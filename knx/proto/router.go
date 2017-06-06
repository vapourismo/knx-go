package proto

import (
	"fmt"
	"io"
	"time"

	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/util"
)

// A RoutingInd indicates to one or more routers that the contents shall be routed.
type RoutingInd struct {
	// L_Data.ind to be routed
	Payload cemi.Message
}

// Service returns the service identifiers for routing indication.
func (RoutingInd) Service() ServiceID {
	return RoutingIndService
}

// Unpack initializes the structure by parsing the given data.
func (ind *RoutingInd) Unpack(data []byte) (uint, error) {
	return cemi.Unpack(data, &ind.Payload)
}

// WriteTo serializes the structure and writes it to the given Writer.
func (ind *RoutingInd) WriteTo(w io.Writer) (int64, error) {
	return cemi.Pack(w, ind.Payload)
}

// DeviceState indicates the state of a device.
type DeviceState uint8

// These are known device states.
const (
	DeviceStateReserved DeviceState = 0xfc
	DeviceStateKNXError DeviceState = 0x01
	DeviceStateIPError  DeviceState = 0x02
)

// String converts the device status to a s tring.
func (status DeviceState) String() string {
	switch status {
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

// Unpack initializes the structure by parsing the given data.
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

// Unpack initializes the structure by parsing the given data.
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
