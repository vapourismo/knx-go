package knx

import (
	"bytes"
	"github.com/vapourismo/knx-go/knx/proto"
	"net"
	"time"
)

// A Socket is a socket, duh.
type Socket interface {
	Send(payload proto.ServiceWriterTo) error
	Inbound() <-chan proto.Service
	Close() error
}

// NewClientSocket creates a new Socket which can used to exchange KNXnet/IP packets with a gateway.
func NewClientSocket(gatewayAddress string) (Socket, error) {
	addr, err := net.ResolveUDPAddr("udp4", gatewayAddress)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return nil, err
	}

	return newUDPSocket(conn, addr), nil
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

	return newUDPSocket(conn, nil), nil
}

// UDP socket for KNXnet/IP packet exchange
type udpSocket struct {
	conn    *net.UDPConn
	inbound <-chan proto.Service
}

// newUDPSocket configures the UDPConn and launches the receiver and sender workers.
func newUDPSocket(conn *net.UDPConn, addr *net.UDPAddr) *udpSocket {
	conn.SetDeadline(time.Time{})

	inbound := make(chan proto.Service)
	go udpSocketReceiver(conn, addr, inbound)

	return &udpSocket{conn, inbound}
}

// Send transmits a KNXnet/IP packet.
func (sock *udpSocket) Send(payload proto.ServiceWriterTo) error {
	buffer := bytes.Buffer{}

	// Packet serialization
	_, err := proto.Pack(&buffer, payload)
	if err != nil {
		return err
	}

	log(sock.conn, "udpSocket", "<- %T %+v", payload, payload)

	// Transmission of the buffer contents
	_, err = sock.conn.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Inbound provides a channel from which you can retrieve incoming packets.
func (sock *udpSocket) Inbound() <-chan proto.Service {
	return sock.inbound
}

// Close shuts the socket down. This will indirectly terminate the associated workers.
func (sock *udpSocket) Close() error {
	return sock.conn.Close()
}

// udpSocketReceiver is the receiver worker for udpSocket.
func udpSocketReceiver(conn *net.UDPConn, addr *net.UDPAddr, inbound chan<- proto.Service) {
	log(conn, "udpSocket", "Started receiver")
	defer log(conn, "udpSocket", "Stopped receiver")

	// A closed inbound channel indicates to its readers that the worker has terminated.
	defer close(inbound)

	buffer := [1024]byte{}

	for {
		len, sender, err := conn.ReadFromUDP(buffer[:])
		if err != nil {
			log(conn, "udpSocket", "Error during read: %v", err)
			return
		}

		// Validate sender origin if necessary
		if addr != nil && (!addr.IP.Equal(sender.IP) || addr.Port != sender.Port) {
			log(conn, "udpSocket", "Origin validation failed: %v (expected %v)", sender, addr)
			continue
		}

		var payload proto.Service
		_, err = proto.Unpack(bytes.NewReader(buffer[:len]), &payload)
		if err != nil {
			log(conn, "udpSocket", "Error during packet parsing: %v", err)
			continue
		}

		log(conn, "udpSocket", "-> %T %+v", payload, payload)

		inbound <- payload
	}
}
