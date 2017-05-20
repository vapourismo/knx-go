package encoding

import (
	"bytes"
	"testing"
	"fmt"
	"reflect"
)

type writeCase struct {
	value interface{}
	result []byte
}

func TestWrite(t *testing.T) {
	u1 := uint8(0xF0)
	u2 := uint16(0xAA55)
	u3 := uint32(0xAA55AA55)
	u4 := uint64(0xAA55AA55AA55AA55)

	i1 := int8(u1)
	i2 := int16(u2)
	i3 := int32(u3)
	i4 := int64(u4)

	cases := []writeCase{
		{u1,  []byte{0xF0}},
		{u2,  []byte{0xAA, 0x55}},
		{u3,  []byte{0xAA, 0x55, 0xAA, 0x55}},
		{u4,  []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},

		{&u1, []byte{0xF0}},
		{&u2, []byte{0xAA, 0x55}},
		{&u3, []byte{0xAA, 0x55, 0xAA, 0x55}},
		{&u4, []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},

		{i1,  []byte{0xF0}},
		{i2,  []byte{0xAA, 0x55}},
		{i3,  []byte{0xAA, 0x55, 0xAA, 0x55}},
		{i4,  []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},

		{&i1, []byte{0xF0}},
		{&i2, []byte{0xAA, 0x55}},
		{&i3, []byte{0xAA, 0x55, 0xAA, 0x55}},
		{&i4, []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},

		{[]uint8{u1, u1},  []byte{0xF0, 0xF0}},
		{[]uint16{u2, u2}, []byte{0xAA, 0x55, 0xAA, 0x55}},
		{[]uint32{u3, u3}, []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},
		{[]uint64{u4, u4}, []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},

		{[]int8{i1, i1},   []byte{0xF0, 0xF0}},
		{[]int16{i2, i2},  []byte{0xAA, 0x55, 0xAA, 0x55}},
		{[]int32{i3, i3},  []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},
		{[]int64{i4, i4},  []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%T", c.value), func (t *testing.T) {
			buffer := &bytes.Buffer{}
			n, err := Write(buffer, c.value)

			if err != nil {
				t.Fatal(err)
			}

			if n != int64(len(c.result)) {
				t.Fatalf("Length mismatch: %v != %v", n, len(c.result))
			}

			if !bytes.Equal(buffer.Bytes(), c.result) {
				t.Fatalf("Result mismatch: %v != %v", buffer.Bytes(), c.result)
			}
		})
	}
}

type readCase struct {
	value interface{}
	input []byte
}

func TestRead(t *testing.T) {
	u1 := uint8(0xF0)
	u2 := uint16(0xAA55)
	u3 := uint32(0xAA55AA55)
	u4 := uint64(0xAA55AA55AA55AA55)

	i1 := int8(u1)
	i2 := int16(u2)
	i3 := int32(u3)
	i4 := int64(u4)

	cases := []readCase{
		{u1,  []byte{0xF0}},
		{u2,  []byte{0xAA, 0x55}},
		{u3,  []byte{0xAA, 0x55, 0xAA, 0x55}},
		{u4,  []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},

		{i1,  []byte{0xF0}},
		{i2,  []byte{0xAA, 0x55}},
		{i3,  []byte{0xAA, 0x55, 0xAA, 0x55}},
		{i4,  []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},

		{[]uint8{u1, u1},  []byte{0xF0, 0xF0}},
		{[]uint16{u2, u2}, []byte{0xAA, 0x55, 0xAA, 0x55}},
		{[]uint32{u3, u3}, []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},
		{[]uint64{u4, u4}, []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},

		{[]int8{i1, i1},   []byte{0xF0, 0xF0}},
		{[]int16{i2, i2},  []byte{0xAA, 0x55, 0xAA, 0x55}},
		{[]int32{i3, i3},  []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},
		{[]int64{i4, i4},  []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%T", c.value), func (t *testing.T) {
			typ := reflect.TypeOf(c.value)

			var n int64
			var err error
			var res interface{}

			if typ.Kind() == reflect.Slice {
				sliceLen := reflect.ValueOf(c.value).Len()

				res = reflect.MakeSlice(typ, sliceLen, sliceLen).Interface()
				n, err = Read(c.input, res)
			} else {
				ptr := reflect.New(typ)

				n, err = Read(c.input, ptr.Interface())
				res = reflect.Indirect(ptr).Interface()
			}

			if err != nil {
				t.Fatal(err)
			}

			if n != int64(len(c.input)) {
				t.Fatalf("Length mismatch: %v != %v", n, len(c.input))
			}

			if !reflect.DeepEqual(res, c.value) {
				t.Fatalf("Value mismatch: %v != %v", res, c.value)
			}
		})
	}
}
