// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"fmt"
)

// A DatapointValue is a value of a datapoint.
type DatapointValue interface {
	Pack() []byte
	Unpack(data []byte) error
}

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

// ValueBrightness is DPT 5.001 (0% - 100%)
type ValueBrightness int32

// Pack the datapoint value.
func (brightness ValueBrightness) Pack() []byte {
	return packI32(int32(brightness))
}

// Unpack the datapoint value from the given data.
func (brightness *ValueBrightness) Unpack(data []byte) error {
	return unpackI32(data, (*int32)(brightness))
}

// String generates a string representation.
func (brightness ValueBrightness) StringBrightness() string {
	return fmt.Sprintf("%d3", uint8(brightness))
}
