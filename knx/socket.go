package knx

import (
	"bytes"
	"net"
	"time"
)

// UDP socket for KNXnet/IP packet exchange
type Socket struct {
	conn *net.UDPConn

	// Inbound relays incoming KNXnet/IP packets.
	// The types of these packets are limited to those returned by ReadPacket.
	Inbound <-chan interface{}

	// Outbound relays outgoing KNXnet/IP packets to the gateway or router.
	Outbound chan<- OutgoingPayload
}

// NewClientSocket creates a new Socket which can used to exchange KNXnet/IP packets with a gateway.
func NewClientSocket(gatewayAddress string) (*Socket, error) {
	addr, err := net.ResolveUDPAddr("udp4", gatewayAddress)
	if err != nil { return nil, err }

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil { return nil, err }

	return makeSocket(conn, addr), nil
}

// NewRoutingSocket creates a new Socket which can be used to exchange KNXnet/IP packets with a
// router.
func NewRoutingSocket(multicastAddress string) (*Socket, error) {
	addr, err := net.ResolveUDPAddr("udp4", multicastAddress)
	if err != nil { return nil, err }

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil { return nil, err }

	return makeSocket(conn, addr), nil
}

// Close shuts the socket down. This will indirectly terminate the associated workers.
func (sock *Socket) Close() error {
	return sock.conn.Close()
}

// makeSocket configures the UDPConn and launches the receiver and sender workers.
func makeSocket(conn *net.UDPConn, addr *net.UDPAddr) *Socket {
	conn.SetDeadline(time.Time{})

	inbound := make(chan interface{})
	go socketReceiver(conn, inbound)

	outbound := make(chan OutgoingPayload)
	go socketSender(conn, outbound)

	return &Socket{conn, inbound, outbound}
}

// isSameUDPAddr tests the given UDPAddrs for equality.
func isSameUDPAddr(a, b *net.UDPAddr) bool {
	return bytes.Equal(a.IP, b.IP) && a.Port == b.Port
}

// socketReceiver is the receiver worker for Socket.
func socketReceiver(conn *net.UDPConn, inbound chan<- interface{}) {
	Logger.Printf("Socket[%v]: Started receiver", conn.RemoteAddr())

	buffer := [1024]byte{}
	reader := bytes.NewReader(buffer[:])

	for {
		len, _, err := conn.ReadFromUDP(buffer[:])
		if err != nil {
			Logger.Printf("Socket[%v]: Error during read: %v", conn.RemoteAddr(), err)
			break
		}

		Logger.Printf("Socket[%v]: Received: %v", conn.RemoteAddr(), buffer[:len])

		reader.Reset(buffer[:len])

		payload, err := ReadPacket(reader)
		if err != nil {
			Logger.Printf("Socket[%v]: Error during packet parsing: %v", conn.RemoteAddr(), err)
			continue
		}

		Logger.Printf("Socket[%v]: Inbound: %+v", conn.RemoteAddr(), payload)

		inbound <- payload
	}

	close(inbound)
	Logger.Printf("Socket[%v]: Stopped receiver", conn.RemoteAddr())
}

// socketSender is the sender worker for Socket.
func socketSender(conn *net.UDPConn, outbound <-chan OutgoingPayload) {
	Logger.Printf("Socket[%v]: Started sender", conn.RemoteAddr())

	buffer := &bytes.Buffer{}

	for payload := range outbound {
		Logger.Printf("Socket[%v]: Outbound: %+v", conn.RemoteAddr(), payload)

		buffer.Reset()

		err := WritePacket(buffer, payload)
		if err != nil {
			Logger.Printf("Socket[%v]: Error during packet generation: %v", conn.RemoteAddr(), err)
			continue
		}

		_, err = conn.Write(buffer.Bytes())
		if err != nil {
			Logger.Printf("Socket[%v]: Error during write: %v", conn.RemoteAddr(), err)
			break
		}

		Logger.Printf("Socket[%v]: Sending: %v", conn.RemoteAddr(), buffer.Bytes())
	}

	Logger.Printf("Socket[%v]: Stopped sender", conn.RemoteAddr())
}
