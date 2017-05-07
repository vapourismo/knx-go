package knx

import (
	"bytes"
	"fmt"
	"errors"
)

// IPv4 address
type Address [4]byte

// Port number
type Port uint16

// Host information
type HostInfo struct {
	Address Address
	Port    Port
}

var hostInfoPreamble = [2]byte{8, 1}

func (info HostInfo) writeTo(w *bytes.Buffer) error {
	return writeSequence(w, hostInfoPreamble, info)
}

func readHostInfo(r *bytes.Reader) (*HostInfo, error) {
	var length, proto byte
	info := &HostInfo{}

	err := readSequence(r, &length, &proto, info)
	if err != nil { return nil, err }

	if length != 8 {
		return nil, errors.New("Host info structure length is invalid")
	}

	if proto != 1 {
		return nil, errors.New("Host info protocol is unknown")
	}

	return info, nil
}

// Connection request
type ConnectionRequest struct {
	Control HostInfo
	Tunnel  HostInfo
}

func (req ConnectionRequest) describe() (serviceIdent, int) {
	return connectionRequestService, 20
}

var connReqInfo = [4]byte{4, 4, 2, 0}

func (req ConnectionRequest) writeTo(w *bytes.Buffer) error {
	err := req.Control.writeTo(w)
	if err != nil { return err }

	err = req.Tunnel.writeTo(w)
	if err != nil { return err }

	_, err = w.Write(connReqInfo[:])
	return err
}

// ConnResStatus is the type of status code carried in a connection response.
type ConnResStatus uint8

// Potential connection response status codes.
var (
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
		return fmt.Sprintf("Unknown status code %#x", status)
	}
}

// Error implements the error Error method.
func (status ConnResStatus) Error() string {
	return status.String()
}

// Connection response
type ConnectionResponse struct {
	Channel uint8
	Status  ConnResStatus
	Host    HostInfo
}

func readConnectionResponse(r *bytes.Reader) (*ConnectionResponse, error) {
	var channel uint8
	var status ConnResStatus

	err := readSequence(r, &channel, &status)
	if err != nil { return nil, err }

	switch status {
		case ConnResOk, ConnResUnsupportedType, ConnResUnsupportedOption, ConnResBusy:

		default:
			return nil, fmt.Errorf("Invalid value for ConnectionResponse.Status: %#x", status)
	}

	host, err := readHostInfo(r)
	if err != nil { return nil, err }

	return &ConnectionResponse{channel, status, *host}, nil
}

// Connection state request
type ConnectionStateRequest struct {
	Channel byte
	Status  byte
	Host    HostInfo
}

func (req ConnectionStateRequest) describe() (serviceIdent, int) {
	return connectionStateRequestService, 10
}

func (req ConnectionStateRequest) writeTo(w *bytes.Buffer) error {
	err := writeSequence(w, req.Channel, req.Status)
	if err != nil { return err }

	return req.Host.writeTo(w)
}

//
type ConnState uint8

//
var (
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
		return fmt.Sprintf("Unknown connection state %#x", state)
	}
}

//
func (state ConnState) Error() string {
	return state.String()
}

// Connection state response
type ConnectionStateResponse struct {
	Channel uint8
	Status  ConnState
}

func readConnectionStateResponse(r *bytes.Reader) (*ConnectionStateResponse, error) {
	res := &ConnectionStateResponse{}

	err := readSequence(r, &res.Channel, &res.Status)
	if err != nil {
		return nil, err
	}

	switch res.Status {
	case ConnStateNormal, ConnStateInactive, ConnStateDataError, ConnStateKNXError:
		return res, nil

	default:
		return nil, fmt.Errorf("Invalid value for ConnectionStateResponse.Status: %#x", res.Status)
	}
}

// Disconnect request
type DisconnectRequest struct {
	Channel byte
	Status  byte
	Host    HostInfo
}

func readDisconnectRequest(r *bytes.Reader) (*DisconnectRequest, error) {
	var channel, status byte

	err := readSequence(r, &channel, &status)
	if err != nil { return nil, err }

	host, err := readHostInfo(r)
	if err != nil { return nil, err }

	return &DisconnectRequest{channel, status, *host}, nil
}

func (req DisconnectRequest) describe() (serviceIdent, int) {
	return disconnectRequestService, 10
}

func (req DisconnectRequest) writeTo(w *bytes.Buffer) error {
	err := writeSequence(w, req.Channel, req.Status)
	if err != nil { return err }

	return req.Host.writeTo(w)
}
