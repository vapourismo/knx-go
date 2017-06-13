// Copyright 2017 Ole Krüger.

package dpt

import (
	"errors"
	"fmt"
)

// ErrInvalidLength is returned when the application data has unexpected length.
var ErrInvalidLength = errors.New("Given application data has invalid length")

// Switch is DPT 1.001.
type Switch bool

// String generates a string representation.
func (sw Switch) String() string {
	if sw {
		return "On"
	}

	return "Off"
}

// Size returns the assembled size.
func (Switch) Size() uint {
	return sizeB1
}

// Pack assembles the datapoint value in the given buffer.
func (sw Switch) Pack(buffer []byte) {
	packB1(buffer, bool(sw))
}

// Unpack extracts the value from the application data.
func (sw *Switch) Unpack(data []byte) (uint, error) {
	if len(data) != sizeB1 {
		return 0, ErrInvalidLength
	}

	unpackB1(data, (*bool)(sw))

	return sizeB1, nil
}

// ValueTemp is DPT 9.0001.
type ValueTemp float32

// String generates a string representation.
func (temp ValueTemp) String() string {
	return fmt.Sprintf("%.2f°C", float32(temp))
}

// Size returns the assembled size.
func (ValueTemp) Size() uint {
	return sizeF16
}

// Pack assembles the datapoint value in the given buffer.
func (temp ValueTemp) Pack(buffer []byte) {
	packF16(buffer, float32(temp))
}

// Unpack extracts the value from the application data.
func (temp *ValueTemp) Unpack(data []byte) (uint, error) {
	if len(data) != sizeF16 {
		return 0, ErrInvalidLength
	}

	unpackF16(data, (*float32)(temp))

	return sizeF16, nil
}
