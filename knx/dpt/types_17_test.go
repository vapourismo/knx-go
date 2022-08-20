package dpt

import "testing"

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
