// Licensed under the MIT license which can be found in the LICENSE file.

package knxnet

import (
	"net"

	"github.com/vapourismo/knx-go/knx/util"
)

// NewSearchReq creates a new SearchReq, addr defines where KNXnet/IP server should send the reponse to.
func NewSearchReq(addr net.Addr) (*SearchReq, error) {
	req := &SearchReq{}

	hostinfo, err := HostInfoFromAddress(addr)
	if err != nil {
		return nil, err
	}
	req.HostInfo = hostinfo

	return req, nil
}

// A SearchReq requests a discovery from all KNXnet/IP servers via multicast.
type SearchReq struct {
	HostInfo
}

// Service returns the service identifier for Search Request.
func (SearchReq) Service() ServiceID {
	return SearchReqService
}

// A SearchRes is a Search Response from a KNXnet/IP server.
type SearchRes struct {
	Control      HostInfo
	DescriptionB DescriptionBlock
}

// Service returns the service identifier for the Search Response.
func (SearchRes) Service() ServiceID {
	return SearchResService
}

// Size returns the packed size.
func (res SearchRes) Size() uint {
	return res.Control.Size() + res.DescriptionB.DeviceHardware.Size() + res.DescriptionB.SupportedServices.Size()
}

// Pack assembles the Search Response structure in the given buffer.
func (res *SearchRes) Pack(buffer []byte) {
	util.PackSome(buffer, res.Control, res.DescriptionB.DeviceHardware, res.DescriptionB.SupportedServices)
}

// Unpack parses the given service payload in order to initialize the Search Response structure.
func (res *SearchRes) Unpack(data []byte) (n uint, err error) {
	return util.UnpackSome(data, &res.Control, &res.DescriptionB.DeviceHardware, &res.DescriptionB.SupportedServices)
}
