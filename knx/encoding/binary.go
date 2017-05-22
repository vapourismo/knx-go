package encoding

import (
	"encoding/binary"
	"io"
)

// Write writes an item to the Writer.
func Write(w io.Writer, item interface{}) (int64, error) {
	if wt, ok := item.(io.WriterTo); ok {
		return wt.WriteTo(w)
	}

	if err := binary.Write(w, binary.BigEndian, item); err != nil {
		return 0, err
	}

	return int64(binary.Size(item)), nil
}

// WriteSome writes multiple items to the Writer.
func WriteSome(w io.Writer, items ...interface{}) (n int64, err error) {
	for _, item := range items {
		len, err := Write(w, item)
		n += len
		if err != nil {
			return n, err
		}
	}

	return
}

// Read reads an item from the Reader.
func Read(r io.Reader, item interface{}) (int64, error) {
	if rf, ok := item.(io.ReaderFrom); ok {
		return rf.ReadFrom(r)
	}

	if err := binary.Read(r, binary.BigEndian, item); err != nil {
		return 0, err
	}

	return int64(binary.Size(item)), nil
}

// ReadSome reads multiple items from the Reader.
func ReadSome(r io.Reader, items ...interface{}) (n int64, err error) {
	for _, item := range items {
		len, err := Read(r, item)
		n += len
		if err != nil {
			return n, err
		}
	}

	return
}
