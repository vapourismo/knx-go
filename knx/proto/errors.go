package proto

import (
	"errors"
)

// General purpose errors
var (
	ErrDataTooShort = errors.New("Given input data is too short")
)