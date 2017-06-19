// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package knxnet

import "fmt"

// What follows are errors codes defined in the KNX standard.
const (
	// NoError indicates a successful operation.
	NoError = 0x00

	// ErrHostProtocolType indicates an unsupported host protocol.
	ErrHostProtocolType = 0x01

	// ErrVersionNotSupported indicates an unsupported KNXnet/IP protocol version.
	ErrVersionNotSupported = 0x02

	// ErrSequenceNumber indicates that an out-of-order sequence number has been received.
	ErrSequenceNumber = 0x04

	// ErrConnectionID indicates that there is no active data connection with given ID.
	ErrConnectionID = 0x21

	// ErrConnectionType indicates an unsupported connection type.
	ErrConnectionType = 0x22

	// ErrConnectionOption indicates an unsupported connection option.
	ErrConnectionOption = 0x23

	// ErrNoMoreConnections is returned by a Tunnelling Server when it cannot accept more
	// connections.
	ErrNoMoreConnections = 0x24

	// ErrNoMoreUniqueConnections is returned by a Tunnelling Server when it has no free Individual
	// Address available that could be used by the connection.
	ErrNoMoreUniqueConnections = 0x25

	// ErrDataConnection indicates an error with a data connection.
	ErrDataConnection = 0x26

	// ErrKNXConnection indicates an error with a KNX connection.
	ErrKNXConnection = 0x27

	// ErrTunnellingLayer indicates an unsupported tunnelling layer.
	ErrTunnellingLayer = 0x29
)

// A ErrCode identifies an error type.
type ErrCode uint8

// String returns a string representation of the error code.
func (err ErrCode) String() string {
	switch err {
	case NoError:
		return "No error"

	case ErrHostProtocolType:
		return "Host protocol is not supported"

	case ErrVersionNotSupported:
		return "KNXnet/IP version is not supported"

	case ErrSequenceNumber:
		return "Sequence number is out-of-order"

	case ErrConnectionID:
		return "No active data connection"

	case ErrConnectionType:
		return "Unsupported connection type"

	case ErrConnectionOption:
		return "Unsupported connection option"

	case ErrNoMoreConnections:
		return "No more connections available"

	case ErrNoMoreUniqueConnections:
		return "No more unique connections available"

	case ErrDataConnection:
		return "Data connection error"

	case ErrKNXConnection:
		return "KNX connection error"

	case ErrTunnellingLayer:
		return "Unsupported tunnelling layer"

	default:
		return fmt.Sprintf("Unknown error code %#x", err)
	}
}

// Error implements the error interface
func (err ErrCode) Error() string {
	return err.String()
}
