// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package util

import (
	"bytes"
	"fmt"
	"io"
)

var (
	stringDecoder = stringCharmap.NewDecoder()
)

// Unpackable is implemented by types that can be initialized by reading from a byte slice.
type Unpackable interface {
	Unpack(data []byte) (uint, error)
}

func unpackUInt16(data []byte, output *uint16) (uint, error) {
	if len(data) < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	*output = uint16(data[1]) | uint16(data[0])<<8

	return 2, nil
}

func unpackUInt32(data []byte, output *uint32) (uint, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}

	*output = uint32(data[3]) | uint32(data[2])<<8 | uint32(data[1])<<16 | uint32(data[0])<<24

	return 4, nil
}

func unpackUInt64(data []byte, output *uint64) (uint, error) {
	if len(data) < 8 {
		return 0, io.ErrUnexpectedEOF
	}

	*output = uint64(data[7]) | uint64(data[6])<<8 | uint64(data[5])<<16 | uint64(data[4])<<24 |
		uint64(data[3])<<32 | uint64(data[2])<<40 | uint64(data[1])<<48 | uint64(data[0])<<56

	return 8, nil
}

// Unpack the byte slice to the given value.
func Unpack(data []byte, output interface{}) (uint, error) {
	switch output := output.(type) {
	case *uint8:
		if len(data) < 1 {
			return 0, io.ErrUnexpectedEOF
		}

		*output = data[0]
		return 1, nil

	case *int8:
		if len(data) < 1 {
			return 0, io.ErrUnexpectedEOF
		}

		*output = int8(data[0])
		return 1, nil

	case *uint16:
		return unpackUInt16(data, output)

	case *int16:
		var u uint16
		n, err := unpackUInt16(data, &u)
		*output = int16(u)
		return n, err

	case *uint32:
		return unpackUInt32(data, output)

	case *int32:
		var u uint32
		n, err := unpackUInt32(data, &u)
		*output = int32(u)
		return n, err

	case *uint64:
		return unpackUInt64(data, output)

	case *int64:
		var u uint64
		n, err := unpackUInt64(data, &u)
		*output = int64(u)
		return n, err

	case []byte:
		if len(output) > len(data) {
			return 0, io.ErrUnexpectedEOF
		}

		return uint(copy(output, data)), nil

	case Unpackable:
		return output.Unpack(data)
	}

	return 0, fmt.Errorf("Can't unpack type %T", output)
}

// UnpackSome unpacks multiple values.
func UnpackSome(data []byte, outputs ...interface{}) (uint, error) {
	var n uint
	for _, output := range outputs {
		m, err := Unpack(data[n:], output)
		n += m

		if err != nil {
			return n, err
		}
	}

	return n, nil
}

// UnpackString unpacks a string
func UnpackString(buffer []byte, len uint, output *string) (uint, error) {
	buffer = buffer[:len]
	buffer = bytes.TrimRight(buffer, string(byte(0x0)))
	buffer, err := stringDecoder.Bytes(buffer)
	if err != nil {
		return 0, fmt.Errorf("Unable to decode string: %s", err)
	}

	*output = string(buffer)
	return len, nil
}
