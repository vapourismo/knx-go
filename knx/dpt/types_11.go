// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"fmt"
	"time"
)

// DPT_11001 represents DPT 11.001 / Date p 34.
// Valid years are limited to 1990 - 2089 for the Year field.
type DPT_11001 struct {
	Year  uint16
	Month uint8
	Day   uint8
}

func (d DPT_11001) Pack() []byte {
	var buf = []byte{0, 0, 0, 0}

	if d.Year >= 1990 && d.Year <= 2089 && d.IsValid() {
		buf[1] = d.Day & 0x1F
		buf[2] = d.Month & 0xF

		if d.Year < 2000 {
			buf[3] = uint8(d.Year - 1900)
		} else {
			buf[3] = uint8(d.Year - 2000)
		}
	}

	buf[3] &= 0x7F

	return buf
}

func (d *DPT_11001) Unpack(data []byte) error {
	if len(data) != 4 {
		return ErrInvalidLength
	}

	d.Day = uint8(data[1] & 0x1F)
	d.Month = uint8(data[2] & 0xF)
	d.Year = uint16(data[3] & 0x7F)

	if d.Year > 99 {
		return fmt.Errorf("payload is out of range")
	}

	if d.Year == 0 && d.Month == 0 && d.Day == 0 {
		d.Year = 90
		d.Month = 1
		d.Day = 1
	}

	if d.Year >= 90 {
		d.Year += 1900
	} else {
		d.Year += 2000
	}

	if !d.IsValid() {
		return fmt.Errorf("payload is out of range")
	}

	return nil
}

func (d DPT_11001) Unit() string {
	return ""
}

func (d DPT_11001) IsValid() bool {
	tm := time.Date(int(d.Year), time.Month(d.Month), int(d.Day), 0, 0, 0, 0, time.UTC)
	if tm.Year() < 1990 || tm.Year() > 2089 {
		return false
	} else {
		return (int(d.Year) == tm.Year() && d.Month == uint8(tm.Month()) && int(d.Day) == tm.Day())
	}
}

func (d DPT_11001) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}
