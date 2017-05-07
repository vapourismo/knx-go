package knx

import (
	"context"
	"testing"
)

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
