package knx

import (
	"errors"
	"time"
	"context"
)

//
type Client struct {
	sock    *Socket
	channel byte
	reaper  context.CancelFunc
	seq     <-chan byte
	ack     <-chan byte

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

var (
	ErrSendClosed   = errors.New("Outbound worker has terminated")
	ErrSendRejected = errors.New("Gateway rejected tunnel request")
	ErrSendTimeout  = errors.New("Gateway did not acknowledge tunnel request in time")
)

//
func (client *Client) Send(data []byte) error {
	seqNumber, open := <-client.seq
	if !open {
		return ErrSendClosed
	}

	req := &TunnelRequest{client.channel, seqNumber, data}
	err := client.sock.Send(req)
	if err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		select {
		case status, open := <-client.ack:
			if !open {
				return ErrSendClosed
			}

			if status == 0 {
				return nil
			} else {
				return ErrSendRejected
			}

		case <-time.After(time.Second):
			err := client.sock.Send(req)
			if err != nil {
				return err
			}
		}
	}

	return ErrSendTimeout
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

	seq := make(chan byte)
	ack := make(chan byte)
	tunRes := make(chan *TunnelResponse)
	go clientOutboundWorker(ctx, reaper, tunRes, seq, ack)

	inbound := make(chan []byte)
	go clientInboundWorker(ctx, reaper, sock, channel, tunRes, inbound)

	return &Client{sock, channel, reaper, seq, ack, inbound}
}

//
func clientInboundWorker(
	ctx     context.Context,
	reaper  context.CancelFunc,
	sock    *Socket,
	channel byte,
	tunRes  chan<- *TunnelResponse,
	inbound chan<- []byte,
) {
	Logger.Printf("Client[%v]: Started inbound worker", sock.conn.RemoteAddr())
	defer Logger.Printf("Client[%v]: Stopped inbound worker", sock.conn.RemoteAddr())

	defer close(inbound)
	defer reaper()

	heartbeatTrigger := make(chan struct{})
	stateResChan := make(chan *ConnectionStateResponse)
	go clientHeartbeatWorker(ctx, reaper, sock, channel, heartbeatTrigger, stateResChan)

	var incomingSeqNumber byte = 0

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
				return
			}

			switch payload.(type) {
			case *ConnectionStateResponse:
				res := payload.(*ConnectionStateResponse)

				if res.Channel == channel {
					select {
					case <-ctx.Done():
						return

					case stateResChan <- res:
					}
				}

			case *TunnelRequest:
				req := payload.(*TunnelRequest)

				if req.Channel == channel {
					// Acknowledge tunnel request
					err := sock.Send(&TunnelResponse{channel, req.SeqNumber, 0})
					if err != nil {
						Logger.Printf("Client[%v]: Error while sending tunnel response: %v",
						              sock.conn.RemoteAddr(), err)
						return
					}

					// Relay to user if it fits the sequence
					if req.SeqNumber == incomingSeqNumber {
						Logger.Printf("Client[%v]: Inbound tunnel request: %v",
						              sock.conn.RemoteAddr(), req.Payload)

						select {
						case <-ctx.Done():
							return

						case inbound <- req.Payload:
						}

						incomingSeqNumber++
					}
				}

			case *TunnelResponse:
				res := payload.(*TunnelResponse)

				if res.Channel == channel {
					select {
					case <-ctx.Done():
						return

					case tunRes <- res:
					}
				}
			}
		}
	}
}

//
func clientOutboundWorker(
	ctx     context.Context,
	reaper  context.CancelFunc,
	tunRes  <-chan *TunnelResponse,
	seq     chan<- byte,
	ack     chan<- byte,
) {
	defer reaper()
	defer close(seq)
	defer close(ack)

	var seqNumber byte = 0

	outerLoop:
	for {
		select {
		case <-ctx.Done():
			return

		// Client requests a sequence number because it wants to send something
		case seq <- seqNumber:

			for {
				select {
				case <-ctx.Done():
					return

				// Await tunnel response
				case res := <-tunRes:
					// We're only interested in the ones that match our sequence number
					if res.SeqNumber == seqNumber {
						select {
						case <-ctx.Done():
							return

						// Send result of the tunnel request to the sender
						case ack <- res.Status:
							seqNumber++
							continue outerLoop
						}
					}
				}
			}

		// Discard out-of-cycle tunnel responses
		case <-tunRes:
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
	resChan <-chan *ConnectionStateResponse, // Really? Why not just 'byte'?
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
