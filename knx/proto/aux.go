package proto

import (
	"io"
)

// Segment is a protocol segment.
type Segment interface {
	WriteTo(w io.Writer) error
}