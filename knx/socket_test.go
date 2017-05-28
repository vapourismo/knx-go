package knx

import (
	"container/list"
	"errors"
	"sync"
	"github.com/vapourismo/knx-go/knx/proto"
)

type dummySocket struct {
	cond    *sync.Cond
	out     *list.List
	in      *list.List
}

func makeDummySockets() (*dummySocket, *dummySocket) {
	cond := sync.NewCond(&sync.Mutex{})
	forGateway := list.New()
	forClient := list.New()

	client := &dummySocket{cond, forGateway, forClient}
	gateway := &dummySocket{cond, forClient, forGateway}

	return client, gateway
}

func (sock *dummySocket) Receive() (proto.Service, error) {
	sock.cond.L.Lock()

	for sock.in != nil && sock.in.Len() < 1 {
		sock.cond.Wait()
	}

	if sock.in == nil {
		sock.cond.L.Unlock()
		return nil, errors.New("Input is closed")
	}

	val := sock.in.Remove(sock.in.Front())

	sock.cond.Broadcast()
	sock.cond.L.Unlock()

	return val.(proto.Service), nil
}

func (sock *dummySocket) Send(payload proto.ServiceWriterTo) error {
	return sock.sendAny(payload)
}

func (sock *dummySocket) sendAny(payload interface{}) error {
	sock.cond.L.Lock()
	defer sock.cond.L.Unlock()

	if sock.out == nil {
		return errors.New("Outbound is closed")
	}

	sock.out.PushBack(payload)
	sock.cond.Broadcast()

	return nil
}

func (sock *dummySocket) Close() error {
	sock.cond.L.Lock()
	defer sock.cond.L.Unlock()

	sock.in = nil
	sock.out = nil

	sock.cond.Broadcast()

	return nil
}

func (sock *dummySocket) closeIn() {
	sock.cond.L.Lock()
	defer sock.cond.L.Unlock()

	sock.in = nil

	sock.cond.Broadcast()
}

func (sock *dummySocket) closeOut() {
	sock.cond.L.Lock()
	defer sock.cond.L.Unlock()

	sock.out = nil

	sock.cond.Broadcast()
}
