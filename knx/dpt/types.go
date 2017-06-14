// Copyright 2017 Ole Krüger.

package dpt

import (
	"fmt"
)

// Switch is DPT 1.001.
type Switch bool

// Pack the datapoint value.
func (sw Switch) Pack() []byte {
	return packB1(bool(sw))
}

// Unpack the datapoint value from the given data.
func (sw *Switch) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(sw))
}

// String generates a string representation.
func (sw Switch) String() string {
	if sw {
		return "On"
	}

	return "Off"
}

// ValueTemp is DPT 9.001.
type ValueTemp float32

// Pack the datapoint value.
func (temp ValueTemp) Pack() []byte {
	return packF16(float32(temp))
}

// Unpack the datapoint value from the given data.
func (temp *ValueTemp) Unpack(data []byte) error {
	return unpackF16(data, (*float32)(temp))
}

// String generates a string representation.
func (temp ValueTemp) String() string {
	return fmt.Sprintf("%.2f°C", float32(temp))
}
