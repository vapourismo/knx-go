// Licensed under the MIT license which can be found in the LICENSE file.

package knx

import (
	"net"
	"time"

	"github.com/vapourismo/knx-go/knx/knxnet"
)

// Discover all KNXnet/IP Server
func Discover(multicastDiscoveryAddress string, searchTimeout time.Duration) ([]*knxnet.SearchRes, error) {
	return DiscoverOnInterface(nil, multicastDiscoveryAddress, searchTimeout)
}

// DiscoverOnInterface discovers all KNXnet/IP Server on a specific interface. If the
// interface is nil, the system-assigned multicast interface is used.
func DiscoverOnInterface(ifi *net.Interface, multicastDiscoveryAddress string, searchTimeout time.Duration) ([]*knxnet.SearchRes, error) {
	socket, err := knxnet.ListenRouterOnInterface(ifi, multicastDiscoveryAddress)
	if err != nil {
		return nil, err
	}
	defer socket.Close()

	req, err := knxnet.NewSearchReq(socket.Addr())
	if err != nil {
		return nil, err
	}

	if err := socket.Send(req); err != nil {
		return nil, err
	}

	results := []*knxnet.SearchRes{}
	timeout := time.After(searchTimeout)

loop:
	for {
		select {
		case msg := <-socket.Inbound():
			searchRes, ok := msg.(*knxnet.SearchRes)
			if !ok {
				continue
			}
			results = append(results, searchRes)

		case <-timeout:
			break loop
		}
	}

	return results, nil
}
