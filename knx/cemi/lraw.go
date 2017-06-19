// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package cemi

// A LRaw is a raw link-layer frame. L_Raw.req, L_Raw.con and L_Raw.ind share this structure.
type LRaw []byte

// Size returns the packed size.
func (lraw LRaw) Size() uint {
	return uint(len(lraw))
}

// Pack the message body into the buffer.
func (lraw LRaw) Pack(buffer []byte) {
	copy(buffer, lraw)
}

// Unpack initializes the structure by parsing the given data.
func (lraw *LRaw) Unpack(data []byte) (n uint, err error) {
	target := []byte(*lraw)

	if len(target) < len(data) {
		target = make([]byte, len(data))
	}

	n = uint(copy(target, data))
	*lraw = LRaw(target)

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
