package encoding

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/vapourismo/knx-go/utilities/testutils"
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
		t.Run(fmt.Sprintf("Types/%T", c.value), func (t *testing.T) {
			buffer := bytes.Buffer{}
			n, err := Write(&buffer, c.value)

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

	t.Run("WriterTo", func (t *testing.T) {
		source := bytes.Buffer{}
		n, err := source.WriteString("Hello World")
		if err != nil {
			t.Fatal(err)
		}

		sourceData := source.Bytes()

		target := bytes.Buffer{}
		m, err := Write(&target, &source)
		if err != nil {
			t.Fatal(err)
		}

		if int64(n) != m {
			t.Errorf("Written length (%d) does not match the read length (%d)", m, n)
		}

		if !bytes.Equal(sourceData, target.Bytes()) {
			t.Errorf(
				"Written contents (%v) do not match read contents (%v)",
				sourceData, target.Bytes(),
			)
		}
	})

	t.Run("BadWriter", func (t *testing.T) {
		len, err := Write(testutils.BadWriter{}, uint8(42))

		if err != testutils.ErrBadWrite {
			t.Errorf("Unexpected error: %v", err)
		}

		if len != 0 {
			t.Errorf("Unexpected length: %d", len)
		}
	})
}

func TestWriteSome(t *testing.T) {
	t.Run("Ok", func (t *testing.T) {
		buffer := bytes.Buffer{}

		len, err := WriteSome(&buffer, uint16(0x1337), uint8(0x42))
		if err != nil {
			t.Fatal("Unexpected error: %v", err)
		}

		if len != 3 {
			t.Fatalf("Unexpected length: ", len)
		}

		cmp := []byte{0x13, 0x37, 0x42}
		if !bytes.Equal(buffer.Bytes(), cmp) {
			t.Fatalf("Written data mismatches expectations: %v != %v", buffer.Bytes(), cmp)
		}
	})

	t.Run("BadWriter", func (t *testing.T) {
		len, err := WriteSome(testutils.BadWriter{}, uint16(0x1337), uint8(0x42))
		if err != testutils.ErrBadWrite {
			t.Fatal("Unexpected error: %v", err)
		}

		if len != 0 {
			t.Fatalf("Unexpected length: ", len)
		}
	})
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
		t.Run(fmt.Sprintf("Types/%T", c.value), func (t *testing.T) {
			typ := reflect.TypeOf(c.value)

			r := bytes.NewReader(c.input)

			var n int64
			var err error
			var res interface{}

			if typ.Kind() == reflect.Slice {
				sliceLen := reflect.ValueOf(c.value).Len()

				res = reflect.MakeSlice(typ, sliceLen, sliceLen).Interface()
				n, err = Read(r, res)
			} else {
				ptr := reflect.New(typ)

				n, err = Read(r, ptr.Interface())
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

	t.Run("ReaderFrom", func (t *testing.T) {
		source := bytes.Buffer{}
		n, err := source.WriteString("Hello World")
		if err != nil {
			t.Fatal(err)
		}

		sourceData := source.Bytes()

		target := bytes.Buffer{}
		m, err := Read(&source, &target)
		if err != nil {
			t.Fatal(err)
		}

		if int64(n) != m {
			t.Errorf("Read length (%d) does not match the written length (%d)", m, n)
		}

		if !bytes.Equal(sourceData, target.Bytes()) {
			t.Errorf(
				"Written contents (%v) do not match read contents (%v)",
				sourceData, target.Bytes(),
			)
		}
	})

	t.Run("BadReader", func (t *testing.T) {
		var target uint8
		len, err := Read(testutils.BadReader{}, &target)

		if err != testutils.ErrBadRead {
			t.Errorf("Unexpected error: %v", err)
		}

		if len != 0 {
			t.Errorf("Unexpected length: %d", len)
		}
	})
}

func TestReadSome(t *testing.T) {
	t.Run("Ok", func (t *testing.T) {
		source := bytes.NewReader([]byte{0x13, 0x37, 0x42})

		var a uint16
		var b uint8

		len, err := ReadSome(source, &a, &b)

		if err != nil {
			t.Fatal(err)
		}

		if len != 3 {
			t.Fatalf("Unexpected length: %d", len)
		}

		if a != 0x1337 {
			t.Errorf("Value 'a' mismatches: %d != 0x1337", a)
		}

		if b != 0x42 {
			t.Errorf("Value 'b' mismatches: %d != 0x42", b)
		}
	})

	t.Run("BadReader", func (t *testing.T) {
		var a uint16
		var b uint8

		len, err := ReadSome(testutils.BadReader{}, &a, &b)

		if err != testutils.ErrBadRead {
			t.Errorf("Unexpected error: %v", err)
		}

		if len != 0 {
			t.Errorf("Unexpected length: %d", len)
		}
	})
}
