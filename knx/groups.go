// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package knx

import (
	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/util"
)

// GroupCommand determines the meaning of a group event.
type GroupCommand uint8

// String generates a string representation of the command.
func (cmd GroupCommand) String() string {
	switch cmd {
	case GroupRead:
		return "Read"

	case GroupResponse:
		return "Response"

	case GroupWrite:
		return "Write"
	}

	return "Unknown"
}

// These are known group commands.
const (
	GroupRead     GroupCommand = 0
	GroupResponse GroupCommand = 1
	GroupWrite    GroupCommand = 2
)

// GroupEvent represents a group communication event.
type GroupEvent struct {
	Command     GroupCommand
	Source      cemi.IndividualAddr
	Destination cemi.GroupAddr
	Data        []byte
}

// A GroupClient is a KNX client which supports group communication.
type GroupClient interface {
	Send(event GroupEvent) error
	Inbound() <-chan GroupEvent
}

// serveGroupInbound serves a group communication.
func serveGroupInbound(inbound <-chan cemi.Message, outbound chan<- GroupEvent) {
	util.Log(inbound, "Started worker")
	defer util.Log(inbound, "Worker exited")

	for msg := range inbound {
		if ind, ok := msg.(*cemi.LDataInd); ok {
			// Filter indications that do not target group addresses.
			if !ind.Control2.IsGroupAddr() {
				util.Log(inbound, "Received L_Data.ind does not target a group address")
				continue
			}

			if app, ok := ind.Data.(*cemi.AppData); ok && app.Command.IsGroupCommand() {
				outbound <- GroupEvent{
					Command:     GroupCommand(app.Command),
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
	Control2: cemi.Control2GroupAddr | cemi.Control2Hops(6),
}

// buildGroupOutbound constructs the L_Data core frame for group communication.
func buildGroupOutbound(event GroupEvent) cemi.LData {
	ldata := defaultGroupLData
	ldata.Data = &cemi.AppData{
		Command: cemi.APCI(event.Command),
		Data:    event.Data,
	}
	ldata.Source = event.Source
	ldata.Destination = uint16(event.Destination)

	if len(event.Data) <= 15 {
		ldata.Control1 |= cemi.Control1StdFrame
	}

	return ldata
}
