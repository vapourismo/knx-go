package cemi

import (
	"math/rand"
	"testing"
)

func makeRandTPDUSegment() []byte {
	n := rand.Int() % 256

	buffer := make([]byte, n+2)
	buffer[0] = byte(n)
	rand.Read(buffer[1:])

	return buffer
}

func TestTPDU_Unpack(t *testing.T) {
	for i := 0; i < 100; i++ {
		data := makeRandTPDUSegment()
		tpdu := TPDU{}

		num, err := tpdu.Unpack(data)

		if err != nil {
			t.Error("Unexpected error:", err, data)
			continue
		}

		if num != uint(len(data)) {
			t.Error("Unexpected length:", num, len(data), data)
		}
	}
}
