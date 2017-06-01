package proto

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	"github.com/vapourismo/knx-go/utilities/testutils"
)

func TestAddress_String(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			buffer := [4]byte{}
			rand.Read(buffer[:])

			addr := Address(buffer)
			result := addr.String()

			sprintResult := fmt.Sprint(addr)
			if sprintResult != result {
				t.Errorf("Sprint gives different result: %v != %v", sprintResult, result)
			}

			expectedResult := fmt.Sprintf("%d.%d.%d.%d", buffer[0], buffer[1], buffer[2], buffer[3])
			if expectedResult != result {
				t.Errorf("Unexpected result: %v != %v", expectedResult, result)
			}
		}
	})

	t.Run("BadReader", func(t *testing.T) {
		var hi HostInfo
		_, err := hi.ReadFrom(testutils.BadReader{})

		if err != testutils.ErrBadRead {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("BadLength", func(t *testing.T) {
		var hi HostInfo
		_, err := hi.ReadFrom(bytes.NewReader([]byte{0, 1, 0, 0, 0, 0, 0, 0}))

		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

	t.Run("BadProtocol", func(t *testing.T) {
		var hi HostInfo
		_, err := hi.ReadFrom(bytes.NewReader([]byte{8, 255, 0, 0, 0, 0, 0, 0}))

		if err == nil {
			t.Fatal("Should not succeed")
		}
	})
}

func makeRandBuffer(size int) []byte {
	buffer := make([]byte, size)
	rand.Read(buffer)
	return buffer
}

func TestHostInfo_ReadFrom(t *testing.T) {
	for i := 0; i < 100; i++ {
		data := makeRandBuffer(6)
		proto := byte(1 + (rand.Int() % 2))
		reader := bytes.NewReader(append([]byte{8, proto}, data...))

		var hi HostInfo
		len, err := hi.ReadFrom(reader)

		if err != nil {
			t.Errorf("Error for data %v: %v", data, err)
			continue
		}

		if len != 8 {
			t.Errorf("Unexpected number of bytes read: %v", len)
		}

		if hi.Protocol != Protocol(proto) {
			t.Errorf("Unexpected protocol: %v != %v", hi.Protocol, data[1])
		}

		if !bytes.Equal(hi.Address[:], data[:4]) {
			var addrData Address
			copy(addrData[:], data[:4])

			t.Errorf("Unexpected address: %v != %v", hi.Address, addrData)
		}

		portData := Port(data[4])<<8 | Port(data[5])
		if hi.Port != portData {
			t.Errorf("Unexpected port: %v != %v", hi.Port, portData)
		}
	}
}

func TestHostInfo_WriteTo(t *testing.T) {
	for i := 0; i < 100; i++ {
		hi := HostInfo{
			Protocol: Protocol(1 + (rand.Int() % 2)),
			Port:     Port(rand.Int()),
		}
		copy(hi.Address[:], makeRandBuffer(4))

		buffer := bytes.Buffer{}
		len, err := hi.WriteTo(&buffer)

		if err != nil {
			t.Errorf("Error for %v: %v", hi, err)
			continue
		}

		if len != 8 {
			t.Errorf("Unexpected number of bytes written: %v", len)
		}

		var hiCmp HostInfo
		_, err = hiCmp.ReadFrom(&buffer)
		if err != nil {
			t.Errorf("Unexpected read error for %v: %v", buffer.Bytes(), err)
			continue
		}

		if !hi.Equals(hiCmp) {
			t.Errorf("Result does not match: %v != %v", hiCmp, hi)
		}
	}
}
