package proto

import (
	"errors"
	"fmt"
	"io"

	"github.com/vapourismo/knx-go/knx/encoding"
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

// Address is a IPv4 address.
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

// Equals checks whether both structures are equal.
func (info HostInfo) Equals(other HostInfo) bool {
	return info.Protocol == other.Protocol &&
		info.Address == other.Address &&
		info.Port == other.Port
}

// ReadFrom initializes the structure by reading from the given Reader.
func (info *HostInfo) ReadFrom(r io.Reader) (n int64, err error) {
	var length uint8
	n, err = encoding.ReadSome(r, &length, &info.Protocol, &info.Address, &info.Port)
	if err != nil {
		return
	}

	if length != 8 {
		return n, errors.New("Host info structure length is invalid")
	}

	switch info.Protocol {
	case UDP4, TCP4:
		return

	default:
		return n, errors.New("Unknown host protocol")
	}
}

// Unpack initializes the structure by parsing the given data.
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

// WriteTo serializes the structure and writes it to the given Writer.
func (info HostInfo) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteSome(w, byte(8), info.Protocol, info.Address, info.Port)
}
