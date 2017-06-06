package cemi

import "io"

// A TPCI is the Transport-layer Protocol Control Information.
type TPCI uint8

// TPCI values
const (
	UnnumberedDataPacket    TPCI = 0
	NumberedDataPacket      TPCI = 1
	UnnumberedControlPacket TPCI = 2
	NumberedControlPacket   TPCI = 3
)

// APCI is the Application-layer Protocol Control Information.
type APCI uint8

// These are usable APCI values.
const (
	GroupValueRead         APCI = 0
	GroupValueResponse     APCI = 1
	GroupValueWrite        APCI = 2
	IndividualAddrWrite    APCI = 3
	IndividualAddrRequest  APCI = 4
	IndividualAddrResponse APCI = 5
	AdcRead                APCI = 6
	AdcResponse            APCI = 7
	MemoryRead             APCI = 8
	MemoryResponse         APCI = 9
	MemoryWrite            APCI = 10
	UserMessage            APCI = 11
	MaskVersionRead        APCI = 12
	MaskVersionResponse    APCI = 13
	Restart                APCI = 14
	Escape                 APCI = 15
)

// An AppData contains application data in a transport unit.
type AppData struct {
	Numbered  bool
	SeqNumber uint8
	Command   APCI
	Data      []byte
}

// Pack the structure.
func (app *AppData) Pack() []byte {
	length := len(app.Data)

	if length > 255 {
		length = 255
	} else if length < 1 {
		length = 1
	}

	buffer := make([]byte, length+2)
	buffer[0] = byte(length)

	if app.Numbered {
		buffer[1] |= 1<<6 | (app.SeqNumber&15)<<2
	}

	buffer[1] |= byte(app.Command>>2) & 3

	copy(buffer[2:], app.Data)

	buffer[2] &= 63
	buffer[2] |= byte(app.Command&3) << 6

	return buffer
}

// A ControlData encodes control information in a transport unit.
type ControlData struct {
	Numbered  bool
	SeqNumber uint8
	Command   uint8
}

// Pack the structure.
func (control *ControlData) Pack() []byte {
	buffer := []byte{0, 1<<7 | (control.Command & 3)}

	if control.Numbered {
		buffer[1] |= 1<<6 | (control.SeqNumber&15)<<2
	}

	return buffer
}

// A TransportUnit is responsive to transport data.
type TransportUnit interface {
	Pack() []byte
}

// UnpackTransportUnit parses the given data in order to extract the transport unit that it encodes.
func UnpackTransportUnit(data []byte, unit *TransportUnit) (uint, error) {
	if len(data) < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	// Does unit contain control information?
	if (data[1] & (1 << 7)) == 1<<7 {
		control := &ControlData{
			Numbered:  (data[1] & (1 << 6)) == 1<<6,
			SeqNumber: (data[1] >> 2) & 15,
			Command:   data[1] & 3,
		}

		*unit = control

		return 2, nil
	}

	dataLength := int(data[0])

	if len(data) < 3 || dataLength+2 < len(data) {
		return 0, io.ErrUnexpectedEOF
	}

	app := &AppData{
		Numbered:  (data[1] & (1 << 6)) == 1<<6,
		SeqNumber: (data[1] >> 2) & 15,
		Command:   APCI((data[1]&3)<<2 | data[2]>>6),
		Data:      make([]byte, dataLength),
	}

	copy(app.Data, data[2:])
	app.Data[0] &= 63

	*unit = app

	return uint(dataLength) + 2, nil
}
