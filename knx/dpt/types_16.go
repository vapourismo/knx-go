// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"unicode"
)

// DPT_16000 represents DPT 16.000 / String ASCII.
// The string should be ASCII and contain at most 14 chars.
// A string longer than 14 chars will be silently truncated.
// Non-ASCII chars will be replaced with a space = 0x20.
type DPT_16000 string

func (d DPT_16000) Pack() []byte {
	var buf = make([]byte, 15)

	r := []rune(d)

	for i := 0; i < len(r) && i < 14; i++ {
		if r[i] > unicode.MaxASCII {
			buf[i+1] = 0x20
		} else {
			buf[i+1] = byte(r[i])
		}
	}

	return buf
}

func (d *DPT_16000) Unpack(data []byte) error {
	if len(data) != 15 {
		return ErrInvalidLength
	}

	var buf = []rune{}

	for i := 1; i < len(data) && data[i] != 0x00; i++ {
		buf = append(buf, rune(data[i]&unicode.MaxASCII))
	}

	*d = DPT_16000(buf)

	return nil
}

func (d DPT_16000) Unit() string {
	return ""
}

func (d DPT_16000) IsValid() bool {
	for _, c := range d {
		if c > unicode.MaxASCII {
			return false
		}
	}

	return len(d) <= 14
}

func (d DPT_16000) String() string {
	return string(d)
}

// DPT_16001 represents DPT 16.001 / String 8859-1.
// The string must be ISO-8859-1 and contain at most 14 chars.
// A string longer than 14 chars will be silently truncated.
// Non-ISO-8859-1 chars will be replaced with a space = 0x20.
type DPT_16001 string

func (d DPT_16001) Pack() []byte {
	buf := make([]byte, 15)

	r := []rune(d)

	for i := 0; i < len(r) && i < 14; i++ {
		if r[i] > unicode.MaxLatin1 {
			buf[i+1] = 0x20
		} else {
			buf[i+1] = byte(r[i])
		}
	}

	return buf
}

func (d *DPT_16001) Unpack(data []byte) error {
	if len(data) != 15 {
		return ErrInvalidLength
	}

	var buf = []rune{}

	for i := 1; i < len(data) && data[i] != 0x00; i++ {
		buf = append(buf, rune(data[i]))
	}

	*d = DPT_16001(buf)

	return nil
}

func (d DPT_16001) Unit() string {
	return ""
}

func (d DPT_16001) IsValid() bool {
	for _, c := range d {
		if c > unicode.MaxLatin1 {
			return false
		}
	}

	return len(d) <= 14
}

func (d DPT_16001) String() string {
	return string(d)
}
