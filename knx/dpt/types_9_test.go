// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"testing"

	"math"
	"math/rand"
)

// Test DPT 9.001 (Temperature) with values within range
func TestDPT_9001(t *testing.T) {
	var buf []byte
	var src, dst DPT_9001

	for i := 1; i <= 10; i++ {
		value := rand.Float32()

		// Scale the random number to the given range
		value *= 670760 - -273
		value += -273

		// Calculate the quantization error we expect
		Q := get_float_quantization_error(value, 0.01, 2047)

		// Pack and unpack to test value
		src = DPT_9001(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_9001! Has value \"%s\".", value, src)
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

// Test DPT 9.004 (Illumination) with values within range
func TestDPT_9004(t *testing.T) {
	var buf []byte
	var src, dst DPT_9004

	for i := 1; i <= 10; i++ {
		value := rand.Float32()

		// Scale the random number to the given range
		value *= 670760

		// Calculate the quantization error we expect
		Q := get_float_quantization_error(value, 0.01, 2047)

		// Pack and unpack to test value
		src = DPT_9004(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_9004! Has value \"%s\".", value, src)
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

// Test DPT 9.005 (Wind Speed) with values within range
func TestDPT_9005(t *testing.T) {
	var buf []byte
	var src, dst DPT_9005

	for i := 1; i <= 10; i++ {
		value := rand.Float32()

		// Scale the random number to the given range
		value *= 670760

		// Calculate the quantization error we expect
		Q := get_float_quantization_error(value, 0.01, 2047)

		// Pack and unpack to test value
		src = DPT_9005(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_9005! Has value \"%s\".", value, src)
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

// Test DPT 9.007 (Humidity) with values within range
func TestDPT_9007(t *testing.T) {
	var buf []byte
	var src, dst DPT_9007

	for i := 1; i <= 10; i++ {
		value := rand.Float32()

		// Scale the random number to the given range
		value *= 670760

		// Calculate the quantization error we expect
		Q := get_float_quantization_error(value, 0.01, 2047)

		// Pack and unpack to test value
		src = DPT_9007(value)
		if abs(float32(src)-value) > epsilon {
			t.Errorf("Assignment of value \"%v\" failed for source of type DPT_9007! Has value \"%s\".", value, src)
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
