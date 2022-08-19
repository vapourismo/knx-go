// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package knxnet

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/vapourismo/knx-go/knx/util"
	"golang.org/x/net/ipv4"
)

// A Socket is a socket, duh.
type Socket interface {
	Send(payload ServicePackable) error
	Inbound() <-chan Service
	Close() error
	LocalAddr() net.Addr
}

// TunnelSocket is a UDP socket for KNXnet/IP packet exchange.
type TunnelSocket struct {
	conn    net.Conn
	inbound <-chan Service
}

// DialTunnelUDP creates a new Socket which can used to exchange KNXnet/IP packets with a single
// endpoint through UDP.
func DialTunnelUDP(address string) (*TunnelSocket, error) {
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return nil, err
	}

	if addr.IP.IsMulticast() {
		return nil, fmt.Errorf("cannot tunnel to multicast address")
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return nil, err
	}

	conn.SetDeadline(time.Time{})

	inbound := make(chan Service)
	go serveUDPSocket(conn, addr, inbound)

	return &TunnelSocket{conn, inbound}, nil
}

// DialTunnelTCP creates a new Socket which can used to exchange KNXnet/IP packets with a single
// endpoint through TCP.
func DialTunnelTCP(address string) (*TunnelSocket, error) {
	addr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, err
	}

	if addr.IP.IsMulticast() {
		return nil, fmt.Errorf("cannot tunnel to multicast address")
	}

	conn, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		return nil, err
	}

	conn.SetDeadline(time.Time{})

	inbound := make(chan Service)
	go serveTCPSocket(conn, addr, inbound)

	return &TunnelSocket{conn, inbound}, nil
}

// Send transmits a KNXnet/IP packet.
func (sock *TunnelSocket) Send(payload ServicePackable) error {
	buffer := make([]byte, Size(payload))
	Pack(buffer, payload)

	// Transmission of the buffer contents.
	_, err := sock.conn.Write(buffer)
	return err
}

// Inbound provides a channel from which you can retrieve incoming packets.
func (sock *TunnelSocket) Inbound() <-chan Service {
	return sock.inbound
}

// Close shuts the socket down. This will indirectly terminate the associated workers.
func (sock *TunnelSocket) Close() error {
	return sock.conn.Close()
}

// LocalAddr returns the local UDP address.
func (sock *TunnelSocket) LocalAddr() net.Addr {
	return sock.conn.LocalAddr()
}

// RouterSocket is a UDP socket for KNXnet/IP packet exchange.
type RouterSocket struct {
	conn    *net.UDPConn
	addr    *net.UDPAddr
	inbound <-chan Service
}

// ListenRouter creates a new Socket which can be used to exchange KNXnet/IP packets with
// multiple endpoints.
func ListenRouter(multicastAddress string) (*RouterSocket, error) {
	return ListenRouterOnInterface(nil, multicastAddress, false)
}

// ListenRouterOnInterface creates a new Socket which can be used to exchange KNXnet/IP packets with
// multiple endpoints. The interface is used to send or listen for KNXnet/IP packets. If the
// interface is nil, the system-assigned multicast interface is used.
func ListenRouterOnInterface(ifi *net.Interface, multicastAddress string, multicastLoopbackEnabled bool) (*RouterSocket, error) {
	addr, err := net.ResolveUDPAddr("udp4", multicastAddress)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return nil, err
	}
	pc := ipv4.NewPacketConn(conn)

	if err := pc.JoinGroup(ifi, addr); err != nil {
		return nil, err
	}

	// Just for logging purposes.
	if loopOn, err := pc.MulticastLoopback(); err == nil {
		util.Log(conn, "MulticastLoopback status: %v", loopOn)
	}
	// Setup interface with Multicast Loopback enabled if desired.
	if err := pc.SetMulticastLoopback(multicastLoopbackEnabled); err != nil {
		util.Log(conn, "SetMulticastLoopback error: %v", err)
	} else {
		util.Log(conn, "MulticastLoopbackEnabled: %t", multicastLoopbackEnabled)
	}

	conn.SetDeadline(time.Time{})

	inbound := make(chan Service)
	go serveUDPSocket(conn, nil, inbound)

	return &RouterSocket{conn, addr, inbound}, nil
}

// Addr returns the multicast destination address.
func (sock *RouterSocket) Addr() *net.UDPAddr {
	return sock.addr
}

// Send transmits a KNXnet/IP packet.
func (sock *RouterSocket) Send(payload ServicePackable) error {
	buffer := make([]byte, Size(payload))
	Pack(buffer, payload)

	// Transmission of the buffer contents
	_, err := sock.conn.WriteToUDP(buffer, sock.addr)
	return err
}

// Inbound provides a channel from which you can retrieve incoming packets.
func (sock *RouterSocket) Inbound() <-chan Service {
	return sock.inbound
}

// Close shuts the socket down. This will indirectly terminate the associated workers.
func (sock *RouterSocket) Close() error {
	return sock.conn.Close()
}

// LocalAddr returns the local UDP address.
func (sock *RouterSocket) LocalAddr() net.Addr {
	return sock.conn.LocalAddr()
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

		// Discard empty frames
		if len == 0 {
			util.Log(conn, "Empty frame discarded")
			continue
		}

		// Validate sender origin if necessary.
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

// serveTCPSocket is the receiver worker for a TCP socket.
func serveTCPSocket(conn *net.TCPConn, addr *net.TCPAddr, inbound chan<- Service) {
	util.Log(conn, "Started worker")
	defer util.Log(conn, "Worker exited")

	// A closed inbound channel indicates to its readers that the worker has terminated.
	defer close(inbound)

	connBuffer := bufio.NewReader(conn)

	for {
		header, err := connBuffer.Peek(6) // KNXnet/IP headers are 6 bytes long
		if err != nil {
			util.Log(conn, "Error during peeking header: %v", err)
			return
		}

		var serviceID ServiceID
		var totalLen uint16

		_, err = UnpackHeader(header, &serviceID, &totalLen)
		if err != nil {
			util.Log(conn, "Error during header inspection: %v", err)
			return
		}

		buffer := make([]byte, totalLen)
		len, err := io.ReadFull(connBuffer, buffer)
		if err != nil {
			util.Log(conn, "Error during ReadFull: %v", err)
			return
		}

		// Discard empty frames
		if len == 0 {
			util.Log(conn, "Empty frame discarded")
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
