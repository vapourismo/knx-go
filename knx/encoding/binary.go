package encoding

import (
	"encoding/binary"
	"io"
)

// WriteSequence writes items to the Writer.
func WriteSequence(w io.Writer, items ...interface{}) error {
	for _, item := range items {
		err := binary.Write(w, binary.BigEndian, item)
		if err != nil { return err }
	}

	return nil
}

// ReadSequence reads items from the Reader.
func ReadSequence(r io.Reader, items ...interface{}) error {
	for _, item := range items {
		err := binary.Read(r, binary.BigEndian, item)
		if err != nil { return err }
	}

	return nil
}

// UInt16 extracts a uint16 from the given data.
func UInt16(data []byte) uint16 {
	return binary.BigEndian.Uint16(data)
}