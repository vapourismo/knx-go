package knx

import (
	"context"
	"errors"
	"sync"
	"testing"
)

type dummySocket struct {
	mu      sync.Mutex
	outOpen bool
	inOpen  bool

	out     chan OutgoingPayload
	in      chan interface{}
}

func makeDummySocket() *dummySocket {
	return &dummySocket{
		sync.Mutex{},
		true,
		true,
		make(chan OutgoingPayload, 10),
		make(chan interface{}, 10),
	}
}

func (sock *dummySocket) gatewaySend(payload interface{}) error {
	for {
		sock.mu.Lock()

		if !sock.inOpen {
			sock.mu.Unlock()
			return errors.New("Socket is closed")
		}

		select {
			case sock.in <- payload:
				sock.mu.Unlock()
				return nil

			default:
				sock.mu.Unlock()
		}
	}
}

func (sock *dummySocket) gatewayInbound() <-chan OutgoingPayload {
	return sock.out
}

func (sock *dummySocket) closeOut() {
	sock.mu.Lock()
	defer sock.mu.Unlock()

	if sock.outOpen {
		close(sock.out)
		sock.outOpen = false
	}
}

func (sock *dummySocket) closeIn() {
	sock.mu.Lock()
	defer sock.mu.Unlock()

	if sock.inOpen {
		close(sock.in)
		sock.inOpen = false
	}
}

func (sock *dummySocket) Send(payload OutgoingPayload) error {
	for {
		sock.mu.Lock()

		if !sock.outOpen {
			sock.mu.Unlock()
			return errors.New("Socket is closed")
		}

		select {
			case sock.out <- payload:
				sock.mu.Unlock()
				return nil

			default:
				sock.mu.Unlock()
		}
	}
}

func (sock *dummySocket) Inbound() <-chan interface{} {
	return sock.in
}

func (sock *dummySocket) Close() error {
	sock.mu.Lock()
	defer sock.mu.Unlock()

	if sock.outOpen {
		close(sock.out)
		sock.outOpen = false
	}

	if sock.inOpen {
		close(sock.in)
		sock.inOpen = false
	}

	return nil
}

type gatewayHelper struct {
	ctx  context.Context
	sock *dummySocket
	test *testing.T
}

func (helper *gatewayHelper) receive() OutgoingPayload {
	select {
	case <-helper.ctx.Done():
		helper.test.Fatalf("While waiting for inbound packet: %v", helper.ctx.Err())
		return nil

	case msg, open := <-helper.sock.gatewayInbound():
		if !open {
			helper.test.Fatal("Inbound socket channel is closed")
			return nil
		}

		return msg
	}
}

func (helper *gatewayHelper) ignore() {
	helper.receive()
}

func (helper *gatewayHelper) send(msg interface{}) {
	err := helper.sock.gatewaySend(msg)
	if err != nil {
		helper.test.Fatalf("While responding: %v", err)
	}
}
