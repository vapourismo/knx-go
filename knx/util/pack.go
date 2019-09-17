// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package util

import (
	"fmt"

	"golang.org/x/text/encoding/charmap"
)

var (
	stringCharmap = charmap.ISO8859_1
	stringEncoder = stringCharmap.NewEncoder()
)

// Packable is implemented by types that can be packed into a byte slice.
type Packable interface {
	Size() uint
	Pack(buffer []byte)
}

// Pack a value into the buffer.
func Pack(buffer []byte, input interface{}) uint {
	switch input := input.(type) {
	case uint8:
		buffer[0] = input
		return 1

	case int8:
		buffer[0] = uint8(input)
		return 1

	case uint16:
		buffer[1] = uint8(input)
		buffer[0] = uint8(input >> 8)
		return 2

	case int16:
		uinput := uint16(input)
		buffer[1] = uint8(uinput)
		buffer[0] = uint8(uinput >> 8)
		return 2

	case uint32:
		buffer[3] = uint8(input)
		buffer[2] = uint8(input >> 8)
		buffer[1] = uint8(input >> 16)
		buffer[0] = uint8(input >> 24)
		return 4

	case int32:
		uinput := uint32(input)
		buffer[3] = uint8(uinput)
		buffer[2] = uint8(uinput >> 8)
		buffer[1] = uint8(uinput >> 16)
		buffer[0] = uint8(uinput >> 24)
		return 4

	case uint64:
		buffer[7] = uint8(input)
		buffer[6] = uint8(input >> 8)
		buffer[5] = uint8(input >> 16)
		buffer[4] = uint8(input >> 24)
		buffer[3] = uint8(input >> 32)
		buffer[2] = uint8(input >> 40)
		buffer[1] = uint8(input >> 48)
		buffer[0] = uint8(input >> 56)
		return 8

	case int64:
		uinput := uint64(input)
		buffer[7] = uint8(uinput)
		buffer[6] = uint8(uinput >> 8)
		buffer[5] = uint8(uinput >> 16)
		buffer[4] = uint8(uinput >> 24)
		buffer[3] = uint8(uinput >> 32)
		buffer[2] = uint8(uinput >> 40)
		buffer[1] = uint8(uinput >> 48)
		buffer[0] = uint8(uinput >> 56)
		return 8

	case []byte:
		return uint(copy(buffer, input))

	case Packable:
		input.Pack(buffer)
		return input.Size()
	}

	panic(fmt.Sprintf("Can't pack type %T", input))
}

// PackSome packs multiple values.
func PackSome(buffer []byte, inputs ...interface{}) {
	var offset uint
	for _, output := range inputs {
		offset += Pack(buffer[offset:], output)
	}
}

// AllocAndPack allocates a buffer and then packs the data inside it.
func AllocAndPack(inputs ...Packable) []byte {
	var size uint
	for _, output := range inputs {
		size += output.Size()
	}

	buffer := make([]byte, size)

	var offset uint
	for _, output := range inputs {
		output.Pack(buffer[offset:])
		offset += output.Size()
	}

	return buffer
}

// PackString packs a string into the buffer
func PackString(buffer []byte, maxLen uint, input string) (uint, error) {
	encoded, err := stringEncoder.Bytes([]byte(input))
	if err != nil {
		return 0, fmt.Errorf("Unable to encode string: %s", err)
	}

	if len(encoded) >= int(maxLen) {
		encoded = encoded[:maxLen]
		encoded[maxLen] = 0x00
	}

	copy(buffer, encoded)
	for i := len(encoded); i < int(maxLen); i++ {
		buffer[i] = 0x00
	}

	return maxLen, nil
}
