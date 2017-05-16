package proto

// import (
// 	"bytes"
// 	"testing"
// )

// func tpduEquals(lhs *TPDU, rhs *TPDU) bool {
// 	if lhs == rhs {
// 		return true
// 	}

// 	if lhs == nil || rhs == nil {
// 		return false
// 	}

// 	return (lhs.PacketType == rhs.PacketType &&
// 		    lhs.SeqNumber == rhs.SeqNumber &&
// 	        lhs.Control == rhs.Control &&
// 	        lhs.Info == rhs.Info &&
// 	        bytes.Equal(lhs.Data, rhs.Data))
// }

// type TPDUCase struct {
// 	in  []byte
// 	res *TPDU
// 	err error
// }

// func TestReadTPDU(t *testing.T) {
// 	cases := []TPDUCase{
// 		{
// 			[]byte{},
// 			nil,
// 			ErrTransportUnitTooShort,
// 		},
// 		{
// 			[]byte{0},
// 			nil,
// 			ErrTransportUnitTooShort,
// 		},
// 		{
// 			[]byte{195},
// 			&TPDU{3, 0, 3, 0, nil},
// 			nil,
// 		},
// 		{
// 			[]byte{255},
// 			&TPDU{3, 15, 3, 0, nil},
// 			nil,
// 		},
// 		{
// 			[]byte{63, 192},
// 			&TPDU{0, 15, 0, 15, []byte{0}},
// 			nil,
// 		},
// 		{
// 			[]byte{63, 255},
// 			&TPDU{0, 15, 0, 15, []byte{63}},
// 			nil,
// 		},
// 		{
// 			[]byte{63, 192, 255},
// 			&TPDU{0, 15, 0, 15, []byte{255}},
// 			nil,
// 		},
// 		{
// 			[]byte{63, 192, 0},
// 			&TPDU{0, 15, 0, 15, []byte{0}},
// 			nil,
// 		},
// 	}

// 	for _, c := range cases {
// 		res, err := ReadTPDU(c.in)
// 		if err != c.err {
// 			t.Errorf("Expected error %v, got %v", c.err, err)
// 		}

// 		if !tpduEquals(c.res, res) {
// 			t.Fatalf("Expected result %v, got %v", c.res, res)
// 		}
// 	}
// }
