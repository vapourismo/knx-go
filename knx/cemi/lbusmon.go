// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package cemi

// A LBusmonInd represents a L_Busmon.ind message.
type LBusmonInd []byte

// MessageCode returns the message code for L_Busmon.ind.
func (LBusmonInd) MessageCode() MessageCode {
	return LBusmonIndCode
}

// Size returns the packed size.
func (lbm LBusmonInd) Size() uint {
	return uint(len(lbm))
}

// Pack the message body into the buffer.
func (lbm LBusmonInd) Pack(buffer []byte) {
	copy(buffer, lbm)
}

// Unpack initializes the structure by parsing the given data.
func (lbm *LBusmonInd) Unpack(data []byte) (n uint, err error) {
	target := []byte(*lbm)

	if len(target) < len(data) {
		target = make([]byte, len(data))
	}

	n = uint(copy(target, data))
	*lbm = LBusmonInd(target)

	return
}
