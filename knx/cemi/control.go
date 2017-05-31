package cemi

// A Priority determines the priority.
type Priority uint8

// These are known priorities.
const (
	PrioritySystem Priority = 0
	PriorityNormal Priority = 1
	PriorityUrgent Priority = 2
	PriorityLow    Priority = 3
)

// ControlField1 contains various control information.
type ControlField1 uint8

// MakeControlField1 generates a control field 1 value.
func MakeControlField1(
	stdFrame bool,
	isRepeated bool,
	sysBroadcast bool,
	prio Priority,
	wantAck bool,
	isErr bool,
) (ret ControlField1) {
	if stdFrame {
		ret |= 1 << 7
	}

	if !isRepeated {
		ret |= 1 << 5
	}

	if !sysBroadcast {
		ret |= 1 << 4
	}

	ret |= ControlField1(prio&3) << 2

	if wantAck {
		ret |= 1 << 1
	}

	if isErr {
		ret |= 1
	}

	return
}

// ControlField2 contains various control information.
type ControlField2 uint8

// MakeControlField2 generates a control field 2 value.
func MakeControlField2(isGroupAddr bool, hopCount uint8, frameFormat uint8) (ret ControlField2) {
	if isGroupAddr {
		ret |= 1 << 7
	}

	ret |= ControlField2(hopCount&7) << 4
	ret |= ControlField2(frameFormat) & 15

	return
}
