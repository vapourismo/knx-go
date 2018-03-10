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

// TrueFalse is DPT 1.002
type TrueFalse bool

// Pack the datapoint value.
func (b TrueFalse) Pack() []byte {
	return packB1(bool(b))
}

// Unpack the datapoint value from the given data.
func (b *TrueFalse) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(b))
}

// String generates a string representation.
func (sw TrueFalse) String() string {
	if sw {
		return "True"
	}
	return "False"
}

// DPT 1.009 OpenClose
type OpenClose bool

// Pack the datapoint value.
func (oc OpenClose) Pack() []byte {
	return packB1(bool(oc))
}

// Unpack the datapoint value from the given data.
func (oc *OpenClose) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(oc))
}

// String generates a string representation.
func (oc OpenClose) String() string {
	if oc {
		return "Close"
	}
	return "Open"
}

// DPT 1.010 StartStop
type StartStop bool

// Pack the datapoint value.
func (ss StartStop) Pack() []byte {
	return packB1(bool(ss))
}

// Unpack the datapoint value from the given data.
func (ss *StartStop) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(ss))
}

// String generates a string representation.
func (ss StartStop) String() string {
	if ss {
		return "Start"
	}
	return "Stop"
}

// DPT 5.001  (0% - 100%)
type Scaling int

// Pack the datapoint value.
func (sc Scaling) Pack() []byte {
	sc = sc * 255 / 100
	return packInt(uint(sc))
}

// Unpack the datapoint value from the given data.
func (sc *Scaling) Unpack(data []byte) error {
	var d uint = (uint(data[1]) * 100) / 255
	data[1] = byte(d)
	return unpackInt(data, (*Scaling)(sc))
}

// String generates a string representation.
func (sc Scaling) String() string {
	return fmt.Sprintf("%d%%", uint(sc))
}

// ValueTemp DPT 9.001
type ValueTemp float32

// Pack the datapoint value.
func (vt ValueTemp) Pack() []byte {
	return packF16(float32(vt))
}

// Unpack the datapoint value from the given data.
func (vt *ValueTemp) Unpack(data []byte) error {
	return unpackF16(data, (*float32)(vt))
}

// String generates a string representation.
func (vt ValueTemp) String() string {
	return fmt.Sprintf("%.2f°C", float32(vt))
}
