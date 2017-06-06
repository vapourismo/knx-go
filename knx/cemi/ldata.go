package cemi

import (
	"errors"
	"io"

	"github.com/vapourismo/knx-go/knx/encoding"
	"github.com/vapourismo/knx-go/knx/util"
)

// A LData is a link-layer data frame. L_Data.req, L_Data.con and L_Data.ind share this structure.
type LData struct {
	Info        Info
	Control1    ControlField1
	Control2    ControlField2
	Source      uint16
	Destination uint16
	Data        TPDU
}

// Unpack initializes the structure by parsing the given data.
func (ldata *LData) Unpack(data []byte) (n uint, err error) {
	return util.UnpackSome(
		data,
		&ldata.Info,
		(*uint8)(&ldata.Control1),
		(*uint8)(&ldata.Control2),
		(*uint16)(&ldata.Source),
		(*uint16)(&ldata.Destination),
		&ldata.Data,
	)
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
		ldata.Info,
		ldata.Control1,
		ldata.Control2,
		ldata.Source,
		ldata.Destination,
		byte(len(ldata.Data)-1),
		ldata.Data,
	)
}

// Copy returns a copy of the LData frame where all struct members are independent from the method
// receiver's members.
func (ldata LData) Copy() LData {
	data := make([]byte, len(ldata.Data))
	copy(data, ldata.Data)
	ldata.Data = data

	return ldata
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
	return LDataConCode
}
