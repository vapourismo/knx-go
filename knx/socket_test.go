// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package knx

import (
	"container/list"
	"errors"
	"net"
	"sync"

	"github.com/vapourismo/knx-go/knx/knxnet"
)

type dummySocket struct {
	cond    *sync.Cond
	out     *list.List
	in      *list.List
	inbound chan knxnet.Service
}

func (sock *dummySocket) serveOne() bool {
	sock.cond.L.Lock()

	for sock.in != nil && sock.in.Len() < 1 {
		sock.cond.Wait()
	}

	if sock.in == nil {
		sock.cond.L.Unlock()
		return false
	}

	val := sock.in.Remove(sock.in.Front())

	sock.cond.Broadcast()
	sock.cond.L.Unlock()

	sock.inbound <- val.(knxnet.Service)

	return true
}

func (sock *dummySocket) serveAll() {
	for sock.serveOne() {
	}
	close(sock.inbound)
}

func (sock *dummySocket) sendAny(payload knxnet.Service) error {
	sock.cond.L.Lock()
	defer sock.cond.L.Unlock()

	if sock.out == nil {
		return errors.New("Outbound is closed")
	}

	sock.out.PushBack(payload)
	sock.cond.Broadcast()

	return nil
}

func (sock *dummySocket) Send(payload knxnet.ServicePackable) error {
	return sock.sendAny(payload)
}

func (sock *dummySocket) Inbound() <-chan knxnet.Service {
	return sock.inbound
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

func (sock *dummySocket) Close() error {
	sock.cond.L.Lock()
	defer sock.cond.L.Unlock()

	sock.in = nil
	sock.out = nil

	sock.cond.Broadcast()

	return nil
}

func (sock *dummySocket) LocalAddr() net.Addr {
	return &net.UDPAddr{
		IP:   net.IPv4(192, 168, 1, 82),
		Port: 4321,
	}
}

func newDummySockets() (*dummySocket, *dummySocket) {
	cond := sync.NewCond(&sync.Mutex{})
	forGateway := list.New()
	forClient := list.New()

	client := &dummySocket{cond, forGateway, forClient, make(chan knxnet.Service)}
	go client.serveAll()

	gateway := &dummySocket{cond, forClient, forGateway, make(chan knxnet.Service)}
	go gateway.serveAll()

	return client, gateway
}
