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

// DPT_20105 represents DPT 20.105 / HVACContrMode.
type DPT_20105 uint8

func (d DPT_20105) Pack() []byte {
	return packU8(uint8(d))
}

func (d *DPT_20105) Unpack(data []byte) error {
	var value uint8

	if err := unpackU8(data, &value); err != nil {
		return err
	}

	*d = DPT_20105(value)
	return nil
}

func (d DPT_20105) Unit() string {
	return ""
}

func (d DPT_20105) String() string {
	switch d {
	case 0:
		return "Auto"
	case 1:
		return "Heat"
	case 2:
		return "Morning Warmup"
	case 3:
		return "Cool"
	case 4:
		return "Night Purge"
	case 5:
		return "Precool"
	case 6:
		return "Off"
	case 7:
		return "Test"
	case 8:
		return "Emergency Heat"
	case 9:
		return "Fan only"
	case 10:
		return "Free Cool"
	case 11:
		return "Ice"
	case 12:
		return "Maximum Heating Mode"
	case 13:
		return "Economic Heat/Cool Mode"
	case 14:
		return "Dehumidification"
	case 15:
		return "Calibration Mode"
	case 16:
		return "Emergency Cool Mode"
	case 17:
		return "Emergency Steam Mode"
	case 20:
		return "NoDem"
	}
	return "reserved"
}
