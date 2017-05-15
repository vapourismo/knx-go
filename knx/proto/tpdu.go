package proto

import (
	"errors"
)

// A TPDU is the transport-layer protocol data unit within a L_Data frame.
type TPDU []byte

var (
	ErrTransportDataTooShort = errors.New("Given TPDU is too short")
	ErrTransportNotData      = errors.New("TPCI does not indicate a data packet")
	ErrTransportNotControl   = errors.New("TPCI does not indicate a control packet")
)

// CheckTPDU validates the length of the given slice.
func CheckTPDU(data []byte) (TPDU, error) {
	if len(data) < 1 {
		return nil, ErrTransportDataTooShort
	}

	tpdu := TPDU(data)

	switch tpdu.PacketType() {
	case UnnumberedDataPacket, NumberedDataPacket:
		if len(data) < 2 {
			return nil, ErrTransportDataTooShort
		}
	}

	return tpdu, nil
}

// A TPCI is the transport-layer protocol control information (TPCI).
type TPCI uint8

const (
	UnnumberedDataPacket    TPCI = 0
	NumberedDataPacket      TPCI = 1
	UnnumberedControlPacket TPCI = 2
	NumberedControlPacket   TPCI = 3
)

// PacketType returns the type of packet in the TPDU.
func (tpdu TPDU) PacketType() TPCI {
	return TPCI((tpdu[0] >> 6) & 3)
}

// SeqNumber retrieves the sequence number.
func (tpdu TPDU) SeqNumber() uint8 {
	return (tpdu[0] >> 2) & 15
}

// AppData parses the application-layer protocol data unit in order to provide control information
// and the actual data.
func (tpdu TPDU) AppData() (APCI, []byte, error) {
	switch tpdu.PacketType() {
	case UnnumberedDataPacket, NumberedDataPacket:
		apci := APCI(((tpdu[0] & 3) << 2) | ((tpdu[1] >> 6) & 3))

		var data []byte

		if len(tpdu) > 2 {
			data = make([]byte, len(tpdu) - 2)
			copy(data, tpdu[2:])
		} else {
			data = make([]byte, 1)
			data[0] = tpdu[1] & 63
		}

		return apci, data, nil

	default:
		return 0, nil, ErrTransportNotData
	}
}

// ControlData retrieves the control data within the data unit.
func (tpdu TPDU) ControlData() (uint8, error) {
	switch tpdu.PacketType() {
	case UnnumberedControlPacket, NumberedControlPacket:
		return tpdu[0] & 3, nil

	default:
		return 0, ErrTransportNotControl
	}
}

// An APCI is the application-layer protocol control information (APCI).
type APCI uint8

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
