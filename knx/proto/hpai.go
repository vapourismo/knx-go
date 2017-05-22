package proto

import (
	"errors"
	"io"
	"github.com/vapourismo/knx-go/knx/encoding"
)

// Address is a IPv4 address.
type Address [4]byte

// Port is a port number.
type Port uint16

// HostInfo contains information about a host.
type HostInfo struct {
	Address Address
	Port    Port
}

// ReadFrom initializes the structure by reading from the given Reader.
func (info *HostInfo) ReadFrom(r io.Reader) (n int64, err error) {
	var length, proto uint8
	n, err = encoding.ReadSome(r, &length, &proto, &info.Address, &info.Port)
	if err != nil {
		return
	}

	if length != 8 {
		return n, errors.New("Host info structure length is invalid")
	}

	if proto != 1 {
		return n, errors.New("Host info protocol is not UDP")
	}

	return
}

// WriteTo serializes the structure and writes it to the given Writer.
func (info *HostInfo) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteSome(w, byte(8), byte(1), info.Address, info.Port)
}
