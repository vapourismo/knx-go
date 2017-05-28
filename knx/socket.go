package knx

import (
	"bytes"
	"errors"
	"net"
	"time"
	"github.com/vapourismo/knx-go/knx/proto"
)

// A Socket is a socket, duh.
type Socket interface {
	Send(payload proto.ServiceWriterTo) error
	Receive() (proto.Service, error)
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

	return makeUDPSocket(conn, addr), nil
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

	return makeUDPSocket(conn, nil), nil
}

// UDP socket for KNXnet/IP packet exchange
type udpSocket struct {
	conn *net.UDPConn
	addr *net.UDPAddr
}

// makeUDPSocket configures the UDPConn and launches the receiver and sender workers.
func makeUDPSocket(conn *net.UDPConn, addr *net.UDPAddr) *udpSocket {
	conn.SetDeadline(time.Time{})

	return &udpSocket{conn, addr}
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

// These are errors that can occur during calls to Receive.
var (
	ErrInvalidOrigin = errors.New("Invalid packet origin")
)

// Receive waits for an incoming KNXnet/IP packet.
func (sock *udpSocket) Receive() (proto.Service, error) {
	buffer := [1024]byte{}

	len, sender, err := sock.conn.ReadFromUDP(buffer[:])
	if err != nil {
		return nil, err
	}

	// Validate sender origin if necessary
	if sock.addr != nil && (!sock.addr.IP.Equal(sender.IP) || sock.addr.Port != sender.Port) {
		return nil, ErrInvalidOrigin
	}

	var payload proto.Service
	if _, err = proto.Unpack(bytes.NewReader(buffer[:len]), &payload); err != nil {
		return nil, err
	}

	log(sock.conn, "udpSocket", "-> %T %+v", payload, payload)

	return payload, nil
}

// Close shuts the socket down. This will indirectly terminate the associated workers.
func (sock *udpSocket) Close() error {
	return sock.conn.Close()
}
