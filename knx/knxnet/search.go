// Licensed under the MIT license which can be found in the LICENSE file.

package knxnet

import (
	"fmt"
	"net"
	"strconv"

	"github.com/vapourismo/knx-go/knx/util"
)

// NewSearchReq creates a new SearchReq, addr defines where ObjectServers should send the reponse to
func NewSearchReq(addr net.Addr) (*SearchReq, error) {
	ipS, portS, err := net.SplitHostPort(addr.String())
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(ipS)
	if ip == nil {
		return nil, fmt.Errorf("Unable to determine IP")
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return nil, fmt.Errorf("Only IPv4 is currently supported")
	}

	port, _ := strconv.ParseUint(portS, 10, 16)
	if port == 0 {
		return nil, fmt.Errorf("Unable to determine port")
	}

	req := &SearchReq{}
	switch addr.Network() {
	case "udp":
		req.Protocol = UDP4
	case "tcp":
		req.Protocol = TCP4
	default:
		return nil, fmt.Errorf("Unsupported network")
	}

	copy(req.Address[:], ipv4)
	req.Port = Port(port)
	return req, nil
}

// A SearchReq requests a discovery from all KNXnet/IP Servers
type SearchReq struct {
	HostInfo
}

// Service returns the service identifier for search request
func (SearchReq) Service() ServiceID {
	return SearchReqService
}

// A SearchRes is a discovery response from a KNXnet/IP Server
type SearchRes struct {
	Control           HostInfo
	DeviceHardware    DeviceInformationBlock
	SupportedServices SupportedServicesDIB
}

// Service returns the service identifier for search response
func (SearchRes) Service() ServiceID {
	return SearchResService
}

// Size returns the packed size.
func (res SearchRes) Size() uint {
	return res.Control.Size() + res.DeviceHardware.Size() + res.SupportedServices.Size()
}

// Pack assembles the search response structure in the given buffer.
func (res *SearchRes) Pack(buffer []byte) {
	util.PackSome(buffer, res.Control, res.DeviceHardware, res.SupportedServices)
}

// Unpack parses the given service payload in order to initialize the structure.
func (res *SearchRes) Unpack(data []byte) (n uint, err error) {
	return util.UnpackSome(data, &res.Control, &res.DeviceHardware, &res.SupportedServices)
}
