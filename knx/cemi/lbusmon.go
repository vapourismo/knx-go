package cemi

import "io"

// LBusmon represents the L_Busmon.ind frame type.
type LBusmon []byte

// WriteTo serializes the LBusmon structure and writes it to the given Writer.
func (lbm LBusmon) WriteTo(w io.Writer) (int64, error) {
	len, err := w.Write([]byte(lbm))
	return int64(len), err
}
