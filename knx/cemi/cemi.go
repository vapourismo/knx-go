// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

// Package cemi provides the functionality to parse and generate KNX CEMI-encoded frames.
package cemi

import (
	"fmt"

	"github.com/vapourismo/knx-go/knx/util"
)

// MessageCode is used to identify the type of message inside a CEMI-encoded frame.
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

// String converts the message code to a string.
func (code MessageCode) String() string {
	switch code {
	case LBusmonIndCode:
		return "LBusmon.ind"

	case LDataReqCode:
		return "LData.req"

	case LDataIndCode:
		return "LData.ind"

	case LDataConCode:
		return "LData.con"

	case LRawReqCode:
		return "LRaw.req"

	case LRawIndCode:
		return "LRaw.ind"

	case LRawConCode:
		return "LRaw.con"

	default:
		return fmt.Sprintf("%#x", uint8(code))
	}
}

// Info is the additional info segment of a CEMI-encoded frame.
type Info []byte

// Size returns the packed size.
func (info Info) Size() uint {
	if len(info) > 255 {
		return 256
	}

	return 1 + uint(len(info))
}

// Pack the info structure into the buffer.
func (info Info) Pack(buffer []byte) {
	if len(info) > 255 {
		buffer[0] = 255
	} else {
		buffer[0] = byte(len(info))
	}

	copy(buffer[1:], info[:buffer[0]])
}

// Unpack initializes the structure by parsing the given data.
func (info *Info) Unpack(data []byte) (n uint, err error) {
	var length uint8

	n, err = util.Unpack(data, &length)
	if err != nil {
		return
	}

	if length > 0 {
		buf := make([]byte, length)
		n += uint(copy(buf, data[n:n+uint(length)]))
		*info = Info(buf)
	} else {
		*info = nil
	}

	return
}

// Message is the body of a CEMI-encoded frame.
type Message interface {
	util.Packable
	MessageCode() MessageCode
}

// An UnsupportedMessage is the raw representation of a message inside a CEMI-encoded frame.
type UnsupportedMessage struct {
	Code MessageCode
	Data []byte
}

// Size returns the packed size.
func (body *UnsupportedMessage) Size() uint {
	return uint(len(body.Data))
}

// MessageCode returns the message code.
func (body *UnsupportedMessage) MessageCode() MessageCode {
	return body.Code
}

// Pack the message body into the buffer.
func (body *UnsupportedMessage) Pack(buffer []byte) {
	copy(buffer, body.Data)
}

// Unpack initializes the structure by parsing the given data.
func (body *UnsupportedMessage) Unpack(data []byte) (uint, error) {
	if len(body.Data) < len(data) {
		body.Data = make([]byte, len(data))
	}

	return uint(copy(body.Data, data)), nil
}

type messageUnpackable interface {
	util.Unpackable
	Message
}

// Unpack a message from a CEMI-encoded frame.
func Unpack(data []byte, message *Message) (n uint, err error) {
	var code MessageCode

	// Read header.
	n, err = util.Unpack(data, (*uint8)(&code))
	if err != nil {
		return
	}

	var body messageUnpackable

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
	m, err := body.Unpack(data[n:])

	if err == nil {
		*message = body
	}

	return n + m, err
}

// Size returns the size for a CEMI-encoded frame with the given message.
func Size(message Message) uint {
	return 1 + message.Size()
}

// Pack assembles a CEMI-encoded frame using the given message.
func Pack(buffer []byte, message Message) {
	util.PackSome(buffer, uint8(message.MessageCode()), message)
}
