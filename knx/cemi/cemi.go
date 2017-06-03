// Package cemi provides the functionality to parse and generate KNX CEMI-encoded frames.
package cemi

import (
	"io"

	"github.com/vapourismo/knx-go/knx/encoding"
)

// MessageCode is used to identify the contents of a CEMI frame.
type MessageCode uint8

const (
	// LBusmonIndCode is the message code for L_Busmon.ind.
	LBusmonIndCode MessageCode = 0x2B

	// LDataReqCode is the message code for L_Data.req.
	LDataReqCode MessageCode = 0x11

	// LDataIndCode is the message code for L_Data.ind.
	LDataIndCode MessageCode = 0x29

	// LDataConCode is the message code for L_Data.con.
	LDataConCode MessageCode = 0x2E

	// LRawReqCode is the message code for L_Raw.req.
	LRawReqCode MessageCode = 0x10

	// LRawIndCode is the message code for L_Raw.ind.
	LRawIndCode MessageCode = 0x2D

	// LRawConCode is the message code for L_Raw.con.
	LRawConCode MessageCode = 0x2F

	// LPollDataReqCode MessageCode = 0x13
	// LPollDataConCode MessageCode = 0x25
)

// Info is the additional info segment of a CEMI-encoded frame.
type Info []byte

// ReadFrom extracts an additional information segment.
func (info *Info) ReadFrom(r io.Reader) (n int64, err error) {
	var length uint8

	n, err = encoding.Read(r, &length)
	if err != nil {
		return
	}

	if length > 0 {
		buf := make([]byte, length)

		m, err := encoding.Read(r, buf)
		n += m
		if err != nil {
			return n, err
		}

		*info = Info(buf)
	} else {
		*info = nil
	}

	return
}

// WriteTo writes an additional information segment.
func (info Info) WriteTo(w io.Writer) (int64, error) {
	length := uint8(len(info))
	return encoding.WriteSome(w, length, []byte(info[:length]))
}

// Message is the body of a Message.
type Message interface {
	io.WriterTo
	MessageCode() MessageCode
}

// An UnsupportedMessage is the raw representation of a message inside a CEMI-encoded frame.
type UnsupportedMessage struct {
	Code MessageCode
	Data []byte
}

// MessageCode returns the message code.
func (body *UnsupportedMessage) MessageCode() MessageCode {
	return body.Code
}

// ReadFrom initializes the structure by reading from the given Reader.
func (body *UnsupportedMessage) ReadFrom(r io.Reader) (n int64, err error) {
	n, body.Data = encoding.ReadAll(r)
	return
}

// WriteTo serializes the structure and writes it to the given Writer.
func (body *UnsupportedMessage) WriteTo(w io.Writer) (int64, error) {
	len, err := w.Write(body.Data)
	return int64(len), err
}

type messageReaderFrom interface {
	io.ReaderFrom
	Message
}

// Unpack extracts the message from a CEMI-encoded frame.
func Unpack(r io.Reader, message *Message) (n int64, err error) {
	var code MessageCode

	// Read header.
	n, err = encoding.Read(r, &code)
	if err != nil {
		return
	}

	var body messageReaderFrom

	// Decide which message is appropriate.
	switch code {
	case LBusmonIndCode:
		body = &LBusmonInd{}

	case LDataReqCode:
		body = &LDataReq{}

	case LDataConCode:
		body = &LDataCon{}

	case LDataIndCode:
		body = &LDataInd{}

	case LRawReqCode:
		body = &LRawReq{}

	case LRawConCode:
		body = &LRawCon{}

	case LRawIndCode:
		body = &LRawInd{}

	default:
		body = &UnsupportedMessage{Code: code}
	}

	// Parse the message.
	m, err := body.ReadFrom(r)

	if err == nil {
		*message = body
	}

	return n + m, err
}

// Pack assembles a CEMI-encoded frame using the given message.
func Pack(w io.Writer, message Message) (int64, error) {
	return encoding.WriteSome(w, message.MessageCode(), message)
}
