package knx

import (
	"bytes"
	"encoding/binary"
)

func writeSequence(w *bytes.Buffer, items ...interface{}) error {
	for _, item := range items {
		err := binary.Write(w, binary.BigEndian, item)
		if err != nil { return err }
	}

	return nil
}

func readSequence(r *bytes.Reader, items ...interface{}) error {
	for _, item := range items {
		err := binary.Read(r, binary.BigEndian, item)
		if err != nil { return err }
	}

	return nil
}
