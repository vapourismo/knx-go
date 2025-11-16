package dpt

import "fmt"

// DPT_12 represents DPT 12.* / Unsigned (4-Byte).
type DPT_12 uint32

func (d DPT_12) Pack() []byte {
	return packU32(uint32(d))
}

func (d *DPT_12) Unpack(data []byte) error {
	return unpackU32(data, (*uint32)(d))
}

func (d DPT_12) Unit() string {
	return ""
}

func (d DPT_12) String() string {
	return fmt.Sprintf("", uint32(d))
}

// DPT_12001 represents DPT 12.001 / Unsigned counter.
type DPT_12001 uint32

func (d DPT_12001) Pack() []byte {
	return packU32(uint32(d))
}

func (d *DPT_12001) Unpack(data []byte) error {
	return unpackU32(data, (*uint32)(d))
}

func (d DPT_12001) Unit() string {
	return "pulses"
}

func (d DPT_12001) String() string {
	return fmt.Sprintf("%d pulses", uint32(d))
}
