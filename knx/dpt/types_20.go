package dpt

const (
	HVACMode_Auto DPT_20102 = iota
	HVACMode_Comfort
	HVACMode_Standby
	HVACMode_Economy
	HVACMode_BuildingProtection
)

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
	case HVACMode_Auto:
		return "Auto"
	case HVACMode_Comfort:
		return "Comfort"
	case HVACMode_Standby:
		return "Standby"
	case HVACMode_Economy:
		return "Economy"
	case HVACMode_BuildingProtection:
		return "Building Protection"
	default:
		return "reserved"
	}
}
