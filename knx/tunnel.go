// Copyright 2017 Ole Krüger.

package knx

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/proto"
)

// TunnelConfig allows you to configure the client's behavior.
type TunnelConfig struct {
	// ResendInterval is how long to wait for a response, until the request is resend. A interval
	// <= 0 can't be used. The default value will be used instead.
	ResendInterval time.Duration

	// HeartbeatInterval specifies the time which has to elapse without any incoming communication,
	// until a heartbeat is triggered. A delay <= 0 will result in the use of a default value.
	HeartbeatInterval time.Duration

	// ResponseTimeout specifies how long to wait for a response. A timeout <= 0 will not be
	// accepted. Instead, the default value will be used.
	ResponseTimeout time.Duration
}

// Default configuration elements
var (
	defaultResendInterval    = 500 * time.Millisecond
	defaultHeartbeatInterval = 10 * time.Second
	defaultResponseTimeout   = 10 * time.Second

	DefaultTunnelConfig = TunnelConfig{
		defaultResendInterval,
		defaultHeartbeatInterval,
		defaultResponseTimeout,
	}
)

// checkTunnelConfig makes sure that the configuration is actually usable.
func checkTunnelConfig(config TunnelConfig) TunnelConfig {
	if config.ResendInterval <= 0 {
		config.ResendInterval = defaultResendInterval
	}

	if config.HeartbeatInterval <= 0 {
		config.HeartbeatInterval = defaultHeartbeatInterval
	}

	if config.ResponseTimeout <= 0 {
		config.ResponseTimeout = defaultResponseTimeout
	}

	return config
}

var (
	errResponseTimeout = errors.New("Response timeout reached")
)

// Tunnel is a handle for a tunnel connection.
type Tunnel struct {
	// Communication methods
	sock   Socket
	config TunnelConfig

	// Connection information
	channel uint8
	control proto.HostInfo

	// For outgoing requests
	seqMu     sync.Mutex
	seqNumber uint8
	ack       chan *proto.TunnelRes

	// Incoming requests
	inbound chan cemi.Message

	// Goroutine controller
	done chan struct{}
	once sync.Once
	wait sync.WaitGroup
}

// requestConn repeatedly sends a connection request through the socket until the configured
// reponse timeout is reached or a response is received. A response that renders the gateway as busy
// will not stop requestConn.
func (conn *Tunnel) requestConn() (err error) {
	conn.control = proto.HostInfo{Protocol: proto.UDP4}

	req := &proto.ConnReq{
		Layer:   proto.TunnelLayerData,
		Control: conn.control,
		Tunnel:  conn.control,
	}

	// Send the initial request.
	err = conn.sock.Send(req)
	if err != nil {
		return
	}

	// Create a resend timer.
	ticker := time.NewTicker(conn.config.ResendInterval)
	defer ticker.Stop()

	// Setup timeout.
	timeout := time.After(conn.config.ResponseTimeout)

	// Cycle until a request gets a response.
	for {
		select {
		// Timeout reached.
		case <-timeout:
			return errResponseTimeout

		// Resend timer triggered.
		case <-ticker.C:
			err = conn.sock.Send(req)
			if err != nil {
				return
			}

		// A message has been received or the channel has been closed.
		case msg, open := <-conn.sock.Inbound():
			if !open {
				return errors.New("Socket's inbound channel has been closed")
			}

			// We're only interested in connection responses.
			if res, ok := msg.(*proto.ConnRes); ok {
				switch res.Status {
				// Conection has been established.
				case proto.NoError:
					conn.channel = res.Channel

					conn.seqMu.Lock()
					conn.seqNumber = 0
					conn.seqMu.Unlock()

					return nil

				// The gateway is busy, but we don't stop yet.
				case proto.ErrNoMoreConnections, proto.ErrNoMoreUniqueConnections:
					continue

				// Connection request has been denied.
				default:
					return res.Status
				}
			}
		}
	}
}

