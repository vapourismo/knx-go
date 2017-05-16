package knx

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ClientConfig allows you to configure the client's behavior.
type ClientConfig struct {
	// ResendInterval is how long to wait for a response, until the request is resend. A interval
	// <= 0 can't be used. The default value will be used instead.
	ResendInterval time.Duration

	// HeartbeatDelay specifies the time which has to elapse without any incoming communication,
	// until a heartbeat is triggered. A delay <= 0 will result in the use of a default value.
	HeartbeatDelay time.Duration

	// ResponseTimeout specifies how long to wait for a response. A timeout <= 0 will not be
	// accepted. Instead, the default value will be used.
	ResponseTimeout time.Duration
}

// Default configuration elements
var (
	defaultResendInterval    = 500 * time.Millisecond
	defaultHeartbeatDelay    = 10 * time.Second
	defaultHeartbeatTimeout  = 10 * time.Second

	DefaultClientConfig = ClientConfig{
		defaultResendInterval,
		defaultHeartbeatDelay,
		defaultHeartbeatTimeout,
	}
)

// checkClientConfig makes sure that the configuration is actually usable.
func checkClientConfig(config ClientConfig) ClientConfig {
	if config.ResendInterval <= 0 {
		config.ResendInterval = defaultResendInterval
	}

	if config.HeartbeatDelay <= 0 {
		config.HeartbeatDelay = defaultHeartbeatDelay
	}

	if config.ResponseTimeout <= 0 {
		config.ResponseTimeout = defaultHeartbeatTimeout
	}

	return config
}

// connHandle is a handle for the client connection.
type connHandle struct {
	sock    Socket
	config  ClientConfig
	channel uint8
}

// requestConnection sends a connection request every 500ms through the socket until the provided
// context gets canceled, or a response is received. A response that renders the gateway as busy
// will not stop requestConnection.
func (conn *connHandle) requestConnection(ctx context.Context) error {
	req := &ConnectionRequest{}

	// Send the initial request.
	err := conn.sock.Send(req)
	if err != nil {
		return err
	}

	// Create a resend timer.
	ticker := time.NewTicker(conn.config.ResendInterval)
	defer ticker.Stop()

	// Cycle until a request gets a response.
	for {
		select {
		// Termination has been requested.
		case <-ctx.Done():
			return ctx.Err()

		// Resend timer triggered.
		case <-ticker.C:
			err := conn.sock.Send(req)
			if err != nil {
				return err
			}

		// A message has been received or the channel has been closed.
		case msg, open := <-conn.sock.Inbound():
			if !open {
				return errors.New("Inbound channel has been closed")
			}

			// We're only interested in connection responses.
			if res, ok := msg.(*ConnectionResponse); ok {
				switch res.Status {
				// Conection has been established.
				case ConnResOk:
					conn.channel = res.Channel
					return nil

				// The gateway is busy, but we don't stop yet.
				case ConnResBusy:
					continue

				// Connection request has been denied.
				default:
					return res.Status
				}
			}
		}
	}
}

// requestConnectionState periodically sends a connection state request to the gateway until it has
// received a response, the context is done, or HeartbeatDelay duration has passed.
func (conn *connHandle) requestConnectionState(
	ctx       context.Context,
	heartbeat <-chan ConnState,
) error {
	req := &ConnectionStateRequest{conn.channel, 0, HostInfo{}}

	// Send first connection state request
	err := conn.sock.Send(req)
	if err != nil {
		return err
	}

	// Start the resend timer.
	ticker := time.NewTicker(conn.config.ResendInterval)
	defer ticker.Stop()

	for {
		select {
		// Termination has been requested.
		case <-ctx.Done():
			return ctx.Err()

		// Resend timer fired.
		case <-ticker.C:
			err := conn.sock.Send(req)
			if err != nil {
				return err
			}

		// Received a connection state response.
		case res, open := <-heartbeat:
			if !open {
				return errors.New("Heartbeat channel is closed")
			}

			// Is connection state positive?
			if res == ConnStateNormal {
				return nil
			}

			return res
		}
	}
}

