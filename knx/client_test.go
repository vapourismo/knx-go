package knx

import (
	"context"
	"testing"
	"time"
)

func TestRequestConnection(t *testing.T) {
	ctx := context.Background()

	t.Run("ClosedSocket", func (t *testing.T) {
		sock := newDummySocket()
		sock.Close()

		_, err := requestConnection(ctx, sock)
		if err == nil {
			t.Fatal("Success on closed socket")
		}
	})

	t.Run("ConnectionOk", func (t *testing.T) {
		sock := newDummySocket()

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			select {
			case <-ctx.Done():
				t.Fatalf("While waiting for inbound packet: %v", ctx.Err())

			case msg := <-sock.gatewayInbound():
				if req, ok := msg.(*ConnectionRequest); ok {
					err := sock.gatewaySend(&ConnectionResponse{1, ConnResOk, req.Control})
					if err != nil {
						t.Fatalf("While sending connection response: %v", nil)
					}
				} else {
					t.Fatalf("Unexpected incoming message type: %T", msg)
				}
			}
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			_, err := requestConnection(ctx, sock)
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("ConnectionTimeout", func (t *testing.T) {
		sock := newDummySocket()

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			select {
			case <-ctx.Done():
				t.Fatalf("While waiting for inbound packet: %v", ctx.Err())

			case <-sock.gatewayInbound():
			}
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(ctx, 200 * time.Millisecond)
			defer cancel()

			_, err := requestConnection(ctx, sock)
			if err != ctx.Err() {
				t.Fatalf("Expected error %v, got %v", ctx.Err(), err)
			}
		})
	})

	t.Run("ConnectionMultipleBusy", func (t *testing.T) {
		sock := newDummySocket()

		t.Run("Gateway", func (t *testing.T) {
			t.Parallel()

			select {
			case <-ctx.Done():
				t.Fatalf("While waiting for inbound packet: %v", ctx.Err())

			case msg := <-sock.gatewayInbound():
				if req, ok := msg.(*ConnectionRequest); ok {
					err := sock.gatewaySend(&ConnectionResponse{0, ConnResBusy, req.Control})
					if err != nil {
						t.Fatalf("While sending connection response: %v", nil)
					}
				} else {
					t.Fatalf("Unexpected incoming message type: %T", msg)
				}
			}

			select {
			case <-ctx.Done():
				t.Fatalf("While waiting for inbound packet: %v", ctx.Err())

			case msg := <-sock.gatewayInbound():
				if req, ok := msg.(*ConnectionRequest); ok {
					err := sock.gatewaySend(&ConnectionResponse{1, ConnResOk, req.Control})
					if err != nil {
						t.Fatalf("While sending connection response: %v", nil)
					}
				} else {
					t.Fatalf("Unexpected incoming message type: %T", msg)
				}
			}
		})

		t.Run("Client", func (t *testing.T) {
			t.Parallel()

			_, err := requestConnection(ctx, sock)
			if err != nil {
				t.Fatal(err)
			}
		})
	})
}
