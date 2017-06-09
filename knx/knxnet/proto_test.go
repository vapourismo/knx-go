// Copyright 2017 Ole Krüger.

package knxnet

import (
	"testing"

	"github.com/vapourismo/knx-go/knx/cemi"
	"github.com/vapourismo/knx-go/knx/util"
)

func BenchmarkPack(b *testing.B) {
	b.ReportAllocs()

	req := &TunnelReq{
		Channel:   1,
		SeqNumber: 0,
		Payload: &cemi.LDataReq{
			LData: cemi.LData{
				Control1: cemi.Control1NoRepeat | cemi.Control1NoSysBroadcast |
					cemi.Control1StdFrame | cemi.Control1WantAck | cemi.Control1Prio(cemi.PrioLow),
				Control2:    cemi.Control2GrpAddr | cemi.Control2Hops(6),
				Source:      0,
				Destination: 0x1337,
				Data: &cemi.AppData{
					Command: cemi.GroupValueWrite,
					Data:    []byte{0, 0x13, 0x37},
				},
			},
		},
	}

	for i := 0; i < b.N; i++ {
		util.AllocAndPack(req)
	}
}
