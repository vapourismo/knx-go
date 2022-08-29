package dpt

import "fmt"

// DPT_6010 represents DPT 6.010 / counter pulses (-128..127).
type DPT_6010 int8

func (d DPT_6010) Pack() []byte {
	return packV8(int8(d))
}

func (d *DPT_6010) Unpack(data []byte) error {
	return unpackV8(data, (*int8)(d))
}

func (d DPT_6010) Unit() string {
	return "counter pulses"
}

func (d DPT_6010) String() string {
	return fmt.Sprintf("%d counter pulses", d)
}
