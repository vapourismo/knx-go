package cemi

import (
	"errors"
)

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

// TPDU is the Transport-layer Protocol Data Unit.
type TPDU []byte

// MakeTPDU generates a TPDU that contains an APDU with the given APCI and data. In order to be able
// to properly format the APDU, the given data must have a certain length.
//
// If len(data) > 0, then you must not utilize the 2 most significant bits of the first byte. They
// will be utilized to store part of the APCI.
//
// Why? Because KNX.
func MakeTPDU(apci APCI, data []byte) TPDU {
	var buffer []byte

	if len(data) > 0 {
		buffer = make([]byte, len(data) + 1)
		copy(buffer[1:], data)
	} else {
		buffer = make([]byte, 2)
	}

	buffer[0] |= byte(UnnumberedDataPacket) << 6
	buffer[0] |= byte(apci >> 2) & 3

	buffer[1] &= 63
	buffer[1] |= byte(apci & 3) << 6

	return TPDU(buffer)
}

// These are errors than can occur when processing the TPDU.
var (
	ErrDataTooShort = errors.New("TPDU is too short")
	ErrNoDataPacket = errors.New("TPCI does not indicate a data packet")
)

// ExtractAPDU parses the APDU section, if one exists.
func (tpdu TPDU) ExtractAPDU() (APCI, []byte, error) {
	if len(tpdu) < 2 {
		return 0, nil, ErrDataTooShort
	}

	packetType := TPCI((tpdu[0] >> 6) & 3)
	switch packetType {
	case UnnumberedDataPacket, NumberedDataPacket:
		apci := APCI((tpdu[0] & 3) << 2 | (tpdu[1] >> 6) & 3)
		data := make([]byte, len(tpdu) - 1)
		copy(data, tpdu[1:])

		data[0] &= 63

		return apci, data, nil
	}

	return 0, nil, ErrNoDataPacket
}
