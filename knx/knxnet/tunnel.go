// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package knxnet

import (
	"errors"

	"github.com/vapourismo/knx-go/knx/cemi"
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

// Size returns the packed size.
func (req *TunnelReq) Size() uint {
	return 4 + cemi.Size(req.Payload)
}

// Pack assembles the service payload in the given buffer.
func (req *TunnelReq) Pack(buffer []byte) {
	buffer[0] = 4
	buffer[1] = req.Channel
	buffer[2] = req.SeqNumber
	buffer[3] = 0
	cemi.Pack(buffer[4:], req.Payload)
}

// Unpack parses the given service payload in order to initialize the structure.
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

// A TunnelRes is a response to a TunnelRequest. It acts as an acknowledgement.
type TunnelRes struct {
	// Communication channel
	Channel uint8

	// Identifies the request that is being acknowledged
	SeqNumber uint8

	// Status code, determines whether the tunneling succeeded or not
	Status ErrCode
}

// Service returns the service identifier for tunnel responses.
func (TunnelRes) Service() ServiceID {
	return TunnelResService
}

// Size returns the packed size.
func (TunnelRes) Size() uint {
	return 4
}

// Pack assembles the service payload in the given buffer.
func (res *TunnelRes) Pack(buffer []byte) {
	buffer[0] = 4
	buffer[1] = res.Channel
	buffer[2] = res.SeqNumber
	buffer[3] = uint8(res.Status)
}

// Unpack parses the given service payload in order to initialize the structure.
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
