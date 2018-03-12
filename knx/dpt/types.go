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

// TrueFalse is DPT 1.002.
type TrueFalse Switch

// String generates a string representation.
func (sw TrueFalse) String() string {
	if sw {
		return "True"
	}

	return "False"
}

// OpenClose is DPT 1.009.
type OpenClose Switch

// String generates a string representation.
func (oc OpenClose) String() string {
	if oc {
		return "Close"
	}

	return "Open"
}

// StartStop is DPT 1.010.
type StartStop Switch

// String generates a string representation.
func (ss StartStop) String() string {
	if ss {
		return "Start"
	}

	return "Stop"
}

// Scaling is DPT 5.001. Scaling goes from 0 to 1.
type Scaling float32

// Pack the datapoint value.
func (sc Scaling) Pack() []byte {
	if sc >= 1 {
		return packU8(255)
	}

	if sc <= 0 {
		return packU8(0)
	}

	return packU8(uint8(sc * 255))
}

// Unpack the datapoint value from the given data.
func (sc *Scaling) Unpack(data []byte) (err error) {
	var value uint8
	err = unpackU8(data, &value)

	if err == nil {
		*sc = Scaling(value) / 255
	}

	return
}

// String generates a string representation.
func (sc Scaling) String() string {
	return fmt.Sprintf("%.2f%%", sc * 100)
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
