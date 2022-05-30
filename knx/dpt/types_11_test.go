// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"testing"
	"time"
)

// Test DPT 11.001 (Date)
func TestDPT_11001(t *testing.T) {
	var src DPT_11001
	var dst DPT_11001

	src.Year = 2001
	src.Month = 2
	src.Day = 29
	buf := src.Pack()
	_ = dst.Unpack(buf)

	if src.IsValid() {
		t.Errorf("Date 2001-02-29 should not be a valid date")
		t.Errorf("%s -> %s should be 2001-02-29\n", src.String(), dst.String())
	}

	src.Month = 0
	src.Day = 1

	for y := 1989; y <= 2090; y++ {
		// year < 1990 or > 2089 will result in date being 1990-01-01
		src.Year = uint16(y)
		for m := 1; m <= 12; m++ {
			src.Month = uint8(m)
			for d := 1; d <= 31; d++ {
				src.Day = uint8(d)
				if src.IsValid() {
					buf := src.Pack()
					dst.Unpack(buf)
					tm, _ := time.Parse("2006-01-02", dst.String())
					r := tm.Format("2006-01-02")
					if src.String() != dst.String() || dst.String() != r {
						t.Errorf("Value \"%s\" is not a time of day! Original value was \"%s\".", dst, src)
					}
				}
			}
		}

	}

}
