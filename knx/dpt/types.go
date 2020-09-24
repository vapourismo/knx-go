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

// DPT_1001 represents DPT 1.001 / Switch.
type DPT_1001 bool

func (d DPT_1001) Pack() []byte {
	return packB1(bool(d))
}

func (d *DPT_1001) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1001) Unit() string {
	return ""
}

func (d DPT_1001) String() string {
	if d {
		return "On"
	} else {
		return "Off"
	}
}

// DPT_1002 represents DPT 1.002 / Bool.
type DPT_1002 bool

func (d DPT_1002) Pack() []byte {
	return packB1(bool(d))
}

func (d *DPT_1002) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1002) Unit() string {
	return ""
}

func (d DPT_1002) String() string {
	if d {
		return "True"
	} else {
		return "False"
	}
}

// DPT_1003 represents DPT 1.003 / Enable.
type DPT_1003 bool

func (d DPT_1003) Pack() []byte {
	return packB1(bool(d))
}

func (d *DPT_1003) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1003) Unit() string {
	return ""
}

func (d DPT_1003) String() string {
	if d {
		return "Enable"
	} else {
		return "Disable"
	}
}

// DPT_1009 represents DPT 1.009 / OpenClose.
type DPT_1009 bool

func (d DPT_1009) Pack() []byte {
	return packB1(bool(d))
}

func (d *DPT_1009) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1009) Unit() string {
	return ""
}

func (d DPT_1009) String() string {
	if d {
		return "Close"
	} else {
		return "Open"
	}
}

// DPT_1010 represents DPT 1.010 / Start.
type DPT_1010 bool

func (d DPT_1010) Pack() []byte {
	return packB1(bool(d))
}

func (d *DPT_1010) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1010) Unit() string {
	return ""
}

