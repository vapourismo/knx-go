package knx

import (
	"context"
	"testing"
)

func TestNewConn(t *testing.T) {
	ctx := context.Background()

	// Socket was closed before anything could be done.
	t.Run("SendFails", func (t *testing.T) {
		client, gateway := makeDummySockets()
		defer gateway.Close()

		client.Close()

		_, err := newConn(ctx, client, DefaultClientConfig)
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	// Context is done.
	t.Run("CancelledContext", func (t *testing.T) {
		client, gateway := makeDummySockets()
		defer client.Close()
		defer gateway.Close()

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		_, err := newConn(ctx, client, DefaultClientConfig)
		if err != ctx.Err() {
			t.Fatalf("Expected error %v, got %v", ctx.Err(), err)
		}
	})

	// Socket is closed before first resend.
	t.Run("ResendFails", func (t *testing.T) {
		client, gateway := makeDummySockets()

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			<-gateway.Inbound()

			client.Close()
			gateway.Close()
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			config := DefaultClientConfig
			config.ResendInterval = 1

			_, err := newConn(ctx, client, config)
			if err == nil {
				t.Fatal("Should not succeed")
			}
		})
	})

	// The gateway responds to the connection request.
	t.Run("Resend", func (t *testing.T) {
		client, gateway := makeDummySockets()

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			<-gateway.Inbound()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*ConnectionRequest); ok {
				gateway.sendAny(&ConnectionResponse{1, ConnResOk, req.Control})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			defer client.Close()

			config := DefaultClientConfig
			config.ResendInterval = 1

			_, err := newConn(ctx, client, config)
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	// Inbound channel is closed.
	t.Run("InboundClosed", func (t *testing.T) {
		client, gatway := makeDummySockets()
		defer gatway.Close()
		defer client.Close()

		client.closeIn()

		_, err := newConn(ctx, client, DefaultClientConfig)
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	// The gateway responds to the connection request.
	t.Run("Ok", func (t *testing.T) {
		client, gateway := makeDummySockets()

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*ConnectionRequest); ok {
				gateway.sendAny(&ConnectionResponse{1, ConnResOk, req.Control})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			defer client.Close()

			_, err := newConn(ctx, client, DefaultClientConfig)
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	// The gateway is only busy for the first attempt.
	t.Run("Busy", func (t *testing.T) {
		client, gateway := makeDummySockets()

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*ConnectionRequest); ok {
				gateway.sendAny(&ConnectionResponse{0, ConnResBusy, req.Control})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}

			msg = <-gateway.Inbound()
			if req, ok := msg.(*ConnectionRequest); ok {
				gateway.sendAny(&ConnectionResponse{1, ConnResOk, req.Control})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			defer client.Close()

			config := DefaultClientConfig
			config.ResendInterval = 1

			_, err := newConn(ctx, client, config)
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	// The gateway doesn't supported the requested connection type.
	t.Run("Unsupported", func (t *testing.T) {
		client, gateway := makeDummySockets()

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*ConnectionRequest); ok {
				gateway.sendAny(&ConnectionResponse{0, ConnResUnsupportedType, req.Control})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			defer client.Close()

			_, err := newConn(ctx, client, DefaultClientConfig)
			if err != ConnResUnsupportedType {
				t.Fatalf("Expected error %v, got %v", ConnResUnsupportedType, err)
			}
		})
	})
}

func TestConnHandle_requestState(t *testing.T) {
	ctx := context.Background()

	t.Run("SendFails", func (t *testing.T) {
		client, gateway := makeDummySockets()
		defer gateway.Close()

		client.Close()

		conn := conn{client, DefaultClientConfig, 1}

		err := conn.requestState(ctx, make(chan ConnState))
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("CancelledContext", func (t *testing.T) {
		client, gateway := makeDummySockets()
		defer client.Close()
		defer gateway.Close()

		conn := conn{client, DefaultClientConfig, 1}

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		err := conn.requestState(ctx, make(chan ConnState))
		if err != ctx.Err() {
			t.Fatalf("Expected error %v, got %v", ctx.Err(), err)
		}
	})

	t.Run("ResendFails", func (t *testing.T) {
		client, gateway := makeDummySockets()

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			<-gateway.Inbound()
			client.closeOut()
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			defer client.Close()

			config := DefaultClientConfig
			config.ResendInterval = 1

			conn := conn{client, config, 1}

			err := conn.requestState(ctx, make(chan ConnState))
			if err == nil {
				t.Fatal("Should not succeed")
			}
		})
	})

	t.Run("Resend", func (t *testing.T) {
		client, gateway := makeDummySockets()

		const channel uint8 = 1
		heartbeat := make(chan ConnState)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			<-gateway.Inbound()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*ConnStateReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.Status != 0 {
					t.Error("Non-null request status")
				}

				heartbeat <- ConnStateNormal
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			defer client.Close()

			config := DefaultClientConfig
			config.ResendInterval = 1

			conn := conn{client, config, channel}

			err := conn.requestState(ctx, heartbeat)
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("InboundClosed", func (t *testing.T) {
		client, gateway := makeDummySockets()
		defer client.Close()
		defer gateway.Close()

		heartbeat := make(chan ConnState)
		close(heartbeat)

		conn := conn{client, DefaultClientConfig, 1}

		err := conn.requestState(ctx, heartbeat)
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("Ok", func (t *testing.T) {
		client, gateway := makeDummySockets()

		const channel uint8 = 1
		heartbeat := make(chan ConnState)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*ConnStateReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.Status != 0 {
					t.Error("Non-null request status")
				}

				heartbeat <- ConnStateNormal
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			defer client.Close()

			conn := conn{client, DefaultClientConfig, channel}

			err := conn.requestState(ctx, heartbeat)
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("Inactive", func (t *testing.T) {
		client, gateway := makeDummySockets()

		const channel uint8 = 1
		heartbeat := make(chan ConnState)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*ConnStateReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.Status != 0 {
					t.Error("Non-null request status")
				}

				heartbeat <- ConnStateInactive
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			defer client.Close()

			conn := conn{client, DefaultClientConfig, channel}

			err := conn.requestState(ctx, heartbeat)
			if err != ConnStateInactive {
				t.Fatal(err)
			}
		})
	})
}

func TestConnHandle_requestTunnel(t *testing.T) {
	ctx := context.Background()

	t.Run("SendFails", func (t *testing.T) {
		client, gateway := makeDummySockets()
		defer gateway.Close()

		client.Close()

		conn := conn{client, DefaultClientConfig, 1}

		err := conn.requestTunnel(ctx, 0, []byte{}, make(chan *TunnelRes))
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("ContextCancelled", func (t *testing.T) {
		client, gateway := makeDummySockets()
		defer client.Close()
		defer gateway.Close()

		conn := conn{client, DefaultClientConfig, 1}

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		err := conn.requestTunnel(ctx, 0, []byte{}, make(chan *TunnelRes))
		if err != ctx.Err() {
			t.Fatalf("Expected %v, got %v", ctx.Err(), err)
		}
	})

	t.Run("ResendFails", func (t *testing.T) {
		client, gateway := makeDummySockets()

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			<-gateway.Inbound()
			client.closeOut()
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			defer client.Close()

			config := DefaultClientConfig
			config.ResendInterval = 1

			conn := conn{client, config, 1}

			err := conn.requestTunnel(ctx, 1, []byte{}, make(chan *TunnelRes))
			if err == nil {
				t.Fatal("Should not succeed")
			}
		})
	})

	t.Run("Resend", func (t *testing.T) {
		client, gateway := makeDummySockets()
		ack := make(chan *TunnelRes)

		const (
			channel   uint8 = 1
			seqNumber uint8 = 0
		)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			<-gateway.Inbound()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*TunnelReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != seqNumber {
					t.Error("Mismatching sequence number")
				}

				ack <- &TunnelRes{req.Channel, req.SeqNumber, 0}
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			defer client.Close()

			config := DefaultClientConfig
			config.ResendInterval = 1

			conn := conn{client, config, channel}

			err := conn.requestTunnel(ctx, seqNumber, []byte{}, ack)
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("ClosedAckChannel", func (t *testing.T) {
		client, gateway := makeDummySockets()
		defer client.Close()
		defer gateway.Close()

		ack := make(chan *TunnelRes)
		close(ack)

		conn := conn{client, DefaultClientConfig, 1}

		err := conn.requestTunnel(ctx, 0, []byte{}, ack)
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("InvalidSeqNumber", func (t *testing.T) {
		client, gateway := makeDummySockets()
		ack := make(chan *TunnelRes)
		ctx, cancel := context.WithCancel(ctx)

		const (
			channel   uint8 = 1
			seqNumber uint8 = 0
		)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*TunnelReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != seqNumber {
					t.Error("Mismatching sequence number")
				}

				ack <- &TunnelRes{req.Channel, req.SeqNumber + 10, 0}
				cancel()
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			defer client.Close()

			conn := conn{client, DefaultClientConfig, channel}

			err := conn.requestTunnel(ctx, seqNumber, []byte{}, ack)
			if err != ctx.Err() {
				t.Fatalf("Expected error %v, got %v", ctx.Err(), err)
			}
		})
	})

	t.Run("BadStatus", func (t *testing.T) {
		client, gateway := makeDummySockets()
		ack := make(chan *TunnelRes)

		const (
			channel   uint8 = 1
			seqNumber uint8 = 0
		)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*TunnelReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != seqNumber {
					t.Error("Mismatching sequence number")
				}

				ack <- &TunnelRes{req.Channel, req.SeqNumber, 1}
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			defer client.Close()

			conn := conn{client, DefaultClientConfig, channel}

			err := conn.requestTunnel(ctx, seqNumber, []byte{}, ack)
			if err == nil {
				t.Fatal("Should not succeed")
			}
		})
	})

	t.Run("Ok", func (t *testing.T) {
		client, gateway := makeDummySockets()
		ack := make(chan *TunnelRes)

		const (
			channel   uint8 = 1
			seqNumber uint8 = 0
		)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*TunnelReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != seqNumber {
					t.Error("Mismatching sequence number")
				}

				ack <- &TunnelRes{req.Channel, req.SeqNumber, 0}
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			defer client.Close()

			conn := conn{client, DefaultClientConfig, channel}

			err := conn.requestTunnel(ctx, seqNumber, []byte{}, ack)
			if err != nil {
				t.Fatal(err)
			}
		})
	})
}

func TestConnHandle_handleTunnelRequest(t *testing.T) {
	ctx := context.Background()

	t.Run("InvalidChannel", func (t *testing.T) {
		client, gateway := makeDummySockets()
		defer client.Close()
		defer gateway.Close()

		var seqNumber uint8

		conn := conn{client, DefaultClientConfig, 1}
		req := &TunnelReq{2, 0, []byte{}}

		err := conn.handleTunnelRequest(ctx, req, &seqNumber, make(chan []byte))
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("InvalidSeqNumber", func (t *testing.T) {
		client, gateway := makeDummySockets()

		const (
			channel       uint8 = 1
			sendSeqNumber uint8 = 0
		)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if res, ok := msg.(*TunnelRes); ok {
				if res.Channel != channel {
					t.Error("Mismatching channel")
				}

				if res.SeqNumber != sendSeqNumber {
					t.Error("Mismatching sequence number")
				}

				if res.Status != 0 {
					t.Error("Invalid response status")
				}
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Worker", func (t *testing.T) {
			t.Parallel()

			defer client.Close()

			seqNumber := sendSeqNumber + 1

			conn := conn{client, DefaultClientConfig, channel}
			req := &TunnelReq{channel, sendSeqNumber, []byte{}}

			err := conn.handleTunnelRequest(ctx, req, &seqNumber, make(chan []byte))
			if err != nil {
				t.Fatal(err)
			}

			if seqNumber != sendSeqNumber + 1 {
				t.Error("Sequence number was modified")
			}
		})
	})

	t.Run("Ok", func (t *testing.T) {
		client, gateway := makeDummySockets()
		inbound := make(chan []byte)

		const (
			channel       uint8 = 1
			sendSeqNumber uint8 = 0
		)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if res, ok := msg.(*TunnelRes); ok {
				if res.Channel != channel {
					t.Error("Mismatching channel")
				}

				if res.SeqNumber != sendSeqNumber {
					t.Error("Mismatching sequence number")
				}

				if res.Status != 0 {
					t.Error("Invalid response status")
				}
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Worker", func (t *testing.T) {
			t.Parallel()

			defer client.Close()

			seqNumber := sendSeqNumber

			conn := conn{client, DefaultClientConfig, channel}
			req := &TunnelReq{channel, sendSeqNumber, []byte{}}

			err := conn.handleTunnelRequest(ctx, req, &seqNumber, inbound)
			if err != nil {
				t.Fatal(err)
			}

			if seqNumber != sendSeqNumber + 1 {
				t.Error("Sequence number has not been increased")
			}
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			<-inbound
		})
	})
}

func TestConnHandle_handleTunnelResponse(t *testing.T) {
	ctx := context.Background()

	t.Run("InvalidChannel", func (t *testing.T) {
		client, gateway := makeDummySockets()
		defer client.Close()
		defer gateway.Close()

		conn := conn{client, DefaultClientConfig, 1}

		res := &TunnelRes{2, 0, 0}
		err := conn.handleTunnelResponse(ctx, res, make(chan *TunnelRes))
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("Ok", func (t *testing.T) {
		client, gateway := makeDummySockets()
		ack := make(chan *TunnelRes)

		t.Run("Worker", func (t *testing.T) {
			t.Parallel()

			defer client.Close()
			defer gateway.Close()

			conn := conn{client, DefaultClientConfig, 1}

			res := &TunnelRes{1, 0, 0}
			err := conn.handleTunnelResponse(ctx, res, ack)
			if err != nil {
				t.Fatal(err)
			}
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			res := <-ack

			if res.Channel != 1 {
				t.Error("Mismatching channel")
			}

			if res.SeqNumber != 0 {
				t.Error("Mismatching sequence number")
			}

			if res.Status != 0 {
				t.Error("Non-zero status")
			}
		})
	})
}
