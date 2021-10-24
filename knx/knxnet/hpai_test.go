// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package knxnet

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"testing"

	"github.com/vapourismo/knx-go/knx/util"
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

	t.Run("BadLength", func(t *testing.T) {
		var hi HostInfo
		_, err := hi.Unpack([]byte{0, 1, 0, 0, 0, 0, 0, 0})

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

func TestHostInfo_Unpack(t *testing.T) {
	for i := 0; i < 100; i++ {
		proto := byte(1 + (rand.Int() % 2))
		data := append([]byte{8, proto}, makeRandBuffer(6)...)

		var hi HostInfo
		num, err := hi.Unpack(data)

		if err != nil {
			t.Errorf("Error for data %v: %v", data, err)
			continue
		}

		if num != uint(len(data)) {
			t.Errorf("Unexpected number of bytes read: %v", num)
		}

		if hi.Protocol != Protocol(proto) {
			t.Errorf("Unexpected protocol: %v != %v", hi.Protocol, data[1])
		}

		if !bytes.Equal(hi.Address[:], data[2:6]) {
			var addrData Address
			copy(addrData[:], data[2:6])

			t.Errorf("Unexpected address: %v != %v", hi.Address, addrData)
		}

		portData := Port(data[6])<<8 | Port(data[7])
		if hi.Port != portData {
			t.Errorf("Unexpected port: %v != %v", hi.Port, portData)
		}
	}
}
func TestHostInfo_Pack(t *testing.T) {
	for i := 0; i < 100; i++ {
		hi := HostInfo{
			Protocol: Protocol(1 + (rand.Int() % 2)),
			Port:     Port(rand.Int()),
		}
		copy(hi.Address[:], makeRandBuffer(4))

		buffer := util.AllocAndPack(&hi)

		var hiCmp HostInfo
		if _, err := hiCmp.Unpack(buffer); err != nil {
			t.Errorf("Unexpected read error for %v: %v", buffer, err)
			continue
		}

		if !hi.Equals(hiCmp) {
			t.Errorf("Result does not match: %v != %v", hiCmp, hi)
		}
	}
}

func TestHostInfoFromAddress(t *testing.T) {
	t.Run("valid UDP4 address", func(t *testing.T) {
		address := net.UDPAddr{
			IP:   net.IPv4(192, 168, 1, 15),
			Port: 1234,
		}

		expected := HostInfo{
			Protocol: UDP4,
			Address:  [4]byte{192, 168, 1, 15},
			Port:     1234,
		}

		info, err := HostInfoFromAddress(&address)

		if err != nil {
			t.Fatal("Expected error to be nil, but it was '", err, "'")
		}

		if !expected.Equals(info) {
			t.Fatal("Expected", expected, "but it was", info)
		}
	})

	t.Run("invalid UDP address", func(t *testing.T) {
		address := net.UDPAddr{}
		_, err := HostInfoFromAddress(&address)

		if err == nil {
			t.Fatal("Should not succeed")
		}
	})

}
