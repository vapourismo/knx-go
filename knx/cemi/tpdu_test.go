package cemi

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/vapourismo/knx-go/knx/util"
)

func TestAppData_Pack(t *testing.T) {
	for i := 0; i < 100; i++ {
		app := AppData{
			Numbered:  rand.Int()%2 == 0,
			SeqNumber: uint8(rand.Int()) % 15,
			Command:   APCI(rand.Int()) % 15,
			Data:      makeRandBuffer(rand.Int() % 300),
		}

		if len(app.Data) > 0 {
			app.Data[0] &= 63
		}

		data := util.AllocAndPack(&app)

		if len(data) < 3 {
			t.Error("Unexpected length:", len(data), app)
			continue
		}

		dataLength := len(app.Data)
		if dataLength > 255 {
			dataLength = 255
		}

		if len(app.Data) > 0 && int(data[0]) != dataLength {
			t.Error("Unexpected unit length:", data[0], app)
		}

		if data[1]>>7 != 0 {
			t.Error("Not a app unit")
		}

		if ((data[1] & (1 << 6)) == (1 << 6)) != app.Numbered {
			t.Error("Unexpected numbered bit:", ((data[1] & (1 << 6)) == (1 << 6)), app.Numbered)
		}

		if app.Numbered && (data[1]>>2)&15 != app.SeqNumber {
			t.Error("Unexpected sequence number", (data[1]>>2)&15, app.SeqNumber)
		}

		apci := APCI((data[1]&3)<<2 | data[2]>>6)
		if apci != app.Command {
			t.Error("Unexpected command:", apci, app.Command)
		}

		if len(app.Data) > 0 && data[2]&63 != app.Data[0]&63 {
			t.Error("First data byte mismatches:", data[2]&63, app.Data[0]&63)
		}

		if len(app.Data) > 1 && !bytes.Equal(data[3:2+dataLength], app.Data[1:dataLength]) {
			t.Error("Trailing data mismatch", data[3:2+dataLength], app.Data[1:dataLength])
		}
	}
}

func TestControlData_Pack(t *testing.T) {
	for i := 0; i < 100; i++ {
		control := ControlData{
			Numbered:  rand.Int()%2 == 0,
			SeqNumber: uint8(rand.Int()) % 15,
			Command:   uint8(rand.Int()) % 3,
		}

		data := util.AllocAndPack(&control)

		if len(data) < 2 {
			t.Error("Unexpected length:", len(data), control)
			continue
		}

		if data[0] != 0 {
			t.Error("Unexpected unit length:", data[0], control)
		}

		if data[1]>>7 != 1 {
			t.Error("Not a control unit")
		}

		if ((data[1] & (1 << 6)) == (1 << 6)) != control.Numbered {
			t.Error("Unexpected numbered bit:", ((data[1] & (1 << 6)) == (1 << 6)), control.Numbered)
		}

		if control.Numbered && (data[1]>>2)&15 != control.SeqNumber {
			t.Error("Unexpected sequence number", (data[1]>>2)&15, control.SeqNumber)
		}

		if data[1]&3 != control.Command {
			t.Error("Unexpected comand:", data[1]&3, control.Command)
		}
	}
}

func TestUnpackTransportUnit(t *testing.T) {
	t.Run("Control", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			data := []byte{0, byte(rand.Int())}
			data[1] |= 1 << 7

			var unit TransportUnit
			num, err := UnpackTransportUnit(data, &unit)

			if err != nil {
				t.Error("Unexpected error:", err, data)
				continue
			}

			if num != 2 {
				t.Error("Unexpected length:", num, data)
				continue
			}

			control, ok := unit.(*ControlData)
			if !ok {
				t.Errorf("Unexpected result type: %T %v", unit, data)
				continue
			}

			if control.Numbered != (data[1]&(1<<6) == 1<<6) {
				t.Error("Unexpected numbered value:", control.Numbered, (data[1]&(1<<6) == 1<<6))
			}

			if control.Numbered && control.SeqNumber != (data[1]>>2)&15 {
				t.Error("Unexpected sequence number:", control.SeqNumber, (data[1]>>2)&15)
			}

			if control.Command != data[1]&3 {
				t.Error("Unexpected command:", control.Command, data[1]&3)
			}
		}
	})

	t.Run("App", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			data := make([]byte, 3+rand.Int()%255)
			rand.Read(data[1:])

			data[0] = byte(len(data) - 2)
			data[1] &= ^(byte(1) << 7)

			var unit TransportUnit
			num, err := UnpackTransportUnit(data, &unit)

			if err != nil {
				t.Error("Unexpected error:", err, data)
				continue
			}

			if num != uint(len(data)) {
				t.Error("Unexpected length:", num, data)
				continue
			}

			app, ok := unit.(*AppData)
			if !ok {
				t.Errorf("Unexpected result type: %T %v", unit, data)
				continue
			}

			if app.Numbered != (data[1]&(1<<6) == 1<<6) {
				t.Error("Unexpected numbered value:", app.Numbered, (data[1]&(1<<6) == 1<<6))
			}

			if app.Numbered && app.SeqNumber != (data[1]>>2)&15 {
				t.Error("Unexpected sequence number:", app.SeqNumber, (data[1]>>2)&15)
			}

			apci := APCI((data[1]&3)<<2 | data[2]>>6)
			if app.Command != apci {
				t.Error("Unexpected command:", app.Command, apci)
			}

			if len(app.Data) > 0 && data[2]&63 != app.Data[0]&63 {
				t.Error("First data byte mismatches:", data[2]&63, app.Data[0]&63)
			}

			if len(app.Data) > 1 && !bytes.Equal(data[3:], app.Data[1:]) {
				t.Error("Trailing data mismatch", data[3:], app.Data[1:])
			}
		}
	})
}
