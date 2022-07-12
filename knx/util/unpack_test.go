// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package util

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type unpackableTester [10]byte

func (u *unpackableTester) Unpack(data []byte) (uint, error) {
	if len(data) < 10 {
		return 0, io.ErrUnexpectedEOF
	}

	return uint(copy(u[:], data)), nil
}

type unpackType struct {
	typ       interface{}
	length    uint
	extractor func(data []byte) interface{}
}

func TestUnpack(t *testing.T) {
	testCases := []unpackType{
		{
			typ:       uint8(0),
			length:    1,
			extractor: func(data []byte) interface{} { return data[0] },
		},
		{
			typ:       int8(0),
			length:    1,
			extractor: func(data []byte) interface{} { return int8(data[0]) },
		},
		{
			typ:       uint16(0),
			length:    2,
			extractor: func(data []byte) interface{} { return binary.BigEndian.Uint16(data) },
		},
		{
			typ:       uint32(0),
			length:    4,
			extractor: func(data []byte) interface{} { return binary.BigEndian.Uint32(data) },
		},
		{
			typ:       uint64(0),
			length:    8,
			extractor: func(data []byte) interface{} { return binary.BigEndian.Uint64(data) },
		},
		{
			typ:       int16(0),
			length:    2,
			extractor: func(data []byte) interface{} { return int16(binary.BigEndian.Uint16(data)) },
		},
		{
			typ:       int32(0),
			length:    4,
			extractor: func(data []byte) interface{} { return int32(binary.BigEndian.Uint32(data)) },
		},
		{
			typ:       int64(0),
			length:    8,
			extractor: func(data []byte) interface{} { return int64(binary.BigEndian.Uint64(data)) },
		},
		{
			typ:    unpackableTester{},
			length: 10,
			extractor: func(data []byte) interface{} {
				var r unpackableTester
				r.Unpack(data)
				return r
			},
		},
	}

	num, err := Unpack(nil, nil)
	if err == nil {
		t.Error("Should not succeed")
	}

	if num != 0 {
		t.Error("Should return 0 in case of error")
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%T", testCase.typ), func(t *testing.T) {
			ptrValue := reflect.New(reflect.TypeOf(testCase.typ))
			ptr := ptrValue.Interface()

			num, err := Unpack(nil, ptr)
			if err != io.ErrUnexpectedEOF {
				t.Fatal("Unexpected error:", err)
			}

			if num != 0 {
				t.Error("Should return 0 in case of error")
			}

			buffer := make([]byte, testCase.length)

			for i := 0; i < 100; i++ {
				if i > 0 {
					rand.Read(buffer)
				}

				num, err := Unpack(buffer, ptr)

				if err != nil {
					t.Errorf("Unexpected error: %v %v", err, buffer)
				}

				if num != testCase.length {
					t.Errorf("Unexpected length: %v %v", num, buffer)
				}

				expected := testCase.extractor(buffer)
				result := reflect.Indirect(ptrValue).Interface()

				if !reflect.DeepEqual(expected, result) {
					t.Errorf("Unexpected result: %v %v %v", result, expected, buffer)
				}
			}
		})
	}
}

type unpackableBad struct{}

var errBadUnpack = errors.New("bad unpack")

func (unpackableBad) Unpack(data []byte) (uint, error) {
	return 0, errBadUnpack
}

func TestUnpackSome(t *testing.T) {
	num, err := UnpackSome(nil, unpackableBad{})
	if err != errBadUnpack {
		t.Error("unexpected error: ", err)
	}

	if num != 0 {
		t.Error("Unexpected length:", num)
	}

	buffer := make([]byte, 25)

	for i := 0; i < 100; i++ {
		rand.Read(buffer)

		var a uint8
		var b uint16
		var c uint32
		var d uint64
		var e unpackableTester

		num, err := UnpackSome(buffer, &a, &b, &c, &d, &e)
		if err != nil {
			t.Errorf("Unexpected error: %v %v", err, buffer)
		}

		if num != uint(len(buffer)) {
			t.Errorf("Unexpected length: %v %v", num, buffer)
		}
	}
}

func TestUnpackString(t *testing.T) {
	testCases := []struct {
		MaxLen   uint
		Data     []byte
		Expected string
	}{
		{
			MaxLen:   30,
			Data:     []byte{0x41, 0x42, 0x42, 0x20, 0x49, 0x50, 0x53, 0x2f, 0x53, 0x32, 0x2e, 0x31, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			Expected: "ABB IPS/S2.1",
		},
	}

	for _, testCase := range testCases {
		var output string
		n, err := UnpackString(testCase.Data, testCase.MaxLen, &output)

		if assert.Nil(t, err) {
			assert.Equal(t, testCase.MaxLen, n, "Consumed bytes not equal")
			assert.Equal(t, testCase.Expected, output, "Unpacked string not equal")
		}
	}

}
