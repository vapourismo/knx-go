package dpt

import (
	"fmt"
)

// DPT_251600 represents DPT 251.600 / Colour RGBW - RGBW value 4x(0..100%) / U8 U8 U8 U8 r8 r4B4
type DPT_251600 struct {
	Red        uint8
	Green      uint8
	Blue       uint8
	White      uint8
	RedValid   bool
	GreenValid bool
	BlueValid  bool
	WhiteValid bool
}

func (d DPT_251600) Pack() []byte {

	validBits := packB4([4]bool{d.WhiteValid, d.BlueValid, d.GreenValid, d.RedValid})

	return []byte{0, d.Red, d.Green, d.Blue, d.White, uint8(0), validBits}
}

func (d *DPT_251600) Unpack(data []byte) error {

	if len(data) != 7 {
		return ErrInvalidLength
	}

	var redValid, greenValid, blueValid, whiteValid bool

	err := unpackB4(data[6], &whiteValid, &blueValid, &greenValid, &redValid)

	if err != nil {
		return ErrInvalidLength
	}

	*d = DPT_251600{
		Red:        uint8(data[1]),
		Green:      uint8(data[2]),
		Blue:       uint8(data[3]),
		White:      uint8(data[4]),
		RedValid:   redValid,
		GreenValid: greenValid,
		BlueValid:  blueValid,
		WhiteValid: whiteValid,
	}
	return nil
}

func (d DPT_251600) Unit() string {
	return ""
}

func (d DPT_251600) String() string {
	return fmt.Sprintf("Red: %d Green: %d Blue: %d White: %d RedValid: %t, GreenValid: %t, BlueValid: %t, WhiteValid: %t", d.Red, d.Green, d.Blue, d.White, d.RedValid, d.GreenValid, d.BlueValid, d.WhiteValid)
}
