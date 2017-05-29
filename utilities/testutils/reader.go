package testutils

import (
	"errors"
)

type BadReader struct {}

var ErrBadRead = errors.New("Bad read")

func (BadReader) Read([]byte) (int, error) {
	return 0, ErrBadRead
}