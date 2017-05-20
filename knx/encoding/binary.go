package encoding

import (
	"fmt"
	"io"
)

func writeUint8(w io.Writer, value uint8) (int64, error) {
	len, err := w.Write([]byte{value})
	return int64(len), err
}

func writeUint16(w io.Writer, value uint16) (int64, error) {
	len, err := w.Write([]byte{
		byte(value >> 8),
		byte(value),
	})
	return int64(len), err
}

func writeUint32(w io.Writer, value uint32) (int64, error) {
	len, err := w.Write([]byte{
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
	return int64(len), err
}

func writeUint64(w io.Writer, value uint64) (int64, error) {
	len, err := w.Write([]byte{
		byte(value >> 56),
		byte(value >> 48),
		byte(value >> 40),
		byte(value >> 32),
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
	return int64(len), err
}

func writeUint8Slice(w io.Writer, values []uint8) (int64, error) {
	len, err := w.Write(values)
	return int64(len), err
}

func writeUint16Slice(w io.Writer, values []uint16) (n int64, err error) {
	for _, value := range values {
		len, err := writeUint16(w, value)
		n += len
		if err != nil {
			return n, err
		}
	}

	return
}

func writeUint32Slice(w io.Writer, values []uint32) (n int64, err error) {
	for _, value := range values {
		len, err := writeUint32(w, value)
		n += len
		if err != nil {
			return n, err
		}
	}

	return
}

func writeUint64Slice(w io.Writer, values []uint64) (n int64, err error) {
	for _, value := range values {
		len, err := writeUint64(w, value)
		n += len
		if err != nil {
			return n, err
		}
	}

	return
}

func writeInt8Slice(w io.Writer, values []int8) (n int64, err error) {
	for _, value := range values {
		len, err := writeUint8(w, uint8(value))
		n += len
		if err != nil {
			return n, err
		}
	}

	return
}

func writeInt16Slice(w io.Writer, values []int16) (n int64, err error) {
	for _, value := range values {
		len, err := writeUint16(w, uint16(value))
		n += len
		if err != nil {
			return n, err
		}
	}

	return
}

func writeInt32Slice(w io.Writer, values []int32) (n int64, err error) {
		for _, value := range values {
		len, err := writeUint32(w, uint32(value))
		n += len
		if err != nil {
			return n, err
		}
	}

	return
}

func writeInt64Slice(w io.Writer, values []int64) (n int64, err error) {
	for _, value := range values {
		len, err := writeUint64(w, uint64(value))
		n += len
		if err != nil {
			return n, err
		}
	}

	return
}

// Write writes an item to the Writer.
func Write(w io.Writer, item interface{}) (int64, error) {
	switch item := item.(type) {
	case uint8:  return writeUint8(w, item)
	case uint16: return writeUint16(w, item)
	case uint32: return writeUint32(w, item)
	case uint64: return writeUint64(w, item)

	case int8:   return writeUint8(w, uint8(item))
	case int16:  return writeUint16(w, uint16(item))
	case int32:  return writeUint32(w, uint32(item))
	case int64:  return writeUint64(w, uint64(item))

	case *uint8:  return writeUint8(w, *item)
	case *uint16: return writeUint16(w, *item)
	case *uint32: return writeUint32(w, *item)
	case *uint64: return writeUint64(w, *item)

	case *int8:   return writeUint8(w, uint8(*item))
	case *int16:  return writeUint16(w, uint16(*item))
	case *int32:  return writeUint32(w, uint32(*item))
	case *int64:  return writeUint64(w, uint64(*item))

	case []uint8:  return writeUint8Slice(w, item)
	case []uint16: return writeUint16Slice(w, item)
	case []uint32: return writeUint32Slice(w, item)
	case []uint64: return writeUint64Slice(w, item)

	case []int8:   return writeInt8Slice(w, item)
	case []int16:  return writeInt16Slice(w, item)
	case []int32:  return writeInt32Slice(w, item)
	case []int64:  return writeInt64Slice(w, item)

	case io.WriterTo:
		return item.WriteTo(w)

	default:
		return 0, fmt.Errorf("Cannot write type %T", item)
	}
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

	return n, err
}

func readUint8(data []byte, ptr *uint8) (int64, error) {
	if len(data) < 1 {
		return 0, io.ErrUnexpectedEOF
	}

	*ptr = data[0]

	return 1, nil
}

func readUint16(data []byte, ptr *uint16) (int64, error) {
	if len(data) < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	*ptr = uint16(data[0]) << 8 | uint16(data[1])

	return 2, nil
}

func readUint32(data []byte, ptr *uint32) (int64, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}

	*ptr = uint32(data[0]) << 24 |
	       uint32(data[1]) << 16 |
	       uint32(data[2]) << 8 |
	       uint32(data[3])

	return 4, nil
}

func readUint64(data []byte, ptr *uint64) (int64, error) {
	if len(data) < 8 {
		return 0, io.ErrUnexpectedEOF
	}

	*ptr = uint64(data[0]) << 56 |
	       uint64(data[1]) << 48 |
	       uint64(data[2]) << 40 |
	       uint64(data[3]) << 32 |
	       uint64(data[4]) << 24 |
	       uint64(data[5]) << 16 |
	       uint64(data[6]) << 8 |
	       uint64(data[7])

	return 8, nil
}

func readInt8(data []byte, ptr *int8) (int64, error) {
	if len(data) < 1 {
		return 0, io.ErrUnexpectedEOF
	}

	*ptr = int8(data[0])

	return 1, nil
}

func readInt16(data []byte, ptr *int16) (int64, error) {
	if len(data) < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	tmp := uint16(data[0]) << 8 | uint16(data[1])
	*ptr = int16(tmp)

	return 2, nil
}

func readInt32(data []byte, ptr *int32) (int64, error) {
	if len(data) < 4 {
		return 0, io.ErrUnexpectedEOF
	}

	tmp := uint32(data[0]) << 24 |
	       uint32(data[1]) << 16 |
	       uint32(data[2]) << 8 |
	       uint32(data[3])
	*ptr = int32(tmp)

	return 4, nil
}

func readInt64(data []byte, ptr *int64) (int64, error) {
	if len(data) < 8 {
		return 0, io.ErrUnexpectedEOF
	}

	tmp := uint64(data[0]) << 56 |
	       uint64(data[1]) << 48 |
	       uint64(data[2]) << 40 |
	       uint64(data[3]) << 32 |
	       uint64(data[4]) << 24 |
	       uint64(data[5]) << 16 |
	       uint64(data[6]) << 8 |
	       uint64(data[7])
	*ptr = int64(tmp)

	return 8, nil
}

func readUint8Slice(data []byte, slice []uint8) (int64, error) {
	if len(data) < len(slice) {
		return 0, io.ErrUnexpectedEOF
	}

	return int64(copy(slice, data)), nil
}

func readUint16Slice(data []byte, slice []uint16) (n int64, err error) {
	for i := range slice {
		len, err := readUint16(data[n:], &slice[i])
		n += len
		if err != nil {
			return n, err
		}
	}

	return
}

func readUint32Slice(data []byte, slice []uint32) (n int64, err error) {
	for i := range slice {
		len, err := readUint32(data[n:], &slice[i])
		n += len
		if err != nil {
			return n, err
		}
	}

	return
}

func readUint64Slice(data []byte, slice []uint64) (n int64, err error) {
	for i := range slice {
		len, err := readUint64(data[n:], &slice[i])
		n += len
		if err != nil {
			return n, err
		}
	}

	return
}

func readInt8Slice(data []byte, slice []int8) (n int64, err error) {
	for i := range slice {
		len, err := readInt8(data[n:], &slice[i])
		n += len
		if err != nil {
			return n, err
		}
	}

	return
}

func readInt16Slice(data []byte, slice []int16) (n int64, err error) {
	for i := range slice {
		len, err := readInt16(data[n:], &slice[i])
		n += len
		if err != nil {
			return n, err
		}
	}

	return
}

func readInt32Slice(data []byte, slice []int32) (n int64, err error) {
	for i := range slice {
		len, err := readInt32(data[n:], &slice[i])
		n += len
		if err != nil {
			return n, err
		}
	}

	return
}

func readInt64Slice(data []byte, slice []int64) (n int64, err error) {
	for i := range slice {
		len, err := readInt64(data[n:], &slice[i])
		n += len
		if err != nil {
			return n, err
		}
	}

	return
}

// Unserialize is implemented by types that can be initialized by reading from a byte slice.
type Unserialize interface {
	ReadFrom(data []byte) (int64, error)
}

// Read parses the given data.
func Read(data []byte, item interface{}) (int64, error) {
	switch item := item.(type) {
	case *uint8:  return readUint8(data, item)
	case *uint16: return readUint16(data, item)
	case *uint32: return readUint32(data, item)
	case *uint64: return readUint64(data, item)
	case *int8:   return readInt8(data, item)
	case *int16:  return readInt16(data, item)
	case *int32:  return readInt32(data, item)
	case *int64:  return readInt64(data, item)

	case []uint8:  return readUint8Slice(data, item)
	case []uint16: return readUint16Slice(data, item)
	case []uint32: return readUint32Slice(data, item)
	case []uint64: return readUint64Slice(data, item)
	case []int8:   return readInt8Slice(data, item)
	case []int16:  return readInt16Slice(data, item)
	case []int32:  return readInt32Slice(data, item)
	case []int64:  return readInt64Slice(data, item)

	case Unserialize:
		return item.ReadFrom(data)

	default:
		return 0, fmt.Errorf("Cannot read type %T", item)
	}
}

// ReadSome reads multiple items from the data.
func ReadSome(data []byte, items ...interface{}) (n int64, err error) {
	for item := range items {
		len, err := Read(data[n:], item)
		n += len
		if err != nil {
			return n, err
		}
	}

	return
}
