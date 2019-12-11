// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

// Package knxnet provides the means to parse and generate frames of the KNXnet/IP protocol.
package knxnet

import (
	"errors"

	"fmt"

	"github.com/vapourismo/knx-go/knx/util"
)

// ServiceID identifies the service that is contained in a packet.
type ServiceID uint16

// String generates a string representation.
func (srv ServiceID) String() string {
	return fmt.Sprintf("%#04x", uint16(srv))
}

// These are supported services.
const (
	SearchReqService    ServiceID = 0x0201
	SearchResService    ServiceID = 0x0202
	ConnReqService      ServiceID = 0x0205
	ConnResService      ServiceID = 0x0206
	ConnStateReqService ServiceID = 0x0207
	ConnStateResService ServiceID = 0x0208
	DiscReqService      ServiceID = 0x0209
	DiscResService      ServiceID = 0x020a
	TunnelReqService    ServiceID = 0x0420
	TunnelResService    ServiceID = 0x0421
	RoutingIndService   ServiceID = 0x0530
	RoutingLostService  ServiceID = 0x0531
	RoutingBusyService  ServiceID = 0x0532
)

// Service describes a KNXnet/IP service.
type Service interface {
	Service() ServiceID
}

// ServicePackable combines Packable and Service.
type ServicePackable interface {
	util.Packable
	Service
}

// UnknownService is the payload of an unknown service.
type UnknownService struct {
	service ServiceID
	Data    []byte
}

// Service returns the service identifier.
func (us *UnknownService) Service() ServiceID {
	return us.service
}

// Size returns the size of the payload.
func (us *UnknownService) Size() uint {
	return uint(len(us.Data))
}

// Pack the payload into the buffer.
func (us *UnknownService) Pack(buffer []byte) {
	copy(buffer, us.Data)
}

// Unpack copies the entire data.
func (us *UnknownService) Unpack(data []byte) (uint, error) {
	us.Data = make([]byte, len(data))
	return uint(copy(us.Data, data)), nil
}

// Size returns the size of a KNXnet/IP packet.
func Size(service ServicePackable) uint {
	return 6 + service.Size()
}

// Pack generates a KNXnet/IP packet. Utilize Size() to determine the required size of the buffer.
func Pack(buffer []byte, srv ServicePackable) {
	buffer[0] = 6
	buffer[1] = 16
	util.Pack(buffer[2:], uint16(srv.Service()))
	util.Pack(buffer[4:], uint16(srv.Size()+6))
	srv.Pack(buffer[6:])
}

// AllocAndPack allocates a buffer and packs the KNXnet/IP packet into it.
func AllocAndPack(srv ServicePackable) []byte {
	buffer := make([]byte, Size(srv))
	Pack(buffer, srv)
	return buffer
}

// These are errors that might occur during unpacking.
var (
	ErrHeaderLength  = errors.New("Header length is not 6")
	ErrHeaderVersion = errors.New("Protocol version is not 16")
)

type serviceUnpackable interface {
	util.Unpackable
	Service
}

// Unpack parses a KNXnet/IP packet and retrieves its service payload.
//
// On success, the variable pointed to by srv will contain a pointer to a service type.
// You can cast it to the matching against service type, like so:
//
// 	var srv Service
//
// 	_, err := Unpack(r, &srv)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	switch srv := srv.(type) {
// 		case *ConnRes:
// 			// ...
//
// 		case *TunnelReq:
// 			// ...
//
// 		// ...
// 	}
//
func Unpack(data []byte, srv *Service) (uint, error) {
	var headerLen, version uint8
	var srvID ServiceID
	var totalLen uint16

	n, err := util.UnpackSome(data, &headerLen, &version, (*uint16)(&srvID), &totalLen)
	if err != nil {
		return n, err
	}

	if headerLen != 6 {
		return n, ErrHeaderLength
	}

	if version != 16 {
		return n, ErrHeaderVersion
	}

	var body serviceUnpackable
	switch srvID {
	case SearchReqService:
		body = &SearchReq{}

	case SearchResService:
		body = &SearchRes{}

	case ConnReqService:
		body = &ConnReq{}

	case ConnResService:
		body = &ConnRes{}

	case ConnStateReqService:
		body = &ConnStateReq{}

	case ConnStateResService:
		body = &ConnStateRes{}

	case DiscReqService:
		body = &DiscReq{}

	case DiscResService:
		body = &DiscRes{}

	case TunnelReqService:
		body = &TunnelReq{}

	case TunnelResService:
		body = &TunnelRes{}

	case RoutingIndService:
		body = &RoutingInd{}

	case RoutingLostService:
		body = &RoutingLost{}

	case RoutingBusyService:
		body = &RoutingBusy{}

	default:
		body = &UnknownService{service: srvID}
	}

	m, err := body.Unpack(data[n:])

	if err == nil {
		*srv = body
	}

	return n + m, err
}
