// Copyright 2017 Ole Krüger.

package knx

// GroupComm represents a group communication.
type GroupComm struct {
	Source      uint16
	Destination uint16
	Data        []byte
}

// A GroupClient is a KNX client which supports group communication.
type GroupClient interface {
	Send(comm GroupComm) error
	Inbound() <-chan GroupComm
}
