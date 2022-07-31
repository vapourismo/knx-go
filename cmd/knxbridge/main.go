// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/vapourismo/knx-go/knx"
	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/knxnet"
	"github.com/vapourismo/knx-go/knx/util"
)

type relay interface {
	relay(data cemi.LData) error

	Inbound() <-chan cemi.Message
	Close()
}

type reqRelay struct {
	*knx.Tunnel
}

func (relay reqRelay) relay(data cemi.LData) error {
	return relay.Send(&cemi.LDataReq{LData: data})
}

type indRelay struct {
	*knx.Router
}

func (relay indRelay) relay(data cemi.LData) error {
	return relay.Send(&cemi.LDataInd{LData: data})
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <gateway addr> <other addr>\n", os.Args[0])
}

func main() {
	if len(os.Args) < 3 {
		printUsage()
		return
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)
	util.Logger = logger

	gatewayAddr := os.Args[1]
	otherAddr := os.Args[2]

	// Loop for ever. Failures don't matter, we'll always retry.
	for {
		br, err := newBridge(gatewayAddr, otherAddr)
		if err != nil {
			logger.Printf("Error while creating: %v\n", err)

			time.Sleep(time.Second)
			continue
		}

		err = br.serve()
		if err != nil {
			logger.Printf("Server terminated with error: %v\n", err)
		}

		br.close()

		time.Sleep(time.Second)
	}
}

type bridge struct {
	tunnel *knx.Tunnel
	other  relay
}

func newBridge(gatewayAddr, otherAddr string) (*bridge, error) {
	// Instantiate tunnel connection.
	tunnel, err := knx.NewTunnel(gatewayAddr, knxnet.TunnelLayerData, knx.DefaultTunnelConfig)
	if err != nil {
		return nil, err
	}

	var other relay

	addr, err := net.ResolveUDPAddr("udp4", otherAddr)
	if err != nil {
		tunnel.Close()
		return nil, err
	}

	if addr.IP.IsMulticast() {
		// Instantiate routing facilities. MulticastLoopback is disabled by default.
		router, err := knx.NewRouter(otherAddr, knx.DefaultRouterConfig)
		if err != nil {
			tunnel.Close()
			return nil, err
		}

		other = indRelay{router}
	} else {
		// Instantiate tunnel connection.
		otherTunnel, err := knx.NewTunnel(otherAddr, knxnet.TunnelLayerData, knx.DefaultTunnelConfig)
		if err != nil {
			tunnel.Close()
			return nil, err
		}

		other = reqRelay{otherTunnel}
	}

	return &bridge{tunnel, other}, nil
}

func (br *bridge) serve() error {
	for {
		select {
		// Receive message from gateway.
		case msg, open := <-br.tunnel.Inbound():
			if !open {
				return errors.New("tunnel channel closed")
			}

			if ind, ok := msg.(*cemi.LDataInd); ok {
				util.Log(br, "%+v", ind)
				if err := br.other.relay(ind.LData); err != nil {
					return err
				}
			}

		// Receive message from router.
		case msg, open := <-br.other.Inbound():
			if !open {
				return errors.New("router channel closed")
			}

			if ind, ok := msg.(*cemi.LDataInd); ok {
				util.Log(br, "%+v", ind)
				if err := br.tunnel.Send(&cemi.LDataReq{LData: ind.LData}); err != nil {
					return err
				}
			}
		}
	}
}

func (br *bridge) close() {
	br.tunnel.Close()
	br.other.Close()
}
