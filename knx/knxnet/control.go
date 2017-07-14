// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package knxnet

import (
	"errors"

	"github.com/vapourismo/knx-go/knx/util"
)

// TunnelLayer identifies the tunnelling layer for a tunnelling connection.
type TunnelLayer uint8

const (
	// TunnelLayerData establishes a data-link layer tunnel. Send and receive L_Data.* messages.
	TunnelLayerData TunnelLayer = 0x02

	// TunnelLayerRaw establishes a raw tunnel. Send and receive L_Raw.* messages.
	TunnelLayerRaw TunnelLayer = 0x04

	// TunnelLayerBusmon establishes a bus monitor tunnel. Receive L_Busmon.ind messages.
	TunnelLayerBusmon TunnelLayer = 0x80
)

// A ConnReq requests a connection to a gateway.
type ConnReq struct {
	Control HostInfo
	Tunnel  HostInfo
	Layer   TunnelLayer
}

// Service returns the service identifier for connection requests.
func (ConnReq) Service() ServiceID {
	return ConnReqService
}

var hostInfoSize = HostInfo{}.Size()

// Size returns the packed size.
func (ConnReq) Size() uint {
	return 2*hostInfoSize + 4
}

// Pack assembles the service payload in the given buffer.
func (req *ConnReq) Pack(buffer []byte) {
	util.PackSome(buffer, &req.Control, &req.Tunnel)

	buffer = buffer[2*hostInfoSize:]
	buffer[0] = 4
	buffer[1] = 4
	buffer[2] = byte(req.Layer)
	buffer[3] = 0
}

// Unpack parses the given service payload in order to initialize the structure.
func (req *ConnReq) Unpack(data []byte) (n uint, err error) {
	var length, connType, reserved uint8

	n, err = util.UnpackSome(
		data, &req.Control, &req.Tunnel, &length, &connType, (*uint8)(&req.Layer), &reserved,
	)
	if err != nil {
		return
	}

	if length != 4 {
		return n, errors.New("Invalid connection request info structure length")
	}

	if connType != 4 {
		return n, errors.New("Invalid connection type")
	}

	return
}

// ConnRes is a response to a connection request.
type ConnRes struct {
	Channel uint8
	Status  ErrCode
	Control HostInfo
}

// Service returns the service identifier for connection responses.
func (ConnRes) Service() ServiceID {
	return ConnResService
}

// Size returns the packed size.
func (res *ConnRes) Size() uint {
	if res.Status == 0 {
		return hostInfoSize + 6
	}

	return 2
}

// Pack assembles the service payload in the given buffer.
func (res *ConnRes) Pack(buffer []byte) {
	if res.Status == 0 {
		util.PackSome(buffer, res.Channel, uint8(0), &res.Control, []byte{4, 4, 0, 0})
	} else {
		util.PackSome(buffer, res.Channel, uint8(res.Status))
	}
}

// Unpack parses the given service payload in order to initialize the structure.
func (res *ConnRes) Unpack(data []byte) (n uint, err error) {
	n, err = util.UnpackSome(data, &res.Channel, (*uint8)(&res.Status))

	if res.Status == 0 {
		var m uint
		m, err = res.Control.Unpack(data[2:])
		n += m
	}

	return
}

// A ConnStateReq requests the connection state from a gateway.
type ConnStateReq struct {
	Channel uint8
	Status  ErrCode
	Control HostInfo
}

// Service returns the service identifier for connection state requests.
func (ConnStateReq) Service() ServiceID {
	return ConnStateReqService
}

// Size returns the packed size.
func (ConnStateReq) Size() uint {
	return 2 + hostInfoSize
}

// Pack assembles the service payload in the given buffer.
func (req *ConnStateReq) Pack(buffer []byte) {
	buffer[0] = req.Channel
	buffer[1] = uint8(req.Status)
	req.Control.Pack(buffer[2:])
}

// Unpack parses the given service payload in order to initialize the structure.
func (req *ConnStateReq) Unpack(data []byte) (uint, error) {
	return util.UnpackSome(data, &req.Channel, (*uint8)(&req.Status), &req.Control)
}

// A ConnStateRes is a response to a connection state request.
type ConnStateRes struct {
	Channel uint8
	Status  ErrCode
}

// Service returns the service identifier for connection state responses.
func (ConnStateRes) Service() ServiceID {
	return ConnStateResService
}

// Size returns the packed size.
func (ConnStateRes) Size() uint {
	return 2
}

// Pack assembles the service payload in the given buffer.
func (res *ConnStateRes) Pack(buffer []byte) {
	buffer[0] = res.Channel
	buffer[1] = uint8(res.Status)
}

// Unpack parses the given service payload in order to initialize the structure.
func (res *ConnStateRes) Unpack(data []byte) (uint, error) {
	return util.UnpackSome(data, &res.Channel, (*uint8)(&res.Status))
}

// A DiscReq requests a connection to be terminated.
type DiscReq struct {
	Channel uint8
	Status  uint8
	Control HostInfo
}

// Service returns the service identifier for disconnect requests.
func (DiscReq) Service() ServiceID {
	return DiscReqService
}

// Size returns the packed size.
func (DiscReq) Size() uint {
	return 2 + hostInfoSize
}

// Pack assembles the service payload in the given buffer.
func (req *DiscReq) Pack(buffer []byte) {
	buffer[0] = req.Channel
	buffer[1] = req.Status
	req.Control.Pack(buffer[2:])
}

// Unpack parses the given service payload in order to initialize the structure.
func (req *DiscReq) Unpack(data []byte) (uint, error) {
	return util.UnpackSome(data, &req.Channel, &req.Status, &req.Control)
}

// A DiscRes is a response to a disconnect request.
type DiscRes struct {
	Channel uint8
	Status  uint8
}

// Service returns the service identifier for disconnect responses.
func (DiscRes) Service() ServiceID {
	return DiscResService
}

// Size returns the packed size.
func (DiscRes) Size() uint {
	return 2
}

// Pack assembles the service payload in the given buffer.
func (res *DiscRes) Pack(data []byte) {
	data[0] = res.Channel
	data[1] = res.Status
}

// Unpack parses the given service payload in order to initialize the structure.
func (res *DiscRes) Unpack(data []byte) (uint, error) {
	return util.UnpackSome(data, &res.Channel, &res.Status)
}
