// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package cemi

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"testing"
)

func makeRandBuffer(n int) []byte {
	buffer := make([]byte, n)
	rand.Read(buffer)
	return buffer
}

func makeRandLData() []byte {
	return bytes.Join([][]byte{
		makeRandInfoSegment(),
		makeRandBuffer(6),
		{0, 1 << 7},
	}, nil)
}

func TestLData_Unpack(t *testing.T) {
	for i := 0; i < 100; i++ {
		data := makeRandLData()
		ldata := LData{}

		num, err := ldata.Unpack(data)
		if err != nil {
			t.Error("Unexpected error:", err, data)
			continue
		}

		if num != uint(len(data)) {
			t.Error("Unexpected length:", num, len(data), data)
			continue
		}

		if int(data[0]) != len(ldata.Info) {
			t.Error("Unexpected info length:", int(data[0]), len(ldata.Info), data)
			continue
		}

		if !bytes.Equal(data[1:1+len(ldata.Info)], ldata.Info) {
			t.Error("Info content mismatch: ", data[1:1+len(ldata.Info)], ldata.Info, data)
			continue
		}

		data = data[1+len(ldata.Info):]

		if ControlField1(data[0]) != ldata.Control1 {
			t.Error("Unexpected control field 1", ControlField1(data[0]), ldata.Control1, data)
		}

		if ControlField2(data[1]) != ldata.Control2 {
			t.Error("Unexpected control field 2", ControlField2(data[1]), ldata.Control2, data)
		}

		if binary.BigEndian.Uint16(data[2:]) != uint16(ldata.Source) {
			t.Error("Unexpected source:", binary.BigEndian.Uint16(data[2:]), ldata.Source, data)
		}

		if binary.BigEndian.Uint16(data[4:]) != ldata.Destination {
			t.Error("Unexpected destination:", binary.BigEndian.Uint16(data[4:]), ldata.Destination, data)
		}
	}
}
