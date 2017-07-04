// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package cemi

import (
	"errors"
	"fmt"
)

// IndividualAddr is an address for a KNX device.
type IndividualAddr uint16

// NewIndividualAddr3 generates a group address with format a/b/c.
func NewIndividualAddr3(a, b, c uint8) IndividualAddr {
	return IndividualAddr(a&15)<<12 | IndividualAddr(b&15)<<8 | IndividualAddr(c)
}

// NewIndividualAddr2 generates a group address with format a/b.
func NewIndividualAddr2(a, b uint8) IndividualAddr {
	return IndividualAddr(a)<<8 | IndividualAddr(b)
}

// NewIndividualAddrString parses the given string as a individual address. Supported formats are
// %d.%d.%d, %d.%d and %d.
func NewIndividualAddrString(addr string) (IndividualAddr, error) {
	var a, b, c uint
	num, _ := fmt.Sscanf(addr, "%d.%d.%d", &a, &b, &c)

	switch {
	case num > 2:
		return NewIndividualAddr3(uint8(a), uint8(b), uint8(c)), nil

	case num > 1:
		return NewIndividualAddr2(uint8(a), uint8(b)), nil

	case num > 0:
		return IndividualAddr(a), nil
	}

	return 0, errors.New("Input is not an individual address")
}

// String generates a string representation.
func (addr IndividualAddr) String() string {
	return fmt.Sprintf("%d.%d.%d", uint8(addr>>12)&15, uint8(addr>>8)&15, uint8(addr))
}

// GroupAddr is an address for a KNX group object.
type GroupAddr uint16

// NewGroupAddr3 generates a group address with format a/b/c.
func NewGroupAddr3(a, b, c uint8) GroupAddr {
	return GroupAddr(a&31)<<11 | GroupAddr(b&7)<<8 | GroupAddr(c)
}

// NewGroupAddr2 generates a group address with format a/b.
func NewGroupAddr2(a, b uint8) GroupAddr {
	return GroupAddr(a)<<8 | GroupAddr(b)
}

// NewGroupAddrString parses the given string as a group address. Supported formats are %d/%d/%d,
// %d/%d and %d.
func NewGroupAddrString(addr string) (GroupAddr, error) {
	var a, b, c uint
	num, _ := fmt.Sscanf(addr, "%d/%d/%d", &a, &b, &c)

	switch {
	case num > 2:
		return NewGroupAddr3(uint8(a), uint8(b), uint8(c)), nil

	case num > 1:
		return NewGroupAddr2(uint8(a), uint8(b)), nil

	case num > 0:
		return GroupAddr(a), nil
	}

	return 0, errors.New("Input is not a group address")
}

// String generates a string representation.
func (addr GroupAddr) String() string {
	return fmt.Sprintf("%d/%d/%d", uint8(addr>>11)&31, uint8(addr>>8)&7, uint8(addr))
}
