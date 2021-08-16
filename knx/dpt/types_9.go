// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"fmt"
)

// DPT_9001 represents DPT 9.001 / Temperature °C.
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
	if value < -273 || value > 670760 {
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

// DPT_9002 represents DPT 9.002 / Temperature K.
type DPT_9002 float32

func (d DPT_9002) Pack() []byte {
	if d <= -670760 {
		return packF16(-670760)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9002) Unpack(data []byte) error {
	var value float32
	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < -670760 || value > 670760 {
		return fmt.Errorf("Temperature \"%.2f\" outside range [-670760, 670760]", value)
	}

	*d = DPT_9002(value)
	return nil
}

func (d DPT_9002) Unit() string {
	return "K"
}

func (d DPT_9002) String() string {
	return fmt.Sprintf("%.2f K", float32(d))
}

// DPT_9003 represents DPT 9.003 / Temperature K/h.
type DPT_9003 float32

func (d DPT_9003) Pack() []byte {
	if d <= -670760 {
		return packF16(-670760)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9003) Unpack(data []byte) error {
	var value float32
	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < -670760 || value > 670760 {
		return fmt.Errorf("Temperature \"%.2f\" outside range [-670760, 670760]", value)
	}

	*d = DPT_9003(value)
	return nil
}

func (d DPT_9003) Unit() string {
	return "K/h"
}

func (d DPT_9003) String() string {
	return fmt.Sprintf("%.2f K/h", float32(d))
}

// DPT_9004 represents DPT 9.004 / Illumination lux.
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
	if value < 0 || value > 670760 {
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

// DPT_9005 represents DPT 9.005 / Wind Speed m/s.
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
	if value < 0 || value > 670760 {
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

// DPT_9006 represents DPT 9.006 / Pressure Pa.
type DPT_9006 float32

func (d DPT_9006) Pack() []byte {
	if d <= 0 {
		return packF16(0)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9006) Unpack(data []byte) error {
	var value float32

	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < 0 || value > 670760 {
		return fmt.Errorf("Pressure \"%.2f\" outside range [0, 670760]", value)
	}

	*d = DPT_9006(value)
	return nil
}

func (d DPT_9006) Unit() string {
	return "Pa"
}

func (d DPT_9006) String() string {
	return fmt.Sprintf("%.2f Pa", float32(d))
}

// DPT_9007 represents DPT 9.007 / Humidity %
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

// DPT_9008 represents DPT 9.008 / Air quality ppm
type DPT_9008 float32

func (d DPT_9008) Pack() []byte {
	if d <= 0 {
		return packF16(0)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9008) Unpack(data []byte) error {
	var value float32

	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < 0 || value > 670760 {
		return fmt.Errorf("Air quality \"%.2f\" outside range [0, 670760]", value)
	}

	*d = DPT_9008(value)

	return nil
}

func (d DPT_9008) Unit() string {
	return "ppm"
}

func (d DPT_9008) String() string {
	return fmt.Sprintf("%.2f ppm", float32(d))
}

// DPT_9010 represents DPT 9.010 / Time s.
type DPT_9010 float32

func (d DPT_9010) Pack() []byte {
	if d <= -670760 {
		return packF16(-670760)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9010) Unpack(data []byte) error {
	var value float32
	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < -670760 || value > 670760 {
		return fmt.Errorf("Time \"%.2f\" outside range [-670760, 670760]", value)
	}

	*d = DPT_9010(value)
	return nil
}

func (d DPT_9010) Unit() string {
	return "s"
}

func (d DPT_9010) String() string {
	return fmt.Sprintf("%.2f s", float32(d))
}

// DPT_9011 represents DPT 9.011 / Time ms.
type DPT_9011 float32

func (d DPT_9011) Pack() []byte {
	if d <= -670760 {
		return packF16(-670760)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9011) Unpack(data []byte) error {
	var value float32
	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < -670760 || value > 670760 {
		return fmt.Errorf("Time \"%.2f\" outside range [-670760, 670760]", value)
	}

	*d = DPT_9011(value)
	return nil
}

func (d DPT_9011) Unit() string {
	return "ms"
}

func (d DPT_9011) String() string {
	return fmt.Sprintf("%.2f ms", float32(d))
}

// DPT_9020 represents DPT 9.020 / Volt mV.
type DPT_9020 float32

func (d DPT_9020) Pack() []byte {
	if d <= -670760 {
		return packF16(-670760)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9020) Unpack(data []byte) error {
	var value float32
	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < -670760 || value > 670760 {
		return fmt.Errorf("Volt \"%.2f\" outside range [-670760, 670760]", value)
	}

	*d = DPT_9020(value)
	return nil
}

func (d DPT_9020) Unit() string {
	return "mV"
}

func (d DPT_9020) String() string {
	return fmt.Sprintf("%.2f mV", float32(d))
}

// DPT_9021 represents DPT 9.021 / Current mA.
type DPT_9021 float32

func (d DPT_9021) Pack() []byte {
	if d <= -670760 {
		return packF16(-670760)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9021) Unpack(data []byte) error {
	var value float32
	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < -670760 || value > 670760 {
		return fmt.Errorf("Current \"%.2f\" outside range [-670760, 670760]", value)
	}

	*d = DPT_9021(value)
	return nil
}

func (d DPT_9021) Unit() string {
	return "mA"
}

func (d DPT_9021) String() string {
	return fmt.Sprintf("%.2f mA", float32(d))
}

// DPT_9022 represents DPT 9.022 / Power Density W/m2.
type DPT_9022 float32

func (d DPT_9022) Pack() []byte {
	if d <= -670760 {
		return packF16(-670760)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9022) Unpack(data []byte) error {
	var value float32
	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < -670760 || value > 670760 {
		return fmt.Errorf("Power Density \"%.2f\" outside range [-670760, 670760]", value)
	}

	*d = DPT_9022(value)
	return nil
}

func (d DPT_9022) Unit() string {
	return "W/m2"
}

func (d DPT_9022) String() string {
	return fmt.Sprintf("%.2f W/m2", float32(d))
}

// DPT_9023 represents DPT 9.023 / Kelvin per Percent K/%.
type DPT_9023 float32

func (d DPT_9023) Pack() []byte {
	if d <= -670760 {
		return packF16(-670760)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9023) Unpack(data []byte) error {
	var value float32
	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < -670760 || value > 670760 {
		return fmt.Errorf("Kelvin per percent \"%.2f\" outside range [-670760, 670760]", value)
	}

	*d = DPT_9023(value)
	return nil
}

func (d DPT_9023) Unit() string {
	return "K/%%"
}

func (d DPT_9023) String() string {
	return fmt.Sprintf("%.2f K/%%", float32(d))
}

// DPT_9024 represents DPT 9.024 / Power kW.
type DPT_9024 float32

func (d DPT_9024) Pack() []byte {
	if d <= -670760 {
		return packF16(-670760)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9024) Unpack(data []byte) error {
	var value float32
	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < -670760 || value > 670760 {
		return fmt.Errorf("Power \"%.2f\" outside range [-670760, 670760]", value)
	}

	*d = DPT_9024(value)
	return nil
}

func (d DPT_9024) Unit() string {
	return "kW"
}

func (d DPT_9024) String() string {
	return fmt.Sprintf("%.2f kW", float32(d))
}

// DPT_9025 represents DPT 9.025 / Volume Flow l/h.
type DPT_9025 float32

func (d DPT_9025) Pack() []byte {
	if d <= -670760 {
		return packF16(-670760)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9025) Unpack(data []byte) error {
	var value float32
	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < -670760 || value > 670760 {
		return fmt.Errorf("Volume flow \"%.2f\" outside range [-670760, 670760]", value)
	}

	*d = DPT_9025(value)
	return nil
}

func (d DPT_9025) Unit() string {
	return "l/h"
}

func (d DPT_9025) String() string {
	return fmt.Sprintf("%.2f l/h", float32(d))
}

// DPT_9026 represents DPT 9.026 / Rain amount l/m^2.
type DPT_9026 float32

func (d DPT_9026) Pack() []byte {
	if d <= -670760 {
		return packF16(-670760)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9026) Unpack(data []byte) error {
	var value float32
	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < -670760 || value > 670760 {
		return fmt.Errorf("Rain amount \"%.2f\" outside range [-670760, 670760]", value)
	}

	*d = DPT_9026(value)
	return nil
}

func (d DPT_9026) Unit() string {
	return "l/m^2"
}

func (d DPT_9026) String() string {
	return fmt.Sprintf("%.2f l/m^2", float32(d))
}

// DPT_9027 represents DPT 9.027 / Temperature °F.
type DPT_9027 float32

func (d DPT_9027) Pack() []byte {
	if d <= -459.6 {
		return packF16(-459.6)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9027) Unpack(data []byte) error {
	var value float32
	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < -459.6 || value > 670760 {
		return fmt.Errorf("Rain amount \"%.2f\" outside range [-670760, 670760]", value)
	}

	*d = DPT_9027(value)
	return nil
}

func (d DPT_9027) Unit() string {
	return "°F"
}

func (d DPT_9027) String() string {
	return fmt.Sprintf("%.2f °F", float32(d))
}

// DPT_9028 represents DPT 9.028 / Wind Speed km/h.
type DPT_9028 float32

func (d DPT_9028) Pack() []byte {
	if d <= 0 {
		return packF16(0)
	} else if d >= 670760 {
		return packF16(670760)
	} else {
		return packF16(float32(d))
	}
}

func (d *DPT_9028) Unpack(data []byte) error {
	var value float32
	if err := unpackF16(data, &value); err != nil {
		return err
	}

	// Check the value for valid range
	if value < 0 || value > 670760 {
		return fmt.Errorf("Wind speed \"%.2f\" outside range [-670760, 670760]", value)
	}

	*d = DPT_9028(value)
	return nil
}

func (d DPT_9028) Unit() string {
	return "km/h"
}

func (d DPT_9028) String() string {
	return fmt.Sprintf("%.2f km/h", float32(d))
}
