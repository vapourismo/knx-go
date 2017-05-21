package cemi

import (
	"io"
	"github.com/vapourismo/knx-go/knx/encoding"
	"bytes"
)

// MessageCode is used to identify the contents of a CEMI frame.
type MessageCode uint8

// Supported message codes
const (
	LDataReq MessageCode = 0x11
	LDataInd MessageCode = 0x29
	LDataCon MessageCode = 0x2e
)

// Segment is a protocol segment.
type Segment interface {
	io.WriterTo
}

// An UnsupportedMessage is the raw representation of a CEMI message body.
type UnsupportedMessage []byte

// WriteTo serializes the message and writes it to the given Writer.
func (data UnsupportedMessage) WriteTo(w io.Writer) (int64, error) {
	len, err := w.Write(data)
	return int64(len), err
}

// CEMI is the Common External Message Interface.
type CEMI struct {
	Code MessageCode
	Info []byte
	Body Segment
}

// ReadFrom initializes the structure by reading from the given Reader.
func (cemi *CEMI) ReadFrom(r io.Reader) (n int64, err error) {
	var infoLen uint8
	n, err = encoding.ReadSome(r, &cemi.Code, &infoLen)
	if err != nil {
		return
	}

	cemi.Info = make([]byte, int(infoLen))
	len, err := encoding.Read(r, cemi.Info)
	n += len

	if err != nil {
		return n, err
	}

	switch cemi.Code {
	case LDataReq, LDataInd, LDataCon:
		ldata := &LData{}
		len, err = ldata.ReadFrom(r)
		n += len

		if err != nil {
			return n, err
		}

		cemi.Body = ldata

		return

	default:
		buffer := bytes.Buffer{}
		len, err = buffer.ReadFrom(r)
		n += len

		if err != nil {
			return n, err
		}

		cemi.Body = UnsupportedMessage(buffer.Bytes())

		return
	}
}

// WriteTo serializes the CEMI frame and writes it to the given Writer.
func (cemi *CEMI) WriteTo(w io.Writer) (int64, error) {
	var infoLen uint8
	var info []byte

	if len(cemi.Info) > 255 {
		infoLen = 255
		info = cemi.Info[:256]
	} else {
		infoLen = uint8(len(info))
		info = cemi.Info
	}

	return encoding.WriteSome(w, cemi.Code, infoLen, info, cemi.Body)
}