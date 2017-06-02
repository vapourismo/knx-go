package knx

import (
	"container/list"

	"sync"

	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/proto"
)

// A RouterConfig determines certain properties of a Router.
type RouterConfig struct {
	// Specify how many sent messages to retain. This is important for when a router indicates that
	// it has lost some messages. If you do not expect to saturate the router, keep this low.
	RetainCount uint
}

// Default configuration elements
var (
	defaultRetainCount uint = 32

	DefaultRouterConfig = RouterConfig{
		defaultRetainCount,
	}
)

// checkRouterConfig validates the given RouterConfig.
func checkRouterConfig(config RouterConfig) RouterConfig {
	if config.RetainCount == 0 {
		config.RetainCount = defaultRetainCount
	}

	return config
}

// tryPushInbound sends the message through the channel. If the sending blocks, it will launch a
// goroutine which will do the sending.
func tryPushInbound(msg cemi.Message, inbound chan<- cemi.Message) {
	select {
	case inbound <- msg:

	default:
		go func() {
			// Since this goroutine decouples from the server goroutine, it might try to send when
			// the server closed the inbound channel. Sending to a closed channel will panic. But we
			// don't care, because cool guys don't look at explosions.
			defer func() { recover() }()
			inbound <- msg
		}()
	}
}

// A Router is a participant in a KNXnet/IP multicast group.
type Router struct {
	sock     Socket
	config   RouterConfig
	inbound  chan cemi.Message
	sendMu   sync.Mutex
	retainer *list.List
}

// serve listens for incoming routing-related packets.
func (router *Router) serve() {
	defer close(router.inbound)

	for msg := range router.sock.Inbound() {
		switch msg := msg.(type) {
		case *proto.RoutingInd:
			// Try to push it to the client without blocking this goroutine to long.
			tryPushInbound(msg.Payload, router.inbound)

		case *proto.RoutingBusy:
			// TODO: Inhibit sending for msg.WaitTime.

		case *proto.RoutingLost:
			// TODO: Resend the last msg.Count frames.
		}
	}
}

// NewRouter creates a new Router that joins the given multicast group.
func NewRouter(multicastAddress string, config RouterConfig) (*Router, error) {
	sock, err := NewRoutingSocket(multicastAddress)
	if err != nil {
		return nil, err
	}

	r := &Router{
		sock:     sock,
		config:   checkRouterConfig(config),
		inbound:  make(chan cemi.Message),
		retainer: list.New(),
	}

	go r.serve()

	return r, nil
}

// Send transmits a packet.
func (router *Router) Send(data cemi.Message) error {
	// We lock this before doing any sending so the server goroutine can adjust the flow control.
	router.sendMu.Lock()
	defer router.sendMu.Unlock()

	err := router.sock.Send(&proto.RoutingInd{Payload: data})

	if err == nil {
		// Store this for potential resending.
		router.retainer.PushBack(data)

		// We don't want to keep more messages than necessary. The overhead needs to be removed.
		for uint(router.retainer.Len()) > router.config.RetainCount {
			router.retainer.Remove(router.retainer.Front())
		}
	}

	return err
}

// Inbound returns the channel which transmits incoming data.
func (router *Router) Inbound() <-chan cemi.Message {
	return router.inbound
}

// Close closes the underlying socket and terminates the Router thereby.
func (router *Router) Close() {
	router.sock.Close()
}
