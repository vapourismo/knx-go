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

	return makeSocket(conn, nil), nil
}

// Close shuts the socket down. This will indirectly terminate the associated workers.
func (sock *Socket) Close() error {
	return sock.conn.Close()
}

// Send transmits a KNXnet/IP packet
func (sock *Socket) Send(payload OutgoingPayload) error {
	buffer := &bytes.Buffer{}

	// Packet serialization
	err := WritePacket(buffer, payload)
	if err != nil { return err }

	Logger.Printf("Socket[%v]: Outbound: %T %+v", sock.conn.RemoteAddr(), payload, payload)

	// Transmission of the buffer contents
	_, err = sock.conn.Write(buffer.Bytes())
	if err != nil { return err }

	// Logger.Printf("Socket[%v]: Sent: %v", sock.conn.RemoteAddr(), buffer.Bytes())

	return nil
}

// makeSocket configures the UDPConn and launches the receiver and sender workers.
func makeSocket(conn *net.UDPConn, addr *net.UDPAddr) *Socket {
	conn.SetDeadline(time.Time{})

	inbound := make(chan interface{})
	go socketReceiver(conn, addr, inbound)

	return &Socket{conn, inbound}
}

// socketReceiver is the receiver worker for Socket.
func socketReceiver(conn *net.UDPConn, addr *net.UDPAddr, inbound chan<- interface{}) {
	Logger.Printf("Socket[%v]: Started receiver", conn.RemoteAddr())

	buffer := [1024]byte{}
	reader := bytes.NewReader(buffer[:])

	for {
		len, sender, err := conn.ReadFromUDP(buffer[:])
		if err != nil {
			Logger.Printf("Socket[%v]: Error during read: %v", conn.RemoteAddr(), err)
			break
		}

		// Validate sender origin if necessary
		if addr != nil && (!addr.IP.Equal(sender.IP) || addr.Port != sender.Port) {
			Logger.Printf("Socket[%v]: Origin validation failed: %v (expected %v)",
			              conn.RemoteAddr(), sender, addr)
			continue
		}

		// Logger.Printf("Socket[%v]: Received: %v", conn.RemoteAddr(), buffer[:len])

		reader.Reset(buffer[:len])

		payload, err := ReadPacket(reader)
		if err != nil {
			Logger.Printf("Socket[%v]: Error during packet parsing: %v", conn.RemoteAddr(), err)
			continue
		}

		Logger.Printf("Socket[%v]: Inbound: %T %+v", conn.RemoteAddr(), payload, payload)

		inbound <- payload
	}

	close(inbound)
	Logger.Printf("Socket[%v]: Stopped receiver", conn.RemoteAddr())
}
