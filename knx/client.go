package knx

import (
	"errors"
	"time"
	"context"
)

//
type Client struct {
	sock    *Socket
	reaper  context.CancelFunc

	//
	Inbound <-chan []byte
}

//
var (
	ErrConnTimeout = errors.New("Gateway did not response to connection request")
	ErrConnRejected = errors.New("Gateway rejected connection request")
)

//
func NewClient(gatewayAddress string) (*Client, error) {
	sock, err := NewClientSocket(gatewayAddress)
	if err != nil {
		return nil, err
	}

	req := &ConnectionRequest{}

	err = sock.Send(req)
	if err != nil {
		return nil, err
	}

	resChan := awaitConnectionResponse(sock)

	// Connection cycle
	for i := 0; i < 5; i++ {
		select {
		case res := <-resChan:
			if res.Status == 0 {
				// Connection is established.

				Logger.Printf("Client[%v]: Connection has been established on channel %v",
				              sock.conn.RemoteAddr(), res.Channel)

				return makeClient(sock, res.Channel), nil
			} else {
				// Connection attempt was rejected.

				Logger.Printf("Client[%v]: Connection attempt was rejected",
				              sock.conn.RemoteAddr())

				sock.Close()
				return nil, ErrConnRejected
			}

		case <-time.After(time.Second):
			// Resend the connection request, if we haven't received a response from the gateway
			// after 1 second.
			err := sock.Send(req)
			if err != nil {
				return nil, err
			}
		}
	}

	Logger.Printf("Client[%v]: Connection attempts timed out", sock.conn.RemoteAddr())

	sock.Close()
	return nil, ErrConnTimeout
}

//
func (client *Client) Close() {
	client.reaper()
	client.sock.Close()
}

//
func awaitConnectionResponse(sock *Socket) <-chan *ConnectionResponse {
	resChan := make(chan *ConnectionResponse)

	go func() {
		for payload := range sock.Inbound {
			res, ok := payload.(*ConnectionResponse)
			if ok {
				resChan <- res
				return
			}
		}
	}()

	return resChan
}

//
func makeClient(sock *Socket, channel byte) *Client {
	ctx, reaper := context.WithCancel(context.Background())

	inbound := make(chan []byte)
	go clientInboundWorker(ctx, reaper, sock, channel, inbound)

	return &Client{sock, reaper, inbound}
}

//
func clientInboundWorker(
	ctx     context.Context,
	reaper  context.CancelFunc,
	sock    *Socket,
	channel byte,
	inbound chan<- []byte,
) {
	Logger.Printf("Client[%v]: Started inbound worker", sock.conn.RemoteAddr())
	defer Logger.Printf("Client[%v]: Stopped inbound worker", sock.conn.RemoteAddr())

	defer close(inbound)

	heartbeatTrigger := make(chan struct{})
	stateResChan := make(chan *ConnectionStateResponse)
	go clientHeartbeatWorker(ctx, reaper, sock, channel, heartbeatTrigger, stateResChan)

	for {
		select {
		// Goroutine exit has been requested
		case <-ctx.Done():
			return

		// 10 seconds without communication, time for a heartbeat
		case <-time.After(10 * time.Second):
			Logger.Printf("Client[%v]: Triggering heartbeat", sock.conn.RemoteAddr())

			select {
			case <-ctx.Done():
				return

			case heartbeatTrigger <- struct{}{}:
			}

		// Incoming packets
		case payload, open := <-sock.Inbound:
			// If the socket inbound channel is closed, this goroutine has no purpose.
			if !open {
				Logger.Printf("Client[%v]: Inbound channel has been closed", sock.conn.RemoteAddr())
				reaper()
				return
			}

			switch payload.(type) {
			case *ConnectionStateResponse:
				res := payload.(*ConnectionStateResponse)

				if res.Channel == channel {
					stateResChan <- res
				}
			}
		}
	}
}

//
func clientHeartbeatWorker(
	ctx     context.Context,
	reaper  context.CancelFunc,
	sock    *Socket,
	channel byte,
	trigger <-chan struct{},
	resChan <-chan *ConnectionStateResponse,
) {
	Logger.Printf("Client[%v]: Started heartbeat worker", sock.conn.RemoteAddr())
	defer Logger.Printf("Client[%v]: Stopped heartbeat worker", sock.conn.RemoteAddr())

	// Make sure we tell the others to exit
	defer reaper()

	outerLoop:
	for {
		select {
		// Gorouting has been asked to exit
		case <-ctx.Done():
			return

		// Inbound worker has triggered a heartbeat
		case <-trigger:
			req := &ConnectionStateRequest{channel, 0, HostInfo{}}

			err := sock.Send(req)
			if err != nil {
				Logger.Printf("Client[%v]: Error while sending heartbeat: %v",
				              sock.conn.RemoteAddr(), err)
				return
			}

			// Heartbeat cycle
			for i := 0; i < 5; i++ {
				select {
				case <-ctx.Done():
					return

				case res := <-resChan:
					if res.Status == 0 {
						Logger.Printf("Client[%v]: Heartbeat successful", sock.conn.RemoteAddr())
						continue outerLoop
					} else {
						Logger.Printf("Client[%v]: Gateway rejected heartbeat",
						              sock.conn.RemoteAddr())
						return
					}

				case <-time.After(time.Second):
					err := sock.Send(req)
					if err != nil {
						Logger.Printf("Client[%v]: Error while sending heartbeat: %v",
						              sock.conn.RemoteAddr(), err)
						return
					}
				}
			}

			// We get here, if the gateway did not respond

			Logger.Printf("Client[%v]: Gateway timed out during heartbeat", sock.conn.RemoteAddr())
			return

		case <-resChan:
			// Discard any connection state response that appears out-of-cycle
		}
	}
}
