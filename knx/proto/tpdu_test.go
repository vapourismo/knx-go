package proto

import (
	"bytes"
	"io"
	"math/rand"
	"testing"
)

func tpduEquals(lhs TPDU, rhs TPDU) bool {
	return (lhs.PacketType == rhs.PacketType &&
	        lhs.SeqNumber == rhs.SeqNumber &&
	        lhs.Control == rhs.Control &&
	        lhs.Info == rhs.Info &&
	        bytes.Equal(lhs.Data, rhs.Data))
}

func TestTPDU_ReadFrom(t *testing.T) {
	t.Run("ReadHeadFails", func (t *testing.T) {
		r := bytes.NewReader([]byte{})
		err := (&TPDU{}).ReadFrom(r)
		if err != io.EOF {
			t.Fatalf("Unexpected error %v", err)
		}
	})

	t.Run("WriteToConfirm", func (t *testing.T) {
		for i := 0; i < 100; i++ {
			buffer := make([]byte, 3)
			rand.Read(buffer)

			r := bytes.NewReader(buffer)
			tpdu := TPDU{}

			err := tpdu.ReadFrom(r)
			if err != nil {
				t.Errorf("ReadFrom error: %v", err)
				continue
			}

			w := &bytes.Buffer{}
			err = tpdu.WriteTo(w)
			if err != nil {
				t.Errorf("%+v WriteTo error: %v", tpdu, err)
				continue
			}

			switch tpdu.PacketType {
			case UnnumberedControlPacket, NumberedControlPacket:
				buffer = buffer[:1]
			}

			if !bytes.Equal(buffer, w.Bytes()) {
				t.Errorf("Mismatch for input %v with value %v and output %v", buffer, tpdu, w.Bytes())
			}
		}
	})
}