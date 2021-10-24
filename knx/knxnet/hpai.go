// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package knxnet

import (
	"errors"
	"fmt"
	"net"

	"github.com/vapourismo/knx-go/knx/util"
)

// Protocol specifies a host protocol to use.
type Protocol uint8

const (
	// UDP4 indicates a communication using UDP over IPv4.
	UDP4 Protocol = 1

	// TCP4 indicates a communication using TCP over IPv4.
	TCP4 Protocol = 2
)

// Address is an IPv4 address.
type Address [4]byte

// String formats the address.
func (addr Address) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", addr[0], addr[1], addr[2], addr[3])
}

// Port is a port number.
type Port uint16

// HostInfo contains information about a host.
type HostInfo struct {
	Protocol Protocol
	Address  Address
	Port     Port
}

// HostInfoFromAddress returns HostInfo from a UDP address
func HostInfoFromAddress(address *net.UDPAddr) (HostInfo, error) {
	ipAddress := address.IP.To4()

	if ipAddress == nil {
		return HostInfo{}, errors.New("unsupported local address")
	}

	var localIP Address
	copy(localIP[:], ipAddress)

	port := (Port)(address.Port)
	return HostInfo{Protocol: UDP4, Address: localIP, Port: port}, nil
}

// Equals checks whether both structures are equal.
func (info HostInfo) Equals(other HostInfo) bool {
	return info.Protocol == other.Protocol &&
		info.Address == other.Address &&
		info.Port == other.Port
}

// Size returns the packed size.
func (HostInfo) Size() uint {
	return 8
}

// Pack assembles the host info structure in the given buffer.
func (info *HostInfo) Pack(buffer []byte) {
	util.PackSome(
		buffer,
		byte(8),
		uint8(info.Protocol),
		info.Address[:],
		uint16(info.Port),
	)
}

// Unpack parses the given data in order to initialize the structure.
func (info *HostInfo) Unpack(data []byte) (n uint, err error) {
	var length uint8

	if n, err = util.UnpackSome(
		data, &length, (*uint8)(&info.Protocol), info.Address[:4], (*uint16)(&info.Port),
	); err != nil {
		return
	}

	if length != 8 {
		return n, errors.New("Host info structure length is invalid")
	}

	return
}
