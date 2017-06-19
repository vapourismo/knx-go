// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package knxnet

import (
	"fmt"
	"time"

	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/util"
)

// A RoutingInd indicates to one or more routers that the contents shall be routed.
type RoutingInd struct {
	Payload cemi.Message
}

// Service returns the service identifiers for routing indication.
func (RoutingInd) Service() ServiceID {
	return RoutingIndService
}

// Size returns the packed size.
func (ind *RoutingInd) Size() uint {
	return cemi.Size(ind.Payload)
}

// Pack assembles the service payload in the given buffer.
func (ind *RoutingInd) Pack(buffer []byte) {
	cemi.Pack(buffer, ind.Payload)
}

// Unpack parses the given service payload in order to initialize the structure.
func (ind *RoutingInd) Unpack(data []byte) (uint, error) {
	return cemi.Unpack(data, &ind.Payload)
}

// DeviceState indicates the state of a device.
type DeviceState uint8

// These are known device states.
const (
	DeviceStateOk       DeviceState = 0x00
	DeviceStateKNXError DeviceState = 0x01
	DeviceStateIPError  DeviceState = 0x02
	DeviceStateReserved DeviceState = 0xfc
)

// String converts the device status to a string.
func (status DeviceState) String() string {
	switch status {
	case DeviceStateOk:
		return "Ok"

	case DeviceStateReserved:
		return "Reserved"

	case DeviceStateKNXError:
		return "KNX error"

	case DeviceStateIPError:
		return "IP error"

	default:
		return fmt.Sprintf("Unknown device status %#x", uint8(status))
	}
}

// A RoutingLost indicates that a packet got lost.
type RoutingLost struct {
	// Device status
	Status DeviceState

	// Number of packets lost
	Count uint16
}

// Service returns the service identifiers for routing lost indication.
func (RoutingLost) Service() ServiceID {
	return RoutingLostService
}

// Unpack parses the given service payload in order to initialize the structure.
func (rl *RoutingLost) Unpack(data []byte) (uint, error) {
	var length uint8
	return util.UnpackSome(data, &length, (*uint8)(&rl.Status), &rl.Count)

	// TODO: Find out if length is supposed to be 4; validate it, if so.
}

// A RoutingBusy indicates that a router is busy.
type RoutingBusy struct {
	// Device status
	Status DeviceState

	// Time to wait
	WaitTime time.Duration

	// ?
	Control uint16
}

// Service returns the service identifiers for routing busy indication.
func (RoutingBusy) Service() ServiceID {
	return RoutingBusyService
}

// Unpack parses the given service payload in order to initialize the structure.
func (rl *RoutingBusy) Unpack(data []byte) (n uint, err error) {
	var length uint8
	var waitTime uint16

	if n, err = util.UnpackSome(
		data, &length, (*uint8)(&rl.Status), &waitTime, &rl.Control,
	); err != nil {
		return
	}

	// TODO: Find out if length is supposed to be 6; validate it, if so.

	rl.WaitTime = time.Duration(waitTime) * time.Millisecond

	return
}
