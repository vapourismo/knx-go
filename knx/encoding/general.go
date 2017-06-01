package encoding

import (
	"io"
)

// ReadAll reads everything that is available.
func ReadAll(r io.Reader) (n int64, result []byte) {
	buffer := [1024]byte{}

	for {
		m, err := r.Read(buffer[:])
		n += int64(m)

		if err != nil {
			return n, result
		}

		result = append(result, buffer[:m]...)
	}
}
