package proto

import (
	"bytes"
	"errors"
	"io"
	"github.com/vapourismo/knx-go/knx/encoding"
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

// Errors returned from ReadTPDU
var (
	ErrDataUnitTooShort = errors.New("Data segment of the TPDU is too short")
)

// ReadFrom parses the given data in order to fill the TPDU struct.
func (tpdu *TPDU) ReadFrom(r io.Reader) error {
	var head uint8

	err := encoding.ReadSequence(r, &head)
	if err != nil {
		return err
	}

	packetType := TPCI((head >> 6) & 3)
	seqNumber := (head >> 2) & 15

	switch packetType {
	case UnnumberedControlPacket, NumberedControlPacket:
		tpdu.PacketType = packetType
		tpdu.SeqNumber = seqNumber
		tpdu.Control = head & 3
		tpdu.Info = 0
		tpdu.Data = nil

		return nil

	case UnnumberedDataPacket, NumberedDataPacket:
		buffer := &bytes.Buffer{}

		// Empty the reader's remaining contents into a buffer.
		len, err := buffer.ReadFrom(r)
		if err != nil {
			return err
		} else if len < 1 {
			return ErrDataUnitTooShort
		}

		data := buffer.Bytes()

		tpdu.PacketType = packetType
		tpdu.SeqNumber = seqNumber
		tpdu.Control = 0
		tpdu.Info = APCI((head & 3) << 2 | (data[0] >> 6) & 3)
		tpdu.Data = data

		// The first 2 bits of the data contain rest of the APCI, we don't need them anymore.
		tpdu.Data[0] &= 63

		return nil
	}

	return errors.New("Unreachable")
}

// WriteTo writes the TPDU structure to the given Writer.
func (tpdu *TPDU) WriteTo(w io.Writer) error {
	headMask := byte(tpdu.PacketType & 3) << 6 | byte(tpdu.SeqNumber & 15) << 2

	switch tpdu.PacketType & 3 {
	case UnnumberedControlPacket, NumberedControlPacket:
		_, err := w.Write([]byte{headMask | byte(tpdu.Control & 3)})
		return err

	case UnnumberedDataPacket, NumberedDataPacket:
		data := make([]byte, 1, 2)
		if len(tpdu.Data) > 0 {
			data = append(data, tpdu.Data...)
		} else {
			data = append(data, 0)
		}

		data[0] |= headMask | byte((tpdu.Info >> 2) & 3)
		data[1] &= 63
		data[1] |= byte(tpdu.Info & 3) << 6

		_, err := w.Write(data)
		return err
	}

	return errors.New("Unreachable")
}
