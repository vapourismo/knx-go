[![Build Status](https://travis-ci.org/vapourismo/knx-go.svg?branch=master)](https://travis-ci.org/vapourismo/knx-go)
[![GoDoc](https://godoc.org/github.com/vapourismo/knx-go?status.svg)](https://godoc.org/github.com/vapourismo/knx-go)

# knx-go

This repository contains a collection of Go packages that provide the means to communicate with KNX
networks.

## Packages

 Package           | Description
-------------------|--------------------------------------------------------------------
 **knx**           | Abstractions to communicate with KNXnet/IP servers
 **knx/knxnet**    | KNXnet/IP protocol services
 **knx/dpt**       | Datapoint types
 **knx/cemi**      | CEMI-encoded frames
 **cmd/knxbridge** | Tool to bridge KNX networks between a KNXnet/IP router and gateway

## Installation

Simply run the following command.

	$ go get -u github.com/vapourismo/knx-go/...

## Examples

### KNXnet/IP Group Client

If you simply want to send and receive group communication, the
[GroupTunnel](https://godoc.org/github.com/vapourismo/knx-go/knx#GroupTunnel) or
[GroupRouter](https://godoc.org/github.com/vapourismo/knx-go/knx#GroupRouter)
might be sufficient to you.

```go
package main

import (
	"log"
	"os"

	"github.com/vapourismo/knx-go/knx"
	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/dpt"
	"github.com/vapourismo/knx-go/knx/util"
)

func main() {
	// Setup logger for auxiliary logging. This enables us to see log messages from internal
	// routines.
	util.Logger = log.New(os.Stdout, "", log.LstdFlags)

	// Connect to the gateway.
	client, err := knx.NewGroupTunnel("10.0.0.7:3671", knx.DefaultTunnelConfig)
	if err != nil {
		util.Logger.Fatal(err)
	}

	// Close upon exiting. Even if the gateway closes the connection, we still have to clean up.
	defer client.Close()

	// Send 20.5°C to group 1/2/3.
	err = client.Send(0, cemi.NewGroupAddr3(1, 2, 3), dpt.ValueTemp(20.5).Pack())
	if err != nil {
		util.Logger.Fatal(err)
	}

	// Receive messages from the gateway. The inbound channel is closed with the connection.
	for msg := range client.Inbound() {
		util.Logger.Printf("%T %+v", msg, msg)
	}
}
```

In case you want to access a KNXnet/IP router instead of a gateway, simply replace

```go
client, err := knx.NewGroupTunnel("10.0.0.7:3671", knx.DefaultTunnelConfig)
```

with

```go
client, err := knx.NewGroupRouter("224.0.23.12:3671", knx.DefaultRouterConfig)
```

### KNXnet/IP CEMI Client

Use [Tunnel](https://godoc.org/github.com/vapourismo/knx-go/knx#Tunnel) or
[Router](https://godoc.org/github.com/vapourismo/knx-go/knx#Router) for finer control over the
communication with a gateway or router.

### KNX Bridge

The **knxbridge** tool (in package `cmd/knxbridge`) has multiple use cases.

Expose a KNX network behind a gateway at `10.0.0.2:3671` on the multicast group `224.0.23.12:3671`.
This allows routers and router clients to access the network.

	$ knxbridge 10.0.0.2:3671 224.0.23.12:3671

Connect two KNX networks through gateways. In this example one gateway is at `10.0.0.2:3671`, the
other is at `10.0.0.3:3671`.

	$ knxbridge 10.0.0.2:3671 10.0.0.3:3671
