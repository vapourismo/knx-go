package knx

import (
	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/proto"
)

// serveRouter listens for incoming routing-related packets.
func serveRouter(sock Socket, inbound chan<- cemi.CEMI) {
	defer close(inbound)

	for msg := range sock.Inbound() {
		if ind, ok := msg.(*proto.RoutingInd); ok {
			inbound <- ind.Payload
		}
	}
}

// A Router is a participant in a KNXnet/IP multicast group.
type Router struct {
	sock    Socket
	inbound <-chan cemi.CEMI
}

// NewRouter creates a new Router that joins the given multicast group.
func NewRouter(multicastAddress string) (*Router, error) {
	sock, err := NewRoutingSocket(multicastAddress)
	if err != nil {
		return nil, err
	}

	inbound := make(chan cemi.CEMI)

	go serveRouter(sock, inbound)

	return &Router{
		sock:    sock,
		inbound: inbound,
	}, nil
}

// Send transmits a packet.
func (router *Router) Send(data cemi.CEMI) error {
	return router.sock.Send(&proto.RoutingInd{Payload: data})
}

// Inbound returns the channel which transmits incoming data.
func (router *Router) Inbound() <-chan cemi.CEMI {
	return router.inbound
}

// Close closes the underlying socket and terminates the Router thereby.
func (router *Router) Close() {
	router.sock.Close()
}
