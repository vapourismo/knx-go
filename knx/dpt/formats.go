// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"errors"

)

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


// int32 to 2-byte array of uint8, leading zero
func packI32(i int32) []byte {
	buffer := []byte{0, 0}
	buffer[1] = uint8(i)
	return buffer
}

// int8 to int32
func unpackI32(data []byte, i *int32) error {
	//fmt.Printf("unpackI32 packed: %+v\n", data)
	if len(data) != 2 {
		return ErrInvalidLength
	}
	//fmt.Printf("unpackI32 data[1]: %+v\n", data[1])
	*i = int32(data[1])
	//fmt.Printf("unpackI32 unpacked: %d\n", *i)
	return nil
}
