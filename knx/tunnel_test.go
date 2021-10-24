// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package knx

import (
	"testing"
	"time"

	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/knxnet"
)

func makeTunnelConn(
	sock knxnet.Socket,
	config TunnelConfig,
	channel uint8,
) *Tunnel {
	return &Tunnel{
		sock:    sock,
		config:  config,
		channel: channel,
		ack:     make(chan *knxnet.TunnelRes),
		inbound: make(chan cemi.Message, 100),
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
			if req, ok := msg.(*knxnet.ConnReq); ok {
				gateway.sendAny(&knxnet.ConnRes{
					Channel: 1,
					Status:  knxnet.NoError,
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
	t.Run("Ok - without local address", func(t *testing.T) {
		client, gateway := newDummySockets()

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*knxnet.ConnReq); ok {

				expectedHostInfo := knxnet.HostInfo{
					Protocol: knxnet.UDP4,
				}

				if !expectedHostInfo.Equals(req.Control) || !expectedHostInfo.Equals(req.Tunnel) {
					t.Fatalf("Unexpected host for request: %+v", req)
				}

				gateway.sendAny(&knxnet.ConnRes{
					Channel: 1,
					Status:  knxnet.NoError,
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

	// The gateway responds to the connection request with local address.
	t.Run("Ok - with local address", func(t *testing.T) {
		client, gateway := newDummySockets()

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*knxnet.ConnReq); ok {

				expectedHostInfo := knxnet.HostInfo{
					Protocol: knxnet.UDP4,
					Address:  [4]byte{192, 168, 1, 82},
					Port:     4321,
				}

				if !expectedHostInfo.Equals(req.Control) || !expectedHostInfo.Equals(req.Tunnel) {
					t.Fatalf("Unexpected host for request: %+v", req)
				}

				gateway.sendAny(&knxnet.ConnRes{
					Channel: 1,
					Status:  knxnet.NoError,
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
				sock: client,
				config: TunnelConfig{
					ResendInterval:    500 * time.Millisecond,
					HeartbeatInterval: 1 * time.Second,
					ResponseTimeout:   1 * time.Second,
					SendLocalAddress:  true,
				},
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
			if req, ok := msg.(*knxnet.ConnReq); ok {
				gateway.sendAny(&knxnet.ConnRes{
					Channel: 0,
					Status:  knxnet.ErrNoMoreConnections,
					Control: req.Control,
				})
			} else {
				t.Fatalf("Unexpected incoming message type: %T", msg)
			}

			msg = <-gateway.Inbound()
			if req, ok := msg.(*knxnet.ConnReq); ok {
				gateway.sendAny(&knxnet.ConnRes{
					Channel: 1,
					Status:  knxnet.NoError,
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
			if req, ok := msg.(*knxnet.ConnReq); ok {
				gateway.sendAny(&knxnet.ConnRes{
					Channel: 0,
					Status:  knxnet.ErrConnectionType,
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
			if err != knxnet.ErrCode(knxnet.ErrConnectionType) {
				t.Fatalf("Expected error %v, got %v", knxnet.ErrConnectionType, err)
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

		_, err := conn.requestConnState(make(chan knxnet.ErrCode))
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

		_, err := conn.requestConnState(make(chan knxnet.ErrCode))
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

			_, err := conn.requestConnState(make(chan knxnet.ErrCode))
			if err == nil {
				t.Fatal("Should not succeed")
			}
		})
	})

	t.Run("Resend", func(t *testing.T) {
		client, gateway := newDummySockets()

		const channel uint8 = 1
		heartbeat := make(chan knxnet.ErrCode)

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			<-gateway.Inbound()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*knxnet.ConnStateReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.Status != 0 {
					t.Error("Non-null request status")
				}

				heartbeat <- knxnet.NoError
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

			if state != knxnet.NoError {
				t.Fatalf("Unexpected connection state: %v", state)
			}
		})
	})

	t.Run("InboundClosed", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer client.Close()
		defer gateway.Close()

		heartbeat := make(chan knxnet.ErrCode)
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
		heartbeat := make(chan knxnet.ErrCode)

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*knxnet.ConnStateReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.Status != 0 {
					t.Error("Non-null request status")
				}

				heartbeat <- knxnet.NoError
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

			if state != knxnet.NoError {
				t.Fatalf("Unexpected connection state: %v", state)
			}
		})
	})

	t.Run("Inactive", func(t *testing.T) {
		client, gateway := newDummySockets()

		const channel uint8 = 1
		heartbeat := make(chan knxnet.ErrCode)

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*knxnet.ConnStateReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.Status != 0 {
					t.Error("Non-null request status")
				}

				heartbeat <- knxnet.ErrConnectionID
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

			if state != knxnet.ErrConnectionID {
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
		ack := make(chan *knxnet.TunnelRes)

		const channel uint8 = 1

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			<-gateway.Inbound()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*knxnet.TunnelReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != 0 {
					t.Error("Expected sequence number 0, got", req.SeqNumber)
				}

				ack <- &knxnet.TunnelRes{Channel: req.Channel, SeqNumber: req.SeqNumber, Status: 0}
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
		ack := make(chan *knxnet.TunnelRes)

		const channel uint8 = 1

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*knxnet.TunnelReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != 0 {
					t.Error("Expected sequence number 0, got", req.SeqNumber)
				}

				ack <- &knxnet.TunnelRes{
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
		ack := make(chan *knxnet.TunnelRes)

		const channel uint8 = 1

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*knxnet.TunnelReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != 0 {
					t.Error("Expected sequence number 0, got", req.SeqNumber)
				}

				ack <- &knxnet.TunnelRes{Channel: req.Channel, SeqNumber: req.SeqNumber, Status: 1}
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
		ack := make(chan *knxnet.TunnelRes)

		const channel uint8 = 1

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if req, ok := msg.(*knxnet.TunnelReq); ok {
				if req.Channel != channel {
					t.Error("Mismatching channel")
				}

				if req.SeqNumber != 0 {
					t.Error("Expected sequence number 0, got", req.SeqNumber)
				}

				ack <- &knxnet.TunnelRes{Channel: req.Channel, SeqNumber: req.SeqNumber, Status: 0}
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
		req := &knxnet.TunnelReq{Channel: 2, SeqNumber: 0, Payload: &cemi.UnsupportedMessage{}}

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
		req := &knxnet.TunnelReq{
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

		const (
			channel       uint8 = 1
			sendSeqNumber uint8 = 0
		)

		t.Run("Gateway", func(t *testing.T) {
			t.Parallel()

			defer gateway.Close()

			msg := <-gateway.Inbound()
			if res, ok := msg.(*knxnet.TunnelRes); ok {
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

			req := &knxnet.TunnelReq{
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

			<-conn.Inbound()
		})
	})
}

func TestTunnelConn_handleTunnelRes(t *testing.T) {
	t.Run("InvalidChannel", func(t *testing.T) {
		client, gateway := newDummySockets()
		defer client.Close()
		defer gateway.Close()

		conn := makeTunnelConn(client, DefaultTunnelConfig, 1)

		res := &knxnet.TunnelRes{Channel: 2, SeqNumber: 0, Status: 0}
		err := conn.handleTunnelRes(res)
		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("Ok", func(t *testing.T) {
		client, gateway := newDummySockets()
		ack := make(chan *knxnet.TunnelRes)

		t.Run("Worker", func(t *testing.T) {
			t.Parallel()

			defer client.Close()
			defer gateway.Close()

			conn := makeTunnelConn(client, DefaultTunnelConfig, 1)
			conn.ack = ack

			res := &knxnet.TunnelRes{Channel: 1, SeqNumber: 0, Status: 0}

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
