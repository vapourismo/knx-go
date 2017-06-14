// Copyright 2017 Ole Krüger.

package dpt

import "errors"

// ErrInvalidLength is returned when the application data has unexpected length.
var ErrInvalidLength = errors.New("Given application data has invalid length")

func packB1(b bool) []byte {
	if b {
		return []byte{1}
	}

	return []byte{0}
}

func unpackB1(data []byte, b *bool) error {
	if len(data) != 1 {
		return ErrInvalidLength
	}

	*b = data[0]&1 == 1

	return nil
}

func packF16(f float32) []byte {
	buffer := []byte{0, 0, 0}

	if f > 670760.96 {
		f = 670760.96
	} else if f < -671088.64 {
		f = -671088.64
	}

	signedMantissa := int(f * 100)
	exp := 0

	for signedMantissa > 2047 || signedMantissa < -2048 {
		signedMantissa /= 2
		exp++
	}

	buffer[1] |= uint8(exp&15) << 3

	if signedMantissa < 0 {
		signedMantissa += 2048
		buffer[1] |= 1 << 7
	}

	mantissa := uint(signedMantissa)

	buffer[1] |= uint8(mantissa>>8) & 7
	buffer[2] |= uint8(mantissa)

	return buffer
}

func unpackF16(data []byte, f *float32) error {
	if len(data) != 3 {
		return ErrInvalidLength
	}

	m := int(data[1]&7)<<8 | int(data[2])
	if data[1]&128 == 128 {
		m -= 2048
	}

	e := (data[1] >> 3) & 15

	*f = 0.01 * float32(m) * float32(uint(1)<<e)
	return nil
}
