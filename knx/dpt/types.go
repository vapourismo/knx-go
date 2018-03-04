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
func (sw Switch) SwitchString() string {
	if sw {
		return "On"
	}
	return "Off"
}

// Bool is DPT 1.002
type Bool bool

// Pack the datapoint value.
func (b Bool) Pack() []byte {
	return packB1(bool(b))
}

// Unpack the datapoint value from the given data.
func (b *Bool) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(b))
}

// String generates a string representation.
func (sw Switch) BoolString() string {
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
func (oc OpenClose) OpenCloseString() string {
	if oc {
		return "Close"
	}
	return "Open"
}

// DPT 1.010 Start
type Start bool

// Pack the datapoint value.
func (st Start) Pack() []byte {
	return packB1(bool(st))
}

// Unpack the datapoint value from the given data.
func (st *Start) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(st))
}

// String generates a string representation.
func (sw Switch) StartString() string {
	if sw {
		return "Start"
	}
	return "Stop"
}



// DPT 5.001  (0% - 100%)
type Scaling int32

// Pack the datapoint value.
func (sc Scaling) Pack() []byte {
	v := packI32(int32(sc))
	//fmt.Printf(" packI32: %+v\n", v)
	return v
}

// Unpack the datapoint value from the given data.
func (sc *Scaling) Unpack(data []byte) error {
	//v := unpackI32(data, (*int32)(sc))
	return unpackI32(data, (*int32)(sc))
}

// String generates a string representation.
func (sc Scaling) StringScaling() string {

	return fmt.Sprintf("%d3", uint8(sc))
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






