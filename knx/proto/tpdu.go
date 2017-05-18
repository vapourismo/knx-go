package proto

import (
	"errors"
)

// A TPCI is the transport-layer protocol control information (TPCI).
type TPCI uint8

// TPCI values
const (
	UnnumberedDataPacket    TPCI = 0
	NumberedDataPacket      TPCI = 1
	UnnumberedControlPacket TPCI = 2
	NumberedControlPacket   TPCI = 3
)

// An APCI is the application-layer protocol control information (APCI).
type APCI uint8

// APCI values
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

// MakeTPDU generates a TPDU that contains an APDU with the given APCI and data. In order to be able
// to properly format the APDU, the given data must have a certain length.
//
// len(data) == 0 indicates no data.
// len(data) == 1 indicates that only the 6 least significant bits are actual data.
// len(data) >  1 indicates that everything but the first byte is data.
//
func (apci APCI) MakeTPDU(data []byte) []byte {
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

	return buffer
}

// Errors from ExtractAPDU
var (
	ErrNoDataPacket = errors.New("Given TPDU is not a data packet")
)

// ExtractAPDUFromTPDU parses the APDU section of a TPDU, if one exists.
func ExtractAPDUFromTPDU(tpdu []byte) (APCI, []byte, error) {
	if len(tpdu) < 2 {
		return 0, nil, ErrDataTooShort
	}

	packetType := TPCI((tpdu[0] >> 6) & 3)
	switch packetType {
	case UnnumberedDataPacket, NumberedDataPacket:
		apci := APCI((tpdu[0] & 3) << 2 | (tpdu[1] >> 6) & 3)
		data := make([]byte, len(tpdu) - 1)
		copy(data, tpdu[1:])

		return apci, data, nil
	}

	return 0, nil, ErrNoDataPacket
}
