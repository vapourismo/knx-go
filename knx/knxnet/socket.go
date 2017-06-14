// Copyright 2017 Ole Krüger.

package knxnet

import (
	"net"
	"time"

	"github.com/vapourismo/knx-go/knx/util"
)

// A Socket is a socket, duh.
type Socket interface {
	Send(payload ServicePackable) error
	Inbound() <-chan Service
	Close() error
}

// UnicastSocket is a UDP socket for KNXnet/IP packet exchange.
type UnicastSocket struct {
	conn    *net.UDPConn
	inbound <-chan Service
}

// NewUnicastSocket creates a new Socket which can used to exchange KNXnet/IP packets with a single
// endpoint.
func NewUnicastSocket(address string) (*UnicastSocket, error) {
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return nil, err
	}

	conn.SetDeadline(time.Time{})

	inbound := make(chan Service)
	go serveUDPSocket(conn, addr, inbound)

	return &UnicastSocket{conn, inbound}, nil
}

// Send transmits a KNXnet/IP packet.
func (sock *UnicastSocket) Send(payload ServicePackable) error {
	buffer := make([]byte, Size(payload))
	Pack(buffer, payload)

	// Transmission of the buffer contents
	_, err := sock.conn.Write(buffer)
	return err
}

// Inbound provides a channel from which you can retrieve incoming packets.
func (sock *UnicastSocket) Inbound() <-chan Service {
	return sock.inbound
}

// Close shuts the socket down. This will indirectly terminate the associated workers.
func (sock *UnicastSocket) Close() error {
	return sock.conn.Close()
}

// MulticastSocket is a UDP socket for KNXnet/IP packet exchange.
type MulticastSocket struct {
	conn    *net.UDPConn
	addr    *net.UDPAddr
	inbound <-chan Service
}

// NewMulticastSocket creates a new Socket which can be used to exchange KNXnet/IP packets with
// multiple endpoints.
func NewMulticastSocket(multicastAddress string) (*MulticastSocket, error) {
	addr, err := net.ResolveUDPAddr("udp4", multicastAddress)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return nil, err
	}

	conn.SetDeadline(time.Time{})

	inbound := make(chan Service)
	go serveUDPSocket(conn, nil, inbound)

	return &MulticastSocket{conn, addr, inbound}, nil
}

// Send transmits a KNXnet/IP packet.
func (sock *MulticastSocket) Send(payload ServicePackable) error {
	buffer := make([]byte, Size(payload))
	Pack(buffer, payload)

	// Transmission of the buffer contents
	_, err := sock.conn.WriteToUDP(buffer, sock.addr)
	return err
}

// Inbound provides a channel from which you can retrieve incoming packets.
func (sock *MulticastSocket) Inbound() <-chan Service {
	return sock.inbound
}

// Close shuts the socket down. This will indirectly terminate the associated workers.
func (sock *MulticastSocket) Close() error {
	return sock.conn.Close()
}

// serveUDPSocket is the receiver worker for a UDP socket.
func serveUDPSocket(conn *net.UDPConn, addr *net.UDPAddr, inbound chan<- Service) {
	util.Log(conn, "Started worker")
	defer util.Log(conn, "Worker exited")

	// A closed inbound channel indicates to its readers that the worker has terminated.
	defer close(inbound)

	buffer := [1024]byte{}

	for {
		len, sender, err := conn.ReadFromUDP(buffer[:])
		if err != nil {
			util.Log(conn, "Error during ReadFromUDP: %v", err)
			return
		}

		// Validate sender origin if necessary
		if addr != nil && (!addr.IP.Equal(sender.IP) || addr.Port != sender.Port) {
			util.Log(conn, "Origin validation failed: %v != %v", addr, sender)
			continue
		}

		var payload Service
		_, err = Unpack(buffer[:len], &payload)
		if err != nil {
			util.Log(conn, "Error during Unpack: %v", err)
			continue
		}

		inbound <- payload
	}
}
