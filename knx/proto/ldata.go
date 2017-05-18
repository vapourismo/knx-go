package proto

import (
	"io"
	"github.com/vapourismo/knx-go/knx/encoding"
)

// A LData is a link-layer data frame.
type LData struct {
	Control1    uint8
	Control2    uint8
	Source      uint16
	Destination uint16
	TPDU        []byte
}

// ReadLData parses the given data in order to extract a LData frame.
func ReadLData(ldata []byte) (*LData, error) {
	if len(ldata) < 7 {
		return nil, ErrDataTooShort
	}

	return &LData{
		Control1:    ldata[0],
		Control2:    ldata[1],
		Source:      encoding.UInt16(ldata[2:]),
		Destination: encoding.UInt16(ldata[4:]),
		TPDU:        ldata[6:],
	}, nil
}

// WriteTo writes the LData structure to the given Writer.
func (ldata *LData) WriteTo(w io.Writer) error {
	return encoding.WriteSequence(
		w, ldata.Control1, ldata.Control2, ldata.Source, ldata.Destination, ldata.TPDU,
	)
}
