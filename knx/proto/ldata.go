package proto

import (
	"io"
	"github.com/vapourismo/knx-go/knx/binary"
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
func ReadLData(r io.Reader) (*LData, error) {
	ldata := &LData{}

	err := binary.ReadSequence(
		r, &ldata.Control1, &ldata.Control2, &ldata.Source, &ldata.Destination,
	)
	if err != nil {
		return nil, err
	}

	err = ldata.Data.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return ldata, nil
}

// WriteTo writes the LData structure to the given Writer.
func (ldata *LData) WriteTo(w io.Writer) error {
	err := binary.WriteSequence(w, ldata.Control1, ldata.Control2, ldata.Source, ldata.Destination)
	if err != nil {
		return err
	}

	return ldata.Data.WriteTo(w)
}
