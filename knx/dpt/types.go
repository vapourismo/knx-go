// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"fmt"
)

// A DatapointValue is a value of a datapoint.
type DatapointValue interface {
	// Pack the datapoint to a byte array.
	Pack() []byte

	// Unpack a the datapoint value from a byte array.
	Unpack(data []byte) error
}

// DatapointMeta gives meta information about a datapoint type.
type DatapointMeta interface {
	// Unit returns the unit of this datapoint type or empty string if it doesn't have a unit.
	Unit() string
}

// DPT_5001 represents DPT 5.001 / Scaling.
type DPT_5001 float32

func (d DPT_5001) Pack() []byte {
	if d <= 0 {
		return packU8(0)
	} else if d >= 100 {
		return packU8(255)
	} else {
		return packU8(uint8(d * 2.55))
	}
}

func (d *DPT_5001) Unpack(data []byte) error {
	var value uint8
	if err := unpackU8(data, &value); err != nil {
		return err
	}

	*d = DPT_5001(value) / 2.55

	return nil
}

func (d DPT_5001) Unit() string {
	return "%"
}

func (d DPT_5001) String() string {
	return fmt.Sprintf("%.2f%%", float32(d))
}

// DPT_5003 represents DPT 5.003 / Angle.
type DPT_5003 float32

func (d DPT_5003) Pack() []byte {
	if d <= 0 {
		return packU8(0)
	} else if d >= 360 {
		return packU8(255)
	} else {
		return packU8(uint8(d * 255 / 360))
	}
}

func (d *DPT_5003) Unpack(data []byte) error {
	var value uint8
	if err := unpackU8(data, &value); err != nil {
		return err
	}

	*d = DPT_5003(value) * 360 / 255

	return nil
}

func (d DPT_5003) Unit() string {
	return "°"
}

func (d DPT_5003) String() string {
	return fmt.Sprintf("%.2f°", float32(d))
}

// DPT_5004 represents DPT 5.004 / Percent_U8.
type DPT_5004 uint8

func (d DPT_5004) Pack() []byte {
	return packU8(uint8(d))
}

func (d *DPT_5004) Unpack(data []byte) error {
	return unpackU8(data, (*uint8)(d))
}

func (d DPT_5004) Unit() string {
	return "%"
}

func (d DPT_5004) String() string {
	return fmt.Sprintf("%.2f%%", float32(d))
}

// DPT_12001 represents DPT 12.001 / Unsigned counter.
type DPT_12001 uint32

func (d DPT_12001) Pack() []byte {
	return packU32(uint32(d))
}

func (d *DPT_12001) Unpack(data []byte) error {
	return unpackU32(data, (*uint32)(d))
}

func (d DPT_12001) Unit() string {
	return "pulses"
}

func (d DPT_12001) String() string {
	return fmt.Sprintf("%d pulses", uint32(d))
}

// DPT_17001 represents DPT 17.001 / Scene Number.
type DPT_17001 uint8

func (d DPT_17001) Pack() []byte {
	if d > 63 {
		return packU8(63)
	} else {
		return packU8(uint8(d))
	}
}

func (d *DPT_17001) Unpack(data []byte) error {
	var value uint8

	if err := unpackU8(data, &value); err != nil {
		return err
	}

	if *d <= 63 {
		*d = DPT_17001(value)
		return nil
	} else {
		*d = DPT_17001(63)
		return nil
	}
}

func (d DPT_17001) Unit() string {
	return ""
}

func (d DPT_17001) String() string {
	return fmt.Sprintf("%d", uint8(d))
}

// DPT_18001 represents DPT 18.001 / Scene Control.
type DPT_18001 uint8

func (d DPT_18001) Pack() []byte {
	if d <= 63 || (d >= 128 && d <= 191) {
		return packU8(uint8(d))
	} else {
		return packU8(63)
	}
}

func (d *DPT_18001) Unpack(data []byte) error {
	var value uint8

	if err := unpackU8(data, &value); err != nil {
		return err
	}

	if *d <= 63 || (*d >= 128 && *d <= 191) {
		*d = DPT_18001(value)
		return nil
	} else {
		*d = DPT_18001(63)
		return nil
	}
}

func (d DPT_18001) Unit() string {
	return ""
}

// KNX Association recommends to display the scene numbers [1..64].
// See note 6 of the KNX Specifications v2.1.
func (d DPT_18001) String() string {
	return fmt.Sprintf("%d", uint8(d))
}
