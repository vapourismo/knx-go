package knx

import (
	"context"
	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/proto"
	"testing"
)

func TestNewConn(t *testing.T) {
	ctx := context.Background()

	// Socket was closed before anything could be done.
	t.Run("SendFails", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer gateway.Close()

		client.Close()

		_, err := newConn(ctx, client, DefaultClientConfig)
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	// Context is done.
	t.Run("CancelledContext", func(t *testing.T) {
		client, gateway := newDummySockets()
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
	t.Run("ResendFails", func(t *testing.T) {
		client, gateway := newDummySockets()

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			<-gateway.Inbound()

			client.Close()
			gateway.Close()
		})

		t.Run("Client", func(t *testing.T) {
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
	t.Run("Resend", func(t *testing.T) {
		client, gateway := newDummySockets()

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			<-gateway.Inbound()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*proto.ConnReq); ok {
				gateway.sendAny(&proto.ConnRes{1, proto.ConnResOk, req.Control})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
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
	t.Run("InboundClosed", func(t *testing.T) {
		client, gatway := newDummySockets()
		defer gatway.Close()
		defer client.Close()

		client.closeIn()

		_, err := newConn(ctx, client, DefaultClientConfig)
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	// The gateway responds to the connection request.
	t.Run("Ok", func(t *testing.T) {
		client, gateway := newDummySockets()

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*proto.ConnReq); ok {
				gateway.sendAny(&proto.ConnRes{1, proto.ConnResOk, req.Control})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			_, err := newConn(ctx, client, DefaultClientConfig)
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	// The gateway is only busy for the first attempt.
	t.Run("Busy", func(t *testing.T) {
		client, gateway := newDummySockets()

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*proto.ConnReq); ok {
				gateway.sendAny(&proto.ConnRes{0, proto.ConnResBusy, req.Control})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}

			msg = <-gateway.Inbound()
			if req, ok := msg.(*proto.ConnReq); ok {
				gateway.sendAny(&proto.ConnRes{1, proto.ConnResOk, req.Control})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
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
	t.Run("Unsupported", func(t *testing.T) {
		client, gateway := newDummySockets()

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*proto.ConnReq); ok {
				gateway.sendAny(&proto.ConnRes{0, proto.ConnResUnsupportedType, req.Control})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			_, err := newConn(ctx, client, DefaultClientConfig)
			if err != proto.ConnResUnsupportedType {
				t.Fatalf("Expected error %v, got %v", proto.ConnResUnsupportedType, err)
			}
		})
	})
}

func TestConnHandle_requestState(t *testing.T) {
	ctx := context.Background()

	t.Run("SendFails", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer gateway.Close()

		client.Close()

		conn := conn{client, DefaultClientConfig, 1}

		err := conn.requestState(ctx, make(chan proto.ConnState))
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("CancelledContext", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer client.Close()
		defer gateway.Close()

		conn := conn{client, DefaultClientConfig, 1}

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		err := conn.requestState(ctx, make(chan proto.ConnState))
		if err != ctx.Err() {
			t.Fatalf("Expected error %v, got %v", ctx.Err(), err)
		}
	})

	t.Run("ResendFails", func(t *testing.T) {
		client, gateway := newDummySockets()

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			<-gateway.Inbound()
			client.closeOut()
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			config := DefaultClientConfig
			config.ResendInterval = 1

			conn := conn{client, config, 1}

			err := conn.requestState(ctx, make(chan proto.ConnState))
			if err == nil {
				t.Fatal("Should not succeed")
			}
		})
	})

	t.Run("Resend", func(t *testing.T) {
		client, gateway := newDummySockets()

		const channel uint8 = 1
		heartbeat := make(chan proto.ConnState)

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			<-gateway.Inbound()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*proto.ConnStateReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.Status != 0 {
					t.Error("Non-null request status")
				}

				heartbeat <- proto.ConnStateNormal
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
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

	t.Run("InboundClosed", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer client.Close()
		defer gateway.Close()

		heartbeat := make(chan proto.ConnState)
		close(heartbeat)

		conn := conn{client, DefaultClientConfig, 1}

		err := conn.requestState(ctx, heartbeat)
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("Ok", func(t *testing.T) {
		client, gateway := newDummySockets()

		const channel uint8 = 1
		heartbeat := make(chan proto.ConnState)

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*proto.ConnStateReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.Status != 0 {
					t.Error("Non-null request status")
				}

				heartbeat <- proto.ConnStateNormal
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			conn := conn{client, DefaultClientConfig, channel}

			err := conn.requestState(ctx, heartbeat)
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("Inactive", func(t *testing.T) {
		client, gateway := newDummySockets()

		const channel uint8 = 1
		heartbeat := make(chan proto.ConnState)

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*proto.ConnStateReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.Status != 0 {
					t.Error("Non-null request status")
				}

				heartbeat <- proto.ConnStateInactive
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			conn := conn{client, DefaultClientConfig, channel}

			err := conn.requestState(ctx, heartbeat)
			if err != proto.ConnStateInactive {
				t.Fatal(err)
			}
		})
	})
}

func TestConnHandle_requestTunnel(t *testing.T) {
	ctx := context.Background()

	t.Run("SendFails", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer gateway.Close()

		client.Close()

		conn := conn{client, DefaultClientConfig, 1}

		err := conn.requestTunnel(ctx, 0, cemi.CEMI{}, make(chan *proto.TunnelRes))
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("ContextCancelled", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer client.Close()
		defer gateway.Close()

		conn := conn{client, DefaultClientConfig, 1}

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		err := conn.requestTunnel(ctx, 0, cemi.CEMI{}, make(chan *proto.TunnelRes))
		if err != ctx.Err() {
			t.Fatalf("Expected %v, got %v", ctx.Err(), err)
		}
	})

	t.Run("ResendFails", func(t *testing.T) {
		client, gateway := newDummySockets()

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			<-gateway.Inbound()
			client.closeOut()
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			config := DefaultClientConfig
			config.ResendInterval = 1

			conn := conn{client, config, 1}

			err := conn.requestTunnel(ctx, 1, cemi.CEMI{}, make(chan *proto.TunnelRes))
			if err == nil {
				t.Fatal("Should not succeed")
			}
		})
	})

	t.Run("Resend", func(t *testing.T) {
		client, gateway := newDummySockets()
		ack := make(chan *proto.TunnelRes)

		const (
			channel   uint8 = 1
			seqNumber uint8 = 0
		)

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			<-gateway.Inbound()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*proto.TunnelReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != seqNumber {
					t.Error("Mismatching sequence number")
				}

				ack <- &proto.TunnelRes{req.Channel, req.SeqNumber, 0}
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			config := DefaultClientConfig
			config.ResendInterval = 1

			conn := conn{client, config, channel}

			err := conn.requestTunnel(ctx, seqNumber, cemi.CEMI{}, ack)
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("ClosedAckChannel", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer client.Close()
		defer gateway.Close()

		ack := make(chan *proto.TunnelRes)
		close(ack)

		conn := conn{client, DefaultClientConfig, 1}

		err := conn.requestTunnel(ctx, 0, cemi.CEMI{}, ack)
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("InvalidSeqNumber", func(t *testing.T) {
		client, gateway := newDummySockets()
		ack := make(chan *proto.TunnelRes)
		ctx, cancel := context.WithCancel(ctx)

		const (
			channel   uint8 = 1
			seqNumber uint8 = 0
		)

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*proto.TunnelReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != seqNumber {
					t.Error("Mismatching sequence number")
				}

				ack <- &proto.TunnelRes{req.Channel, req.SeqNumber + 10, 0}
				cancel()
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			conn := conn{client, DefaultClientConfig, channel}

			err := conn.requestTunnel(ctx, seqNumber, cemi.CEMI{}, ack)
			if err != ctx.Err() {
				t.Fatalf("Expected error %v, got %v", ctx.Err(), err)
			}
		})
	})

	t.Run("BadStatus", func(t *testing.T) {
		client, gateway := newDummySockets()
		ack := make(chan *proto.TunnelRes)

		const (
			channel   uint8 = 1
			seqNumber uint8 = 0
		)

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*proto.TunnelReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != seqNumber {
					t.Error("Mismatching sequence number")
				}

				ack <- &proto.TunnelRes{req.Channel, req.SeqNumber, 1}
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			conn := conn{client, DefaultClientConfig, channel}

			err := conn.requestTunnel(ctx, seqNumber, cemi.CEMI{}, ack)
			if err == nil {
				t.Fatal("Should not succeed")
			}
		})
	})

	t.Run("Ok", func(t *testing.T) {
		client, gateway := newDummySockets()
		ack := make(chan *proto.TunnelRes)

		const (
			channel   uint8 = 1
			seqNumber uint8 = 0
		)

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*proto.TunnelReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != seqNumber {
					t.Error("Mismatching sequence number")
				}

				ack <- &proto.TunnelRes{req.Channel, req.SeqNumber, 0}
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			conn := conn{client, DefaultClientConfig, channel}

			err := conn.requestTunnel(ctx, seqNumber, cemi.CEMI{}, ack)
			if err != nil {
				t.Fatal(err)
			}
		})
	})
}

func TestConnHandle_handleTunnelRequest(t *testing.T) {
	ctx := context.Background()

	t.Run("InvalidChannel", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer client.Close()
		defer gateway.Close()

		var seqNumber uint8

		conn := conn{client, DefaultClientConfig, 1}
		req := &proto.TunnelReq{2, 0, cemi.CEMI{}}

		err := conn.handleTunnelRequest(ctx, req, &seqNumber, make(chan *cemi.CEMI))
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("InvalidSeqNumber", func(t *testing.T) {
		client, gateway := newDummySockets()

		const (
			channel       uint8 = 1
			sendSeqNumber uint8 = 0
		)

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if res, ok := msg.(*proto.TunnelRes); ok {
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

		t.Run("Worker", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			seqNumber := sendSeqNumber + 1

			conn := conn{client, DefaultClientConfig, channel}
			req := &proto.TunnelReq{channel, sendSeqNumber, cemi.CEMI{}}

			err := conn.handleTunnelRequest(ctx, req, &seqNumber, make(chan *cemi.CEMI))
			if err != nil {
				t.Fatal(err)
			}

			if seqNumber != sendSeqNumber+1 {
				t.Error("Sequence number was modified")
			}
		})
	})

	t.Run("Ok", func(t *testing.T) {
		client, gateway := newDummySockets()
		inbound := make(chan *cemi.CEMI)

		const (
			channel       uint8 = 1
			sendSeqNumber uint8 = 0
		)

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if res, ok := msg.(*proto.TunnelRes); ok {
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

		t.Run("Worker", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			seqNumber := sendSeqNumber

			conn := conn{client, DefaultClientConfig, channel}
			req := &proto.TunnelReq{channel, sendSeqNumber, cemi.CEMI{}}

			err := conn.handleTunnelRequest(ctx, req, &seqNumber, inbound)
			if err != nil {
				t.Fatal(err)
			}

			if seqNumber != sendSeqNumber+1 {
				t.Error("Sequence number has not been increased")
			}
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			<-inbound
		})
	})
}

func TestConnHandle_handleTunnelResponse(t *testing.T) {
	ctx := context.Background()

	t.Run("InvalidChannel", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer client.Close()
		defer gateway.Close()

		conn := conn{client, DefaultClientConfig, 1}

		res := &proto.TunnelRes{2, 0, 0}
		err := conn.handleTunnelResponse(ctx, res, make(chan *proto.TunnelRes))
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("Ok", func(t *testing.T) {
		client, gateway := newDummySockets()
		ack := make(chan *proto.TunnelRes)

		t.Run("Worker", func(t *testing.T) {
			t.Parallel()

			defer client.Close()
			defer gateway.Close()

			conn := conn{client, DefaultClientConfig, 1}

			res := &proto.TunnelRes{1, 0, 0}
			err := conn.handleTunnelResponse(ctx, res, ack)
			if err != nil {
				t.Fatal(err)
			}
		})

		t.Run("Client", func(t *testing.T) {
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
