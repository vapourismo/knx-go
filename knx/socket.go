package knx

import (
	"bytes"
	"errors"
	"net"
	"sync"
	"time"
)

//
type Socket interface {
	Send(payload OutgoingPayload) error
	Inbound() <-chan interface{}
	Close() error
}

// NewClientSocket creates a new Socket which can used to exchange KNXnet/IP packets with a gateway.
func NewClientSocket(gatewayAddress string) (Socket, error) {
	addr, err := net.ResolveUDPAddr("udp4", gatewayAddress)
	if err != nil { return nil, err }

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil { return nil, err }

	return makeUdpSocket(conn, addr), nil
}

// NewRoutingSocket creates a new Socket which can be used to exchange KNXnet/IP packets with a
// router.
func NewRoutingSocket(multicastAddress string) (Socket, error) {
	addr, err := net.ResolveUDPAddr("udp4", multicastAddress)
	if err != nil { return nil, err }

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil { return nil, err }

	return makeUdpSocket(conn, nil), nil
}

// UDP socket for KNXnet/IP packet exchange
type udpSocket struct {
	conn    *net.UDPConn
	inbound <-chan interface{}
}

// makeUdpSocket configures the UDPConn and launches the receiver and sender workers.
func makeUdpSocket(conn *net.UDPConn, addr *net.UDPAddr) *udpSocket {
	conn.SetDeadline(time.Time{})

	inbound := make(chan interface{})
	go udpSocketReceiver(conn, addr, inbound)

	return &udpSocket{conn, inbound}
}

// Send transmits a KNXnet/IP packet.
func (sock *udpSocket) Send(payload OutgoingPayload) error {
	buffer := &bytes.Buffer{}

	// Packet serialization
	err := WritePacket(buffer, payload)
	if err != nil { return err }

	Logger.Printf("udpSocket[%v]: <- %T %+v", sock.conn.RemoteAddr(), payload, payload)

	// Transmission of the buffer contents
	_, err = sock.conn.Write(buffer.Bytes())
	if err != nil { return err }

	// Logger.Printf("udpSocket[%v]: Sent: %v", sock.conn.RemoteAddr(), buffer.Bytes())

	return nil
}

// Inbound provides a channel from which you can retrieve incoming packets.
func (sock *udpSocket) Inbound() <-chan interface{} {
	return sock.inbound
}

// Close shuts the socket down. This will indirectly terminate the associated workers.
func (sock *udpSocket) Close() error {
	return sock.conn.Close()
}

// udpSocketReceiver is the receiver worker for udpSocket.
func udpSocketReceiver(conn *net.UDPConn, addr *net.UDPAddr, inbound chan<- interface{}) {
	Logger.Printf("udpSocket[%v]: Started receiver", conn.RemoteAddr())
	defer Logger.Printf("udpSocket[%v]: Stopped receiver", conn.RemoteAddr())

	defer close(inbound)

	buffer := [1024]byte{}
	reader := bytes.NewReader(buffer[:])

	for {
		len, sender, err := conn.ReadFromUDP(buffer[:])
		if err != nil {
			Logger.Printf("udpSocket[%v]: Error during read: %v", conn.RemoteAddr(), err)
			return
		}

		// Validate sender origin if necessary
		if addr != nil && (!addr.IP.Equal(sender.IP) || addr.Port != sender.Port) {
			Logger.Printf("udpSocket[%v]: Origin validation failed: %v (expected %v)",
			              conn.RemoteAddr(), sender, addr)
			continue
		}

		// Logger.Printf("udpSocket[%v]: Received: %v", conn.RemoteAddr(), buffer[:len])

		reader.Reset(buffer[:len])

		payload, err := ReadPacket(reader)
		if err != nil {
			Logger.Printf("udpSocket[%v]: Error during packet parsing: %v", conn.RemoteAddr(), err)
			continue
		}

		Logger.Printf("udpSocket[%v]: -> %T %+v", conn.RemoteAddr(), payload, payload)

		inbound <- payload
	}
}

//
type dummySocket struct {
	mu   sync.Mutex
	open bool

	out  chan OutgoingPayload
	in   chan interface{}
}

func newDummySocket() *dummySocket {
	return &dummySocket{sync.Mutex{}, true, make(chan OutgoingPayload), make(chan interface{})}
}

//
func (sock *dummySocket) gatewaySend(payload interface{}) error {
	sock.mu.Lock()
	defer sock.mu.Unlock()

	if !sock.open {
		return errors.New("Socket is closed")
	}

	sock.in <- payload

	return nil
}

//
func (sock *dummySocket) gatewayInbound() <-chan OutgoingPayload {
	return sock.out
}

//
func (sock *dummySocket) Send(payload OutgoingPayload) error {
	sock.mu.Lock()
	defer sock.mu.Unlock()

	if !sock.open {
		return errors.New("Socket is closed")
	}

	sock.out <- payload

	return nil
}

//
func (sock *dummySocket) Inbound() <-chan interface{} {
	return sock.in
}

//
func (sock *dummySocket) Close() error {
	sock.mu.Lock()
	defer sock.mu.Unlock()

	if !sock.open {
		return errors.New("Socket already closed")
	}

	sock.open = false

	close(sock.out)
	close(sock.in)

	return nil
}
