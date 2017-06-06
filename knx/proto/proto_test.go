package proto

import (
	"bytes"
	"testing"

	"github.com/vapourismo/knx-go/knx/cemi"
)

func BenchmarkPack1(b *testing.B) {
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
		buffer := bytes.NewBuffer(make([]byte, 0, 32))

		num, err := Pack(buffer, req)

		if err != nil {
			b.Fatal(err)
		}

		if num != int64(len(buffer.Bytes())) {
			b.Fatal("Length mismatch", num, len(buffer.Bytes()))
		}
	}
}

type bytesWriter []byte

func (w *bytesWriter) Write(data []byte) (int, error) {
	*w = append(*w, data...)
	return len(data), nil
}

func BenchmarkPack2(b *testing.B) {
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
		buffer := bytesWriter(make([]byte, 0, 32))

		num, err := Pack(&buffer, req)

		if err != nil {
			b.Fatal(err)
		}

		if num != int64(len(buffer)) {
			b.Fatal("Length mismatch", num, len(buffer))
		}
	}
}
