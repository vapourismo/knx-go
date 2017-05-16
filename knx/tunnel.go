package knx

import (
	"bytes"
	"errors"
	"io"
	"github.com/vapourismo/knx-go/knx/encoding"
)

// A TunnelRequest asks a gateway to transmit data.
type TunnelRequest struct {
	Channel   byte
	SeqNumber byte
	Payload   []byte
}

func readTunnelRequest(r *bytes.Reader) (*TunnelRequest, error) {
	var length, channel, seq byte

	err := encoding.ReadSequence(r, &length, &channel, &seq)
	if err != nil { return nil, err }

	if length != 4 {
		return nil, errors.New("Invalid structure length")
	}

	_, err = r.Seek(1, io.SeekCurrent)
	if err != nil { return nil, err }

	buffer := &bytes.Buffer{}

	_, err = r.WriteTo(buffer)
	if err != nil { return nil, err }

	return &TunnelRequest{channel, seq, buffer.Bytes()}, nil
}

func (req TunnelRequest) describe() (serviceIdent, int) {
	return tunnelRequestService, 4 + len(req.Payload)
}

func (req TunnelRequest) writeTo(w *bytes.Buffer) error {
	return encoding.WriteSequence(w, byte(4), req.Channel, req.SeqNumber, byte(0), req.Payload)
}

// A TunnelResponse is a response to a TunnelRequest.
type TunnelResponse struct {
	Channel   byte
	SeqNumber byte
	Status    byte
}

func readTunnelResponse(r *bytes.Reader) (*TunnelResponse, error) {
	var length, channel, seq, status byte

	err := encoding.ReadSequence(r, &length, &channel, &seq, &status)
	if err != nil { return nil, err }

	if length != 4 {
		return nil, errors.New("Invalid structure length")
	}

	return &TunnelResponse{channel, seq, status}, nil
}

func (res TunnelResponse) describe() (serviceIdent, int) {
	return tunnelResponseService, 4
}

func (res TunnelResponse) writeTo(w *bytes.Buffer) error {
	return encoding.WriteSequence(w, byte(4), res.Channel, res.SeqNumber, res.Status)
}
