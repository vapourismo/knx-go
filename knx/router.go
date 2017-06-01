package knx

import (
	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/proto"
)

// tryPushInbound sends the message through the channel. If the sending blocks, it will launch a
// goroutine which will do the sending.
func tryPushInbound(msg cemi.Message, inbound chan<- cemi.Message) {
	select {
	case inbound <- msg:

	default:
		go func() { inbound <- msg }()
	}
}

// serveRouter listens for incoming routing-related packets.
func serveRouter(sock Socket, inbound chan<- cemi.Message) {
	defer close(inbound)

	for msg := range sock.Inbound() {
		switch msg := msg.(type) {
		case *proto.RoutingInd:
			tryPushInbound(msg.Payload, inbound)

		case *proto.RoutingBusy:
			// TODO: Inhibit sending for msg.WaitTime.

		case *proto.RoutingLost:
			// TODO: Resend the last msg.Count frames.
		}
	}
}

// A Router is a participant in a KNXnet/IP multicast group.
type Router struct {
	sock    Socket
	inbound <-chan cemi.Message
}

// NewRouter creates a new Router that joins the given multicast group.
func NewRouter(multicastAddress string) (*Router, error) {
	sock, err := NewRoutingSocket(multicastAddress)
	if err != nil {
		return nil, err
	}

	inbound := make(chan cemi.Message)

	go serveRouter(sock, inbound)

	return &Router{
		sock:    sock,
		inbound: inbound,
	}, nil
}

// Send transmits a packet.
func (router *Router) Send(data cemi.Message) error {
	return router.sock.Send(&proto.RoutingInd{Payload: data})
}

// Inbound returns the channel which transmits incoming data.
func (router *Router) Inbound() <-chan cemi.Message {
	return router.inbound
}

// Close closes the underlying socket and terminates the Router thereby.
func (router *Router) Close() {
	router.sock.Close()
}
