package dpt

import (
	"testing"

	"math"
	"math/rand"
)

func getRandValue() (float32, float32) {
	value := rand.Float32()

	// Scale the random number to the given range
	value *= 670760

	// Calculate the quantization error we expect
	Q := get_float_quantization_error(value, 0.01, 2047)
	return value, Q
}

// Test DPT 14.000 (Acceleration) with values within range
func TestDPT_14000(t *testing.T) {
	var buf []byte
	var src, dst DPT_14000

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14000(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14000! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.001 (Acceleration Angular) with values within range
func TestDPT_14001(t *testing.T) {
	var buf []byte
	var src, dst DPT_14001

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14001(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14001! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.002 (ActivationEnergy) with values within range
func TestDPT_14002(t *testing.T) {
	var buf []byte
	var src, dst DPT_14002

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14002(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14002! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.003 (Activity) with values within range
func TestDPT_14003(t *testing.T) {
	var buf []byte
	var src, dst DPT_14003

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14003(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14003! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.004 (Mol) with values within range
func TestDPT_14004(t *testing.T) {
	var buf []byte
	var src, dst DPT_14004

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14004(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14004! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.005 (Amplitude) with values within range
func TestDPT_14005(t *testing.T) {
	var buf []byte
	var src, dst DPT_14005

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14005(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14005! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.006 (AngleRad) with values within range
func TestDPT_14006(t *testing.T) {
	var buf []byte
	var src, dst DPT_14006

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14006(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14006! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.007 (AngleDeg) with values within range
func TestDPT_14007(t *testing.T) {
	var buf []byte
	var src, dst DPT_14007

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14007(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14007! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.008 (Angular Momentum) with values within range
func TestDPT_14008(t *testing.T) {
	var buf []byte
	var src, dst DPT_14008

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14008(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14008! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.009 (Angular Velocity) with values within range
func TestDPT_14009(t *testing.T) {
	var buf []byte
	var src, dst DPT_14009

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14009(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14009! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.010 (Area) with values within range
func TestDPT_14010(t *testing.T) {
	var buf []byte
	var src, dst DPT_14010

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14010(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14010! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.011 (Capacitance) with values within range
func TestDPT_14011(t *testing.T) {
	var buf []byte
	var src, dst DPT_14011

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14011(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14011! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.012 (Charge DensitySurface) with values within range
func TestDPT_14012(t *testing.T) {
	var buf []byte
	var src, dst DPT_14012

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14012(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14012! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.013 (Charge DensityVolume) with values within range
func TestDPT_14013(t *testing.T) {
	var buf []byte
	var src, dst DPT_14013

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14013(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14013! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.014 (Compressibility) with values within range
func TestDPT_14014(t *testing.T) {
	var buf []byte
	var src, dst DPT_14014

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14014(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14014! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.015 (Conductance) with values within range
func TestDPT_14015(t *testing.T) {
	var buf []byte
	var src, dst DPT_14015

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14015(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14015! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.016 (Electrical Conductivity) with values within range
func TestDPT_14016(t *testing.T) {
	var buf []byte
	var src, dst DPT_14016

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14016(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14016! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.017 (Density) with values within range
func TestDPT_14017(t *testing.T) {
	var buf []byte
	var src, dst DPT_14017

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14017(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14017! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.018 (Electric Charge) with values within range
func TestDPT_14018(t *testing.T) {
	var buf []byte
	var src, dst DPT_14018

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14018(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14018! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.019 (Electric Current) with values within range
func TestDPT_14019(t *testing.T) {
	var buf []byte
	var src, dst DPT_14019

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14019(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14019! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.020 (Electric CurrentDensity) with values within range
func TestDPT_14020(t *testing.T) {
	var buf []byte
	var src, dst DPT_14020

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14020(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14020! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.021 (Electric DipoleMoment) with values within range
func TestDPT_14021(t *testing.T) {
	var buf []byte
	var src, dst DPT_14021

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14021(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14021! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.022 (Electric Displacement) with values within range
func TestDPT_14022(t *testing.T) {
	var buf []byte
	var src, dst DPT_14022

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14022(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14022! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.023 (Electric FieldStrength) with values within range
func TestDPT_14023(t *testing.T) {
	var buf []byte
	var src, dst DPT_14023

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14023(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14023! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.024 (Electric Flux) with values within range
func TestDPT_14024(t *testing.T) {
	var buf []byte
	var src, dst DPT_14024

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14024(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14024! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.025 (Electric FluxDensity) with values within range
func TestDPT_14025(t *testing.T) {
	var buf []byte
	var src, dst DPT_14025

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14025(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14025! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.026 (Electric Polarization) with values within range
func TestDPT_14026(t *testing.T) {
	var buf []byte
	var src, dst DPT_14026

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14026(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14026! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.027 (Electric Potential) with values within range
func TestDPT_14027(t *testing.T) {
	var buf []byte
	var src, dst DPT_14027

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14027(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14027! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.028 (Electric PotentialDifference) with values within range
func TestDPT_14028(t *testing.T) {
	var buf []byte
	var src, dst DPT_14028

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14028(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14028! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.029 (ElectromagneticMoment) with values within range
func TestDPT_14029(t *testing.T) {
	var buf []byte
	var src, dst DPT_14029

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14029(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14029! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.030 (Electromotive_Force) with values within range
func TestDPT_14030(t *testing.T) {
	var buf []byte
	var src, dst DPT_14030

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14030(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14030! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.031 (Energy) with values within range
func TestDPT_14031(t *testing.T) {
	var buf []byte
	var src, dst DPT_14031

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14031(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14031! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.032 (Force) with values within range
func TestDPT_14032(t *testing.T) {
	var buf []byte
	var src, dst DPT_14032

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14032(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14032! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.033 (Frequency) with values within range
func TestDPT_14033(t *testing.T) {
	var buf []byte
	var src, dst DPT_14033

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14033(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14033! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.034 (Angular Frequency) with values within range
func TestDPT_14034(t *testing.T) {
	var buf []byte
	var src, dst DPT_14034

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14034(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14034! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.035 (Heat Capacity) with values within range
func TestDPT_14035(t *testing.T) {
	var buf []byte
	var src, dst DPT_14035

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14035(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14035! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.036 (Heat Flow Rate) with values within range
func TestDPT_14036(t *testing.T) {
	var buf []byte
	var src, dst DPT_14036

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14036(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14036! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.037 (Heat Quantity) with values within range
func TestDPT_14037(t *testing.T) {
	var buf []byte
	var src, dst DPT_14037

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14037(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14037! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.038 (Impedance) with values within range
func TestDPT_14038(t *testing.T) {
	var buf []byte
	var src, dst DPT_14038

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14038(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14038! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.039 (Length) with values within range
func TestDPT_14039(t *testing.T) {
	var buf []byte
	var src, dst DPT_14039

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14039(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14039! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.040 (Light_Quantity) with values within range
func TestDPT_14040(t *testing.T) {
	var buf []byte
	var src, dst DPT_14040

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14040(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14040! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.041 (Luminance) with values within range
func TestDPT_14041(t *testing.T) {
	var buf []byte
	var src, dst DPT_14041

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14041(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14041! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.042 (Luminous Flux) with values within range
func TestDPT_14042(t *testing.T) {
	var buf []byte
	var src, dst DPT_14042

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14042(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14042! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.043 (Luminous Intensity) with values within range
func TestDPT_14043(t *testing.T) {
	var buf []byte
	var src, dst DPT_14043

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14043(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14043! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.044 (Magnetic FieldStrength) with values within range
func TestDPT_14044(t *testing.T) {
	var buf []byte
	var src, dst DPT_14044

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14044(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14044! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.045 (Magnetic Flux) with values within range
func TestDPT_14045(t *testing.T) {
	var buf []byte
	var src, dst DPT_14045

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14045(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14045! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.046 (Magnetic FluxDensity) with values within range
func TestDPT_14046(t *testing.T) {
	var buf []byte
	var src, dst DPT_14046

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14046(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14046! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.047 (Magnetic Moment) with values within range
func TestDPT_14047(t *testing.T) {
	var buf []byte
	var src, dst DPT_14047

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14047(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14047! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.048 (Magnetic Polarization) with values within range
func TestDPT_14048(t *testing.T) {
	var buf []byte
	var src, dst DPT_14048

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14048(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14048! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.049 (Magnetization) with values within range
func TestDPT_14049(t *testing.T) {
	var buf []byte
	var src, dst DPT_14049

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14049(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14049! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.050 (MagnetomotiveForce) with values within range
func TestDPT_14050(t *testing.T) {
	var buf []byte
	var src, dst DPT_14050

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14050(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14050! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.051 (Mass) with values within range
func TestDPT_14051(t *testing.T) {
	var buf []byte
	var src, dst DPT_14051

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14051(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14051! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.052 (MassFlux) with values within range
func TestDPT_14052(t *testing.T) {
	var buf []byte
	var src, dst DPT_14052

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14052(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14052! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.053 (Momentum) with values within range
func TestDPT_14053(t *testing.T) {
	var buf []byte
	var src, dst DPT_14053

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14053(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14053! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.054 (Phase Angle, Radiant) with values within range
func TestDPT_14054(t *testing.T) {
	var buf []byte
	var src, dst DPT_14054

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14054(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14054! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.055 (Phase Angle, Degree) with values within range
func TestDPT_14055(t *testing.T) {
	var buf []byte
	var src, dst DPT_14055

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14055(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14055! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.056 (Power) with values within range
func TestDPT_14056(t *testing.T) {
	var buf []byte
	var src, dst DPT_14056

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14056(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14056! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.057 (Power Factor) with values within range
func TestDPT_14057(t *testing.T) {
	var buf []byte
	var src, dst DPT_14057

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14057(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14057! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.058 (Pressure) with values within range
func TestDPT_14058(t *testing.T) {
	var buf []byte
	var src, dst DPT_14058

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14058(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14058! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.059 (Reactance) with values within range
func TestDPT_14059(t *testing.T) {
	var buf []byte
	var src, dst DPT_14059

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14059(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14059! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.060 (Resistance) with values within range
func TestDPT_14060(t *testing.T) {
	var buf []byte
	var src, dst DPT_14060

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14060(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14060! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.061 (Resistivity) with values within range
func TestDPT_14061(t *testing.T) {
	var buf []byte
	var src, dst DPT_14061

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14061(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14061! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.062 (SelfInductance) with values within range
func TestDPT_14062(t *testing.T) {
	var buf []byte
	var src, dst DPT_14062

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14062(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14062! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.063 (SolidAngle) with values within range
func TestDPT_14063(t *testing.T) {
	var buf []byte
	var src, dst DPT_14063

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14063(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14063! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.064 (Sound Intensity) with values within range
func TestDPT_14064(t *testing.T) {
	var buf []byte
	var src, dst DPT_14064

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14064(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14064! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.065 (Speed) with values within range
func TestDPT_14065(t *testing.T) {
	var buf []byte
	var src, dst DPT_14065

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14065(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14065! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.066 (Stress) with values within range
func TestDPT_14066(t *testing.T) {
	var buf []byte
	var src, dst DPT_14066

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14066(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14066! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.067 (Surface Tension) with values within range
func TestDPT_14067(t *testing.T) {
	var buf []byte
	var src, dst DPT_14067

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14067(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14067! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.068 (Common Temperature) with values within range
func TestDPT_14068(t *testing.T) {
	var buf []byte
	var src, dst DPT_14068

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14068(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14068! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.069 (Absolute Temperature) with values within range
func TestDPT_14069(t *testing.T) {
	var buf []byte
	var src, dst DPT_14069

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14069(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14069! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.070 (Temperature Difference) with values within range
func TestDPT_14070(t *testing.T) {
	var buf []byte
	var src, dst DPT_14070

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14070(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14070! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.071 (Thermal Capacity) with values within range
func TestDPT_14071(t *testing.T) {
	var buf []byte
	var src, dst DPT_14071

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14071(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14071! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.072 (Thermal Conductivity) with values within range
func TestDPT_14072(t *testing.T) {
	var buf []byte
	var src, dst DPT_14072

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14072(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14072! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.073 (Thermoelectric Power) with values within range
func TestDPT_14073(t *testing.T) {
	var buf []byte
	var src, dst DPT_14073

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14073(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14073! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.074 (Time) with values within range
func TestDPT_14074(t *testing.T) {
	var buf []byte
	var src, dst DPT_14074

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14074(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14074! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.075 (Torque) with values within range
func TestDPT_14075(t *testing.T) {
	var buf []byte
	var src, dst DPT_14075

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14075(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14075! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.076 (Volume) with values within range
func TestDPT_14076(t *testing.T) {
	var buf []byte
	var src, dst DPT_14076

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14076(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14076! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.077 (Volume Flux) with values within range
func TestDPT_14077(t *testing.T) {
	var buf []byte
	var src, dst DPT_14077

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14077(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14077! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.078 (Weight) with values within range
func TestDPT_14078(t *testing.T) {
	var buf []byte
	var src, dst DPT_14078

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14078(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14078! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}

// Test DPT 14.079 (Work) with values within range
func TestDPT_14079(t *testing.T) {
	var buf []byte
	var src, dst DPT_14079

	for i := 1; i <= 10; i++ {
		value, Q := getRandValue()

		// Pack and unpack to test value
		src = DPT_14079(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_14079! Has value \"%s\".", value, src)
		}
		buf = src.Pack()
		_ = dst.Unpack(buf)
		if math.IsNaN(float64(dst)) {
			t.Errorf("Value \"%s\" is not a valid number! Original value was \"%v\".", dst, value)
		}
		if abs(float32(dst)-value) > (Q + epsilon) {
			t.Errorf("Value \"%s\" after pack/unpack above quantization noise! Original value was \"%v\", noise is \"%f\"", dst, value, Q)
		}
	}
}
