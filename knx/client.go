package knx

import (
	"context"
	"errors"
	"time"
)

// ClientConfig allows you to configure the client's behavior.
type ClientConfig struct {
	// ConnectionTimeout determines how long to wait for a connection response.
	ConnectionTimeout time.Duration

	// ResendInterval is how long to wait for a response, until the request is resend. A interval
	// <= 0 can't be used. The default value will be used instead.
	ResendInterval time.Duration

	// HeartbeatDelay specifies the time which has to elapse without any incoming communication,
	// until a heartbeat is triggered. A delay <= 0 will result in the use of a default value.
	HeartbeatDelay time.Duration

	// HeartbeatTimeout specifies how long to wait for a connection state response. A timeout <= 0
	// will not be accepted. Instead, the default value will be used.
	HeartbeatTimeout time.Duration
}

// Default configuration elements
var (
	defaultConnectionTimeout = 10 * time.Second
	defaultResendInterval    = 500 * time.Millisecond
	defaultHeartbeatDelay    = 10 * time.Second
	defaultHeartbeatTimeout  = 10 * time.Second

	DefaultClientConfig = ClientConfig{
		defaultConnectionTimeout,
		defaultResendInterval,
		defaultHeartbeatDelay,
		defaultHeartbeatTimeout,
	}
)

// checkClientConfig makes sure that the configuration is actually usable.
func checkClientConfig(config ClientConfig) ClientConfig {
	if config.ConnectionTimeout <= 0 {
		config.ConnectionTimeout = defaultConnectionTimeout
	}

	if config.ResendInterval <= 0 {
		config.ResendInterval = defaultResendInterval
	}

	if config.HeartbeatDelay <= 0 {
		config.HeartbeatDelay = defaultHeartbeatDelay
	}

	if config.HeartbeatTimeout <= 0 {
		config.HeartbeatTimeout = defaultHeartbeatTimeout
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
			} else {
				return res
			}
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
	childCtx, cancel := context.WithTimeout(ctx, conn.config.HeartbeatTimeout)
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
	ctx     context.Context,
	inbound chan<- []byte,
) {

	defer close(inbound)

	heartbeat := make(chan ConnState)
	timeout := make(chan struct{})

	var seqNumber uint8 = 0

	for {
		select {
		// Termination has been requested.
		case <-ctx.Done():
			log(conn, "connHandle", "Exiting inbound server due to context error: %v", ctx.Err())
			return

		// Heartbeat worker signals a result.
		case <-timeout:
			log(conn, "connHandle", "Exiting inbound server due to heartbeat timeout")
			return

		// There were no incoming packets for some time.
		case <-time.After(conn.config.HeartbeatDelay):
			go conn.performHeartbeat(ctx, heartbeat, timeout)

		// A message has been received or the channel is closed.
		case msg, open := <-conn.sock.Inbound():
			if !open {
				log(conn, "connHandle", "Exiting inbound server due to closed socket's inbound channel")
				return
			}

			// Determine what to do with the message.
			switch msg.(type) {
			case *TunnelRequest:
				req := msg.(*TunnelRequest)
				err := conn.handleTunnelRequest(ctx, req, &seqNumber, inbound)
				if err != nil {
					log(conn, "connHandle", "Error while handling tunnel request %v: %v", req, err)
				}

			case *ConnectionStateResponse:
				res := msg.(*ConnectionStateResponse)
				err := conn.handleConnectionStateResponse(ctx, res, heartbeat)
				if err != nil {
					log(conn, "connHandle", "Error while handling connection state response: %v", err)
				}
			}
		}
	}
}

//
type Client struct {
	ctx     context.Context
	cancel  context.CancelFunc

	Inbound <-chan []byte
}

//
func NewClient(gatewayAddr string, config ClientConfig) (*Client, error) {
	sock, err := NewClientSocket(gatewayAddr)
	if err != nil {
		return nil, err
	}

	connHandle := connHandle{sock, checkClientConfig(config), 0}

	connectCtx, cancelConnect := context.WithTimeout(context.Background(), config.ConnectionTimeout)
	defer cancelConnect()

	err = connHandle.requestConnection(connectCtx)
	if err != nil {
		return nil, err
	}

	inbound := make(chan []byte)

	ctx, cancel := context.WithCancel(context.Background())
	go connHandle.serveInbound(ctx, inbound)

	return &Client{ctx, cancel, inbound}, nil
}

//
func (client Client) Close() {
	client.cancel()
}
