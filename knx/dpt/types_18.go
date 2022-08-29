package dpt

import "fmt"

// DPT_18001 represents DPT 18.001 / Scene Control.
type DPT_18001 uint8

func (d DPT_18001) Pack() []byte {
	if d <= 63 || (d >= 128 && d <= 191) {
		return packU8(uint8(d))
	} else {
		return packU8(63)
	}
}

func (d *DPT_18001) Unpack(data []byte) error {
	var value uint8

	if err := unpackU8(data, &value); err != nil {
		return err
	}

	if value <= 63 || (value >= 128 && value <= 191) {
		*d = DPT_18001(value)
		return nil
	} else {
		*d = DPT_18001(63)
		return nil
	}
}

func (d DPT_18001) Unit() string {
	return ""
}

// KNX Association recommends to display the scene numbers [1..64].
// See note 6 of the KNX Specifications v2.1.
func (d DPT_18001) String() string {
	return fmt.Sprintf("%d", uint8(d))
}