// requestConnState periodically sends a connection state request to the gateway until it has
// received a response or the response timeout is reached.
func (conn *Tunnel) requestConnState(
	heartbeat <-chan proto.ErrCode,
) (proto.ErrCode, error) {
	req := &proto.ConnStateReq{Channel: conn.channel, Status: 0, Control: conn.control}

	// Send first connection state request
	err := conn.sock.Send(req)
	if err != nil {
		return proto.ErrConnectionID, err
	}

	// Start the resend timer.
	ticker := time.NewTicker(conn.config.ResendInterval)
	defer ticker.Stop()

	// Setup timeout timer.
	timeout := time.After(conn.config.ResponseTimeout)

	for {
		select {
		// Reached timeout
		case <-timeout:
			return proto.ErrConnectionID, errResponseTimeout

		// Resend timer fired.
		case <-ticker.C:
			err := conn.sock.Send(req)
			if err != nil {
				return proto.ErrConnectionID, err
			}

		// Received a connection state response.
		case res, open := <-heartbeat:
			if !open {
				return proto.ErrConnectionID, errors.New("Connection server has terminated")
			}

			return res, nil
		}
	}
}

// requestDisc sends a disconnect request to the gateway.
func (conn *Tunnel) requestDisc() error {
	return conn.sock.Send(&proto.DiscReq{
		Channel: conn.channel,
		Status:  0,
		Control: conn.control,
	})
}

// requestTunnel sends a tunnel request to the gateway and waits for an appropriate acknowledgement.
func (conn *Tunnel) requestTunnel(data cemi.Message) error {
	// Sequence numbers cannot be reused, therefore we must protect against that.
	conn.seqMu.Lock()
	defer conn.seqMu.Unlock()

	req := &proto.TunnelReq{
		Channel:   conn.channel,
		SeqNumber: conn.seqNumber,
		Payload:   data,
	}

	// Send initial request.
	err := conn.sock.Send(req)
	if err != nil {
		return err
	}

	// Start the resend timer.
	ticker := time.NewTicker(conn.config.ResendInterval)
	defer ticker.Stop()

	// Setup timeout.
	timeout := time.After(conn.config.ResponseTimeout)

	for {
		select {
		// Timeout reached.
		case <-timeout:
			return errResponseTimeout

		// Resend timer fired.
		case <-ticker.C:
			err := conn.sock.Send(req)
			if err != nil {
				return err
			}

		// Received a tunnel response.
		case res, open := <-conn.ack:
			if !open {
				return errors.New("Connection server has terminated")
			}

			// Ignore mismatching sequence numbers.
			if res.SeqNumber != conn.seqNumber {
				continue
			}

			// Gateway has received the request, therefore we can increase on our side.
			conn.seqNumber++

			// Check if the response confirms the tunnel request.
			if res.Status == 0 {
				return nil
			}

			return fmt.Errorf("tunnelConn request has been rejected with status %#x", res.Status)
		}
	}
}

// performHeartbeat uses requestConnState to determine if the gateway is still alive.
func (conn *Tunnel) performHeartbeat(
	heartbeat <-chan proto.ErrCode,
	timeout chan<- struct{},
) {
	// Request the connction state.
	state, err := conn.requestConnState(heartbeat)
	if err != nil || state != proto.NoError {
		if err != nil {
			log(conn, "conn", "Error while requesting connection state: %v", err)
		} else {
			log(conn, "conn", "Bad connection state: %v", state)
		}

		// Write to timeout as an indication that the heartbeat has failed.
		select {
		case <-conn.done:
		case timeout <- struct{}{}:
		}
	}
}

// handleDiscReq validates the request.
func (conn *Tunnel) handleDiscReq(req *proto.DiscReq) error {
	// Validate the request channel.
	if req.Channel != conn.channel {
		return errors.New("Invalid communication channel in disconnect request")
	}

	// We don't need to check if this errors or not. It doesn't matter.
	conn.sock.Send(&proto.DiscRes{Channel: req.Channel, Status: 0})

	return nil
}

