package knx

import (
	"bytes"
	"net"
	"time"
)

// UDP socket for KNXnet/IP packet exchange
type Socket struct {
	conn *net.UDPConn

	// Inbound is a channel that relays incoming KNXnet/IP packets.
	// The types of these packets are limited to those returned by ReadPacket.
	Inbound <-chan interface{}
}

// NewClientSocket creates a new Socket which can used to exchange KNXnet/IP packets with a gateway.
func NewClientSocket(gatewayAddress string) (*Socket, error) {
	addr, err := net.ResolveUDPAddr("udp4", gatewayAddress)
	if err != nil { return nil, err }

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil { return nil, err }

	return makeSocket(conn)
}

// NewRoutingSocket creates a new Socket which can be used to exchange KNXnet/IP packets with a
// router.
func NewRoutingSocket(multicastAddress string) (*Socket, error) {
	addr, err := net.ResolveUDPAddr("udp4", multicastAddress)
	if err != nil { return nil, err }

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil { return nil, err }

	return makeSocket(conn)
}

// Close closes the socket.
func (sock *Socket) Close() error {
	return sock.conn.Close()
}

// Send transmits a packet.
func (sock *Socket) Send(payload OutgoingPayload) error {
	buffer := bytes.NewBuffer(make([]byte, 0, 64))

	err := WritePacket(buffer, payload)
	if err != nil { return err }

	_, err = sock.conn.Write(buffer.Bytes())
	return err
}

func makeSocket(conn *net.UDPConn) (*Socket, error) {
	conn.SetDeadline(time.Time{})

	inbound := make(chan interface{}, 10)
	go socketWorker(conn, inbound)

	return &Socket{conn, inbound}, nil
}

func socketWorker(conn *net.UDPConn, inbound chan<- interface{}) {
	buffer := make([]byte, 1024)
	reader := bytes.NewReader(buffer)

	for {
		len, err := conn.Read(buffer[:1024])

		if err != nil {
			logf("socket[%v]: Error during read: %v", conn.RemoteAddr(), err)
			break
		}

		reader.Reset(buffer[:len])
		payload, err := ReadPacket(reader)

		if err != nil {
			logf("socket[%v]: Error during parsing: %v", conn.RemoteAddr(), err)
			continue
		}

		inbound <- payload
	}

	close(inbound)
}
