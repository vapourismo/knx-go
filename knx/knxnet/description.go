// Copyright (c) 2022 mobilarte.
// Licensed under the MIT license which can be found in the LICENSE file.

package knxnet

import (
	"net"

	"github.com/vapourismo/knx-go/knx/util"
)

// NewDescriptionReq creates a new Description Request, addr defines where
// KNXnet/IP server should send the response to.
func NewDescriptionReq(addr net.Addr) (*DescriptionReq, error) {
	req := &DescriptionReq{}

	hostinfo, err := HostInfoFromAddress(addr)
	if err != nil {
		return nil, err
	}
	req.HostInfo = hostinfo

	return req, nil
}

// A DescriptionReq requests a description from a particular KNXnet/IP server via unicast.
type DescriptionReq struct {
	HostInfo
}

// Service returns the service identifier for a Description Request.
func (DescriptionReq) Service() ServiceID {
	return DescrReqService
}

// A DescriptionRes is a Description Response from a KNXnet/IP server.
type DescriptionRes DescriptionBlock

// Service returns the service identifier for Description Response.
func (DescriptionRes) Service() ServiceID {
	return DescrResService
}

// Size returns the packed size of a Description Response.
func (res DescriptionRes) Size() uint {
	return res.DeviceHardware.Size() + res.SupportedServices.Size()
}

// Pack assembles the Description Response structure in the given buffer.
func (res *DescriptionRes) Pack(buffer []byte) {
	util.PackSome(buffer, res.DeviceHardware, res.SupportedServices)
}

// Unpack parses the given service payload in order to initialize the Description Response.
func (res *DescriptionRes) Unpack(data []byte) (n uint, err error) {
	return (*DescriptionBlock)(res).Unpack(data)
}
