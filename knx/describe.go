// Copyright (c) 2022 mobilarte.
// Licensed under the MIT license which can be found in the LICENSE file.

package knx

import (
	"time"

	"github.com/vapourismo/knx-go/knx/knxnet"
)

// Describe a single KNXnet/IP server. Uses unicast UDP, address format is "ip:port".
func DescribeTunnel(address string, searchTimeout time.Duration) (*knxnet.DescriptionRes, error) {
	// Uses a UDP socket.
	socket, err := knxnet.DialTunnelUDP(address)
	if err != nil {
		return nil, err
	}
	defer socket.Close()

	addr := socket.LocalAddr()

	req, err := knxnet.NewDescriptionReq(addr)
	if err != nil {
		return nil, err
	}

	if err := socket.Send(req); err != nil {
		return nil, err
	}

	timeout := time.After(searchTimeout)

	for {
		select {
		case msg := <-socket.Inbound():
			descriptionRes, ok := msg.(*knxnet.DescriptionRes)
			if ok {
				return descriptionRes, nil
			}

		case <-timeout:
			return nil, nil
		}
	}
}
