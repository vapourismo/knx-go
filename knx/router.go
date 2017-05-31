package knx

import (
	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/proto"
)

func serveRouter(sock Socket, inbound chan<- cemi.CEMI) {
	defer close(inbound)

	for msg := range sock.Inbound() {
		if ind, ok := msg.(*proto.RoutingInd); ok {
			inbound <- ind.Payload
		}
	}
}

type Router struct {
	sock    Socket
	inbound <-chan cemi.CEMI
}

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

func (router *Router) Send(data cemi.CEMI) error {
	return router.sock.Send(&proto.RoutingInd{Payload: data})
}

func (router *Router) Inbound() <-chan cemi.CEMI {
	return router.inbound
}

func (router *Router) Close() {
	router.sock.Close()
}
