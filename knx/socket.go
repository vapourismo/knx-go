package knx

import (
	"bytes"
	"net"
	"time"

	"github.com/vapourismo/knx-go/knx/proto"
)

// A Socket is a socket, duh.
type Socket interface {
	Send(payload proto.ServiceWriterTo) error
	Inbound() <-chan proto.Service
	Close() error
}

// NewTunnelSocket creates a new Socket which can used to exchange KNXnet/IP packets with a gateway.
func NewTunnelSocket(gatewayAddress string) (Socket, error) {
	addr, err := net.ResolveUDPAddr("udp4", gatewayAddress)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return nil, err
	}

	conn.SetDeadline(time.Time{})

	inbound := make(chan proto.Service)
	go serveUDPSocket(conn, addr, inbound)

	return &tunnelSock{conn, inbound}, nil
}

// NewRoutingSocket creates a new Socket which can be used to exchange KNXnet/IP packets with a
// router.
func NewRoutingSocket(multicastAddress string) (Socket, error) {
	addr, err := net.ResolveUDPAddr("udp4", multicastAddress)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return nil, err
	}

	conn.SetDeadline(time.Time{})

	inbound := make(chan proto.Service)
	go serveUDPSocket(conn, nil, inbound)

	return &routerSock{conn, addr, inbound}, nil
}

// tunnelSock is a UDP socket for KNXnet/IP packet exchange.
type tunnelSock struct {
	conn    *net.UDPConn
	inbound <-chan proto.Service
}

// Send transmits a KNXnet/IP packet.
func (sock *tunnelSock) Send(payload proto.ServiceWriterTo) error {
	buffer := bytes.Buffer{}

	// Packet serialization
	_, err := proto.Pack(&buffer, payload)
	if err != nil {
		return err
	}

	log(sock.conn, "Socket", "<- %T %+v", payload, payload)

	// Transmission of the buffer contents
	_, err = sock.conn.Write(buffer.Bytes())
	return err
}

// Inbound provides a channel from which you can retrieve incoming packets.
func (sock *tunnelSock) Inbound() <-chan proto.Service {
	return sock.inbound
}

// Close shuts the socket down. This will indirectly terminate the associated workers.
func (sock *tunnelSock) Close() error {
	return sock.conn.Close()
}

// routerSock is a UDP socket for KNXnet/IP packet exchange.
type routerSock struct {
	conn    *net.UDPConn
	addr    *net.UDPAddr
	inbound <-chan proto.Service
}

// Send transmits a KNXnet/IP packet.
func (sock *routerSock) Send(payload proto.ServiceWriterTo) error {
	buffer := bytes.Buffer{}

	// Packet serialization
	_, err := proto.Pack(&buffer, payload)
	if err != nil {
		return err
	}

	log(sock.conn, "Socket", "<- %T %+v", payload, payload)

	// Transmission of the buffer contents
	_, err = sock.conn.WriteToUDP(buffer.Bytes(), sock.addr)
	return err
}

// Inbound provides a channel from which you can retrieve incoming packets.
func (sock *routerSock) Inbound() <-chan proto.Service {
	return sock.inbound
}

// Close shuts the socket down. This will indirectly terminate the associated workers.
func (sock *routerSock) Close() error {
	return sock.conn.Close()
}

// serveUDPSocket is the receiver worker for a UDP socket.
func serveUDPSocket(conn *net.UDPConn, addr *net.UDPAddr, inbound chan<- proto.Service) {
	log(conn, "Socket", "Started receiver")
	defer log(conn, "Socket", "Stopped receiver")

	// A closed inbound channel indicates to its readers that the worker has terminated.
	defer close(inbound)

	buffer := [1024]byte{}

	for {
		len, sender, err := conn.ReadFromUDP(buffer[:])
		if err != nil {
			log(conn, "Socket", "Error during read: %v", err)
			return
		}

		// Validate sender origin if necessary
		if addr != nil && (!addr.IP.Equal(sender.IP) || addr.Port != sender.Port) {
			log(conn, "Socket", "Origin validation failed: %v (expected %v)", sender, addr)
			continue
		}

		var payload proto.Service
		_, err = proto.Unpack(buffer[:len], &payload)
		if err != nil {
			log(conn, "Socket", "Error during packet parsing: %v", err)
			continue
		}

		log(conn, "Socket", "-> %T %+v", payload, payload)

		inbound <- payload
	}
}
