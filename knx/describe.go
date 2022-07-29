// Copyright (c) 2022 mobilarte.
// Licensed under the MIT license which can be found in the LICENSE file.

package knx

import (
	"time"

	"github.com/vapourismo/knx-go/knx/knxnet"
)

// Describe a single KNXnet/IP server.
func Describe(Address string, searchTimeout time.Duration) (*knxnet.DescriptionRes, error) {
	// Uses a UDP socket.
	socket, err := knxnet.DialTunnel(Address)
	if err != nil {
		return nil, err
	}
	defer socket.Close()

	addr, _ := socket.LocalAddr()

	req, err := knxnet.NewDescriptionReq(addr)
	if err != nil {
		return nil, err
	}

	if err := socket.Send(req); err != nil {
		return nil, err
	}

	timeout := time.After(searchTimeout)

loop:
	for {
		select {
		case msg := <-socket.Inbound():
			descriptionRes, ok := msg.(*knxnet.DescriptionRes)
			if ok {
				return descriptionRes, nil
			}

		case <-timeout:
			break loop
		}
	}

	return nil, nil
}
