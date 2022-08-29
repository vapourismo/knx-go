package dpt

import "testing"

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
