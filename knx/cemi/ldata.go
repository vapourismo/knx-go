package cemi

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

// ReadFrom initializes the LData structure by reading from the given Reader.
func (ldata *LData) ReadFrom(r io.Reader) (n int64, err error) {
	var tpduLen8 uint8
	len, err := encoding.ReadSome(
		r,
		&ldata.Control1,
		&ldata.Control2,
		&ldata.Source,
		&ldata.Destination,
		&tpduLen8,
	)
	n += len

	if err != nil {
		return n, err
	}

	tpdu := make([]byte, int(tpduLen8) + 1)
	len, err = encoding.Read(r, tpdu)
	n += len

	if err != nil {
		return n, err
	}

	ldata.Data = TPDU(tpdu)

	return
}

// WriteTo serializes the LData structure and writes it to the given Writer.
func (ldata *LData) WriteTo(w io.Writer) (int64, error) {
	if len(ldata.Data) < 1 {
		return 0, errors.New("TPDU length has be 1 or more")
	} else if len(ldata.Data) > 256 {
		return 0, errors.New("TPDU is too large")
	}

	return encoding.WriteSome(
		w,
		ldata.Control1,
		ldata.Control2,
		ldata.Source,
		ldata.Destination,
		byte(len(ldata.Data) - 1),
		ldata.Data,
	)
}
