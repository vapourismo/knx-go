// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

// DPT_1001 represents DPT 1.001 (G) / DPT_Switch.
type DPT_1001 bool

func (d DPT_1001) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1001) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1001) Unit() string {
	return ""
}

func (d DPT_1001) String() string {
	if d {
		return "On"
	} else {
		return "Off"
	}
}

// DPT_1002 represents DPT 1.002 (G) / DPT_Bool.
type DPT_1002 bool

func (d DPT_1002) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1002) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1002) Unit() string {
	return ""
}

func (d DPT_1002) String() string {
	if d {
		return "True"
	} else {
		return "False"
	}
}

// DPT_1003 represents DPT 1.003 (G) / DPT_Enable.
type DPT_1003 bool

func (d DPT_1003) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1003) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1003) Unit() string {
	return ""
}

func (d DPT_1003) String() string {
	if d {
		return "Enable"
	} else {
		return "Disable"
	}
}

// DPT_1004 represents DPT 1.004 (FB) / DPT_Ramp.
type DPT_1004 bool

func (d DPT_1004) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1004) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1004) Unit() string {
	return ""
}

func (d DPT_1004) String() string {
	if d {
		return "Ramp"
	} else {
		return "No ramp"
	}
}

// DPT_1005 represents DPT 1.005 (FB) / DPT_Alarm.
type DPT_1005 bool

func (d DPT_1005) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1005) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1005) Unit() string {
	return ""
}

func (d DPT_1005) String() string {
	if d {
		return "Alarm"
	} else {
		return "No alarm"
	}
}

// DPT_1006 represents DPT 1.006 (FB) / DPT_BinaryValue.
type DPT_1006 bool

func (d DPT_1006) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1006) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1006) Unit() string {
	return ""
}

func (d DPT_1006) String() string {
	if d {
		return "High"
	} else {
		return "Low"
	}
}

// DPT_1007 represents DPT 1.007 (FB) / DPT_Step.
type DPT_1007 bool

func (d DPT_1007) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1007) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1007) Unit() string {
	return ""
}

func (d DPT_1007) String() string {
	if d {
		return "Increase"
	} else {
		return "Decrease"
	}
}

// DPT_1008 represents DPT 1.008 (G) / DPT_UpDown.
type DPT_1008 bool

func (d DPT_1008) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1008) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1008) Unit() string {
	return ""
}

func (d DPT_1008) String() string {
	if d {
		return "Down"
	} else {
		return "Up"
	}
}

// DPT_1009 represents DPT 1.009 (G) / DPT_OpenClose.
type DPT_1009 bool

func (d DPT_1009) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1009) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1009) Unit() string {
	return ""
}

func (d DPT_1009) String() string {
	if d {
		return "Close"
	} else {
		return "Open"
	}
}

// DPT_1010 represents DPT 1.010 (G) / DPT_Start.
type DPT_1010 bool

func (d DPT_1010) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1010) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1010) Unit() string {
	return ""
}

func (d DPT_1010) String() string {
	if d {
		return "Start"
	} else {
		return "Stop"
	}
}

// DPT_1011 represents DPT 1.011 (FB) / DPT_State.
type DPT_1011 bool

func (d DPT_1011) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1011) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1011) Unit() string {
	return ""
}

func (d DPT_1011) String() string {
	if d {
		return "Active"
	} else {
		return "Inactive"
	}
}

// DPT_1012 represents DPT 1.012 (FB) / DPT_Invert.
type DPT_1012 bool

func (d DPT_1012) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1012) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1012) Unit() string {
	return ""
}

func (d DPT_1012) String() string {
	if d {
		return "Inverted"
	} else {
		return "Not inverted"
	}
}

// DPT_1013 represents DPT 1.013 (FB) / DPT_DimSendStyle.
type DPT_1013 bool

func (d DPT_1013) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1013) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1013) Unit() string {
	return ""
}

func (d DPT_1013) String() string {
	if d {
		return "Cyclically"
	} else {
		return "Start/stop"
	}
}

// DPT_1014 represents DPT 1.014 (FB) / DPT_InputSource.
type DPT_1014 bool

func (d DPT_1014) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1014) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1014) Unit() string {
	return ""
}

func (d DPT_1014) String() string {
	if d {
		return "Calculated"
	} else {
		return "Fixed"
	}
}

// DPT_1015 represents DPT 1.015 (G) / DPT_Reset.
type DPT_1015 bool

func (d DPT_1015) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1015) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1015) Unit() string {
	return ""
}

func (d DPT_1015) String() string {
	if d {
		return "reset command"
	} else {
		return "no action"
	}
}

// DPT_1016 represents DPT 1.016 (G) / DPT_Ack.
type DPT_1016 bool

func (d DPT_1016) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1016) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1016) Unit() string {
	return ""
}

func (d DPT_1016) String() string {
	if d {
		return "acknowledge command"
	} else {
		return "no action"
	}
}

// DPT_1017 represents DPT 1.017 (G) / DPT_Trigger.
type DPT_1017 bool

func (d DPT_1017) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1017) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1017) Unit() string {
	return ""
}

func (d DPT_1017) String() string {
	if d {
		return "trigger"
	} else {
		return "trigger"
	}
}

// DPT_1018 represents DPT 1.018 (G) / DPT_Occupancy.
type DPT_1018 bool

func (d DPT_1018) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1018) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1018) Unit() string {
	return ""
}

func (d DPT_1018) String() string {
	if d {
		return "occupied"
	} else {
		return "not occupied"
	}
}

// DPT_1019 represents DPT 1.019 (G) / DPT_Window_Door.
type DPT_1019 bool

func (d DPT_1019) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1019) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1019) Unit() string {
	return ""
}

func (d DPT_1019) String() string {
	if d {
		return "open"
	} else {
		return "closed"
	}
}

// DPT_1021 represents DPT 1.021 (FB) / DPT_LogicalFunction.
type DPT_1021 bool

func (d DPT_1021) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1021) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1021) Unit() string {
	return ""
}

func (d DPT_1021) String() string {
	if d {
		return "AND"
	} else {
		return "OR"
	}
}

// DPT_1022 represents DPT 1.022 (FB) / DPT_Scene_AB.
type DPT_1022 bool

func (d DPT_1022) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1022) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1022) Unit() string {
	return ""
}

func (d DPT_1022) String() string {
	if d {
		return "scene B"
	} else {
		return "scene A"
	}
}

// DPT_1023 represents DPT 1.023 (FB) / DPT_ShutterBlinds_Mode.
type DPT_1023 bool

func (d DPT_1023) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1023) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1023) Unit() string {
	return ""
}

func (d DPT_1023) String() string {
	if d {
		return "move Up/Down + StepStop mode"
	} else {
		return "only move Up/Down mode"
	}
}

// DPT_1024 represents DPT 1.024 (G) / DPT_DayNight.
type DPT_1024 bool

func (d DPT_1024) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1024) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1024) Unit() string {
	return ""
}

func (d DPT_1024) String() string {
	if d {
		return "Night"
	} else {
		return "Day"
	}
}

// DPT_1100 represents DPT 1.100 (FB) / DPT_Heat/Cool.
type DPT_1100 bool

func (d DPT_1100) Pack() []byte {
	return []byte{packB1(bool(d))}
}

func (d *DPT_1100) Unpack(data []byte) error {
	return unpackB1(data, (*bool)(d))
}

func (d DPT_1100) Unit() string {
	return ""
}

func (d DPT_1100) String() string {
	if d {
		return "heating"
	} else {
		return "cooling"
	}
}
