// Copyright 2017 Ole Krüger.

package knx

import (
	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/util"
)

// GroupComm represents a group communication.
type GroupComm struct {
	Source      cemi.IndividualAddr
	Destination cemi.GroupAddr
	Data        []byte
}

// A GroupClient is a KNX client which supports group communication.
type GroupClient interface {
	Send(src cemi.IndividualAddr, dest cemi.GroupAddr, data []byte) error
	Inbound() <-chan GroupComm
}

// serveGroupInbound serves a group communication.
func serveGroupInbound(inbound <-chan cemi.Message, outbound chan<- GroupComm) {
	util.Log(inbound, "Started worker")
	defer util.Log(inbound, "Worker exited")

	for msg := range inbound {
		if ind, ok := msg.(*cemi.LDataInd); ok {
			if app, ok := ind.Data.(*cemi.AppData); ok && (app.Command == cemi.GroupValueResponse || app.Command == cemi.GroupValueWrite) {
				outbound <- GroupComm{
					Source:      ind.Source,
					Destination: cemi.GroupAddr(ind.Destination),
					Data:        app.Data,
				}
			} else {
				util.Log(inbound, "Received L_Data.ind frame does not contain application data")
			}
		} else {
			util.Log(inbound, "Received frame is not a L_Data.ind frame")
		}
	}

	close(outbound)
}

var defaultGroupLData = cemi.LData{
	Control1: cemi.Control1NoRepeat | cemi.Control1NoSysBroadcast | cemi.Control1WantAck | cemi.Control1Prio(cemi.PrioLow),
	Control2: cemi.Control2GrpAddr | cemi.Control2Hops(6),
}

// buildGroupOutbound constructs the L_Data core frame for group communication.
func buildGroupOutbound(src cemi.IndividualAddr, dest cemi.GroupAddr, data []byte) cemi.LData {
	ldata := defaultGroupLData
	ldata.Data = &cemi.AppData{Command: cemi.GroupValueWrite, Data: data}
	ldata.Source = src
	ldata.Destination = uint16(dest)

	if len(data) <= 15 {
		ldata.Control1 |= cemi.Control1StdFrame
	}

	return ldata
}
