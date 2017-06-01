package cemi

import (
	"io"

	"github.com/vapourismo/knx-go/knx/encoding"
)

// A LBusmonInd represents a L_Busmon.ind message.
type LBusmonInd []byte

// MessageCode returns the message code for L_Busmon.ind.
func (LBusmonInd) MessageCode() MessageCode {
	return LBusmonIndCode
}

// WriteTo serializes the structure and writes it to the given Writer.
func (lbm *LBusmonInd) WriteTo(w io.Writer) (int64, error) {
	len, err := w.Write([]byte(*lbm))
	return int64(len), err
}

// ReadFrom initializes the structure by reading from the given Reader.
func (lbm *LBusmonInd) ReadFrom(r io.Reader) (n int64, err error) {
	n, *lbm = encoding.ReadAll(r)
	return
}
