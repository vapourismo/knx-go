// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package cemi

import (
	"bytes"
	"math/rand"
	"testing"
)

func makeRandInfoSegment() []byte {
	n := rand.Int() % 256

	buffer := make([]byte, n+1)
	buffer[0] = byte(n)
	rand.Read(buffer[1:])

	return buffer
}

func TestInfo_Unpack(t *testing.T) {
	for i := 0; i < 100; i++ {
		data := makeRandInfoSegment()
		info := Info{}

		num, err := info.Unpack(data)

		if err != nil {
			t.Error("Unexpected error:", err, data)
			continue
		}

		if num != uint(len(data)) {
			t.Error("Unexpected length:", num, len(data), data)
		}

		if !bytes.Equal([]byte(info), data[1:]) {
			t.Error("Unexpected result:", info, data)
		}
	}
}

func TestUnpack(t *testing.T) {
	ldataCodes := []MessageCode{LDataReqCode, LDataConCode, LDataIndCode}

	for i := 0; i < 100; i++ {
		code := ldataCodes[rand.Int()%3]
		data := append([]byte{byte(code)}, makeRandLData()...)

		var msg Message
		num, err := Unpack(data, &msg)

		if err != nil {
			t.Error("Unexpected error:", err, data)
			continue
		}

		if num != uint(len(data)) {
			t.Error("Unexpected length:", num, len(data), data)
		}
	}
}
