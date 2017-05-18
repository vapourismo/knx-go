package proto

import (
	"io"
	"github.com/vapourismo/knx-go/knx/encoding"
)

// MessageCode identifies the message body of a CEMI frame.
type MessageCode uint8

// Supported message codes
const (
	LDataReq MessageCode = 0x11
	LDataInd MessageCode = 0x29
	LDataCon MessageCode = 0x2e
)

// Segment is a protocol segment.
type Segment interface {
	WriteTo(w io.Writer) error
}

// A UnsupportedMessage is the raw representation of a CEMI message body.
type UnsupportedMessage []byte

// WriteTo writes the contents to a Writer.
func (data UnsupportedMessage) WriteTo(w io.Writer) error {
	_, err := w.Write(data)
	return err
}

// CEMI is a common external message interface.
type CEMI struct {
	Code MessageCode
	Info []byte
	Body Segment
}

// ReadCEMI extract a CEMI frame from the given data.
func ReadCEMI(cemi []byte) (*CEMI, error) {
	if len(cemi) < 2 {
		return nil, ErrDataTooShort
	}

	code := MessageCode(cemi[0])
	infoLen := int(cemi[1])

	if infoLen > len(cemi) - 2 {
		return nil, ErrDataIncomplete
	}

	info := cemi[2:2 + infoLen]
	data := cemi[2 + infoLen:]

	switch code {
	case LDataReq, LDataInd, LDataCon:
		body, err := ReadLData(data)
		if err != nil {
			return nil, err
		}

		return &CEMI{code, info, body}, nil

	default:
		return &CEMI{code, info, UnsupportedMessage(data)}, nil
	}
}

// WriteTo writes the CEMI frame to the Writer.
func (cemi *CEMI) WriteTo(w io.Writer) error {
	var infoLen uint8
	var info []byte

	if len(cemi.Info) > 255 {
		infoLen = 255
		info = cemi.Info[:256]
	} else {
		infoLen = uint8(len(info))
		info = cemi.Info
	}

	err := encoding.WriteSequence(w, cemi.Code, infoLen, info)
	if err != nil {
		return err
	}

	return cemi.Body.WriteTo(w)
}