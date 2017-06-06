package proto

import (
	"errors"
	"fmt"

	"github.com/vapourismo/knx-go/knx/util"
)

// TunnelLayer identifies the tunnelling layer for a tunnelling connection.
type TunnelLayer uint8

const (
	// TunnelLayerData establishes a data-link layer tunnel.
	TunnelLayerData TunnelLayer = 0x02

	// TunnelLayerRaw establishes a raw tunnel.
	TunnelLayerRaw TunnelLayer = 0x04

	// TunnelLayerBusmon establishes a bus monitor tunnel.
	TunnelLayerBusmon TunnelLayer = 0x80
)

// A ConnReq requests a connection to a gateway.
type ConnReq struct {
	Control HostInfo
	Tunnel  HostInfo
	Layer   TunnelLayer
}

// Service returns the service identifier for connection requests.
func (ConnReq) Service() ServiceID {
	return ConnReqService
}

var hostInfoSize = HostInfo{}.Size()

// Size returns the packed size.
func (ConnReq) Size() uint {
	return 2*hostInfoSize + 4
}

// Pack the structure into the given buffer.
func (req *ConnReq) Pack(buffer []byte) {
	util.PackSome(buffer, &req.Control, &req.Tunnel)

	buffer = buffer[2*hostInfoSize:]
	buffer[0] = 4
	buffer[1] = 4
	buffer[2] = byte(req.Layer)
	buffer[3] = 0
}

// Unpack initializes the structure by parsing the given data.
func (req *ConnReq) Unpack(data []byte) (n uint, err error) {
	var length, connType, reserved uint8

	n, err = util.UnpackSome(
		data, &req.Control, &req.Tunnel, &length, &connType, (*uint8)(&req.Layer), &reserved,
	)
	if err != nil {
		return
	}

	if length != 4 {
		return n, errors.New("Invalid connection request info structure length")
	}

	if connType != 4 {
		return n, errors.New("Invalid connection type")
	}

	return
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

// Service returns the service identifier for connection responses.
func (ConnRes) Service() ServiceID {
	return ConnResService
}

// Unpack initializes the structure by parsing the given data.
func (res *ConnRes) Unpack(data []byte) (uint, error) {
	return util.UnpackSome(data, &res.Channel, (*uint8)(&res.Status), &res.Control)
}

// A ConnStateReq requests the the connection state from a gateway.
type ConnStateReq struct {
	Channel uint8
	Status  uint8
	Control HostInfo
}

// Service returns the service identifier for connection state requests.
func (ConnStateReq) Service() ServiceID {
	return ConnStateReqService
}

// Size returns the packed size.
func (ConnStateReq) Size() uint {
	return 2 + hostInfoSize
}

// Pack the structure into the buffer.
func (req *ConnStateReq) Pack(buffer []byte) {
	buffer[0] = req.Channel
	buffer[1] = req.Status
	req.Control.Pack(buffer[2:])
}

// Unpack initializes the structure by parsing the given data.
func (req *ConnStateReq) Unpack(data []byte) (uint, error) {
	return util.UnpackSome(data, &req.Channel, &req.Status, &req.Control)
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

// A ConnStateRes is a response to a ConnStateReq.
type ConnStateRes struct {
	Channel uint8
	Status  ConnState
}

// Service returns the service identifier for connection state responses.
func (ConnStateRes) Service() ServiceID {
	return ConnStateResService
}

// Size returns the packed size.
func (ConnStateRes) Size() uint {
	return 2
}

// Pack the structure into the buffer.
func (res *ConnStateRes) Pack(buffer []byte) {
	buffer[0] = res.Channel
	buffer[1] = byte(res.Status)
}

// Unpack initializes the structure by parsing the given data.
func (res *ConnStateRes) Unpack(data []byte) (uint, error) {
	return util.UnpackSome(data, &res.Channel, (*uint8)(&res.Status))
}

// A DiscReq requests a connection to be terminated.
type DiscReq struct {
	Channel uint8
	Status  uint8
	Control HostInfo
}

// Service returns the service identifier for disconnect requests.
func (DiscReq) Service() ServiceID {
	return DiscReqService
}

// Size returns the packed size.
func (DiscReq) Size() uint {
	return 2 + hostInfoSize
}

// Pack the structure into the buffer.
func (req *DiscReq) Pack(buffer []byte) {
	buffer[0] = req.Channel
	buffer[1] = req.Status
	req.Control.Pack(buffer[2:])
}

// Unpack initializes the structure by parsing the given data.
func (req *DiscReq) Unpack(data []byte) (uint, error) {
	return util.UnpackSome(data, &req.Channel, &req.Status, &req.Control)
}

// A DiscRes is a response to a DiscReq..
type DiscRes struct {
	Channel uint8
	Status  uint8
}

// Service returns the service identifier for disconnect responses.
func (DiscRes) Service() ServiceID {
	return DiscResService
}

// Size returns the packed size.
func (DiscRes) Size() uint {
	return 2
}

// Pack the structure into the buffer.
func (res *DiscRes) Pack(data []byte) {
	data[0] = res.Channel
	data[1] = res.Status
}

// Unpack initializes the structure by parsing the given data.
func (res *DiscRes) Unpack(data []byte) (uint, error) {
	return util.UnpackSome(data, &res.Channel, &res.Status)
}
