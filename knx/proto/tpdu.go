package proto

import (
	"io"
	"errors"
)

// A TPCI is the transport-layer protocol control information (TPCI).
type TPCI uint8

//
const (
	UnnumberedDataPacket    TPCI = 0
	NumberedDataPacket      TPCI = 1
	UnnumberedControlPacket TPCI = 2
	NumberedControlPacket   TPCI = 3
)

// An APCI is the application-layer protocol control information (APCI).
type APCI uint8

//
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

// A TPDU is the transport-layer protocol data unit within a L_Data frame.
type TPDU struct {
	PacketType TPCI
	SeqNumber  uint8
	Control    uint8
	Info       APCI
	Data       []byte
}

//
var (
	ErrTransportUnitTooShort = errors.New("Given TPDU is too short")
)

// ReadTPDU parses the given data in order to produce a TPDU struct.
func ReadTPDU(data []byte) (*TPDU, error) {
	if len(data) < 1 {
		return nil, ErrTransportUnitTooShort
	}

	packetType := TPCI((data[0] >> 6) & 3)
	seqNumber := (data[0] >> 2) & 15

	switch packetType {
	case UnnumberedControlPacket, NumberedControlPacket:
		return &TPDU{packetType, seqNumber, data[0] & 3, 0, nil}, nil

	case UnnumberedDataPacket, NumberedDataPacket:
		if len(data) < 2 {
			return nil, ErrTransportUnitTooShort
		}

		info := APCI((data[0] & 3) << 2 | (data[1] >> 6) & 3)

		var appData []byte
		if len(data) > 2 {
			appData = make([]byte, len(data) - 2)
			copy(appData, data[2:])
		} else {
			appData = []byte{data[1] & 63}
		}

		return &TPDU{packetType, seqNumber, 0, info, appData}, nil
	}

	return nil, errors.New("Unreachable")
}

// WriteTo writes the TPDU structure to the given Writer.
func (tpdu TPDU) WriteTo(w io.Writer) (err error) {
	buffer := []byte{
		byte(tpdu.PacketType & 3) << 6 | byte(tpdu.SeqNumber & 15) << 2,
	}

	switch tpdu.PacketType {
	case UnnumberedControlPacket, NumberedControlPacket:
		buffer[0] |= byte(tpdu.Control & 3)

	case UnnumberedDataPacket, NumberedDataPacket:
		buffer[0] |= byte(tpdu.Info >> 2) & 3

		if len(tpdu.Data) > 0 {
			buffer = append(buffer, tpdu.Data...)
			buffer[1] &= 63
			buffer[1] |= byte(tpdu.Info & 3) << 6
		} else {
			buffer = []byte{buffer[0], byte(tpdu.Info & 3) << 6}
		}
	}

	_, err = w.Write(buffer)
	return
}
