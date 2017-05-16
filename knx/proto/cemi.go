package proto

import (
	"errors"
	"io"
	"github.com/vapourismo/knx-go/knx/binary"
)

// MessageCode identifies the message body of a CEMI frame.
type MessageCode uint8

// Supported message codes
const (
	LDataReq MessageCode = 0x11
	LDataInd MessageCode = 0x29
	LDataCon MessageCode = 0x2e
)

// CEMI is a common external message interface.
type CEMI struct {
	Code MessageCode
	Info []byte
	Body Segment
}

// Errors from ReadCEMI
var (
	ErrUnsupportedMessageCode = errors.New("CEMI message code is unsupported")
)

// ReadCEMI extract a CEMI frame from the given data.
func ReadCEMI(r io.Reader) (*CEMI, error) {
	var code MessageCode
	var infoLen uint8

	err := binary.ReadSequence(r, &code, &infoLen)
	if err != nil {
		return nil, err
	}

	info := make([]byte, int(infoLen))

	len, err := r.Read(info)
	if err != nil {
		return nil, err
	} else if len < int(infoLen) {
		return nil, io.ErrUnexpectedEOF
	}

	var body Segment

	switch code {
	case LDataReq, LDataInd, LDataCon:
		body, err = ReadLData(r)
		if err != nil {
			return nil, err
		}

	default:
		return nil, ErrUnsupportedMessageCode
	}

	return &CEMI{code, info, body}, nil
}

// WriteTo writes the CEMI frame to the Writer.
func (cemi *CEMI) WriteTo(w io.Writer) error {
	var infoLen uint8
	var info []byte

	if len(cemi.Info) > 255 {
		infoLen = 255
		info = cemi.Info[:256]
	} else {
		info = cemi.Info
		infoLen = uint8(len(info))
	}

	err := binary.WriteSequence(w, cemi.Code, infoLen, info)
	if err != nil {
		return err
	}

	return cemi.Body.WriteTo(w)
}