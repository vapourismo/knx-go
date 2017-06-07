// Copyright 2017 Ole Krüger.

package knx

import (
	"testing"

	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/proto"
)

func makeTunnelConn(
	sock Socket,
	config TunnelConfig,
	channel uint8,
) *Tunnel {
	return &Tunnel{
		sock:    sock,
		config:  config,
		channel: channel,
		ack:     make(chan *proto.TunnelRes),
		inbound: make(chan cemi.Message),
	}
}

func TestTunnelConn_requestConn(t *testing.T) {
	// Socket was closed before anything could be done.
	t.Run("SendFails", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer gateway.Close()

		client.Close()

		conn := Tunnel{
			sock:   client,
			config: DefaultTunnelConfig,
		}

		err := conn.requestConn()
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	// Context is done.
	t.Run("Timeout", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer client.Close()
		defer gateway.Close()

		config := DefaultTunnelConfig
		config.ResponseTimeout = 1

		conn := Tunnel{
			sock:   client,
			config: config,
		}

		err := conn.requestConn()
		if err != errResponseTimeout {
			t.Fatalf("Expected error %v, got %v", errResponseTimeout, err)
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

			config := DefaultTunnelConfig
			config.ResendInterval = 1

			conn := Tunnel{
				sock:   client,
				config: config,
			}

			err := conn.requestConn()
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
				gateway.sendAny(&proto.ConnRes{
					Channel: 1,
					Status:  proto.ConnResOk,
					Control: req.Control,
				})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			config := DefaultTunnelConfig
			config.ResendInterval = 1

			conn := Tunnel{
				sock:   client,
				config: config,
			}

			err := conn.requestConn()
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

		conn := Tunnel{
			sock:   client,
			config: DefaultTunnelConfig,
		}

		err := conn.requestConn()
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
				gateway.sendAny(&proto.ConnRes{
					Channel: 1,
					Status:  proto.ConnResOk,
					Control: req.Control,
				})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			conn := Tunnel{
				sock:   client,
				config: DefaultTunnelConfig,
			}

			err := conn.requestConn()
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
				gateway.sendAny(&proto.ConnRes{Channel: 0, Status: proto.ConnResBusy, Control: req.Control})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}

			msg = <-gateway.Inbound()
			if req, ok := msg.(*proto.ConnReq); ok {
				gateway.sendAny(&proto.ConnRes{
					Channel: 1,
					Status:  proto.ConnResOk,
					Control: req.Control,
				})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			config := DefaultTunnelConfig
			config.ResendInterval = 1

			conn := Tunnel{
				sock:   client,
				config: config,
			}

			err := conn.requestConn()
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
				gateway.sendAny(&proto.ConnRes{
					Channel: 0,
					Status:  proto.ConnResUnsupportedType,
					Control: req.Control,
				})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			conn := Tunnel{
				sock:   client,
				config: DefaultTunnelConfig,
			}

			err := conn.requestConn()
			if err != proto.ConnResUnsupportedType {
				t.Fatalf("Expected error %v, got %v", proto.ConnResUnsupportedType, err)
			}
		})
	})
}

