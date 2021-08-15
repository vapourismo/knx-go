// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"fmt"
)

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
