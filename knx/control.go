package knx

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/vapourismo/knx-go/knx/encoding"
)

// Address is a IPv4 address.
type Address [4]byte

// Port is a port number.
type Port uint16

// A HostInfo contains information about a host.
type HostInfo struct {
	Address Address
	Port    Port
}

var hostInfoPreamble = [2]byte{8, 1}

func (info HostInfo) writeTo(w *bytes.Buffer) error {
	return encoding.WriteSequence(w, hostInfoPreamble, info)
}

func readHostInfo(r *bytes.Reader) (*HostInfo, error) {
	var length, proto byte
	info := &HostInfo{}

	err := encoding.ReadSequence(r, &length, &proto, info)
	if err != nil {
		return nil, err
	}

	if length != 8 {
		return nil, errors.New("Host info structure length is invalid")
	}

	if proto != 1 {
		return nil, errors.New("Host info protocol is unknown")
	}

	return info, nil
}

// A ConnectionRequest requests a connection to a gateway.
type ConnectionRequest struct {
	Control HostInfo
	Tunnel  HostInfo
}

func (req ConnectionRequest) describe() (ServiceID, int) {
	return ConnReqService, 20
}

var connReqInfo = [4]byte{4, 4, 2, 0}

func (req ConnectionRequest) writeTo(w *bytes.Buffer) error {
	err := req.Control.writeTo(w)
	if err != nil {
		return err
	}

	err = req.Tunnel.writeTo(w)
	if err != nil {
		return err
	}

	_, err = w.Write(connReqInfo[:])
	return err
}

// ConnResStatus is the type of status code carried in a connection response.
type ConnResStatus uint8

// Potential connection response status codes.
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

// ConnectionResponse is a response to a ConnectionRequest.
type ConnectionResponse struct {
	Channel uint8
	Status  ConnResStatus
	Host    HostInfo
}

func readConnectionResponse(r *bytes.Reader) (*ConnectionResponse, error) {
	var channel uint8
	var status ConnResStatus

	err := encoding.ReadSequence(r, &channel, &status)
	if err != nil {
		return nil, err
	}

	host, err := readHostInfo(r)
	if err != nil {
		return nil, err
	}

	return &ConnectionResponse{channel, status, *host}, nil
}

// A ConnStateReq requests the the connection state from a gateway.
type ConnStateReq struct {
	Channel byte
	Status  byte
	Host    HostInfo
}

func (req ConnStateReq) describe() (ServiceID, int) {
	return ConnStateReqService, 10
}

func (req ConnStateReq) writeTo(w *bytes.Buffer) error {
	err := encoding.WriteSequence(w, req.Channel, req.Status)
	if err != nil {
		return err
	}

	return req.Host.writeTo(w)
}

// A ConnState represents the state of a connection.
type ConnState uint8

//
const (
	ConnStateNormal    ConnState = 0x00
	ConnStateInactive  ConnState = 0x21
	ConnStateDataError ConnState = 0x26
	ConnStateKNXError  ConnState = 0x27
)

//
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

//
func (state ConnState) Error() string {
	return state.String()
}

// A ConnStateRes is a response to a ConnectionStateRequest.
type ConnStateRes struct {
	Channel uint8
	Status  ConnState
}

func readConnStateRes(r *bytes.Reader) (*ConnStateRes, error) {
	res := &ConnStateRes{}

	err := encoding.ReadSequence(r, &res.Channel, &res.Status)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// A DiscReq requests a connection to be terminated.
type DiscReq struct {
	Channel byte
	Status  byte
	Host    HostInfo
}

func readDiscReq(r *bytes.Reader) (*DiscReq, error) {
	var channel, status byte

	err := encoding.ReadSequence(r, &channel, &status)
	if err != nil {
		return nil, err
	}

	host, err := readHostInfo(r)
	if err != nil {
		return nil, err
	}

	return &DiscReq{channel, status, *host}, nil
}

func (req DiscReq) describe() (ServiceID, int) {
	return DiscReqService, 10
}

func (req DiscReq) writeTo(w *bytes.Buffer) error {
	err := encoding.WriteSequence(w, req.Channel, req.Status)
	if err != nil {
		return err
	}

	return req.Host.writeTo(w)
}

// A DiscRes is a response to a DisconnectRequest.
type DiscRes struct {
	Channel uint8
	Status  uint8
}

func (res DiscRes) describe() (ServiceID, int) {
	return DiscResService, 2
}

func (res DiscRes) writeTo(w *bytes.Buffer) error {
	return encoding.WriteSequence(w, res.Channel, res.Status)
}

func readDiscRes(r *bytes.Reader) (*DiscRes, error) {
	res := &DiscRes{}

	err := encoding.ReadSequence(r, &res.Channel, &res.Status)
	if err != nil {
		return nil, err
	}

	return res, nil
}

