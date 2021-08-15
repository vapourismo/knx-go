// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"testing"

	"math"
	"math/rand"
)

// Test DPT 7.001 (pulses) with values within range
func TestDPT_7001(t *testing.T) {
	var buf []byte
	var src, dst DPT_7001

	for i := 1; i <= 10; i++ {
		value := uint16(rand.Uint32() % math.MaxUint16)

		// Pack and unpack to test value
		src = DPT_7001(value)
		if uint16(src) != value {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_7001! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if uint16(dst) != value {
			t.Errorf("Value \"%s\" after pack/unpack different from Original value. Was \"%v\"", dst, value)
		}
	}
}

// Test DPT 7.002 (ms) with values within range
func TestDPT_7002(t *testing.T) {
	var buf []byte
	var src, dst DPT_7002

	for i := 1; i <= 10; i++ {
		value := uint16(rand.Uint32() % math.MaxUint16)

		// Pack and unpack to test value
		src = DPT_7002(value)
		if uint16(src) != value {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_7002! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if uint16(dst) != value {
			t.Errorf("Value \"%s\" after pack/unpack different from Original value. Was \"%v\"", dst, value)
		}
	}
}

// Test DPT 7.003 (s) with values within range
func TestDPT_7003(t *testing.T) {
	var buf []byte
	var src, dst DPT_7003

	for i := 1; i <= 10; i++ {
		value := uint16(rand.Uint32() % math.MaxUint16)

		// Pack and unpack to test value
		src = DPT_7003(value)
		if uint16(src) != value {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_7003! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if uint16(dst) != value {
			t.Errorf("Value \"%s\" after pack/unpack different from Original value. Was \"%v\"", dst, value)
		}
	}
}

// Test DPT 7.004 (s) with values within range
func TestDPT_7004(t *testing.T) {
	var buf []byte
	var src, dst DPT_7004

	for i := 1; i <= 10; i++ {
		value := uint16(rand.Uint32() % math.MaxUint16)

		// Pack and unpack to test value
		src = DPT_7004(value)
		if uint16(src) != value {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_7004! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if uint16(dst) != value {
			t.Errorf("Value \"%s\" after pack/unpack different from Original value. Was \"%v\"", dst, value)
		}
	}
}

// Test DPT 7.005 (s) with values within range
func TestDPT_7005(t *testing.T) {
	var buf []byte
	var src, dst DPT_7005

	for i := 1; i <= 10; i++ {
		value := uint16(rand.Uint32() % math.MaxUint16)

		// Pack and unpack to test value
		src = DPT_7005(value)
		if uint16(src) != value {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_7005! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if uint16(dst) != value {
			t.Errorf("Value \"%s\" after pack/unpack different from Original value. Was \"%v\"", dst, value)
		}
	}
}

// Test DPT 7.006 (m) with values within range
func TestDPT_7006(t *testing.T) {
	var buf []byte
	var src, dst DPT_7006

	for i := 1; i <= 10; i++ {
		value := uint16(rand.Uint32() % math.MaxUint16)

		// Pack and unpack to test value
		src = DPT_7006(value)
		if uint16(src) != value {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_7006! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if uint16(dst) != value {
			t.Errorf("Value \"%s\" after pack/unpack different from Original value. Was \"%v\"", dst, value)
		}
	}
}

// Test DPT 7.007 (h) with values within range
func TestDPT_7007(t *testing.T) {
	var buf []byte
	var src, dst DPT_7007

	for i := 1; i <= 10; i++ {
		value := uint16(rand.Uint32() % math.MaxUint16)

		// Pack and unpack to test value
		src = DPT_7007(value)
		if uint16(src) != value {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_7007! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if uint16(dst) != value {
			t.Errorf("Value \"%s\" after pack/unpack different from Original value. Was \"%v\"", dst, value)
		}
	}
}

// Test DPT 7.010 with values within range
func TestDPT_7010(t *testing.T) {
	var buf []byte
	var src, dst DPT_7010

	for i := 1; i <= 10; i++ {
		value := uint16(rand.Uint32() % math.MaxUint16)

		// Pack and unpack to test value
		src = DPT_7010(value)
		if uint16(src) != value {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_7010! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if uint16(dst) != value {
			t.Errorf("Value \"%s\" after pack/unpack different from Original value. Was \"%v\"", dst, value)
		}
	}
}

// Test DPT 7.011 (mm) with values within range
func TestDPT_7011(t *testing.T) {
	var buf []byte
	var src, dst DPT_7011

	for i := 1; i <= 10; i++ {
		value := uint16(rand.Uint32() % math.MaxUint16)

		// Pack and unpack to test value
		src = DPT_7011(value)
		if uint16(src) != value {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_7011! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if uint16(dst) != value {
			t.Errorf("Value \"%s\" after pack/unpack different from Original value. Was \"%v\"", dst, value)
		}
	}
}

// Test DPT 7.012 (mA) with values within range
func TestDPT_7012(t *testing.T) {
	var buf []byte
	var src, dst DPT_7012

	for i := 1; i <= 10; i++ {
		value := uint16(rand.Uint32() % math.MaxUint16)

		// Pack and unpack to test value
		src = DPT_7012(value)
		if uint16(src) != value {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_7012! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if uint16(dst) != value {
			t.Errorf("Value \"%s\" after pack/unpack different from Original value. Was \"%v\"", dst, value)
		}
	}
}

// Test DPT 7.013 (lux) with values within range
func TestDPT_7013(t *testing.T) {
	var buf []byte
	var src, dst DPT_7013

	for i := 1; i <= 10; i++ {
		value := uint16(rand.Uint32() % math.MaxUint16)

		// Pack and unpack to test value
		src = DPT_7013(value)
		if uint16(src) != value {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_7013! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		dst.Unpack(buf)
		if uint16(dst) != value {
			t.Errorf("Value \"%s\" after pack/unpack different from Original value. Was \"%v\"", dst, value)
		}
	}
}
