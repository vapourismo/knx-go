package knx

import (
	"context"
	"testing"
)

func TestConnHandle_requestConnection(t *testing.T) {
	ctx := context.Background()

	// Socket was closed before anything could be done.
	t.Run("SendFails", func (t *testing.T) {
		conn := connHandle{makeDummySocket(), DefaultClientConfig, 0}
		conn.sock.Close()

		err := conn.requestConnection(ctx)
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	// Context is done.
	t.Run("CancelledContext", func (t *testing.T) {
		sock := makeDummySocket()
		defer sock.Close()

		conn := connHandle{sock, DefaultClientConfig, 0}

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		err := conn.requestConnection(ctx)
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

			conn := connHandle{sock, config, 0}

			err := conn.requestConnection(ctx)
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

			conn := connHandle{sock, config, 0}

			err := conn.requestConnection(ctx)
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

		conn := connHandle{sock, DefaultClientConfig, 0}

		err := conn.requestConnection(ctx)
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

			conn := connHandle{sock, DefaultClientConfig, 0}

			err := conn.requestConnection(ctx)
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

			conn := connHandle{sock, config, 0}

			err := conn.requestConnection(ctx)
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

			conn := connHandle{sock, DefaultClientConfig, 0}

			err := conn.requestConnection(ctx)
			if err != ConnResUnsupportedType {
				t.Fatalf("Expected error %v, got %v", ConnResUnsupportedType, err)
			}
		})
	})
}

