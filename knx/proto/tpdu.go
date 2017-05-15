package proto

import (
	"errors"
)

// A TransportData is the transport-layer protocol data unit (TPDU) within a L_Data frame.
type TransportData []byte

var (
	ErrTransportDataTooShort = errors.New("Given TPDU is too short")
)

// CheckTransportData validates the length of the given slice.
func CheckTransportData(data []byte) (TransportData, error) {
	if len(data) < 1 {
		return nil, ErrTransportDataTooShort
	}

	tpdu := TransportData(data)

	switch tpdu.ControlInfo() {
	case UnnumberedDataPacket, NumberedDataPacket:
		if len(data) < 2 {
			return nil, ErrTransportDataTooShort
		}
	}

	return tpdu, nil
}

// A TransportControlInfo is the transport-layer protocol control information (TPCI).
type TransportControlInfo uint8

const (
	UnnumberedDataPacket    TransportControlInfo = 0
	NumberedDataPacket      TransportControlInfo = 1
	UnnumberedControlPacket TransportControlInfo = 2
	NumberedControlPacket   TransportControlInfo = 3
)

// ControlInfo returns the type of packet in the TPDU.
func (tpdu TransportData) ControlInfo() TransportControlInfo {
	return TransportControlInfo((tpdu[0] >> 6) & 3)
}

// SeqNumber retrieves the sequence number.
func (tpdu TransportData) SeqNumber() uint8 {
	return (tpdu[0] >> 2) & 15
}

// Data parses the application-layer protocol data unit in order to provide control information
// and the actual data.
func (tpdu TransportData) Data() (AppControlInfo, []byte, error) {
	switch tpdu.ControlInfo() {
	case UnnumberedDataPacket, NumberedDataPacket:
		apci := AppControlInfo(((tpdu[0] & 3) << 2) | ((tpdu[1] >> 6) & 3))

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
		return 0, nil, errors.New("TransportControlInfo does not indicate a data packet")
	}
}

// An AppControlInfo is the application-layer protocol control information (APCI).
type AppControlInfo uint8

const (
	GroupValueRead         AppControlInfo = 0
	GroupValueResponse     AppControlInfo = 1
	GroupValueWrite        AppControlInfo = 2
	IndividualAddrWrite    AppControlInfo = 3
	IndividualAddrRequest  AppControlInfo = 4
	IndividualAddrResponse AppControlInfo = 5
	AdcRead                AppControlInfo = 6
	AdcResponse            AppControlInfo = 7
	MemoryRead             AppControlInfo = 8
	MemoryResponse         AppControlInfo = 9
	MemoryWrite            AppControlInfo = 10
	UserMessage            AppControlInfo = 11
	MaskVersionRead        AppControlInfo = 12
	MaskVersionResponse    AppControlInfo = 13
	Restart                AppControlInfo = 14
	Escape                 AppControlInfo = 15
)
