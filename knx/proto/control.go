package proto

import (
	"errors"
	"fmt"
	"io"
	"github.com/vapourismo/knx-go/knx/encoding"
)

// A ConnReq requests a connection to a gateway.
type ConnReq struct {
	Control HostInfo
	Tunnel  HostInfo
}

var connReqInfo = [4]byte{4, 4, 2, 0}

// WriteTo serializes the structure and writes it to the given Writer.
func (req *ConnReq) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteSome(w, &req.Control, &req.Tunnel, connReqInfo)
}

// ConnResStatus is the type of status code carried in a connection response.
type ConnResStatus uint8

// Therese are known connection response status codes.
const (
	ConnResOk                ConnResStatus = 0x00
	ConnResUnsupportedType   ConnResStatus = 0x22
	ConnResUnsupportedOption ConnResStatus = 0x23
	ConnResBusy              ConnResStatus = 0x24
)

// String describes the status code.
func (status ConnResStatus) String() string {
	switch status {
	case ConnResOk:
		return "Connection established"

	case ConnResUnsupportedType:
		return "Requested connection type is unsupported"

	case ConnResUnsupportedOption:
		return "One of the requested options is unsupported"

	case ConnResBusy:
		return "No data channel is available"

	default:
		return fmt.Sprintf("Unknown status code %#x", uint8(status))
	}
}

// Error implements the error Error method.
func (status ConnResStatus) Error() string {
	return status.String()
}

// ConnRes is a response to a ConnReq.
type ConnRes struct {
	Channel uint8
	Status  ConnResStatus
	Control HostInfo
}

// ReadFrom initializes the structure by reading from the given Reader.
func (res *ConnRes) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadSome(r, &res.Channel, &res.Status, &res.Control)
}

// A ConnStateReq requests the the connection state from a gateway.
type ConnStateReq struct {
	Channel uint8
	Status  uint8
	Control HostInfo
}

// WriteTo serializes the structure and writes it to the given Writer.
func (req *ConnStateReq) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteSome(w, req.Channel, req.Status, &req.Control)
}

// A ConnState represents the state of a connection.
type ConnState uint8

// These are known connection states.
const (
	ConnStateNormal    ConnState = 0x00
	ConnStateInactive  ConnState = 0x21
	ConnStateDataError ConnState = 0x26
	ConnStateKNXError  ConnState = 0x27
)

// String converts the connection state to a string.
func (state ConnState) String() string {
	switch state {
	case ConnStateNormal:
		return "Connection is intact"

	case ConnStateInactive:
		return "Connection is inactive"

	case ConnStateDataError:
		return "Gateway encountered a data error"

	case ConnStateKNXError:
		return "Gateway encountered a KNX error"

	default:
		return fmt.Sprintf("Unknown connection state %#x", uint8(state))
	}
}

// Error implements the error Error method.
func (state ConnState) Error() string {
	return state.String()
}

// A ConnStateRes is a response to a ConnStateReq.
type ConnStateRes struct {
	Channel uint8
	Status  ConnState
}

// ReadFrom initializes the structure by reading from the given Reader.
func (res *ConnStateRes) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadSome(r, &res.Channel, &res.Status)
}

// A DiscReq requests a connection to be terminated.
type DiscReq struct {
	Channel uint8
	Status  uint8
	Control HostInfo
}

// ReadFrom initializes the structure by reading from the given Reader.
func (req *DiscReq) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadSome(r, &req.Channel, &req.Status, &req.Control)
}

// WriteTo serializes the structure and writes it to the given Writer.
func (req *DiscReq) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteSome(w, req.Channel, req.Status, &req.Control)
}

// A DiscRes is a response to a DiscReq..
type DiscRes struct {
	Channel uint8
	Status  uint8
}

// ReadFrom initializes the structure by reading from the given Reader.
func (res *DiscRes) ReadFrom(r io.Reader) (int64, error) {
	return encoding.ReadSome(r, &res.Channel, &res.Status)
}

// WriteTo serializes the structure and writes it to the given Writer.
func (res *DiscRes) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteSome(w, res.Channel, res.Status)
}
