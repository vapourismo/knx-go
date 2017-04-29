package knx

import (
	"errors"
	"time"
)

var (
	ErrConnRejected = errors.New("Gateway rejected connection")
	ErrConnClosed   = errors.New("Socket closed before connection was established")
	ErrConnTimeout  = errors.New("Connection timed out")
)

// Gateway client
type Client struct {
	sock    *Socket
	Inbound <-chan []byte
}

// NewClient establishes a connection to the given gateway.
func NewClient(gatewayAddr string) (*Client, error) {
	sock, err := NewClientSocket(gatewayAddr)
	if err != nil { return nil, err }

	channel, err := establishConnection(sock)
	if err != nil { return nil, err }

	inbound := make(chan []byte, 10)
	heartbeat := make(chan struct{}, 1)

	go clientWorker(sock, channel, inbound, heartbeat)
	go clientHeartbeat(sock, channel, heartbeat)

	return &Client{sock, inbound}, nil
}

// Close terminates the connection.
func (client *Client) Close() {
	client.sock.Close()
}

func attemptConnection(sock *Socket) (byte, error) {
	sock.Send(&ConnectionRequest{})
	timeout := time.After(time.Second)

	for {
		select {
			case <-timeout:
				return 0, ErrConnTimeout

			case msg, open := <-sock.Inbound:
				if !open {
					return 0, ErrConnClosed
				}

				// We are only interested in connection responses
				res, ok := msg.(*ConnectionResponse)
				if !ok { continue }

				if res.Status == 0 {
					return res.Channel, nil
				} else {
					return 0, ErrConnRejected
				}
		}
	}
}

func establishConnection(sock *Socket) (byte, error) {
	var channel byte
	var err error

	for i := 0; i < 5; i++ {
		channel, err = attemptConnection(sock)

		// Unless it timed out, there is no point in continuing
		if err != ErrConnTimeout {
			break
		}
	}

	return channel, err
}

func clientWorker(sock *Socket, channel byte, inbound chan<- []byte, heartbeat chan<- struct{}) {
	for payload := range sock.Inbound {
		switch payload.(type) {
			case *DisconnectRequest:
				sock.Close()
				break

			case *ConnectionStateResponse:
				res := payload.(*ConnectionStateResponse)

				// Make sure we only process the response if it is actually meant for us
				if res.Channel == channel {
					if res.Status == 0 {
						heartbeat <- struct{}{}
					} else {
						sock.Close()
						break
					}
				}

			case *TunnelRequest:
				req := payload.(*TunnelRequest)

				// Make sure we only process the response if it is actually meant for us
				if req.Channel == channel {
					sock.Send(&TunnelResponse{channel, req.SeqNumber, 0})
					inbound <- req.Payload
				}
		}
	}

	close(heartbeat)
	close(inbound)
}

func clientHeartbeat(sock *Socket, channel byte, heartbeat <-chan struct{}) {
	tick := time.Tick(time.Second * 20)

	for {
		select {
			case <-tick:
				if !performHeartbeat(sock, channel, heartbeat) {
					return
				}

			case _, open := <-heartbeat:
				// Discard heartbeats that occur outside of the heartbeat cycle. That way, the
				// gateway can't pre-send the connection state response.

				if !open {
					// Heartbeat channel is closed, which means that the main worker has terminated
					// and with it the client connection.
					return
				}
		}
	}
}

func performHeartbeat(sock *Socket, channel byte, heartbeat <-chan struct{}) bool {
	for i := 0; i < 5; i++ {
		sock.Send(&ConnectionStateRequest{channel, 0, HostInfo{}})

		select {
			case _, open := <-heartbeat:
				return open

			case <-time.After(time.Second):
		}
	}

	// Gateway timed out
	sock.Close()

	return false
}
