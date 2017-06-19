// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package cemi

import "github.com/vapourismo/knx-go/knx/util"

// A LData is a link-layer data frame. L_Data.req, L_Data.con and L_Data.ind share this structure.
type LData struct {
	Info        Info
	Control1    ControlField1
	Control2    ControlField2
	Source      IndividualAddr
	Destination uint16
	Data        TransportUnit
}

// Unpack initializes the structure by parsing the given data.
func (ldata *LData) Unpack(data []byte) (n uint, err error) {
	if n, err = util.UnpackSome(
		data,
		&ldata.Info,
		(*uint8)(&ldata.Control1),
		(*uint8)(&ldata.Control2),
		(*uint16)(&ldata.Source),
		&ldata.Destination,
	); err != nil {
		return
	}

	m, err := unpackTransportUnit(data[n:], &ldata.Data)
	n += m

	return
}

// Size returns the packed size.
func (ldata *LData) Size() uint {
	return ldata.Info.Size() + 6 + ldata.Data.Size()
}

// Pack the message body into the buffer.
func (ldata *LData) Pack(buffer []byte) {
	util.PackSome(
		buffer,
		ldata.Info,
		uint8(ldata.Control1),
		uint8(ldata.Control2),
		uint16(ldata.Source),
		ldata.Destination,
		ldata.Data,
	)
}

// A LDataReq represents a L_Data.req message body.
type LDataReq struct {
	LData
}

// MessageCode returns the message code for L_Data.req.
func (LDataReq) MessageCode() MessageCode {
	return LDataReqCode
}

// A LDataCon represents a L_Data.con message body.
type LDataCon struct {
	LData
}

// MessageCode returns the message code for L_Data.con.
func (LDataCon) MessageCode() MessageCode {
	return LDataConCode
}

// A LDataInd represents a L_Data.ind message body.
type LDataInd struct {
	LData
}

// MessageCode returns the message code for L_Data.ind.
func (LDataInd) MessageCode() MessageCode {
	return LDataIndCode
}
