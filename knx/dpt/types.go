// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import "fmt"

// A DatapointValue is a value of a datapoint.
type DatapointValue interface {
	// Pack the datapoint to a byte array.
	Pack() []byte

	// Unpack a the datapoint value from a byte array.
	Unpack(data []byte) error
}

// DatapointMeta gives meta information about a datapoint type.
type DatapointMeta interface {
	// Unit returns the unit of this datapoint type or empty string if it doesn't have a unit.
	Unit() string

	// fmt.Stringer provides a string representation of the datapoint.
	fmt.Stringer
}

// Datapoint represents a datapoint with both its value and metadata.
type Datapoint interface {
	DatapointValue
	DatapointMeta
}