func TestConnHandle_requestConnectionState(t *testing.T) {
	ctx := context.Background()

	t.Run("SendFails", func (t *testing.T) {
		sock := makeDummySocket()
		sock.Close()

		conn := connHandle{sock, DefaultClientConfig, 1}

		err := conn.requestConnectionState(ctx, make(chan ConnState))
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("CancelledContext", func (t *testing.T) {
		sock := makeDummySocket()
		defer sock.Close()

		conn := connHandle{sock, DefaultClientConfig, 1}

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		err := conn.requestConnectionState(ctx, make(chan ConnState))
		if err != ctx.Err() {
			t.Fatal("Expected error %v, got %v", ctx.Err(), err)
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

			conn := connHandle{sock, config, 1}

			err := conn.requestConnectionState(ctx, make(chan ConnState))
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
			if req, ok := msg.(*ConnectionStateRequest); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.Status != 0 {
					t.Error("Non-null request status")
				}

				heartbeat <- ConnStateNormal
			} else {
				t.Fatal("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			defer sock.Close()
			t.Parallel()

			config := DefaultClientConfig
			config.ResendInterval = 1

			conn := connHandle{sock, config, channel}

			err := conn.requestConnectionState(ctx, heartbeat)
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("InboundClosed", func (t *testing.T) {
		sock := makeDummySocket()

		heartbeat := make(chan ConnState)
		close(heartbeat)

		conn := connHandle{sock, DefaultClientConfig, 1}

		err := conn.requestConnectionState(ctx, heartbeat)
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
			if req, ok := msg.(*ConnectionStateRequest); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.Status != 0 {
					t.Error("Non-null request status")
				}

				heartbeat <- ConnStateNormal
			} else {
				t.Fatal("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			defer sock.Close()
			t.Parallel()

			conn := connHandle{sock, DefaultClientConfig, channel}

			err := conn.requestConnectionState(ctx, heartbeat)
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
			if req, ok := msg.(*ConnectionStateRequest); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.Status != 0 {
					t.Error("Non-null request status")
				}

				heartbeat <- ConnStateInactive
			} else {
				t.Fatal("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			defer sock.Close()
			t.Parallel()

			conn := connHandle{sock, DefaultClientConfig, channel}

			err := conn.requestConnectionState(ctx, heartbeat)
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

		conn := connHandle{sock, DefaultClientConfig, 1}

		err := conn.requestTunnel(ctx, 0, []byte{}, make(chan *TunnelResponse))
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("ContextCancelled", func (t *testing.T) {
		sock := makeDummySocket()
		defer sock.Close()

		conn := connHandle{sock, DefaultClientConfig, 1}

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		err := conn.requestTunnel(ctx, 0, []byte{}, make(chan *TunnelResponse))
		if err != ctx.Err() {
			t.Fatal("Expected %v, got %v", ctx.Err(), err)
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

			conn := connHandle{sock, config, 1}

			err := conn.requestTunnel(ctx, 1, []byte{}, make(chan *TunnelResponse))
			if err == nil {
				t.Fatal("Should not succeed")
			}
		})
	})

	t.Run("Resend", func (t *testing.T) {
		sock := makeDummySocket()
		ack := make(chan *TunnelResponse)

		const (
			channel uint8 = 1
			seqNumber uint8 = 0
		)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			gw.ignore()

			msg := gw.receive()
			if req, ok := msg.(*TunnelRequest); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != seqNumber {
					t.Error("Mismatching sequence number")
				}

				ack <- &TunnelResponse{req.Channel, req.SeqNumber, 0}
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			defer sock.Close()
			t.Parallel()

			config := DefaultClientConfig
			config.ResendInterval = 1

			conn := connHandle{sock, config, channel}

			err := conn.requestTunnel(ctx, seqNumber, []byte{}, ack)
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("ClosedAckChannel", func (t *testing.T) {
		sock := makeDummySocket()
		defer sock.Close()

		ack := make(chan *TunnelResponse)
		close(ack)

		conn := connHandle{sock, DefaultClientConfig, 1}

		err := conn.requestTunnel(ctx, 0, []byte{}, ack)
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("InvalidSeqNumber", func (t *testing.T) {
		sock := makeDummySocket()
		ack := make(chan *TunnelResponse)

		ctx, cancel := context.WithCancel(ctx)

		const (
			channel uint8 = 1
			seqNumber uint8 = 0
		)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			msg := gw.receive()
			if req, ok := msg.(*TunnelRequest); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != seqNumber {
					t.Error("Mismatching sequence number")
				}

				ack <- &TunnelResponse{req.Channel, req.SeqNumber + 10, 0}
				cancel()
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			defer sock.Close()
			t.Parallel()

			conn := connHandle{sock, DefaultClientConfig, channel}

			err := conn.requestTunnel(ctx, seqNumber, []byte{}, ack)
			if err != ctx.Err() {
				t.Fatalf("Expected error %v, got %v", ctx.Err(), err)
			}
		})
	})

	t.Run("BadStatus", func (t *testing.T) {
		sock := makeDummySocket()
		ack := make(chan *TunnelResponse)

		const (
			channel uint8 = 1
			seqNumber uint8 = 0
		)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			msg := gw.receive()
			if req, ok := msg.(*TunnelRequest); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != seqNumber {
					t.Error("Mismatching sequence number")
				}

				ack <- &TunnelResponse{req.Channel, req.SeqNumber, 1}
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			defer sock.Close()
			t.Parallel()

			conn := connHandle{sock, DefaultClientConfig, channel}

			err := conn.requestTunnel(ctx, seqNumber, []byte{}, ack)
			if err == nil {
				t.Fatal("Should not succeed")
			}
		})
	})

	t.Run("Ok", func (t *testing.T) {
		sock := makeDummySocket()
		ack := make(chan *TunnelResponse)

		const (
			channel uint8 = 1
			seqNumber uint8 = 0
		)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			msg := gw.receive()
			if req, ok := msg.(*TunnelRequest); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != seqNumber {
					t.Error("Mismatching sequence number")
				}

				ack <- &TunnelResponse{req.Channel, req.SeqNumber, 0}
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func (t *testing.T) {
			defer sock.Close()
			t.Parallel()

			conn := connHandle{sock, DefaultClientConfig, channel}

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

		var seqNumber uint8 = 0

		conn := connHandle{sock, DefaultClientConfig, 1}
		req := &TunnelRequest{2, 0, []byte{}}

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
			if res, ok := msg.(*TunnelResponse); ok {
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
				t.Fatal("Unexpected type %T", msg)
			}
		})

		t.Run("Worker", func (t *testing.T) {
			t.Parallel()

			var seqNumber uint8 = sendSeqNumber + 1

			conn := connHandle{sock, DefaultClientConfig, channel}
			req := &TunnelRequest{channel, sendSeqNumber, []byte{}}

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
			if res, ok := msg.(*TunnelResponse); ok {
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
				t.Fatal("Unexpected type %T", msg)
			}
		})

		t.Run("Worker", func (t *testing.T) {
			t.Parallel()

			var seqNumber uint8 = sendSeqNumber

			conn := connHandle{sock, DefaultClientConfig, channel}
			req := &TunnelRequest{channel, sendSeqNumber, []byte{}}

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
