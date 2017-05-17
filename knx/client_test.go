package knx

import (
	"context"
	"testing"
)

func TestNewConn(t *testing.T) {
	ctx := context.Background()

	// Socket was closed before anything could be done.
	t.Run("SendFails", func (t *testing.T) {
		sock := makeDummySocket()
		sock.Close()

		_, err := newConn(ctx, sock, DefaultClientConfig)
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	// Context is done.
	t.Run("CancelledContext", func (t *testing.T) {
		sock := makeDummySocket()
		defer sock.Close()

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		_, err := newConn(ctx, sock, DefaultClientConfig)
		if err != ctx.Err() {
			t.Fatalf("Expected error %v, got %v", ctx.Err(), err)
		}
	})

	// Socket is closed before first resend.
	t.Run("ResendFails", func (t *testing.T) {
		sock := makeDummySocket()

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}
			gw.ignore()

			sock.closeOut()
		})

		t.Run("Client", func (t *testing.T) {
			defer sock.Close()
			t.Parallel()

			config := DefaultClientConfig
			config.ResendInterval = 1

			_, err := newConn(ctx, sock, config)
			if err == nil {
				t.Fatal("Should not succeed")
			}
		})
	})

	// The gateway responds to the connection request.
	t.Run("Resend", func (t *testing.T) {
		sock := makeDummySocket()

		const channel uint8 = 1

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			gw.ignore()

			msg := gw.receive()
			if req, ok := msg.(*ConnectionRequest); ok {
				gw.send(&ConnectionResponse{channel, ConnResOk, req.Control})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			defer sock.Close()
			t.Parallel()

			config := DefaultClientConfig
			config.ResendInterval = 1

			conn, err := newConn(ctx, sock, config)
			if err != nil {
				t.Fatal(err)
			}

			if conn.channel != channel {
				t.Error("Mismatching channel")
			}
		})
	})

	// Inbound channel is closed.
	t.Run("InboundClosed", func (t *testing.T) {
		sock := makeDummySocket()
		sock.closeIn()
		defer sock.Close()

		_, err := newConn(ctx, sock, DefaultClientConfig)
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	// The gateway responds to the connection request.
	t.Run("Ok", func (t *testing.T) {
		sock := makeDummySocket()

		const channel uint8 = 1

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			msg := gw.receive()
			if req, ok := msg.(*ConnectionRequest); ok {
				gw.send(&ConnectionResponse{channel, ConnResOk, req.Control})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			defer sock.Close()
			t.Parallel()

			conn, err := newConn(ctx, sock, DefaultClientConfig)
			if err != nil {
				t.Fatal(err)
			}

			if conn.channel != channel {
				t.Error("Mismatching channel")
			}
		})
	})

	// The gateway is only busy for the first attempt.
	t.Run("Busy", func (t *testing.T) {
		sock := makeDummySocket()

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			msg := gw.receive()
			if req, ok := msg.(*ConnectionRequest); ok {
				gw.send(&ConnectionResponse{0, ConnResBusy, req.Control})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}

			msg = gw.receive()
			if req, ok := msg.(*ConnectionRequest); ok {
				gw.send(&ConnectionResponse{1, ConnResOk, req.Control})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			defer sock.Close()
			t.Parallel()

			config := DefaultClientConfig
			config.ResendInterval = 1

			_, err := newConn(ctx, sock, config)
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	// The gateway doesn't supported the requested connection type.
	t.Run("Unsupported", func (t *testing.T) {
		sock := makeDummySocket()

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			msg := gw.receive()
			if req, ok := msg.(*ConnectionRequest); ok {
				gw.send(&ConnectionResponse{0, ConnResUnsupportedType, req.Control})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			defer sock.Close()
			t.Parallel()

			_, err := newConn(ctx, sock, DefaultClientConfig)
			if err != ConnResUnsupportedType {
				t.Fatalf("Expected error %v, got %v", ConnResUnsupportedType, err)
			}
		})
	})
}

func TestConnHandle_requestState(t *testing.T) {
	ctx := context.Background()

	t.Run("SendFails", func (t *testing.T) {
		sock := makeDummySocket()
		sock.Close()

		conn := conn{sock, DefaultClientConfig, 1}

		err := conn.requestState(ctx, make(chan ConnState))
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("CancelledContext", func (t *testing.T) {
		sock := makeDummySocket()
		defer sock.Close()

		conn := conn{sock, DefaultClientConfig, 1}

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		err := conn.requestState(ctx, make(chan ConnState))
		if err != ctx.Err() {
			t.Fatalf("Expected error %v, got %v", ctx.Err(), err)
		}
	})

	t.Run("ResendFails", func (t *testing.T) {
		sock := makeDummySocket()

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}
			gw.ignore()

			sock.closeOut()
		})

		t.Run("Client", func (t *testing.T) {
			defer sock.Close()
			t.Parallel()

			config := DefaultClientConfig
			config.ResendInterval = 1

			conn := conn{sock, config, 1}

			err := conn.requestState(ctx, make(chan ConnState))
			if err == nil {
				t.Fatal("Should not succeed")
			}
		})
	})

	t.Run("Resend", func (t *testing.T) {
		sock := makeDummySocket()

		const channel uint8 = 1
		heartbeat := make(chan ConnState)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			gw.ignore()

			msg := gw.receive()
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
			defer sock.Close()
			t.Parallel()

			config := DefaultClientConfig
			config.ResendInterval = 1

			conn := conn{sock, config, channel}

			err := conn.requestState(ctx, heartbeat)
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("InboundClosed", func (t *testing.T) {
		sock := makeDummySocket()

		heartbeat := make(chan ConnState)
		close(heartbeat)

		conn := conn{sock, DefaultClientConfig, 1}

		err := conn.requestState(ctx, heartbeat)
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("Ok", func (t *testing.T) {
		sock := makeDummySocket()

		const channel uint8 = 1
		heartbeat := make(chan ConnState)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			msg := gw.receive()
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
			defer sock.Close()
			t.Parallel()

			conn := conn{sock, DefaultClientConfig, channel}

			err := conn.requestState(ctx, heartbeat)
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("Inactive", func (t *testing.T) {
		sock := makeDummySocket()

		const channel uint8 = 1
		heartbeat := make(chan ConnState)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			msg := gw.receive()
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
			defer sock.Close()
			t.Parallel()

			conn := conn{sock, DefaultClientConfig, channel}

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
		sock := makeDummySocket()
		sock.Close()

		conn := conn{sock, DefaultClientConfig, 1}

		err := conn.requestTunnel(ctx, 0, []byte{}, make(chan *TunnelRes))
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("ContextCancelled", func (t *testing.T) {
		sock := makeDummySocket()
		defer sock.Close()

		conn := conn{sock, DefaultClientConfig, 1}

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		err := conn.requestTunnel(ctx, 0, []byte{}, make(chan *TunnelRes))
		if err != ctx.Err() {
			t.Fatalf("Expected %v, got %v", ctx.Err(), err)
		}
	})

	t.Run("ResendFails", func (t *testing.T) {
		sock := makeDummySocket()

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			gw.ignore()
			sock.closeOut()
		})

		t.Run("Client", func (t *testing.T) {
			defer sock.Close()
			t.Parallel()

			config := DefaultClientConfig
			config.ResendInterval = 1

			conn := conn{sock, config, 1}

			err := conn.requestTunnel(ctx, 1, []byte{}, make(chan *TunnelRes))
			if err == nil {
				t.Fatal("Should not succeed")
			}
		})
	})

	t.Run("Resend", func (t *testing.T) {
		sock := makeDummySocket()
		ack := make(chan *TunnelRes)

		const (
			channel   uint8 = 1
			seqNumber uint8 = 0
		)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			gw.ignore()

			msg := gw.receive()
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
			defer sock.Close()
			t.Parallel()

			config := DefaultClientConfig
			config.ResendInterval = 1

			conn := conn{sock, config, channel}

			err := conn.requestTunnel(ctx, seqNumber, []byte{}, ack)
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("ClosedAckChannel", func (t *testing.T) {
		sock := makeDummySocket()
		defer sock.Close()

		ack := make(chan *TunnelRes)
		close(ack)

		conn := conn{sock, DefaultClientConfig, 1}

		err := conn.requestTunnel(ctx, 0, []byte{}, ack)
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("InvalidSeqNumber", func (t *testing.T) {
		sock := makeDummySocket()
		ack := make(chan *TunnelRes)

		ctx, cancel := context.WithCancel(ctx)

		const (
			channel   uint8 = 1
			seqNumber uint8 = 0
		)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			msg := gw.receive()
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
			defer sock.Close()
			t.Parallel()

			conn := conn{sock, DefaultClientConfig, channel}

			err := conn.requestTunnel(ctx, seqNumber, []byte{}, ack)
			if err != ctx.Err() {
				t.Fatalf("Expected error %v, got %v", ctx.Err(), err)
			}
		})
	})

	t.Run("BadStatus", func (t *testing.T) {
		sock := makeDummySocket()
		ack := make(chan *TunnelRes)

		const (
			channel   uint8 = 1
			seqNumber uint8 = 0
		)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			msg := gw.receive()
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
			defer sock.Close()
			t.Parallel()

			conn := conn{sock, DefaultClientConfig, channel}

			err := conn.requestTunnel(ctx, seqNumber, []byte{}, ack)
			if err == nil {
				t.Fatal("Should not succeed")
			}
		})
	})

	t.Run("Ok", func (t *testing.T) {
		sock := makeDummySocket()
		ack := make(chan *TunnelRes)

		const (
			channel   uint8 = 1
			seqNumber uint8 = 0
		)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			msg := gw.receive()
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
			defer sock.Close()
			t.Parallel()

			conn := conn{sock, DefaultClientConfig, channel}

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
		sock := makeDummySocket()
		defer sock.Close()

		var seqNumber uint8

		conn := conn{sock, DefaultClientConfig, 1}
		req := &TunnelReq{2, 0, []byte{}}

		err := conn.handleTunnelRequest(ctx, req, &seqNumber, make(chan []byte))
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("InvalidSeqNumber", func (t *testing.T) {
		sock := makeDummySocket()

		const (
			channel       uint8 = 1
			sendSeqNumber uint8 = 0
		)

		t.Run("Gateway", func (t *testing.T) {
			defer sock.Close()
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			msg := gw.receive()
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

			seqNumber := sendSeqNumber + 1

			conn := conn{sock, DefaultClientConfig, channel}
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
		sock := makeDummySocket()
		inbound := make(chan []byte)

		const (
			channel       uint8 = 1
			sendSeqNumber uint8 = 0
		)

		t.Run("Gateway", func (t *testing.T) {
			defer sock.Close()
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			msg := gw.receive()
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

			seqNumber := sendSeqNumber

			conn := conn{sock, DefaultClientConfig, channel}
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
		sock := makeDummySocket()
		defer sock.Close()

		conn := conn{sock, DefaultClientConfig, 1}

		res := &TunnelRes{2, 0, 0}
		err := conn.handleTunnelResponse(ctx, res, make(chan *TunnelRes))
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("Ok", func (t *testing.T) {
		sock := makeDummySocket()
		ack := make(chan *TunnelRes)

		t.Run("Worker", func (t *testing.T) {
			t.Parallel()

			conn := conn{sock, DefaultClientConfig, 1}

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
