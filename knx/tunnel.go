// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package knx

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/knxnet"
	"github.com/vapourismo/knx-go/knx/util"
)

// TunnelConfig allows you to configure the tunnel client's behavior.
type TunnelConfig struct {
	// ResendInterval is the interval with which requests will be resend if no response is received.
	ResendInterval time.Duration

	// HeartbeatInterval specifies the time interval which triggers a heartbeat check.
	HeartbeatInterval time.Duration

	// ResponseTimeout specifies how long to wait for a response.
	ResponseTimeout time.Duration

	// SendLocalAddress specifies if local address should be sent on connection request.
	SendLocalAddress bool
}

// DefaultTunnelConfig is a good default configuration for a Tunnel client.
var DefaultTunnelConfig = TunnelConfig{
	ResendInterval:    500 * time.Millisecond,
	HeartbeatInterval: 10 * time.Second,
	ResponseTimeout:   10 * time.Second,
	SendLocalAddress:  false,
}

// checkTunnelConfig makes sure that the configuration is actually usable.
func checkTunnelConfig(config TunnelConfig) TunnelConfig {
	if config.ResendInterval <= 0 {
		config.ResendInterval = DefaultTunnelConfig.ResendInterval
	}

	if config.HeartbeatInterval <= 0 {
		config.HeartbeatInterval = DefaultTunnelConfig.HeartbeatInterval
	}

	if config.ResponseTimeout <= 0 {
		config.ResponseTimeout = DefaultTunnelConfig.ResponseTimeout
	}

	return config
}

var (
	errResponseTimeout = errors.New("Response timeout reached")
)

// A Tunnel provides methods to communicate with a KNXnet/IP gateway.
type Tunnel struct {
	// Communication methods
	sock   knxnet.Socket
	config TunnelConfig

	// Connection information
	layer   knxnet.TunnelLayer
	channel uint8
	control knxnet.HostInfo

	// For outgoing requests
	seqMu     sync.Mutex
	seqNumber uint8
	ack       chan *knxnet.TunnelRes

	// Incoming requests
	inbound chan cemi.Message

	// Goroutine controller
	done chan struct{}
	once sync.Once
	wait sync.WaitGroup
}

func (conn *Tunnel) hostInfo() (knxnet.HostInfo, error) {
	if conn.config.SendLocalAddress {
		localAddr, err := conn.sock.LocalAddr()

		if err != nil {
			return knxnet.HostInfo{}, err
		}

		return knxnet.HostInfoFromAddress(localAddr)
	} else {
		return knxnet.HostInfo{Protocol: knxnet.UDP4}, nil
	}
}