// handleDiscRes validates the response.
func (conn *Tunnel) handleDiscRes(res *proto.DiscRes) error {
	// Validate the response channel.
	if res.Channel != conn.channel {
		return errors.New("Invalid communication channel in disconnect response")
	}

	return nil
}

// pushInbound sends the message through the inbound channel. If the sending blocks, it will launch
// a goroutine which will do the sending.
func (conn *Tunnel) pushInbound(msg cemi.Message) {
	select {
	case conn.inbound <- msg:

	default:
		go func() {
			// Since this goroutine decouples from the server goroutine, it might try to send when
			// the server closed the inbound channel. Sending to a closed channel will panic. But we
			// don't care, because cool guys don't look at explosions.
			defer func() { recover() }()
			conn.inbound <- msg
		}()
	}
}

// handleTunnelReq validates the request, pushes the data to the client and acknowledges the
// request for the gateway.
func (conn *Tunnel) handleTunnelReq(req *proto.TunnelReq, seqNumber *uint8) error {
	// Validate the request channel.
	if req.Channel != conn.channel {
		return errors.New("Invalid communication channel in tunnel request")
	}

	expected := *seqNumber

	// Is the sequence number what we expected?
	if req.SeqNumber == expected {
		*seqNumber++

		// Send tunnel data to the client without blocking this goroutine to long.
		conn.pushInbound(req.Payload)
	} else if req.SeqNumber != expected-1 {
		// The sequence number is out of the range which we would have to acknowledge.
		return errors.New("Out of sequence tunnel acknowledgement")
	}

	// Send the acknowledgement.
	return conn.sock.Send(&proto.TunnelRes{
		Channel:   conn.channel,
		SeqNumber: req.SeqNumber,
		Status:    0,
	})
}

// handleTunnelRes validates the response and relays it to a sender that is awaiting an
// acknowledgement.
func (conn *Tunnel) handleTunnelRes(res *proto.TunnelRes) error {
	// Validate the request channel.
	if res.Channel != conn.channel {
		return errors.New("Invalid communication channel in connection state response")
	}

	// Send to client.
	go func() {
		// Ack channel might be closed, but we don't care. Just catch the panic that occurs when
		// writing to a closed channel here, and be done with it.
		defer func() { recover() }()

		select {
		case <-conn.done:
		case <-time.After(conn.config.ResendInterval):
		case conn.ack <- res:
		}
	}()

	return nil
}

// handleConnStateRes validates the response and sends it to the heartbeat routine, if there is a
// waiting one.
func (conn *Tunnel) handleConnStateRes(
	res *proto.ConnStateRes,
	heartbeat chan<- proto.ErrCode,
) error {
	// Validate the request channel.
	if res.Channel != conn.channel {
		return errors.New("Invalid communication channel in connection state response")
	}

	// Send connection state to the heartbeat goroutine.
	go func() {
		// Heartbeat channel might be closed, but we don't care. Just catch the panic that occurs
		// when writing to a closed channel here, and be done with it.
		defer func() { recover() }()

		select {
		case <-conn.done:
		case <-time.After(conn.config.ResendInterval):
		case heartbeat <- res.Status:
		}
	}()

	return nil
}

var (
	errHeartbeatFailed = errors.New("Heartbeat did not succeed")
	errInboundClosed   = errors.New("Socket's inbound channel is closed")
	errDisconnected    = errors.New("Gateway terminated the connection")
)

