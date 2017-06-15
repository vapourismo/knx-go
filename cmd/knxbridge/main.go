// Copyright 2017 Ole Krüger.

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"time"

	"github.com/vapourismo/knx-go/knx"
	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/knxnet"
	"github.com/vapourismo/knx-go/knx/util"
)

var (
	flagRouter = flag.String("router", "224.0.23.12:3671", "KNXnet/IP router multicast group")
	flagHelp   = flag.Bool("help", false, "Display usage")
)

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <gateway addr>\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	if *flagHelp || len(flag.Args()) < 1 {
		printUsage()
		return
	}

	util.Logger = log.New(os.Stdout, "", log.LstdFlags)

	// Loop for ever. Failures don't matter, we'll always retry.
	for {
		br, err := newBridge(flag.Arg(0), *flagRouter)
		if err != nil {
			util.Log(br, "Error while creating: %v", err)

			time.Sleep(time.Second)
			continue
		}

		err = br.serve()
		if err != nil {
			util.Log(br, "Server terminated with error: %v", err)
		}

		br.close()

		time.Sleep(time.Second)
	}
}

type bridge struct {
	tunnel *knx.Tunnel
	router *knx.Router
}

func newBridge(gatewayAddr, routerAddr string) (*bridge, error) {
	// Instantiate tunnel connection.
	tunnel, err := knx.NewTunnel(gatewayAddr, knxnet.TunnelLayerData, knx.DefaultTunnelConfig)
	if err != nil {
		return nil, err
	}

	// Instantiate routing facilities.
	router, err := knx.NewRouter(routerAddr, knx.DefaultRouterConfig)
	if err != nil {
		tunnel.Close()
		return nil, err
	}

	return &bridge{tunnel, router}, nil
}

func (br *bridge) serve() error {
	for {
		select {
		// Receive message from gateway.
		case msg, open := <-br.tunnel.Inbound():
			if !open {
				return errors.New("Tunnel channel closed")
			}

			if ind, ok := msg.(*cemi.LDataInd); ok {
				util.Log(br, "To router: %v", ind)
				if err := br.router.Send(ind); err != nil {
					return err
				}
			}

		// Receive message from router.
		case msg, open := <-br.router.Inbound():
			if !open {
				return errors.New("Router channel closed")
			}

			if ind, ok := msg.(*cemi.LDataInd); ok {
				util.Log(br, "To tunnel: %v", ind)
				if err := br.tunnel.Send(&cemi.LDataReq{LData: ind.LData}); err != nil {
					return err
				}
			}
		}
	}
}

func (br *bridge) close() {
	br.tunnel.Close()
	br.router.Close()
}
