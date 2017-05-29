package testutils

import (
	"errors"
)

type BadWriter struct{}

var ErrBadWrite = errors.New("Bad write")

func (BadWriter) Write([]byte) (int, error) {
	return 0, ErrBadWrite
}
