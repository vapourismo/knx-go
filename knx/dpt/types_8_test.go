// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"fmt"
	"testing"

	"math"
	"math/rand"
)

// Test DPT 8.001 with values within range
func TestDPT_8001(t *testing.T) {
	var buf []byte
	var src, dst DPT_8001

	for i := 1; i <= 10; i++ {
		value := int16(rand.Uint32() % math.MaxInt16)

		// Pack and unpack to test value
		src = DPT_8001(value)
		if int16(src) != value {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_8001! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int16(dst) != value {
			t.Errorf("Value \"%s\" after pack/unpack different from Original value. Was \"%v\"", dst, value)
		}
	}
}

// Test DPT 8.002 with values within range
func TestDPT_8002(t *testing.T) {
	var buf []byte
	var src, dst DPT_8002

	for i := 1; i <= 10; i++ {
		value := int16(rand.Uint32() % math.MaxInt16)

		// Pack and unpack to test value
		src = DPT_8002(value)
		if int16(src) != value {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_8002! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int16(dst) != value {
			t.Errorf("Value \"%s\" after pack/unpack different from Original value. Was \"%v\"", dst, value)
		}
	}
}

// Test DPT 8.003 with values within range
func TestDPT_8003(t *testing.T) {
	var buf []byte
	var src, dst DPT_8003

	for i := 1; i <= 10; i++ {

		value := float32(int16(rand.Uint32()%math.MaxInt16)) / 100

		src = DPT_8003(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_8003! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if fmt.Sprintf("%.2f", dst) != fmt.Sprintf("%.2f", value) {
			t.Errorf("Value \"%f\" after pack/unpack different from Original value. Was \"%f\"", dst, value)
		}
	}
}

// Test DPT 8.004 with values within range
func TestDPT_8004(t *testing.T) {
	var buf []byte
	var src, dst DPT_8004

	for i := 1; i <= 10; i++ {

		value := float32(int16(rand.Uint32()%math.MaxInt16)) / 10

		src = DPT_8004(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_8004! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if fmt.Sprintf("%.1f", dst) != fmt.Sprintf("%.1f", value) {
			t.Errorf("Value \"%f\" after pack/unpack different from Original value. Was \"%f\"", dst, value)
		}
	}
}

// Test DPT 8.005 with values within range
func TestDPT_8005(t *testing.T) {
	var buf []byte
	var src, dst DPT_8005

	for i := 1; i <= 10; i++ {
		value := int16(rand.Uint32() % math.MaxInt16)

		// Pack and unpack to test value
		src = DPT_8005(value)
		if int16(src) != value {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_8005! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int16(dst) != value {
			t.Errorf("Value \"%s\" after pack/unpack different from Original value. Was \"%v\"", dst, value)
		}
	}
}

// Test DPT 8.006 with values within range
func TestDPT_8006(t *testing.T) {
	var buf []byte
	var src, dst DPT_8006

	for i := 1; i <= 10; i++ {
		value := int16(rand.Uint32() % math.MaxInt16)

		// Pack and unpack to test value
		src = DPT_8006(value)
		if int16(src) != value {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_8006! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int16(dst) != value {
			t.Errorf("Value \"%s\" after pack/unpack different from Original value. Was \"%v\"", dst, value)
		}
	}
}

// Test DPT 8.007 with values within range
func TestDPT_8007(t *testing.T) {
	var buf []byte
	var src, dst DPT_8007

	for i := 1; i <= 10; i++ {
		value := int16(rand.Uint32() % math.MaxInt16)

		// Pack and unpack to test value
		src = DPT_8007(value)
		if int16(src) != value {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_8007! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int16(dst) != value {
			t.Errorf("Value \"%s\" after pack/unpack different from Original value. Was \"%v\"", dst, value)
		}
	}
}

// Test DPT 8.010 with values within range
func TestDPT_8010(t *testing.T) {
	var buf []byte
	var src, dst DPT_8010

	for i := 1; i <= 10; i++ {

		iValue := int16(rand.Uint32() % math.MaxInt16)
		value := float32(iValue) / 100

		src = DPT_8010(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_8010! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if fmt.Sprintf("%.2f", dst) != fmt.Sprintf("%.2f", value) {
			t.Errorf("Value \"%f\" after pack/unpack different from Original value. Was \"%f\"", dst, value)
		}
		fmt.Printf("Original value: %d, dpt.String(): %s\n", iValue, dst.String())
	}
}

// Test DPT 8.011 (°) with values within range
func TestDPT_8011(t *testing.T) {
	var buf []byte
	var src, dst DPT_8011

	for i := 1; i <= 10; i++ {
		value := int16(rand.Uint32() % math.MaxInt16)

		// Pack and unpack to test value
		src = DPT_8011(value)
		if int16(src) != value {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_8011! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if int16(dst) != value {
			t.Errorf("Value \"%s\" after pack/unpack different from Original value. Was \"%v\"", dst, value)
		}
	}
}
