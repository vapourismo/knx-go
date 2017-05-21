package proto

import (
	"errors"
	"io"
	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/encoding"
)

// A TunnelReq asks a gateway to transmit data.
type TunnelReq struct {
	// Communication channel
	Channel uint8

	// Sequential number, used to track acknowledgements
	SeqNumber uint8

	// Data to be tunneled
	Payload cemi.CEMI
}

// ReadFrom initializes the structure by reading from the given Reader.
func (req *TunnelReq) ReadFrom(r io.Reader) (n int64, err error) {
	var length, reserved uint8

	n, err = encoding.ReadSome(r, &length, &req.Channel, &req.SeqNumber, &reserved, &req.Payload)
	if err != nil {
		return
	}

	if length != 4 {
		return n, errors.New("Length header is not 4")
	}

	return
}

// WriteTo serializes the structure and writes it to the given Writer.
func (req *TunnelReq) WriteTo(w io.Writer) (int64, error) {
	return encoding.WriteSome(w, byte(4), req.Channel, req.SeqNumber, byte(0), req.Payload)
}

// A TunnelRes is a response to a TunnelRequest. It acts as an acknowledgement.
type TunnelRes struct {
	// Communication channel
	Channel uint8

	// Identifies the request that is being acknowledged
	SeqNumber uint8

	// Status code, determines whether the tunneling succeeded or not
	Status uint8
}

// ReadFrom initializes the structure by reading from the given Reader.
func (res *TunnelRes) ReadFrom(r io.Reader) (n int64, err error) {
	var length uint8

	n, err = encoding.ReadSome(r, &length, &res.Channel, &res.SeqNumber, &res.Status)
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
