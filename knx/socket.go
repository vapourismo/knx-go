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

	return makeSocket(conn, addr)
}

// NewRoutingSocket creates a new Socket which can be used to exchange KNXnet/IP packets with a
// router.
func NewRoutingSocket(multicastAddress string) (*Socket, error) {
	addr, err := net.ResolveUDPAddr("udp4", multicastAddress)
	if err != nil { return nil, err }

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil { return nil, err }

	return makeSocket(conn, addr)
}

// Close shuts the socket down. This will indirectly terminate the associated workers.
func (sock *Socket) Close() error {
	return sock.conn.Close()
}

// makeSocket configures the UDPConn and launches the receiver and sender workers.
func makeSocket(conn *net.UDPConn, addr *net.UDPAddr) (*Socket, error) {
	conn.SetDeadline(time.Time{})

	inbound := make(chan interface{})
	go socketReceiver(conn, addr, inbound)

	outbound := make(chan OutgoingPayload)
	go socketSender(conn, addr, outbound)

	return &Socket{conn, inbound, outbound}, nil
}

// isSameUDPAddr tests the given UDPAddrs for equality.
func isSameUDPAddr(a, b *net.UDPAddr) bool {
	return bytes.Equal(a.IP, b.IP) && a.Port == b.Port
}

// socketReceiver is the receiver worker for Socket.
func socketReceiver(conn *net.UDPConn, addr *net.UDPAddr, inbound chan<- interface{}) {
	Logger.Printf("socket[%v]: Started receiver", conn.RemoteAddr())

	buffer := [1024]byte{}
	reader := bytes.NewReader(buffer[:])

	for {
		len, sender, err := conn.ReadFromUDP(buffer[:])
		if err != nil {
			Logger.Printf("socket[%v]: Error during read: %v", conn.RemoteAddr(), err)
			break
		} else if !isSameUDPAddr(sender, addr) {
			Logger.Printf("socket[%v]: Packet from invalid sender: %+v", conn.RemoteAddr(), sender)
			continue
		}

		reader.Reset(buffer[:len])

		payload, err := ReadPacket(reader)
		if err != nil {
			Logger.Printf("socket[%v]: Error during packet parsing: %v\n" +
			              "            Buffer was: %v",
			              conn.RemoteAddr(), err, buffer[:len])
			continue
		}

		Logger.Printf("socket[%v]: Received: %+v", conn.RemoteAddr(), payload)

		inbound <- payload
	}

	close(inbound)
	Logger.Printf("socket[%v]: Stopped receiver", conn.RemoteAddr())
}

// socketSender is the sender worker for Socket.
func socketSender(conn *net.UDPConn, addr *net.UDPAddr, outbound <-chan OutgoingPayload) {
	Logger.Printf("socket[%v]: Started sender", conn.RemoteAddr())

	buffer := &bytes.Buffer{}

	for payload := range outbound {
		buffer.Reset()

		err := WritePacket(buffer, payload)
		if err != nil {
			Logger.Printf("socket[%v]: Error during packet generation: %v", conn.RemoteAddr(), err)
			continue
		}

		Logger.Printf("socket[%v]: Sending: %+v", conn.RemoteAddr(), payload)

		_, err = conn.WriteToUDP(buffer.Bytes(), addr)
		if err != nil {
			Logger.Printf("socket[%v]: Error during write: %v", conn.RemoteAddr(), err)
			break
		}
	}

	Logger.Printf("socket[%v]: Stopped sender", conn.RemoteAddr())
}
