package dpt

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestDPT_5005(t *testing.T) {
	knxValue := []byte{0, 42}
	dptValue := DPT_5005(42)

	var tmpDPT DPT_5005
	assert.NoError(t, tmpDPT.Unpack(knxValue))
	assert.Equal(t, dptValue, tmpDPT)

	assert.Equal(t, knxValue, dptValue.Pack())

	assert.Equal(t, "42", dptValue.String())
}
