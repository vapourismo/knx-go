package proto

import (
	"errors"
	"io"
	"github.com/vapourismo/knx-go/knx/encoding"
)

// A LData is a link-layer data frame.
type LData struct {
	Control1    uint8
	Control2    uint8
	Source      uint16
	Destination uint16
	Data        TPDU
}

// ReadLData parses the given data in order to extract a LData frame.
func ReadLData(ldata []byte) (*LData, error) {
	if len(ldata) < 8 {
		return nil, ErrDataTooShort
	}

	tpduLen := int(ldata[6])

	if tpduLen > len(ldata) - 8 {
		return nil, ErrDataIncomplete
	}

	return &LData{
		Control1:    ldata[0],
		Control2:    ldata[1],
		Source:      encoding.UInt16(ldata[2:]),
		Destination: encoding.UInt16(ldata[4:]),
		Data:        TPDU(ldata[7:8 + tpduLen]),
	}, nil
}

// WriteTo writes the LData structure to the given Writer.
func (ldata *LData) WriteTo(w io.Writer) error {
	if len(ldata.Data) < 1 {
		return errors.New("TPDU length has be 1 or more")
	} else if len(ldata.Data) > 256 {
		return errors.New("TPDU is too large")
	}

	dataLen := byte(len(ldata.Data) - 1)

	return encoding.WriteSequence(
		w, ldata.Control1, ldata.Control2, ldata.Source, ldata.Destination, dataLen, ldata.Data,
	)
}
