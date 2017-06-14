// Copyright 2017 Ole Krüger.

package knx

import (
	"container/list"
	"errors"
	"sync"
	"time"

	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/knxnet"
	"github.com/vapourismo/knx-go/knx/util"
)

// A RouterConfig determines certain properties of a Router.
type RouterConfig struct {
	// Specify how many sent messages to retain. This is important for when a router indicates that
	// it has lost some messages. If you do not expect to saturate the router, keep this low.
	RetainCount uint
}

// DefaultRouterConfig is a good default configuration for a Router client.
var DefaultRouterConfig = RouterConfig{
	RetainCount: 32,
}

// checkRouterConfig validates the given RouterConfig.
func checkRouterConfig(config RouterConfig) RouterConfig {
	if config.RetainCount == 0 {
		config.RetainCount = DefaultRouterConfig.RetainCount
	}

	return config
}

// A Router provides the means to communicate with KNXnet/IP routers in a IP multicast group.
// It supports sending and receiving CEMI-encoded frames, aswell as basic flow control.
type Router struct {
	sock     knxnet.Socket
	config   RouterConfig
	inbound  chan cemi.Message
	sendMu   sync.Mutex
	retainer *list.List
}

// sendMultiple sends each message from the slice. Doesn't matter if one fails, all will be tried.
func (router *Router) sendMultiple(messages []cemi.Message) {
	for _, message := range messages {
		router.Send(message)
	}
}

// resendLost resends the last count messages.
func (router *Router) resendLost(count uint16) {
	router.sendMu.Lock()
	defer router.sendMu.Unlock()

	// Make sure not to overflow our retainer list.
	if int(count) > router.retainer.Len() {
		count = uint16(router.retainer.Len())
	}

	messages := make([]cemi.Message, count)

	// Retrieve the messages in reverse. This enables us to resend them in the order in which the
	// have been sent initially.
	for i := len(messages) - 1; i >= 0; i-- {
		messages[i] = router.retainer.Remove(router.retainer.Back()).(cemi.Message)
	}

	go router.sendMultiple(messages)
}

// pushInbound sends the message through the inbound channel. If the sending blocks, it will launch
// a goroutine which will do the sending.
func (router *Router) pushInbound(msg cemi.Message) {
	select {
	case router.inbound <- msg:

	default:
		go func() {
			// Since this goroutine decouples from the server goroutine, it might try to send when
			// the server closed the inbound channel. Sending to a closed channel will panic. But we
			// don't care, because cool guys don't look at explosions.
			defer func() { recover() }()
			router.inbound <- msg
		}()
	}
}

// serve listens for incoming routing-related packets.
func (router *Router) serve() {
	util.Log(router, "Started worker")
	defer util.Log(router, "Worker exited")

	defer close(router.inbound)

	for msg := range router.sock.Inbound() {
		switch msg := msg.(type) {
		case *knxnet.RoutingInd:
			// Try to push it to the client without blocking this goroutine to long.
			router.pushInbound(msg.Payload)

		case *knxnet.RoutingBusy:
			// Inhibit sending for the given time.
			router.sendMu.Lock()
			time.AfterFunc(msg.WaitTime, router.sendMu.Unlock)

			// TODO: Slow down pace after busy indication.

		case *knxnet.RoutingLost:
			// Resend the last msg.Count messages.
			router.resendLost(msg.Count)
		}
	}
}

// NewRouter creates a new Router that joins the given multicast group. You may pass a
// zero-initialized value as parameter config, the default values will be set up.
func NewRouter(multicastAddress string, config RouterConfig) (*Router, error) {
	sock, err := knxnet.NewMulticastSocket(multicastAddress)
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
	if data == nil {
		return errors.New("Nil-pointers are not sendable")
	}

	// We lock this before doing any sending so the server goroutine can adjust the flow control.
	router.sendMu.Lock()
	defer router.sendMu.Unlock()

	err := router.sock.Send(&knxnet.RoutingInd{Payload: data})

	if err == nil {
		// Store this for potential resending.
		// TODO: Ensure that the retained value is independent from the parameter, i.e. not modified
		//       when the user changes a member of data.
		router.retainer.PushBack(data)

		// We don't want to keep more messages than necessary. The overhead needs to be removed.
		for uint(router.retainer.Len()) > router.config.RetainCount {
			router.retainer.Remove(router.retainer.Front())
		}
	}

	return err
}

// Inbound returns the channel which transmits incoming data. The channel will be closed when the
// underlying Socket closes its inbound channel (which happens on read errors or upon closing it).
func (router *Router) Inbound() <-chan cemi.Message {
	return router.inbound
}

// Close closes the underlying socket and terminates the Router thereby.
func (router *Router) Close() {
	router.sock.Close()
}

// GroupRouter is a Router that provides only a group communication interface.
type GroupRouter struct {
	*Router
	inbound chan GroupComm
}

// NewGroupRouter creates a new Router for group communication.
func NewGroupRouter(multicastAddress string, config RouterConfig) (gr GroupRouter, err error) {
	gr.Router, err = NewRouter(multicastAddress, config)

	if err == nil {
		gr.inbound = make(chan GroupComm)
		go gr.serve()
	}

	return
}

func (gr *GroupRouter) serve() {
	util.Log(gr, "Started worker")
	defer util.Log(gr, "Worker exited")

	for msg := range gr.Router.Inbound() {
		if ind, ok := msg.(*cemi.LDataInd); ok {
			if app, ok := ind.Data.(*cemi.AppData); ok && (app.Command == cemi.GroupValueResponse || app.Command == cemi.GroupValueWrite) {
				gr.inbound <- GroupComm{
					Source:      ind.Source,
					Destination: ind.Destination,
					Data:        app.Data,
				}
			} else {
				util.Log(gr, "Received L_Data.ind frame does not contain application data")
			}
		} else {
			util.Log(gr, "Received frame is not a L_Data.ind frame")
		}
	}

	close(gr.inbound)
}

var defaultInd = cemi.LDataInd{
	LData: cemi.LData{
		Control1: cemi.Control1NoRepeat | cemi.Control1NoSysBroadcast | cemi.Control1WantAck | cemi.Control1Prio(cemi.PrioLow),
		Control2: cemi.Control2GrpAddr | cemi.Control2Hops(6),
	},
}

// Send a group communication.
func (gr *GroupRouter) Send(comm GroupComm) error {
	ind := defaultInd
	ind.Data = &cemi.AppData{Command: cemi.GroupValueWrite, Data: comm.Data}
	ind.Source = comm.Source
	ind.Destination = comm.Destination

	if len(comm.Data) <= 15 {
		ind.Control1 |= cemi.Control1StdFrame
	}

	return gr.Router.Send(&ind)
}

// Inbound returns the channel on which group communication can be received.
func (gr *GroupRouter) Inbound() <-chan GroupComm {
	return gr.inbound
}