func TestTunnelConn_requestState(t *testing.T) {
	t.Run("SendFails", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer gateway.Close()

		client.Close()

		conn := makeTunnelConn(client, DefaultTunnelConfig, 1)

		_, err := conn.requestConnState(make(chan proto.ConnStateResStatus))
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("CancelledContext", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer client.Close()
		defer gateway.Close()

		config := DefaultTunnelConfig
		config.ResponseTimeout = 1

		conn := makeTunnelConn(client, config, 1)

		_, err := conn.requestConnState(make(chan proto.ConnStateResStatus))
		if err != errResponseTimeout {
			t.Fatalf("Expected error %v, got %v", errResponseTimeout, err)
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

			config := DefaultTunnelConfig
			config.ResendInterval = 1

			conn := makeTunnelConn(client, config, 1)

			_, err := conn.requestConnState(make(chan proto.ConnStateResStatus))
			if err == nil {
				t.Fatal("Should not succeed")
			}
		})
	})

	t.Run("Resend", func(t *testing.T) {
		client, gateway := newDummySockets()

		const channel uint8 = 1
		heartbeat := make(chan proto.ConnStateResStatus)

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

			config := DefaultTunnelConfig
			config.ResendInterval = 1

			conn := makeTunnelConn(client, config, channel)

			state, err := conn.requestConnState(heartbeat)

			if err != nil {
				t.Fatal(err)
			}

			if state != proto.ConnStateNormal {
				t.Fatalf("Unexpected connection state: %v", state)
			}
		})
	})

	t.Run("InboundClosed", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer client.Close()
		defer gateway.Close()

		heartbeat := make(chan proto.ConnStateResStatus)
		close(heartbeat)

		conn := makeTunnelConn(client, DefaultTunnelConfig, 1)

		_, err := conn.requestConnState(heartbeat)
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("Ok", func(t *testing.T) {
		client, gateway := newDummySockets()

		const channel uint8 = 1
		heartbeat := make(chan proto.ConnStateResStatus)

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

			conn := makeTunnelConn(client, DefaultTunnelConfig, channel)

			state, err := conn.requestConnState(heartbeat)

			if err != nil {
				t.Fatal(err)
			}

			if state != proto.ConnStateNormal {
				t.Fatalf("Unexpected connection state: %v", state)
			}
		})
	})

	t.Run("Inactive", func(t *testing.T) {
		client, gateway := newDummySockets()

		const channel uint8 = 1
		heartbeat := make(chan proto.ConnStateResStatus)

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

			conn := makeTunnelConn(client, DefaultTunnelConfig, channel)

			state, err := conn.requestConnState(heartbeat)

			if err != nil {
				t.Fatal(err)
			}

			if state != proto.ConnStateInactive {
				t.Fatalf("Unexpected connection state: %v", state)
			}
		})
	})
}

func TestTunnelConn_requestTunnel(t *testing.T) {
	t.Run("SendFails", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer gateway.Close()

		client.Close()

		conn := makeTunnelConn(client, DefaultTunnelConfig, 1)

		err := conn.requestTunnel(&cemi.UnsupportedMessage{})
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("Timeout", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer client.Close()
		defer gateway.Close()

		config := DefaultTunnelConfig
		config.ResponseTimeout = 1

		conn := makeTunnelConn(client, config, 1)

		err := conn.requestTunnel(&cemi.UnsupportedMessage{})
		if err != errResponseTimeout {
			t.Fatalf("Expected %v, got %v", errResponseTimeout, err)
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

			config := DefaultTunnelConfig
			config.ResendInterval = 1

			conn := makeTunnelConn(client, config, 1)

			err := conn.requestTunnel(&cemi.UnsupportedMessage{})
			if err == nil {
				t.Fatal("Should not succeed")
			}
		})
	})

	t.Run("Resend", func(t *testing.T) {
		client, gateway := newDummySockets()
		ack := make(chan *proto.TunnelRes)

		const channel uint8 = 1

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			<-gateway.Inbound()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*proto.TunnelReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != 0 {
					t.Error("Expected sequence number 0, got", req.SeqNumber)
				}

				ack <- &proto.TunnelRes{Channel: req.Channel, SeqNumber: req.SeqNumber, Status: 0}
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			config := DefaultTunnelConfig
			config.ResendInterval = 1

			conn := makeTunnelConn(client, config, channel)
			conn.ack = ack

			err := conn.requestTunnel(&cemi.UnsupportedMessage{})
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("ClosedAckChannel", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer client.Close()
		defer gateway.Close()

		conn := makeTunnelConn(client, DefaultTunnelConfig, 1)
		close(conn.ack)

		err := conn.requestTunnel(&cemi.UnsupportedMessage{})
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("InvalidSeqNumber", func(t *testing.T) {
		client, gateway := newDummySockets()
		ack := make(chan *proto.TunnelRes)

		const channel uint8 = 1

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*proto.TunnelReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != 0 {
					t.Error("Expected sequence number 0, got", req.SeqNumber)
				}

				ack <- &proto.TunnelRes{
					Channel:   req.Channel,
					SeqNumber: req.SeqNumber + 10,
					Status:    0,
				}
				close(ack)
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			conn := makeTunnelConn(client, DefaultTunnelConfig, channel)
			conn.ack = ack

			err := conn.requestTunnel(&cemi.UnsupportedMessage{})
			if err == nil {
				t.Fatal("Should not succeed")
			}
		})
	})

	t.Run("BadStatus", func(t *testing.T) {
		client, gateway := newDummySockets()
		ack := make(chan *proto.TunnelRes)

		const channel uint8 = 1

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*proto.TunnelReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != 0 {
					t.Error("Expected sequence number 0, got", req.SeqNumber)
				}

				ack <- &proto.TunnelRes{Channel: req.Channel, SeqNumber: req.SeqNumber, Status: 1}
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			conn := makeTunnelConn(client, DefaultTunnelConfig, channel)
			conn.ack = ack

			err := conn.requestTunnel(&cemi.UnsupportedMessage{})
			if err == nil {
				t.Fatal("Should not succeed")
			}
		})
	})

	t.Run("Ok", func(t *testing.T) {
		client, gateway := newDummySockets()
		ack := make(chan *proto.TunnelRes)

		const channel uint8 = 1

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*proto.TunnelReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != 0 {
					t.Error("Expected sequence number 0, got", req.SeqNumber)
				}

				ack <- &proto.TunnelRes{Channel: req.Channel, SeqNumber: req.SeqNumber, Status: 0}
			} else {
				t.Fatalf("Unexpected type %T", msg)
			}
		})

		t.Run("Client", func(t *testing.T) {
			t.Parallel()

			defer client.Close()

			conn := makeTunnelConn(client, DefaultTunnelConfig, channel)
			conn.ack = ack

			err := conn.requestTunnel(&cemi.UnsupportedMessage{})
			if err != nil {
				t.Fatal(err)
			}
		})
	})
}

