// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"testing"

	"math"
	"math/rand"
)

// Test DPT 13.001 (counter value (pulses))
func TestDPT_13001(t *testing.T) {
	var buf []byte
	var src, dst DPT_13001

	// Corner cases
	for _, value := range []int32{math.MinInt32, math.MaxInt32} {
		src = DPT_13001(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Positive
	for i := 1; i <= 10; i++ {
		value := rand.Int31()

		src = DPT_13001(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Negative
	for i := 1; i <= 10; i++ {
		value := -rand.Int31()

		src = DPT_13001(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}
}

// Test DPT 13.002 (flow rate (m^3/h))
func TestDPT_13002(t *testing.T) {
	var buf []byte
	var src, dst DPT_13002

	// Corner cases
	for _, value := range []int32{math.MinInt32, math.MaxInt32} {
		src = DPT_13002(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Positive
	for i := 1; i <= 10; i++ {
		value := rand.Int31()

		src = DPT_13002(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Negative
	for i := 1; i <= 10; i++ {
		value := -rand.Int31()

		src = DPT_13002(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}
}

// Test DPT 13.010 (active energy (Wh))
func TestDPT_13010(t *testing.T) {
	var buf []byte
	var src, dst DPT_13010

	// Corner cases
	for _, value := range []int32{math.MinInt32, math.MaxInt32} {
		src = DPT_13010(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Positive
	for i := 1; i <= 10; i++ {
		value := rand.Int31()

		src = DPT_13010(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Negative
	for i := 1; i <= 10; i++ {
		value := -rand.Int31()

		src = DPT_13010(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}
}

// Test DPT 13.011 (apparant energy (VAh))
func TestDPT_13011(t *testing.T) {
	var buf []byte
	var src, dst DPT_13011

	// Corner cases
	for _, value := range []int32{math.MinInt32, math.MaxInt32} {
		src = DPT_13011(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Positive
	for i := 1; i <= 10; i++ {
		value := rand.Int31()

		src = DPT_13011(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Negative
	for i := 1; i <= 10; i++ {
		value := -rand.Int31()

		src = DPT_13011(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}
}

// Test DPT 13.012 (reactive energy (VARh))
func TestDPT_13012(t *testing.T) {
	var buf []byte
	var src, dst DPT_13012

	// Corner cases
	for _, value := range []int32{math.MinInt32, math.MaxInt32} {
		src = DPT_13012(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Positive
	for i := 1; i <= 10; i++ {
		value := rand.Int31()

		src = DPT_13012(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Negative
	for i := 1; i <= 10; i++ {
		value := -rand.Int31()

		src = DPT_13012(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}
}

// Test DPT 13.013 (active energy (kWh))
func TestDPT_13013(t *testing.T) {
	var buf []byte
	var src, dst DPT_13013

	// Corner cases
	for _, value := range []int32{math.MinInt32, math.MaxInt32} {
		src = DPT_13013(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Positive
	for i := 1; i <= 10; i++ {
		value := rand.Int31()

		src = DPT_13013(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Negative
	for i := 1; i <= 10; i++ {
		value := -rand.Int31()

		src = DPT_13013(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}
}

// Test DPT 13.014 (apparant energy (kVAh))
func TestDPT_13014(t *testing.T) {
	var buf []byte
	var src, dst DPT_13014

	// Corner cases
	for _, value := range []int32{math.MinInt32, math.MaxInt32} {
		src = DPT_13014(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Positive
	for i := 1; i <= 10; i++ {
		value := rand.Int31()

		src = DPT_13014(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Negative
	for i := 1; i <= 10; i++ {
		value := -rand.Int31()

		src = DPT_13014(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}
}

// Test DPT 13.015 (reactive energy (kVARh))
func TestDPT_13015(t *testing.T) {
	var buf []byte
	var src, dst DPT_13015

	// Corner cases
	for _, value := range []int32{math.MinInt32, math.MaxInt32} {
		src = DPT_13015(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Positive
	for i := 1; i <= 10; i++ {
		value := rand.Int31()

		src = DPT_13015(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Negative
	for i := 1; i <= 10; i++ {
		value := -rand.Int31()

		src = DPT_13015(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}
}

// Test DPT 13.016 (apparant energy (MWh))
func TestDPT_13016(t *testing.T) {
	var buf []byte
	var src, dst DPT_13016

	// Corner cases
	for _, value := range []int32{math.MinInt32, math.MaxInt32} {
		src = DPT_13016(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Positive
	for i := 1; i <= 10; i++ {
		value := rand.Int31()

		src = DPT_13016(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Negative
	for i := 1; i <= 10; i++ {
		value := -rand.Int31()

		src = DPT_13016(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}
}

// Test DPT 13.100 (delta time (s))
func TestDPT_13100(t *testing.T) {
	var buf []byte
	var src, dst DPT_13100

	// Corner cases
	for _, value := range []int32{math.MinInt32, math.MaxInt32} {
		src = DPT_13100(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Positive
	for i := 1; i <= 10; i++ {
		value := rand.Int31()

		src = DPT_13100(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}

	// Negative
	for i := 1; i <= 10; i++ {
		value := -rand.Int31()

		src = DPT_13100(value)
		if int32(src) != value {
			t.Errorf("Assignment of value \"%v\" failed! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int32(dst) != value {
			t.Errorf("Wrong value \"%s\" after pack/unpack! Original value was \"%v\".", dst, value)
		}
	}
}
