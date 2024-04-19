package dpt

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDPT_232600(t *testing.T) {
	var buf []byte
	var dst DPT_232600
	sources := []DPT_232600{
		{Red: 255, Green: 96, Blue: 0},
		{Red: 0, Green: 128, Blue: 128},
		{Red: 0, Green: 128, Blue: 255},
		{Red: 0, Green: 0, Blue: 255},
	}

	for _, src := range sources {
		buf = src.Pack()
		_ = dst.Unpack(buf)

		if !reflect.DeepEqual(src, dst) {
			fmt.Printf("%+v\n", src)
			fmt.Printf("%+v\n", dst)
			t.Errorf("Value [%s] after pack/unpack for DPT_232600 differs. Original value was [%v].", dst, src)
		}
	}
}
