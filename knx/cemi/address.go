// Copyright 2017 Ole Krüger.
// Copyright 2022 Martin Müller.
// Licensed under the MIT license which can be found in the LICENSE file.

package cemi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// IndividualAddr is an individual address for a KNX device. Its format is defined in
// 03_03_02 Data Link Layer General v01.02.02 AS.pdf
// 1.4.2 Individual Address
// Individual address zero (0.0.0) is not allowed.
type IndividualAddr uint16

// NewIndividualAddr3 generates an individual address from an
// "a.b.c" representation, where a is the Area Address [0..15],
// b is the Line Address [0..15] and c is the Device Address [0..255].
func NewIndividualAddr3(a, b, c uint8) IndividualAddr {
	return IndividualAddr(a&0xF)<<12 | IndividualAddr(b&0xF)<<8 | IndividualAddr(c)
}

// NewIndividualAddr2 generates an individual address from an
// "a.b" representation, where a is the Subnetwork Address [0..255],
// b is the Device Address [0..255].
func NewIndividualAddr2(a, b uint8) IndividualAddr {
	return IndividualAddr(a)<<8 | IndividualAddr(b)
}

// NewIndividualAddrString parses the given string to an individual address.
// Supported formats are
// %d.%d.%d ([0..15], [0..15], [0..255]),
// %d.%d ([0..255], [0..255]) and
// %d ([0..65535]). Validity is checked.
func NewIndividualAddrString(addr string) (IndividualAddr, error) {
	var nums []int

	numstrings := strings.Split(addr, ".")

	for _, s := range numstrings {
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		nums = append(nums, i)
	}

	switch len(nums) {
	case 3:
		if nums[0] < 0 || nums[0] > 15 ||
			nums[1] < 0 || nums[1] > 15 ||
			nums[2] < 0 || nums[2] > 255 {
			return 0, fmt.Errorf("invalid area, line or device address in %s", addr)
		}
		if nums[0] == 0 && nums[1] == 0 && nums[2] == 0 {
			return 0, errors.New("invalid individual address 0.0.0")
		}
		return NewIndividualAddr3(uint8(nums[0]), uint8(nums[1]), uint8(nums[2])), nil
	case 2:
		if nums[0] < 0 || nums[0] > 255 || nums[1] < 0 || nums[1] > 255 {
			return 0, fmt.Errorf("invalid subnetwork or device address in %s", addr)
		}
		if nums[0] == 0 && nums[1] == 0 {
			return 0, errors.New("invalid individual address 0.0")
		}
		return NewIndividualAddr2(uint8(nums[0]), uint8(nums[1])), nil
	case 1:
		if nums[0] <= 0 || nums[0] > 65535 {
			return 0, fmt.Errorf("invalid raw individual address in %s", addr)
		}
		return IndividualAddr(nums[0]), nil
	}

	return 0, errors.New("string cannot be parsed to an individual address")
}

// String generates a string representation "a.b.c" where
// a = Area Address = 4 bits, b = Line Address = 4 bits,
// c = Device Address = 1 byte.
func (addr IndividualAddr) String() string {
	return fmt.Sprintf("%d.%d.%d", uint8(addr>>12)&0xF, uint8(addr>>8)&0xF, uint8(addr))
}

// GroupAddr is an address for a KNX group object. Group address
// zero (0/0/0) is not allowed.
type GroupAddr uint16

// NewGroupAddr3 generates a group address from an "a/b/c"
// representation, where a is the Main Group [0..31], b is
// the Middle Group [0..7], c is the Sub Group [0..255].
func NewGroupAddr3(a, b, c uint8) GroupAddr {
	return GroupAddr(a&0x1F)<<11 | GroupAddr(b&0x7)<<8 | GroupAddr(c)
}

// NewGroupAddr2 generates a group address from and "a/b"
// representation, where a is the Main Group [0..31] and b is
// the Sub Group [0..2047].
func NewGroupAddr2(a uint8, b uint16) GroupAddr {
	return GroupAddr(a)<<8 | GroupAddr(b&0x7FF)
}

// NewGroupAddrString parses the given string to a group address.
// Supported formats are:
// %d/%d/%d ([0..31], [0..7], [0..255]),
// %d/%d ([0..31], [0..2047]) and
// %d ([0..65535]). Validity is checked.
func NewGroupAddrString(addr string) (GroupAddr, error) {
	var nums []int

	numstrings := strings.Split(addr, "/")

	for _, s := range numstrings {
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("invalid parsing %s", err)
		}
		nums = append(nums, i)
	}

	switch len(nums) {
	case 3:
		if nums[0] < 0 || nums[0] > 31 ||
			nums[1] < 0 || nums[1] > 7 ||
			nums[2] < 0 || nums[2] > 255 {
			return 0, fmt.Errorf("invalid main, middle or sub group address in %s", addr)
		}
		if nums[0] == 0 && nums[1] == 0 && nums[2] == 0 {
			return 0, errors.New("invalid group address 0/0/0")
		}
		return NewGroupAddr3(uint8(nums[0]), uint8(nums[1]), uint8(nums[2])), nil
	case 2:
		if nums[0] < 0 || nums[0] > 31 ||
			nums[1] < 0 || nums[1] > 2047 {
			return 0, fmt.Errorf("invalid main or sub group address in %s", addr)
		}
		if nums[0] == 0 && nums[1] == 0 {
			return 0, errors.New("invalid group address 0/0")
		}
		return NewGroupAddr2(uint8(nums[0]), uint16(nums[1])), nil
	case 1:
		if nums[0] <= 0 || nums[0] > 65535 {
			return 0, fmt.Errorf("invalid raw group address in %s", addr)
		}
		return GroupAddr(nums[0]), nil
	}

	return 0, errors.New("string cannot be parsed to a group address")
}

// String generates a string representation with groups "a/b/c" where
// a = Main Group = 5 bits, b = Middle Group = 3 bits, c = Sub Group = 1 byte.
func (addr GroupAddr) String() string {
	return fmt.Sprintf("%d/%d/%d", uint8(addr>>11)&0x1F, uint8(addr>>8)&0x7, uint8(addr))
}
