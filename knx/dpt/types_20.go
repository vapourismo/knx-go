package dpt

// DPT_20102 represents DPT 20.102 / HVAC Mode.
type DPT_20102 uint8

func (d DPT_20102) Pack() []byte {
	return packU8(uint8(d))
}

func (d *DPT_20102) Unpack(data []byte) error {
	return unpackU8(data, (*uint8)(d))
}

func (d DPT_20102) Unit() string {
	return ""
}

func (d DPT_20102) String() string {
	switch d {
	case 0:
		return "Auto"
	case 1:
		return "Comfort"
	case 2:
		return "Standby"
	case 3:
		return "Economy"
	case 4:
		return "Building Protection"
	default:
		return ""
	}
}
