package dpt

import (
	"fmt"
)

// DPT_232600 represents DPT 232.600 (G) / DPT_Colour_RGB.
// Colour RGB - RGB value 4x(0..100%) / U8 U8 U8
type DPT_232600 struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func (d DPT_232600) Pack() []byte {
	return []byte{0, d.Red, d.Green, d.Blue}
}

func (d *DPT_232600) Unpack(data []byte) error {

	if len(data) != 4 {
		return ErrInvalidLength
	}

	d.Red = uint8(data[1])
	d.Green = uint8(data[2])
	d.Blue = uint8(data[3])

	return nil
}

func (d DPT_232600) Unit() string {
	return ""
}

func (d DPT_232600) String() string {
	return fmt.Sprintf("#%02X%02X%02X", d.Red, d.Green, d.Blue)
}