// requestConn repeatedly sends a connection request through the socket until the configured
// reponse timeout is reached or a response is received. A response that renders the gateway as busy
// will not stop requestConn.
func (conn *Tunnel) requestConn() (err error) {

	hostInfo, err := conn.hostInfo()
	if err != nil {
		return err
	}

	conn.control = hostInfo

	req := &knxnet.ConnReq{
		Layer:   conn.layer,
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
			if res, ok := msg.(*knxnet.ConnRes); ok {
				switch res.Status {
				// Conection has been established.
				case knxnet.NoError:
					conn.channel = res.Channel

					conn.seqMu.Lock()
					conn.seqNumber = 0
					conn.seqMu.Unlock()

					return nil

				// The gateway is busy, but we don't stop yet.
				case knxnet.ErrNoMoreConnections, knxnet.ErrNoMoreUniqueConnections:
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
	heartbeat <-chan knxnet.ErrCode,
) (knxnet.ErrCode, error) {
	req := &knxnet.ConnStateReq{Channel: conn.channel, Status: 0, Control: conn.control}

	// Send first connection state request
	err := conn.sock.Send(req)
	if err != nil {
		return knxnet.ErrConnectionID, err
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
			return knxnet.ErrConnectionID, errResponseTimeout

		// Resend timer fired.
		case <-ticker.C:
			err := conn.sock.Send(req)
			if err != nil {
				return knxnet.ErrConnectionID, err
			}

		// Received a connection state response.
		case res, open := <-heartbeat:
			if !open {
				return knxnet.ErrConnectionID, errors.New("Connection server has terminated")
			}

			return res, nil
		}
	}
}

// requestDisc sends a disconnect request to the gateway.
func (conn *Tunnel) requestDisc() error {
	return conn.sock.Send(&knxnet.DiscReq{
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

	req := &knxnet.TunnelReq{
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
	heartbeat <-chan knxnet.ErrCode,
	timeout chan<- struct{},
) {
	// Request the connction state.
	state, err := conn.requestConnState(heartbeat)
	if err != nil || state != knxnet.NoError {
		if err != nil {
			util.Log(conn, "Error while requesting connection state: %v", err)
		} else {
			util.Log(conn, "Bad connection state: %v", state)
		}

		// Write to timeout as an indication that the heartbeat has failed.
		select {
		case <-conn.done:
		case timeout <- struct{}{}:
		}
	}
}

// handleDiscReq validates the request.
func (conn *Tunnel) handleDiscReq(req *knxnet.DiscReq) error {
	// Validate the request channel.
	if req.Channel != conn.channel {
		return errors.New("Invalid communication channel in disconnect request")
	}

	// We don't need to check if this errors or not. It doesn't matter.
	conn.sock.Send(&knxnet.DiscRes{Channel: req.Channel, Status: 0})

	return nil
}

// handleDiscRes validates the response.
func (conn *Tunnel) handleDiscRes(res *knxnet.DiscRes) error {
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
func (conn *Tunnel) handleTunnelReq(req *knxnet.TunnelReq, seqNumber *uint8) error {
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
	return conn.sock.Send(&knxnet.TunnelRes{
		Channel:   conn.channel,
		SeqNumber: req.SeqNumber,
		Status:    0,
	})
}

// handleTunnelRes validates the response and relays it to a sender that is awaiting an
// acknowledgement.
func (conn *Tunnel) handleTunnelRes(res *knxnet.TunnelRes) error {
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
	res *knxnet.ConnStateRes,
	heartbeat chan<- knxnet.ErrCode,
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
	heartbeat := make(chan knxnet.ErrCode)
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
			case *knxnet.DiscReq:
				err := conn.handleDiscReq(msg)
				if err == nil {
					return errDisconnected
				}

				util.Log(conn, "Error while handling disconnect request %v: %v", msg, err)

			case *knxnet.DiscRes:
				err := conn.handleDiscRes(msg)
				if err == nil {
					return nil
				}

				util.Log(conn, "Error while handling disconnect response %v: %v", msg, err)

			case *knxnet.TunnelReq:
				err := conn.handleTunnelReq(msg, &seqNumber)
				if err != nil {
					util.Log(conn, "Error while handling tunnel request %v: %v", msg, err)
				}

			case *knxnet.TunnelRes:
				err := conn.handleTunnelRes(msg)
				if err != nil {
					util.Log(conn, "Error while handling tunnel response %v: %v", msg, err)
				}

			case *knxnet.ConnStateRes:
				err := conn.handleConnStateRes(msg, heartbeat)
				if err != nil {
					util.Log(
						conn,
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
	util.Log(conn, "Started worker")
	defer util.Log(conn, "Worker exited")

	defer close(conn.ack)
	defer close(conn.inbound)
	defer conn.wait.Done()

	for {
		err := conn.process()

		if err != nil {
			util.Log(conn, "Server terminated with error: %v", err)
		}

		// Check if we can try again.
		if err == errDisconnected || err == errHeartbeatFailed {
			util.Log(conn, "Attempting reconnect")

			reconnErr := conn.requestConn()

			if reconnErr == nil {
				util.Log(conn, "Reconnect succeeded")
				continue
			}

			util.Log(conn, "Reconnect failed: %v", reconnErr)
		}

		return
	}
}

// NewTunnel establishes a connection to a gateway. You can pass a zero initialized ClientConfig;
// the function will take care of filling in the default values.
func NewTunnel(gatewayAddr string, layer knxnet.TunnelLayer, config TunnelConfig) (*Tunnel, error) {
	// Create socket which will be used for communication.
	sock, err := knxnet.DialTunnel(gatewayAddr)
	if err != nil {
		return nil, err
	}

	// Initialize the Client structure.
	client := &Tunnel{
		sock:    sock,
		config:  checkTunnelConfig(config),
		layer:   layer,
		ack:     make(chan *knxnet.TunnelRes),
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

// Close will terminate the connection and wait for the server routine to exit. Although a
// disconnect request is sent, it does not wait for a disconnect response.
func (conn *Tunnel) Close() {
	conn.once.Do(func() {
		conn.requestDisc()

		close(conn.done)
		conn.wait.Wait()

		conn.sock.Close()
	})
}

// Inbound retrieves the channel which transmits incoming data. The channel is closed when the
// underlying Socket closes its inbound channel or when the connection is terminated.
func (conn *Tunnel) Inbound() <-chan cemi.Message {
	return conn.inbound
}

// Send relays a tunnel request to the gateway with the given contents.
func (conn *Tunnel) Send(data cemi.Message) error {
	return conn.requestTunnel(data)
}

// GroupTunnel is a Tunnel that provides only a group communication interface.
type GroupTunnel struct {
	*Tunnel
	inbound chan GroupEvent
}

// NewGroupTunnel creates a new Tunnel for group communication.
func NewGroupTunnel(gatewayAddr string, config TunnelConfig) (gt GroupTunnel, err error) {
	gt.Tunnel, err = NewTunnel(gatewayAddr, knxnet.TunnelLayerData, config)

	if err == nil {
		gt.inbound = make(chan GroupEvent)
		go serveGroupInbound(gt.Tunnel.Inbound(), gt.inbound)
	}

	return
}

// Send a group communication.
func (gt *GroupTunnel) Send(event GroupEvent) error {
	return gt.Tunnel.Send(&cemi.LDataReq{LData: buildGroupOutbound(event)})
}

// Inbound returns the channel on which group communication can be received.
func (gt *GroupTunnel) Inbound() <-chan GroupEvent {
	return gt.inbound
}