func (d DPT_1010) String() string {
	if d {
		return "Start"
	} else {
		return "Stop"
	}
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

// DPT_9001 represents DPT 9.001 / Temperature.
type DPT_9001 float32

func (d DPT_9001) Pack() []byte {
	if d <= -273 {
		return packF16(-273)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9001) Unpack(data []byte) error {
	var value float32
	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < -273 {
		return fmt.Errorf("Temperature \"%.2f\" outside range [-273, 670760]", value)
	} else if value > 670760 {
		return fmt.Errorf("Temperature \"%.2f\" outside range [-273, 670760]", value)
	}

	*d = DPT_9001(value)

	return nil
}

func (d DPT_9001) Unit() string {
	return "°C"
}

func (d DPT_9001) String() string {
	return fmt.Sprintf("%.2f °C", float32(d))
}

// DPT_9004 represents DPT 9.004 / Illumination.
type DPT_9004 float32

func (d DPT_9004) Pack() []byte {
	if d <= 0 {
		return packF16(0)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9004) Unpack(data []byte) error {
	var value float32
	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < 0 {
		return fmt.Errorf("Illumination \"%.2f\" outside range [0, 670760]", value)
	} else if value > 670760 {
		return fmt.Errorf("Illumination \"%.2f\" outside range [0, 670760]", value)
	}

	*d = DPT_9004(value)

	return nil
}

func (d DPT_9004) Unit() string {
	return "lux"
}

func (d DPT_9004) String() string {
	return fmt.Sprintf("%.2f lux", float32(d))
}

// DPT_9005 represents DPT 9.005 / Wind Speed.
type DPT_9005 float32

func (d DPT_9005) Pack() []byte {
	if d <= 0 {
		return packF16(0)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9005) Unpack(data []byte) error {
	var value float32

	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < 0 {
		return fmt.Errorf("Wind speed \"%.2f\" outside range [0, 670760]", value)
	} else if value > 670760 {
		return fmt.Errorf("Wind speed \"%.2f\" outside range [0, 670760]", value)
	}

	*d = DPT_9005(value)
	return nil
}

func (d DPT_9005) Unit() string {
	return "m/s"
}

func (d DPT_9005) String() string {
	return fmt.Sprintf("%.2f m/s", float32(d))
}

// DPT_9007 represents DPT 9.007 / Humidity
type DPT_9007 float32

func (d DPT_9007) Pack() []byte {
	if d <= 0 {
		return packF16(0)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9007) Unpack(data []byte) error {
	var value float32

	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < 0 || value > 670760 {
		return fmt.Errorf("Humidity \"%.2f\" outside range [0, 670760]", value)
	}

	*d = DPT_9007(value)

	return nil
}

func (d DPT_9007) Unit() string {
	return "%"
}

func (d DPT_9007) String() string {
	return fmt.Sprintf("%.2f %%", float32(d))
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

// DPT_13001 represents DPT 13.001 / counter value.
type DPT_13001 int32

func (d DPT_13001) Pack() []byte {
	return packV32(int32(d))
}

func (d *DPT_13001) Unpack(data []byte) error {
	return unpackV32(data, (*int32)(d))
}

func (d DPT_13001) Unit() string {
	return "pulses"
}

func (d DPT_13001) String() string {
	return fmt.Sprintf("%d pulses", int32(d))
}

// DPT_13002 represents DPT 13.002 / flow rate.
type DPT_13002 int32

func (d DPT_13002) Pack() []byte {
	return packV32(int32(d))
}

func (d *DPT_13002) Unpack(data []byte) error {
	return unpackV32(data, (*int32)(d))
}

func (d DPT_13002) Unit() string {
	return "m^3/h"
}

func (d DPT_13002) String() string {
	return fmt.Sprintf("%d m^3/h", int32(d))
}

// DPT_13010 represents DPT 13.010 / active energy.
type DPT_13010 int32

func (d DPT_13010) Pack() []byte {
	return packV32(int32(d))
}

func (d *DPT_13010) Unpack(data []byte) error {
	return unpackV32(data, (*int32)(d))
}

func (d DPT_13010) Unit() string {
	return "Wh"
}

func (d DPT_13010) String() string {
	return fmt.Sprintf("%d Wh", int32(d))
}

// DPT_13011 represents DPT 13.011 / apparant energy.
type DPT_13011 int32

func (d DPT_13011) Pack() []byte {
	return packV32(int32(d))
}

func (d *DPT_13011) Unpack(data []byte) error {
	return unpackV32(data, (*int32)(d))
}

func (d DPT_13011) Unit() string {
	return "VAh"
}

func (d DPT_13011) String() string {
	return fmt.Sprintf("%d VAh", int32(d))
}

// DPT_13012 represents DPT 13.012 / reactive energy.
type DPT_13012 int32

func (d DPT_13012) Pack() []byte {
	return packV32(int32(d))
}

func (d *DPT_13012) Unpack(data []byte) error {
	return unpackV32(data, (*int32)(d))
}

func (d DPT_13012) Unit() string {
	return "VARh"
}

func (d DPT_13012) String() string {
	return fmt.Sprintf("%d VARh", int32(d))
}

// DPT_13013 represents DPT 13.010 / active energy (kWh).
type DPT_13013 int32

func (d DPT_13013) Pack() []byte {
	return packV32(int32(d))
}

func (d *DPT_13013) Unpack(data []byte) error {
	return unpackV32(data, (*int32)(d))
}

func (d DPT_13013) Unit() string {
	return "kWh"
}

func (d DPT_13013) String() string {
	return fmt.Sprintf("%d kWh", int32(d))
}

// DPT_13014 represents DPT 13.014 / apparant energy (kVAh).
type DPT_13014 int32

func (d DPT_13014) Pack() []byte {
	return packV32(int32(d))
}

func (d *DPT_13014) Unpack(data []byte) error {
	return unpackV32(data, (*int32)(d))
}

func (d DPT_13014) Unit() string {
	return "kVAh"
}

func (d DPT_13014) String() string {
	return fmt.Sprintf("%d kVAh", int32(d))
}

// DPT_13015 represents DPT 13.015 / reactive energy (kVARh).
type DPT_13015 int32

func (d DPT_13015) Pack() []byte {
	return packV32(int32(d))
}

func (d *DPT_13015) Unpack(data []byte) error {
	return unpackV32(data, (*int32)(d))
}

func (d DPT_13015) Unit() string {
	return "kVARh"
}

func (d DPT_13015) String() string {
	return fmt.Sprintf("%d kVARh", int32(d))
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
