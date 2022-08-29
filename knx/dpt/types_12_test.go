package dpt

import (
	"math/rand"
	"testing"
)

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
