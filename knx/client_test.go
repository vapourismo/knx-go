package knx

import (
	"context"
	"testing"
	"time"
)

var clientConfig = ClientConfig{
	2 * time.Second,
	defaultResendInterval,
	2 * time.Second,
	2 * time.Second,
}

func TestConnHandle_RequestConnection(t *testing.T) {
	ctx := context.Background()

	// Socket was closed before anything could be done.
	t.Run("SendFails", func (t *testing.T) {
		conn := connHandle{makeDummySocket(), clientConfig, 0}
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

		conn := connHandle{sock, clientConfig, 0}

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

			conn := connHandle{sock, clientConfig, 0}

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

		conn := connHandle{sock, clientConfig, 0}

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

			conn := connHandle{sock, clientConfig, 0}

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

			conn := connHandle{sock, clientConfig, 0}

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

		conn := connHandle{sock, clientConfig, 1}

		err := conn.requestConnectionState(ctx, make(chan ConnState))
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("CancelledContext", func (t *testing.T) {
		sock := makeDummySocket()
		defer sock.Close()

		conn := connHandle{sock, clientConfig, 1}

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

		conn := connHandle{sock, clientConfig, 1}

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

			conn := connHandle{sock, clientConfig, channel}

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

			conn := connHandle{sock, clientConfig, channel}

			err := conn.requestConnectionState(ctx, heartbeat)
			if err != ConnStateInactive {
				t.Fatal(err)
			}
		})
	})
}

func TestConnHandle_handleTunnelRequest(t *testing.T) {
	ctx := context.Background()

	// Proper tunnel request from gateway.
	t.Run("Ok", func (t *testing.T) {
		sock := makeDummySocket()
		inbound := make(chan []byte)

		const (
			channel uint8 = 1
			seqNumber uint8 = 0
		)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			msg := gw.receive()
			if res, ok := msg.(*TunnelResponse); ok {
				if res.Channel != channel {
					t.Error("Mismatching channels")
				}

				if res.SeqNumber != seqNumber {
					t.Error("Mismatching sequence numbers")
				}

				if res.Status != 0 {
					t.Error("Invalid response status")
				}
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}
		})

		t.Run("Worker", func (t *testing.T) {
			t.Parallel()

			var seqNo uint8 = seqNumber
			handle := &connHandle{sock, clientConfig, 1}

			req := &TunnelRequest{channel, seqNumber, []byte{}}
			err := handle.handleTunnelRequest(ctx, req, &seqNo, inbound)
			if err != nil {
				t.Fatal(err)
			}

			if seqNo != req.SeqNumber + 1 {
				t.Error("Sequence number was not increased")
			}
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			select {
			case <-ctx.Done():
				t.Fatalf("While waiting for inbound packet: %v", ctx.Err())

			case _, open := <-inbound:
				if !open {
					t.Fatal("Inbound channel was closed")
				}
			}
		})
	})

	// Out-of-sequence tunnel request from the gateway.
	t.Run("OutOfSequence", func (t *testing.T) {
		sock := makeDummySocket()
		inbound := make(chan []byte)

		const (
			channel uint8 = 1
			seqNumber uint8 = 1
		)

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			gw := gatewayHelper{ctx, sock, t}

			msg := gw.receive()
			if res, ok := msg.(*TunnelResponse); ok {
				if res.Channel != channel {
					t.Error("Mismatching channels")
				}

				if res.SeqNumber != seqNumber {
					t.Error("Mismatching sequence numbers")
				}

				if res.Status != 0 {
					t.Error("Invalid response status")
				}
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}
		})

		t.Run("Worker", func (t *testing.T) {
			t.Parallel()

			var seqNo uint8 = seqNumber - 1
			handle := &connHandle{sock, clientConfig, channel}

			req := &TunnelRequest{channel, seqNumber, []byte{}}
			err := handle.handleTunnelRequest(ctx, req, &seqNo, inbound)
			if err != nil {
				t.Fatal(err)
			}

			if seqNo != seqNumber - 1 {
				t.Error("Sequence number was changed by an out-of-sequence tunnel request")
			}
		})
	})

	// Tunnel request on incorrect channel.
	t.Run("WrongChannel", func (t *testing.T) {
		sock := makeDummySocket()
		defer sock.Close()

		inbound := make(chan []byte)

		const (
			channel uint8 = 1
			seqNumber uint8 = 1
		)

		var seqNo uint8 = 0
		handle := &connHandle{sock, clientConfig, channel + 1}

		req := &TunnelRequest{channel, seqNumber, []byte{}}
		err := handle.handleTunnelRequest(ctx, req, &seqNo, inbound)
		if err == nil {
			t.Fatal("Tunnel request with wrong channel has been successful")
		}

		if seqNo != 0 {
			t.Error("Sequence number was changed by an invalid tunnel request")
		}
	})
}
