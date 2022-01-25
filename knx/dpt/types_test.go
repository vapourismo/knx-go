// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"testing"

	"math"
	"math/rand"
)

// Define epsilon constant for floating point checks
const epsilon = 1e-3

func abs(x float32) float32 {
	if x < 0.0 {
		return -x
	}
	return x
}

func get_float_quantization_error(value, resolution float32, mantis int) float32 {
	// Calculate the exponent for the value given the mantis and value resolution
	value_m := value / (resolution * float32(mantis))
	value_exp := math.Ceil(math.Log2(float64(value_m)))

	// Calculate the worst quantization error by assuming the
	// mantis to be off by one
	q := math.Pow(2, value_exp)

	// Scale back the quantization error with the given resolution
	return float32(q) / resolution
}

// Test DPT 1.001 (Switch) with values within range
func TestDPT_1001(t *testing.T) {
	var buf []byte
	var src, dst DPT_1001

	for _, value := range []bool{true, false} {
		src = DPT_1001(value)
		if bool(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if bool(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}
}

// Test DPT 1.002 (Bool) with values within range
func TestDPT_1002(t *testing.T) {
	var buf []byte
	var src, dst DPT_1002

	for _, value := range []bool{true, false} {
		src = DPT_1002(value)
		if bool(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if bool(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}
}

// Test DPT 1.003 (Enable) with values within range
func TestDPT_1003(t *testing.T) {
	var buf []byte
	var src, dst DPT_1003

	for _, value := range []bool{true, false} {
		src = DPT_1003(value)
		if bool(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if bool(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}
}

// Test DPT 1.008 (Up/Down) with values within range
func TestDPT_1008(t *testing.T) {
	var buf []byte
	var src, dst DPT_1008

	for _, value := range []bool{true, false} {
		src = DPT_1008(value)
		if bool(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if bool(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}
}

// Test DPT 1.009 (OpenClose) with values within range
func TestDPT_1009(t *testing.T) {
	var buf []byte
	var src, dst DPT_1009

	for _, value := range []bool{true, false} {
		src = DPT_1009(value)
		if bool(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if bool(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}
}

// Test DPT 1.010 (Start) with values within range
func TestDPT_1010(t *testing.T) {
	var buf []byte
	var src, dst DPT_1010

	for _, value := range []bool{true, false} {
		src = DPT_1010(value)
		if bool(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if bool(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}
}

// Test DPT 1.100 (Heat/Cool) with values within range
func TestDPT_1100(t *testing.T) {
	var buf []byte
	var src, dst DPT_1100

	for _, value := range []bool{true, false} {
		src = DPT_1100(value)
		if bool(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if bool(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}
}

// Test DPT 5.001 (Scaling) with values within range
func TestDPT_5001(t *testing.T) {
	var buf []byte
	var src, dst DPT_5001

	// Calculate the quantization error we expect
	const Q = float32(100) / 255

	for i := 1; i <= 10; i++ {
		value := rand.Float32()

		// Scale the random number to the given range
		value *= 100

		// Pack and unpack to test value
		src = DPT_5001(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_5001! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 5.003 (Angle) with values within range
func TestDPT_5003(t *testing.T) {
	var buf []byte
	var src, dst DPT_5003

	// Calculate the quantization error we expect
	const Q = float32(360) / 255

	for i := 1; i <= 10; i++ {
		value := rand.Float32()

		// Scale the random number to the given range
		value *= 360

		// Pack and unpack to test value
		src = DPT_5003(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_5003! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 12.001 (Unsigned counter) with values within range
func TestDPT_12001(t *testing.T) {
	var buf []byte
	var src, dst DPT_12001

	for i := 1; i <= 10; i++ {
		value := rand.Uint32()

		// Pack and unpack to test value
		src = DPT_12001(value)
		if uint32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_12001! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if uint32(dst) != value {
			t.Errorf("Value \"%s\" after pack/unpack different from Original value. Was \"%v\"", dst, value)
		}
	}
}

// Test DPT 17.001 (scene number)
func TestDPT_17001(t *testing.T) {
	var buf []byte
	var src, dst DPT_17001

	for i := 0; i <= 255; i++ {
		value := uint8(i)
		src = DPT_17001(value)
		buf = src.Pack()
		dst.Unpack(buf)
		if value <= 63 && uint8(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, src)
		} else if value > 63 && uint8(dst) != 63 {
			t.Errorf("Wrong value \"%s\" (undefined range) after pack/unpack! Original value was \"%v\".", dst, src)
		}
	}
}

// Test DPT 18.001 (scene control)
func TestDPT_18001(t *testing.T) {
	var buf []byte
	var src, dst DPT_18001

	for i := 0; i <= 255; i++ {
		value := uint8(i)
		src = DPT_18001(value)
		buf = src.Pack()
		dst.Unpack(buf)
		if value <= 63 && uint8(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, src)
		}
		if value > 63 && value < 128 && uint8(dst) != 63 {
			t.Errorf("Wrong value \"%s\" (undefined lower range) after pack/unpack! Original value was \"%v\".", dst, src)
		}
		if value >= 128 && value <= 191 && uint8(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, src)
		}
		if value > 191 && uint8(dst) != 63 {
			t.Errorf("Wrong value \"%s\" (undefined upper range) after pack/unpack! Original value was \"%v\".", dst, src)
		}
	}
}
