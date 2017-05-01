package knx

import (
	"errors"
	"time"
)

// These errors can occur during a connection attempt.
var (
	ErrConnRejected = errors.New("Gateway rejected connection")
	ErrConnClosed   = errors.New("Socket closed before connection was established")
	ErrConnTimeout  = errors.New("Connection timed out")
)

// Gateway client
type Client struct {
	sock     *Socket
	Inbound  <-chan []byte
	Outbound chan<- []byte
}

// NewClient establishes a connection to the given gateway.
func NewClient(gatewayAddr string) (*Client, error) {
	sock, err := NewClientSocket(gatewayAddr)
	if err != nil { return nil, err }

	channel, err := attemptConnection(sock)
	if err != nil { return nil, err }

	inbound := make(chan []byte, 10)
	outbound := make(chan []byte, 10)

	go clientReceiver(sock, channel, inbound, outbound)

	return &Client{sock, inbound, outbound}, nil
}

// Send transmits data via a tunnel request.
func (client *Client) Send(data []byte) {
	client.Outbound <- data
}

// Close terminates the connection.
func (client *Client) Close() {
	client.sock.Close()
}

func attemptConnection(sock *Socket) (byte, error) {
	for i := 0; i < 5; i++ {
		err := sock.Send(&ConnectionRequest{})
		if err != nil { return 0, err }

		select {
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

			case <-time.After(time.Second):
		}
	}

	return 0, ErrConnTimeout
}

func clientReceiver(sock *Socket, channel byte, inbound chan<- []byte, outbound <-chan []byte) {
	Logger.Printf("tunnel[%v]: Started main worker", sock.conn.RemoteAddr())

	reaper, scythe := makeReaper()

	heartbeat := make(chan *ConnectionStateResponse, 1)
	ack := make(chan *TunnelResponse, 1)

	go clientHeartbeat(sock, channel, heartbeat, scythe)
	go clientSender(sock, channel, outbound, ack, scythe)

	var seqNumber byte = 0

	loop:
	for {
		select {
			case <-reaper:
				break loop

			case payload, open := <-sock.Inbound:
				if !open || !clientProcessPayload(sock, channel, heartbeat, ack, inbound,
				                                  &seqNumber, payload) {
					break loop
				}
		}
	}

	close(heartbeat)
	close(ack)
	close(inbound)

	Logger.Printf("tunnel[%v]: Stopped main worker", sock.conn.RemoteAddr())
}

func clientProcessPayload(
	sock      *Socket,
	channel   byte,
	heartbeat chan<- *ConnectionStateResponse,
	ack       chan<- *TunnelResponse,
	inbound   chan<- []byte,
	seqNumber *byte,
	payload   interface{},
) bool {
	switch payload.(type) {
		case *DisconnectRequest:
			req := payload.(*DisconnectRequest)

			if req.Channel == channel {
				Logger.Printf("tunnel[%v]: Received disconnect request", sock.conn.RemoteAddr())
				return false
			}

		case *ConnectionStateResponse:
			res := payload.(*ConnectionStateResponse)

			// Make sure we only process the response if it is meant for us
			if res.Channel == channel {
				heartbeat <- res
			}

		case *TunnelRequest:
			req := payload.(*TunnelRequest)

			// Make sure we only process the request if it is meant for us
			if req.Channel == channel {
				sock.Send(&TunnelResponse{channel, req.SeqNumber, 0})

				if *seqNumber == req.SeqNumber {
					Logger.Printf("tunnel[%v]: Received tunnel request: seqNumber=%v, data=%v",
					              sock.conn.RemoteAddr(), req.SeqNumber, req.Payload)

					inbound <- req.Payload
					(*seqNumber)++
				} else {
					Logger.Printf(
						"tunnel[%v]: Received out-of-sequence tunnel request: seqNumber=%v (expected  %v)",
						sock.conn.RemoteAddr(), req.SeqNumber, *seqNumber,
					)
				}
			}

		case *TunnelResponse:
			res := payload.(*TunnelResponse)

			// Make sure we only process the request if it is meant for us
			if res.Channel == channel {
				ack <- res
			}
	}

	return true
}

var (
	errHeartbeatTimeout  = errors.New("Gateway timed out during heartbeat check")
	errHeartbeatRejected = errors.New("Gateway closed the connection")
	errHeartbeatClosed   = errors.New("Heartbeat channel is closed")
)