func TestTunnelConn_handleTunnelReq(t *testing.T) {
	t.Run("InvalidChannel", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer client.Close()
		defer gateway.Close()

		var seqNumber uint8

		conn := makeTunnelConn(client, DefaultTunnelConfig, 1)
		req := &proto.TunnelReq{Channel: 2, SeqNumber: 0, Payload: &cemi.UnsupportedMessage{}}

		err := conn.handleTunnelReq(req, &seqNumber)
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

		defer client.Close()
		defer gateway.Close()

		seqNumber := sendSeqNumber

		conn := makeTunnelConn(client, DefaultTunnelConfig, channel)
		req := &proto.TunnelReq{
			Channel:   channel,
			SeqNumber: seqNumber + 1,
			Payload:   &cemi.UnsupportedMessage{},
		}

		err := conn.handleTunnelReq(req, &seqNumber)
		if err == nil {
			t.Error("Should not succeed")
		}

		if seqNumber != sendSeqNumber {
			t.Error("Sequence number was modified")
		}
	})

	t.Run("Ok", func(t *testing.T) {
		client, gateway := newDummySockets()
		inbound := make(chan cemi.Message)

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

			conn := makeTunnelConn(client, DefaultTunnelConfig, channel)
			conn.inbound = inbound

			req := &proto.TunnelReq{
				Channel:   channel,
				SeqNumber: sendSeqNumber,
				Payload:   &cemi.UnsupportedMessage{},
			}

			err := conn.handleTunnelReq(req, &seqNumber)
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

func TestTunnelConn_handleTunnelRes(t *testing.T) {
	t.Run("InvalidChannel", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer client.Close()
		defer gateway.Close()

		conn := makeTunnelConn(client, DefaultTunnelConfig, 1)

		res := &proto.TunnelRes{Channel: 2, SeqNumber: 0, Status: 0}
		err := conn.handleTunnelRes(res)
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

			conn := makeTunnelConn(client, DefaultTunnelConfig, 1)
			conn.ack = ack

			res := &proto.TunnelRes{Channel: 1, SeqNumber: 0, Status: 0}

			err := conn.handleTunnelRes(res)
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