//
func (conn *connHandle) requestTunnel(
	ctx       context.Context,
	seqNumber uint8,
	data      []byte,
	ack       <-chan *TunnelResponse,
) error {
	req := &TunnelRequest{conn.channel, seqNumber, data}

	// Send initial request.
	err := conn.sock.Send(req)
	if err != nil {
		return err
	}

	// Start the resend timer.
	ticker := time.NewTicker(conn.config.ResendInterval)
	defer ticker.Stop()

	for {
		select {
		// Termination has been requested.
		case <-ctx.Done():
			return ctx.Err()

		// Resend timer fired.
		case <-ticker.C:
			err := conn.sock.Send(req)
			if err != nil {
				return err
			}

		// Received a tunnel response.
		case res, open := <-ack:
			if !open {
				return errors.New("Ack channel is closed")
			}

			// Ignore mismatching sequence numbers.
			if res.SeqNumber != seqNumber {
				continue
			}

			// Check if the response confirms the tunnel request.
			if res.Status == 0 {
				return nil
			}

			return fmt.Errorf("Tunnel request has been rejected with status %#x", res.Status)
		}
	}
}

// performHeartbeat uses requestConnectionState to determine if the gateway is still alive.
func (conn *connHandle) performHeartbeat(
	ctx       context.Context,
	heartbeat <-chan ConnState,
	timeout   chan<- struct{},
) {
	// Setup a child context which will time out with the given heartbeat timeout.
	childCtx, cancel := context.WithTimeout(ctx, conn.config.ResponseTimeout)
	defer cancel()

	// Request the connction state.
	err := conn.requestConnectionState(childCtx, heartbeat)
	if err != nil {
		log(conn, "connHandle", "Error while requesting connection state: %v", err)

		// Write to timeout as an indication that the heartbeat has failed.
		select {
		case <-ctx.Done():
		case timeout <- struct{}{}:
		}
	}
}

// handleDisconnectRequest validates the request.
func (conn *connHandle) handleDisconnectRequest(
	ctx context.Context,
	req *DisconnectRequest,
) error {
	// Validate the request channel.
	if req.Channel != conn.channel {
		return errors.New("Invalid communication channel in disconnect request")
	}

	return nil
}

// handleTunnelRequest validates the request, pushes the data to the client and acknowledges the
// request for the gateway.
func (conn *connHandle) handleTunnelRequest(
	ctx       context.Context,
	req       *TunnelRequest,
	seqNumber *uint8,
	inbound   chan<- []byte,
) error {
	// Validate the request channel.
	if req.Channel != conn.channel {
		return errors.New("Invalid communication channel in tunnel request")
	}

	// Is the sequence number what we expected?
	if req.SeqNumber == *seqNumber {
		*seqNumber++

		// Send tunnel data to the client.
		go func () {
			select {
			case <-ctx.Done():
			case inbound <- req.Payload:
			}
		}()
	}

	// Send the acknowledgement.
	return conn.sock.Send(&TunnelResponse{conn.channel, req.SeqNumber, 0})
}

// handleTunnelResponse validates the response and relays it to a sender that is awaiting an
// acknowledgement.
func (conn *connHandle) handleTunnelResponse(
	ctx context.Context,
	res *TunnelResponse,
	ack chan<- *TunnelResponse,
) error {
	// Validate the request channel.
	if res.Channel != conn.channel {
		return errors.New("Invalid communication channel in connection state response")
	}

	// Send to client.
	go func () {
		select {
		case <-ctx.Done():
		case <-time.After(conn.config.ResendInterval):
		case ack <- res:
		}
	}()

	return nil
}

// handleConnectionStateResponse validates the response and sends it to the heartbeat routine, if
// there is a waiting one.
func (conn *connHandle) handleConnectionStateResponse(
	ctx       context.Context,
	res       *ConnectionStateResponse,
	heartbeat chan<- ConnState,
) error {
	// Validate the request channel.
	if res.Channel != conn.channel {
		return errors.New("Invalid communication channel in connection state response")
	}

	// Send connection state to the heartbeat goroutine.
	go func () {
		select {
		case <-ctx.Done():
		case <-time.After(conn.config.ResendInterval):
		case heartbeat <- res.Status:
		}
	}()

	return nil
}

