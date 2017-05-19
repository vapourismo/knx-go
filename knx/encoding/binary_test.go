package encoding

import (
	"bytes"
	"testing"
)

func TestWriteSequence(t *testing.T) {
	var b1 uint8 = 0x01
	var b2 uint16 = 0x2345
	var b3 uint32 = 0x6789ABCD
	var b4 uint64 = 0xEF0123456789ABCD

	b5 := []uint8{b1, b1}
	b6 := []uint16{b2, b2}
	b7 := []uint32{b3, b3}
	b8 := []uint64{b4, b4}

	buffer := &bytes.Buffer{}

	len, err := WriteSequence(
		buffer,
		b1, b2, b3, b4,
		&b1, &b2, &b3, &b4,
		b5, b6, b7, b8,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len != 60 {
		t.Fatalf("Length mismatch: 60 != %d", len)
	}

	cmp := []byte{
		0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF, 0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD,
		0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF, 0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD,
		0x01, 0x01,
		0x23, 0x45, 0x23, 0x45,
		0x67, 0x89, 0xAB, 0xCD, 0x67, 0x89, 0xAB, 0xCD,
		0xEF, 0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF, 0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD,
	}

	if !bytes.Equal(buffer.Bytes(), cmp) {
		t.Fatalf("Result mismatches: %v != %v", cmp, buffer.Bytes())
	}

	anotherBuffer := &bytes.Buffer{}
	len, err = WriteSequence(anotherBuffer, buffer)

	if err != nil {
		t.Fatal(err)
	}

	if len != 60 {
		t.Fatalf("Length mismatch: 60 != %d", len)
	}

	if !bytes.Equal(anotherBuffer.Bytes(), cmp) {
		t.Fatalf("Result mismatches: %v != %v", cmp, anotherBuffer.Bytes())
	}
}
