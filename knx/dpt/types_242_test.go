package dpt

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDPT_242600(t *testing.T) {
	var buf []byte
	var dst DPT_242600
	sources := []DPT_242600{
		{X: 0, Y: 0, YBrightness: 1, ColorValid: true, BrightnessValid: true},
		{X: 1, Y: 1, YBrightness: 255, ColorValid: true, BrightnessValid: true},
		{X: 0.5, Y: 0.5, YBrightness: 127, ColorValid: true, BrightnessValid: true},
		{X: 0.1, Y: 0.9, YBrightness: 127, ColorValid: true, BrightnessValid: true},
	}

	for _, src := range sources {
		buf = src.Pack()
		fmt.Printf("==> %v\n", buf)
		err := dst.Unpack(buf)

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(src, dst) {
			fmt.Printf("%+v\n", src)
			fmt.Printf("%+v\n", dst)
			t.Errorf("Value \"%s\" after pack/unpack for DPT_242600 differs. Original value was \"%v\"!", dst, src)
		}
	}
}
