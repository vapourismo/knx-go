// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"fmt"
	"testing"
	"time"
)

// Test DPT 10.001 (Time of Day)
func TestDPT_10001(t *testing.T) {
	var src DPT_10001
	var dst DPT_10001

	// Test all weekdays, all hours, all minutes and seconds.
	for w := 0; w <= 7; w++ {
		src.Weekday = uint8(w)
		for m := 0; m < 24; m++ {
			src.Hour = uint8(m)
			for s := 0; s < 60; s++ {
				src.Minutes = uint8(s)
				src.Seconds = uint8(s)
				buf := src.Pack()
				dst.Unpack(buf)
				if dst.String() != src.String() {
					t.Errorf("Value \"%s\" is not a time of day! Original value was \"%s\".", dst, src)
				}
			}
		}
	}

	src.Minutes = 0
	weekday := []string{"", time.Monday.String(), time.Tuesday.String(), time.Wednesday.String(),
		time.Thursday.String(), time.Friday.String(), time.Saturday.String(), time.Sunday.String()}

	// Test print with weekday
	for w := 1; w <= 7; w++ {
		src.Weekday = uint8(w)
		buf := src.Pack()
		dst.Unpack(buf)
		r := fmt.Sprintf("%s %02d:%02d:%02d", weekday[src.Weekday], src.Hour, src.Minutes, src.Seconds)
		if dst.String() != r {
			t.Errorf("Value \"%s\" is not a time of day! Original value was \"%s\".", dst, r)
		}
	}
}
