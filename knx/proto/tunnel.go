package proto

import (
	"errors"
	"io"

	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/encoding"
	"github.com/vapourismo/knx-go/knx/util"
)

// A TunnelReq asks a gateway to transmit data.
type TunnelReq struct {
	// Communication channel
	Channel uint8

	// Sequential number, used to track acknowledgements
	SeqNumber uint8

	// Data to be tunneled
	Payload cemi.Message
}

// Service returns the service identifiers for tunnel requests.
func (TunnelReq) Service() ServiceID {
	return TunnelReqService
}

// Unpack initializes the structure by parsing the given data.
func (req *TunnelReq) Unpack(data []byte) (n uint, err error) {
	var length, reserved uint8

	if n, err = util.UnpackSome(
		data, &length, &req.Channel, &req.SeqNumber, &reserved,
	); err != nil {
		return
	}

	if length != 4 {
		return n, errors.New("Length header is not 4")
	}

	m, err := cemi.Unpack(data[n:], &req.Payload)
	n += m

	return
}

// WriteTo serializes the structure and writes it to the given Writer.
func (req *TunnelReq) WriteTo(w io.Writer) (n int64, err error) {
	if n, err = encoding.WriteSome(w, byte(4), req.Channel, req.SeqNumber, byte(0)); err != nil {
		return
	}

	m, err := cemi.Pack(w, req.Payload)
	n += m

	return
}

// A TunnelResStatus is the status in a tunnel response.
type TunnelResStatus uint8

const (
	// TunnelResUnsupported indicates that the CEMI-encoded frame inside the tunnel request was not
	// understood or is not supported.
	TunnelResUnsupported TunnelResStatus = 0x29
)

// A TunnelRes is a response to a TunnelRequest. It acts as an acknowledgement.
type TunnelRes struct {
	// Communication channel
	Channel uint8

	// Identifies the request that is being acknowledged
	SeqNumber uint8

	// Status code, determines whether the tunneling succeeded or not
	Status TunnelResStatus
}

// Service returns the service identifier for tunnel responses.
func (TunnelRes) Service() ServiceID {
	return TunnelResService
}

// Unpack initializes the structure by parsing the given data.
func (res *TunnelRes) Unpack(data []byte) (n uint, err error) {
	var length uint8

	n, err = util.UnpackSome(data, &length, &res.Channel, &res.SeqNumber, (*uint8)(&res.Status))
	if err != nil {
		return
	}

	if length != 4 {
		return n, errors.New("Length header is not 4")
	}

	return
}

// WriteTo serializes the structure and writes it to the given Writer.
func (res *TunnelRes) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteSome(w, byte(4), res.Channel, res.SeqNumber, res.Status)
}
