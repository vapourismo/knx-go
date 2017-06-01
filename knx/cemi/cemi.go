package cemi

import (
	"io"

	"github.com/vapourismo/knx-go/knx/encoding"
)

// MessageCode is used to identify the contents of a Message.
type MessageCode uint8

const (
	// LBusmonIndCode is the message code for L_Busmon.ind.
	LBusmonIndCode MessageCode = 0x2B

	// LDataReqCode is the message code for L_Data.Req.
	LDataReqCode MessageCode = 0x11

	// LDataIndCode is the message code for L_Data.Ind.
	LDataIndCode MessageCode = 0x29

	// LDataConCode is the message code for L_Data.Con.
	LDataConCode MessageCode = 0x2E

	// LRawReqCode      MessageCode = 0x10
	// LRawIndCode      MessageCode = 0x2D
	// LRawConCode      MessageCode = 0x2F
	// LPollDataReqCode MessageCode = 0x13
	// LPollDataConCode MessageCode = 0x25
)

// MessageBody is the body of a Message.
type MessageBody interface {
	io.WriterTo
	MessageCode() MessageCode
}

// An UnsupportedMessageBody is the raw representation of a CEMI message body.
type UnsupportedMessageBody struct {
	Code MessageCode
	Data []byte
}

// MessageCode returns the message code.
func (body *UnsupportedMessageBody) MessageCode() MessageCode {
	return body.Code
}

// ReadFrom initializes the structure by reading from the given Reader.
func (body *UnsupportedMessageBody) ReadFrom(r io.Reader) (n int64, err error) {
	n, body.Data = encoding.ReadAll(r)
	return
}

// WriteTo serializes the structure and writes it to the given Writer.
func (body *UnsupportedMessageBody) WriteTo(w io.Writer) (int64, error) {
	len, err := w.Write(body.Data)
	return int64(len), err
}

// Message represents the Common External Message Interface.
type Message struct {
	Info []byte
	Body MessageBody
}

type messageBodyReaderFrom interface {
	io.ReaderFrom
	MessageBody
}

// ReadFrom initializes the structure by reading from the given Reader.
func (cemi *Message) ReadFrom(r io.Reader) (n int64, err error) {
	var infoLen uint8
	var code MessageCode

	// Retrieve CEMI header.
	n, err = encoding.ReadSome(r, &code, &infoLen)
	if err != nil {
		return
	}

	cemi.Info = make([]byte, infoLen)

	// Extract the additional info segment.
	len, err := encoding.Read(r, cemi.Info)
	n += len

	if err != nil {
		return
	}

	var body messageBodyReaderFrom

	// Decide which message body is appropriate.
	switch code {
	case LBusmonIndCode:
		body = &LBusmonInd{}

	case LDataReqCode:
		body = &LDataReq{}

	case LDataConCode:
		body = &LDataCon{}

	case LDataIndCode:
		body = &LDataInd{}

	default:
		body = &UnsupportedMessageBody{Code: code}
	}

	// Parse the message body.
	len, err = body.ReadFrom(r)

	if err == nil {
		cemi.Body = body
	}

	return n + len, err
}

// WriteTo serializes the CEMI frame and writes it to the given Writer.
func (cemi *Message) WriteTo(w io.Writer) (int64, error) {
	var infoLen uint8
	var info []byte

	if len(cemi.Info) > 255 {
		infoLen = 255
		info = cemi.Info[:256]
	} else {
		infoLen = uint8(len(info))
		info = cemi.Info
	}

	return encoding.WriteSome(w, cemi.Body.MessageCode(), infoLen, info, cemi.Body)
}