// process incoming packets.
func (conn *Tunnel) process() error {
	heartbeat := make(chan proto.ErrCode)
	defer close(heartbeat)

	timeout := make(chan struct{})

	var seqNumber uint8

	heartbeatInterval := time.NewTicker(conn.config.HeartbeatInterval)
	defer heartbeatInterval.Stop()

	for {
		select {
		// Termination has been requested.
		case <-conn.done:
			return nil

		// Heartbeat worker signals a result.
		case <-timeout:
			return errHeartbeatFailed

		// Heartbeat check is due.
		case <-heartbeatInterval.C:
			go conn.performHeartbeat(heartbeat, timeout)

		// A message has been received or the channel is closed.
		case msg, open := <-conn.sock.Inbound():
			if !open {
				return errInboundClosed
			}

			// Determine what to do with the message.
			switch msg := msg.(type) {
			case *proto.DiscReq:
				err := conn.handleDiscReq(msg)
				if err == nil {
					return errDisconnected
				}

				log(conn, "conn", "Error while handling disconnect request %v: %v", msg, err)

			case *proto.DiscRes:
				err := conn.handleDiscRes(msg)
				if err == nil {
					return nil
				}

				log(conn, "conn", "Error while handling disconnect response %v: %v", msg, err)

			case *proto.TunnelReq:
				err := conn.handleTunnelReq(msg, &seqNumber)
				if err != nil {
					log(conn, "conn", "Error while handling tunnel request %v: %v", msg, err)
				}

			case *proto.TunnelRes:
				err := conn.handleTunnelRes(msg)
				if err != nil {
					log(conn, "conn", "Error while handling tunnel response %v: %v", msg, err)
				}

			case *proto.ConnStateRes:
				err := conn.handleConnStateRes(msg, heartbeat)
				if err != nil {
					log(
						conn, "conn",
						"Error while handling connection state response: %v", err,
					)
				}
			}
		}
	}
}

// serve serves the tunnel connection. It can sustain certain failures. This method will try to
// reconnect in case of a heartbeat failure or disconnect.
func (conn *Tunnel) serve() {
	defer close(conn.ack)
	defer close(conn.inbound)
	defer conn.wait.Done()

	for {
		err := conn.process()

		if err != nil {
			log(conn, "conn", "Server terminated with error: %v", err)
		}

		// Check if we can try again.
		if err == errDisconnected || err == errHeartbeatFailed {
			log(conn, "conn", "Attempting reconnect")

			reconnErr := conn.requestConn()

			if reconnErr == nil {
				log(conn, "conn", "Reconnect succeeded")
				continue
			}

			log(conn, "conn", "Reconnect failed: %v", reconnErr)
		}

		return
	}
}

// NewTunnel establishes a connection to a gateway. You can pass a zero initialized ClientConfig;
// the function will take care of filling in the default values.
func NewTunnel(gatewayAddr string, config TunnelConfig) (*Tunnel, error) {
	// Create socket which will be used for communication.
	sock, err := NewTunnelSocket(gatewayAddr)
	if err != nil {
		return nil, err
	}

	// Initialize the Client structure.
	client := &Tunnel{
		sock:    sock,
		config:  checkTunnelConfig(config),
		ack:     make(chan *proto.TunnelRes),
		inbound: make(chan cemi.Message),
		done:    make(chan struct{}),
	}

	// Connect to the gateway.
	err = client.requestConn()
	if err != nil {
		sock.Close()
		return nil, err
	}

	client.wait.Add(1)
	go client.serve()

	return client, nil
}

// Close will terminate the connection and wait for the server routine to exit.
func (conn *Tunnel) Close() {
	conn.once.Do(func() {
		conn.requestDisc()

		close(conn.done)
		conn.wait.Wait()

		conn.sock.Close()
	})
}

// Inbound retrieves the channel which transmits incoming data.
func (conn *Tunnel) Inbound() <-chan cemi.Message {
	return conn.inbound
}

// Send relays a tunnel request to the gateway with the given contents.
func (conn *Tunnel) Send(data cemi.Message) error {
	// Send the tunnel reqest.
	return conn.requestTunnel(data)
}
