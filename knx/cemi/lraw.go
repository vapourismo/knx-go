package cemi

import (
	"io"

	"github.com/vapourismo/knx-go/knx/encoding"
)

// A LRaw is a raw link-layer frame. L_Raw.req, L_Raw.con and L_Raw.ind share this structure.
type LRaw []byte

// WriteTo serializes the structure and writes it to the given Writer.
func (lbm *LRaw) WriteTo(w io.Writer) (int64, error) {
	len, err := w.Write([]byte(*lbm))
	return int64(len), err
}

// ReadFrom initializes the structure by reading from the given Reader.
func (lbm *LRaw) ReadFrom(r io.Reader) (n int64, err error) {
	n, *lbm = encoding.ReadAll(r)
	return
}

// A LRawReq represents a L_Raw.req message body.
type LRawReq struct {
	LRaw
}

// MessageCode returns the message code for L_Raw.req.
func (LRawReq) MessageCode() MessageCode {
	return LRawReqCode
}

// A LRawCon represents a L_Raw.con message body.
type LRawCon struct {
	LRaw
}

// MessageCode returns the message code for L_Raw.con.
func (LRawCon) MessageCode() MessageCode {
	return LRawConCode
}

// A LRawInd represents a L_Raw.ind message body.
type LRawInd struct {
	LRaw
}

// MessageCode returns the message code for L_Raw.ind.
func (LRawInd) MessageCode() MessageCode {
	return LRawConCode
}