func clientHeartbeat(
	sock       *Socket,
	channel    byte,
	heartbeat  <-chan *ConnectionStateResponse,
	killParent scythe,
) {
	Logger.Printf("tunnel[%v]: Started heartbeat worker", sock.conn.RemoteAddr())
	defer Logger.Printf("tunnel[%v]: Stopped heartbeat worker", sock.conn.RemoteAddr())

	tick := time.Tick(time.Second * 20)

	for {
		select {
			case <-tick:
				Logger.Printf("tunnel[%v]: Starting heartbeat cycle", sock.conn.RemoteAddr())
				err := clientPerformHeartbeat(sock, channel, heartbeat)

				switch err {
					// Main worker terminated
					case errHeartbeatClosed:
						return

					// Gateway died or closed the connection
					case errHeartbeatTimeout, errHeartbeatRejected:
						Logger.Printf("tunnel[%v]: Error during heartbeat: %v",
						              sock.conn.RemoteAddr(), err)

						// Tell main worker to shutdown
						killParent()

						return

					// Nothing
					case nil:
						Logger.Printf("tunnel[%v]: Heartbeat succeeded", sock.conn.RemoteAddr())

					// Unknown error
					default:
						Logger.Printf("tunnel[%v]: Error during heartbeat: %v",
						              sock.conn.RemoteAddr(), err)
				}

			case _, open := <-heartbeat:
				// Discard heartbeats that occur outside of the heartbeat cycle. That way, the
				// gateway can't pre-send the connection state response. Also the connection cannot
				// be terminated by a random connection state response.

				if !open {
					// Main worker has terminated
					return
				} else {
					Logger.Printf("tunnel[%v]: Encountered out of cycle heartbeat",
					              sock.conn.RemoteAddr())
				}
		}
	}
}

func clientPerformHeartbeat(
	sock      *Socket,
	channel   byte,
	heartbeat <-chan *ConnectionStateResponse,
) error {
	req := &ConnectionStateRequest{channel, 0, HostInfo{}}

	for i := 0; i < 5; i++ {
		err := sock.Send(req)
		if err != nil { return err }

		select {
			case res, open := <-heartbeat:
				switch {
					// Main worker terminated
					case !open:
						return errHeartbeatClosed

					// Successful heartbeat
					case res.Status == 0:
						return nil

					// Gateway closed the connection
					default:
						return errHeartbeatRejected
				}

			case <-time.After(time.Second):
		}
	}

	return errHeartbeatTimeout
}

var (
	errSendTimeout  = errors.New("Gateway did not acknowledge tunnel request")
	errSendClosed   = errors.New("Acknowledge channel is closed")
	errSendWrongAck = errors.New("Gateway sent incorrect tunnel response")
	errSendRejected = errors.New("Gateway rejected the tunnel request")
)

func clientSender(
	sock       *Socket,
	channel    byte,
	outbound   <-chan []byte,
	ack        <-chan *TunnelResponse,
	killParent scythe,
) {
	Logger.Printf("tunnel[%v]: Started send worker", sock.conn.RemoteAddr())
	defer Logger.Printf("tunnel[%v]: Stopped send worker", sock.conn.RemoteAddr())

	var seqNumber byte = 0

	for {
		select {
			case data, open := <-outbound:
				if !open {
					// Sender channel has been closed
					return
				}

				Logger.Printf("tunnel[%v]: Sending tunnel request: seqNumber=%v, data=%v",
				              sock.conn.RemoteAddr(), seqNumber, data)

				err := clientPerformSend(sock, &TunnelRequest{channel, seqNumber, data}, ack)
				seqNumber++

				switch err {
					// Main worker terminated
					case errSendClosed:
						return

					// Gateway timed out
					case errSendTimeout:
						Logger.Printf("tunnel[%v]: Error during send: %v",
						              sock.conn.RemoteAddr(), err)

						// Tell main worker to shut down
						killParent()

						return

					// Nothing
					case nil:
						Logger.Printf("tunnel[%v]: Send succeeded", sock.conn.RemoteAddr())

					// Other errors
					default:
						Logger.Printf("tunnel[%v]: Error during send: %v",
						              sock.conn.RemoteAddr(), err)
				}

			case _, open := <-ack:
				// Discard random tunnel responses

				if !open {
					// Main worker has terminated
					return
				} else {
					Logger.Printf("tunnel[%v]: Encountered out of cycle tunnel response",
					              sock.conn.RemoteAddr())
				}
		}
	}
}

func clientPerformSend(sock *Socket, req *TunnelRequest, ack <-chan *TunnelResponse) error {
	for i := 0; i < 5; i++ {
		err := sock.Send(req)
		if err != nil { return err }

		select {
			case res, open := <-ack:
				switch {
					// Main worker has terminated
					case !open:
						return errSendClosed

					// Non-sequential SeqNumber
					case res.SeqNumber != req.SeqNumber:
						return errSendWrongAck

					// Successful acknowledgement
					case res.Status == 0:
						return nil

					// Gateway rejected the tunnel request
					default:
						return errSendRejected
				}

			case <-time.After(time.Second):
		}
	}

	return errSendTimeout
}