// serveInbound processes incoming packets.
func (conn *connHandle) serveInbound(
	ctx      context.Context,
	inbound  chan<- []byte,
	ack      chan<- *TunnelResponse,
) error {
	defer close(ack)
	defer close(inbound)

	heartbeat := make(chan ConnState)
	timeout := make(chan struct{})

	var seqNumber uint8

	for {
		select {
		// Termination has been requested.
		case <-ctx.Done():
			return ctx.Err()

		// Heartbeat worker signals a result.
		case <-timeout:
			return errors.New("Heartbeat did not succeed")

		// There were no incoming packets for some time.
		case <-time.After(conn.config.HeartbeatDelay):
			go conn.performHeartbeat(ctx, heartbeat, timeout)

		// A message has been received or the channel is closed.
		case msg, open := <-conn.sock.Inbound():
			if !open {
				return errors.New("Socket's inbound channel is closed")
			}

			// Determine what to do with the message.
			switch msg.(type) {
			case *DisconnectRequest:
				req := msg.(*DisconnectRequest)

				err := conn.handleDisconnectRequest(ctx, req)
				if err == nil {
					return nil
				}

				log(conn, "connHandle", "Error while handling disconnect request %v: %v", req, err)

			case *TunnelRequest:
				req := msg.(*TunnelRequest)

				err := conn.handleTunnelRequest(ctx, req, &seqNumber, inbound)
				if err != nil {
					log(conn, "connHandle", "Error while handling tunnel request %v: %v", req, err)
				}

			case *TunnelResponse:
				res := msg.(*TunnelResponse)

				err := conn.handleTunnelResponse(ctx, res, ack)
				if err != nil {
					log(conn, "connHandle", "Error while handling tunnel response %v: %v", res, err)
				}

			case *ConnectionStateResponse:
				res := msg.(*ConnectionStateResponse)

				err := conn.handleConnectionStateResponse(ctx, res, heartbeat)
				if err != nil {
					log(conn, "connHandle",
					    "Error while handling connection state response: %v", err)
				}
			}
		}
	}
}

// Client represents the client endpoint in a connection with a gateway.
type Client struct {
	ctx       context.Context
	cancel    context.CancelFunc

	conn      *connHandle

	mu        sync.Mutex
	seqNumber uint8
	ack       chan *TunnelResponse

	inbound   chan []byte
}

// Connect establishes a connection with a gateway.
func Connect(gatewayAddr string, config ClientConfig) (*Client, error) {
	// Create socket which will be used for communication.
	sock, err := NewClientSocket(gatewayAddr)
	if err != nil {
		return nil, err
	}

	// Initialize the connection handle.
	conn := &connHandle{sock, checkClientConfig(config), 0}

	// Prepare a context, so that the connection request cannot run forever.
	connectCtx, cancelConnect := context.WithTimeout(context.Background(), config.ResponseTimeout)
	defer cancelConnect()

	// Connect to the gateway.
	err = conn.requestConnection(connectCtx)
	if err != nil {
		return nil, err
	}

	// Prepare a context for the inbound server.
	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		ctx,
		cancel,
		conn,
		sync.Mutex{},
		0,
		make(chan *TunnelResponse),
		make(chan []byte),
	}, nil
}

// Serve starts the internal connection server, which is needed to process incoming packets.
func (client *Client) Serve() error {
	return client.conn.serveInbound(client.ctx, client.inbound, client.ack)
}

// Close will terminate the connection.
func (client *Client) Close() {
	client.cancel()
}

// Inbound retrieves the channel which transmits incoming data.
func (client *Client) Inbound() <-chan []byte {
	return client.inbound
}

// Send relays a tunnel request to the gateway with the given contents.
func (client *Client) Send(data []byte) error {
	// Establish a lock so that nobody else can modify the sequence number.
	client.mu.Lock()
	defer client.mu.Unlock()

	// Prepare a context, so that we won't wait forever for a tunnel response.
	ctx, cancel := context.WithTimeout(client.ctx, client.conn.config.ResponseTimeout)
	defer cancel()

	// Send the tunnel reqest.
	err := client.conn.requestTunnel(ctx, client.seqNumber, data, client.ack)
	if err != nil {
		return err
	}

	// We are able to increase the sequence number of success.
	client.seqNumber++

	return nil
}
