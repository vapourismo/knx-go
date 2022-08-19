// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package knxnet

import (
	"errors"
	"fmt"
	"net"
	"strconv"

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

// HostInfoFromAddress returns HostInfo from an address.
func HostInfoFromAddress(address net.Addr) (HostInfo, error) {
	hostinfo := HostInfo{}

	ipS, portS, err := net.SplitHostPort(address.String())
	if err != nil {
		return hostinfo, err
	}

	ip := net.ParseIP(ipS)
	if ip == nil {
		return hostinfo, fmt.Errorf("unable to determine IP")
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return hostinfo, fmt.Errorf("only IPv4 is currently supported")
	}

	port, _ := strconv.ParseUint(portS, 10, 16)
	if port == 0 {
		return hostinfo, fmt.Errorf("unable to determine port")
	}

	copy(hostinfo.Address[:], ipv4)
	hostinfo.Port = Port(port)

	switch address.Network() {
	case "udp":
		hostinfo.Protocol = UDP4
	case "tcp":
		hostinfo.Protocol = TCP4
	default:
		return hostinfo, fmt.Errorf("unsupported network")
	}

	return hostinfo, nil
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

// Pack assembles the Host Info structure in the given buffer.
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
		return n, errors.New("host info structure length is invalid")
	}

	return
}
