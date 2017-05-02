package knx

import (
	"errors"
	"time"
)

//
type Client struct {
	sock *Socket

	//
	Inbound <-chan []byte

	//
	Outbound chan<- []byte
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

	sock.Outbound <- &ConnectionRequest{}

	resChan := clientReceiveConnectionResponse(sock)

	for i := 0; i < 5; i++ {
		select {
			case res := <-resChan:
				if res.Status == 0 {
					// Connection is established.

					Logger.Printf("Client[%v]: Connection has been established on channel %v",
					              sock.conn.RemoteAddr(), res.Channel)

					return clientLaunch(sock, res.Channel), nil
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
				sock.Outbound <- &ConnectionRequest{}
		}
	}

	Logger.Printf("Client[%v]: Connection attempts timed out", sock.conn.RemoteAddr())

	sock.Close()
	return nil, ErrConnTimeout
}

//
func clientReceiveConnectionResponse(sock *Socket) <-chan *ConnectionResponse {
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
func clientLaunch(sock *Socket, channel byte) *Client {
	inbound := make(chan []byte)
	outbound := make(chan []byte)

	return &Client{sock, inbound, outbound}
}
