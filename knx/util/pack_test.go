// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackString(t *testing.T) {
	testCases := []struct {
		MaxLen   uint
		Data     string
		Expected []byte
	}{
		{
			MaxLen:   30,
			Data:     "ABB IPS/S2.1",
			Expected: []byte{0x41, 0x42, 0x42, 0x20, 0x49, 0x50, 0x53, 0x2f, 0x53, 0x32, 0x2e, 0x31, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
	}

	for _, testCase := range testCases {
		buffer := make([]byte, testCase.MaxLen)
		for i := range buffer {
			buffer[i] = 0xff
		}

		n, err := PackString(buffer, testCase.MaxLen, testCase.Data)
		if assert.Nil(t, err) {
			assert.Equal(t, testCase.MaxLen, n, "Produced bytes not equal")
			assert.Equal(t, testCase.Expected, buffer, "Packed burffer not equal")
		}
	}

}
