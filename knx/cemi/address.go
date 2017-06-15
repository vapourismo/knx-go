// Copyright 2017 Ole Krüger.

package cemi

import (
	"fmt"
)

// IndividualAddr is an address for a KNX device.
type IndividualAddr uint16

// String generates a string representation.
func (addr IndividualAddr) String() string {
	return fmt.Sprintf("%d.%d.%d", uint8(addr>>12)&15, uint8(addr>>8)&15, uint8(addr))
}

// GroupAddr is an address for a KNX group object.
type GroupAddr uint16

// String generates a string representation.
func (addr GroupAddr) String() string {
	return fmt.Sprintf("%d/%d/%d", uint8(addr>>11)&31, uint8(addr>>8)&7, uint8(addr))
}
