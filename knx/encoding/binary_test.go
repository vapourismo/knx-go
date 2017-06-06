package encoding

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/vapourismo/knx-go/knx/util"
)

type writeCase struct {
	value  interface{}
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
		{u1, []byte{0xF0}},
		{u2, []byte{0xAA, 0x55}},
		{u3, []byte{0xAA, 0x55, 0xAA, 0x55}},
		{u4, []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},

		{&u1, []byte{0xF0}},
		{&u2, []byte{0xAA, 0x55}},
		{&u3, []byte{0xAA, 0x55, 0xAA, 0x55}},
		{&u4, []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},

		{i1, []byte{0xF0}},
		{i2, []byte{0xAA, 0x55}},
		{i3, []byte{0xAA, 0x55, 0xAA, 0x55}},
		{i4, []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},

		{&i1, []byte{0xF0}},
		{&i2, []byte{0xAA, 0x55}},
		{&i3, []byte{0xAA, 0x55, 0xAA, 0x55}},
		{&i4, []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},

		{[]uint8{u1, u1}, []byte{0xF0, 0xF0}},
		{[]uint16{u2, u2}, []byte{0xAA, 0x55, 0xAA, 0x55}},
		{[]uint32{u3, u3}, []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},
		{[]uint64{u4, u4}, []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},

		{[]int8{i1, i1}, []byte{0xF0, 0xF0}},
		{[]int16{i2, i2}, []byte{0xAA, 0x55, 0xAA, 0x55}},
		{[]int32{i3, i3}, []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},
		{[]int64{i4, i4}, []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Types/%T", c.value), func(t *testing.T) {
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

	t.Run("WriterTo", func(t *testing.T) {
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

	t.Run("BadWriter", func(t *testing.T) {
		len, err := Write(util.BadWriter{}, uint8(42))

		if err != util.ErrBadWrite {
			t.Error("Unexpected error:", err)
		}

		if len != 0 {
			t.Error("Unexpected length:", len)
		}
	})
}

func TestWriteSome(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		buffer := bytes.Buffer{}

		len, err := WriteSome(&buffer, uint16(0x1337), uint8(0x42))
		if err != nil {
			t.Fatal("Unexpected error:", err)
		}

		if len != 3 {
			t.Fatal("Unexpected length:", len)
		}

		cmp := []byte{0x13, 0x37, 0x42}
		if !bytes.Equal(buffer.Bytes(), cmp) {
			t.Fatalf("Written data mismatches expectations: %v != %v", buffer.Bytes(), cmp)
		}
	})

	t.Run("BadWriter", func(t *testing.T) {
		len, err := WriteSome(util.BadWriter{}, uint16(0x1337), uint8(0x42))
		if err != util.ErrBadWrite {
			t.Fatal("Unexpected error:", err)
		}

		if len != 0 {
			t.Fatal("Unexpected length:", len)
		}
	})
}
