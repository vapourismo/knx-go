package dpt

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDPT_251600(t *testing.T) {
	var buf []byte
	var dst DPT_251600
	sources := []DPT_251600{
		{Red: 255, Green: 96, Blue: 0, White: 18, RedValid: true, GreenValid: true, BlueValid: true, WhiteValid: true},
		{Red: 255, Green: 96, Blue: 0, White: 18, RedValid: false, GreenValid: false, BlueValid: false, WhiteValid: false},

		{Red: 255, Green: 96, Blue: 0, White: 18, RedValid: false, GreenValid: true, BlueValid: true, WhiteValid: true},
		{Red: 255, Green: 96, Blue: 0, White: 18, RedValid: true, GreenValid: false, BlueValid: true, WhiteValid: true},
		{Red: 255, Green: 96, Blue: 0, White: 18, RedValid: true, GreenValid: true, BlueValid: false, WhiteValid: true},
		{Red: 255, Green: 96, Blue: 0, White: 18, RedValid: true, GreenValid: true, BlueValid: true, WhiteValid: false},
	}

	for _, src := range sources {
		buf = src.Pack()
		_ = dst.Unpack(buf)

		if !reflect.DeepEqual(src, dst) {
			fmt.Printf("%+v\n", src)
			fmt.Printf("%+v\n", dst)
			t.Errorf("Value \"%s\" after pack/unpack for DPT_251600 differs. Original value was \"%v\"!", dst, src)
		}
	}
}
