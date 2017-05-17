package knx

import (
	"bytes"
	"errors"
	"github.com/vapourismo/knx-go/knx/encoding"
	"io"
)

// A TunnelReq asks a gateway to transmit data.
type TunnelReq struct {
	Channel   byte
	SeqNumber byte
	Payload   []byte
}

func readTunnelReq(r *bytes.Reader) (*TunnelReq, error) {
	var length, channel, seq byte

	err := encoding.ReadSequence(r, &length, &channel, &seq)
	if err != nil {
		return nil, err
	}

	if length != 4 {
		return nil, errors.New("Invalid structure length")
	}

	_, err = r.Seek(1, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	buffer := &bytes.Buffer{}

	_, err = r.WriteTo(buffer)
	if err != nil {
		return nil, err
	}

	return &TunnelReq{channel, seq, buffer.Bytes()}, nil
}

func (req TunnelReq) describe() (ServiceID, int) {
	return TunnelReqService, 4 + len(req.Payload)
}

func (req TunnelReq) writeTo(w *bytes.Buffer) error {
	return encoding.WriteSequence(w, byte(4), req.Channel, req.SeqNumber, byte(0), req.Payload)
}

// A TunnelRes is a response to a TunnelRequest.
type TunnelRes struct {
	Channel   byte
	SeqNumber byte
	Status    byte
}

func readTunnelRes(r *bytes.Reader) (*TunnelRes, error) {
	var length, channel, seq, status byte

	err := encoding.ReadSequence(r, &length, &channel, &seq, &status)
	if err != nil {
		return nil, err
	}

	if length != 4 {
		return nil, errors.New("Invalid structure length")
	}

	return &TunnelRes{channel, seq, status}, nil
}

func (res TunnelRes) describe() (ServiceID, int) {
	return TunnelResService, 4
}

func (res TunnelRes) writeTo(w *bytes.Buffer) error {
	return encoding.WriteSequence(w, byte(4), res.Channel, res.SeqNumber, res.Status)
}
