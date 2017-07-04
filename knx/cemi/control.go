// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package cemi

// A Priority determines the priority.
type Priority uint8

const (
	// PrioSystem is for high priority, configuration and management frames.
	PrioSystem Priority = 0

	// PrioNormal indicates normal priority. Ideal for short frames.
	PrioNormal Priority = 1

	// PrioUrgent indicates urgent priority.
	PrioUrgent Priority = 2

	// PrioLow indicates low priority. Ideal for long frames.
	PrioLow Priority = 3
)

// ControlField1 contains various control information.
type ControlField1 uint8

const (
	// Control1StdFrame indicate that the frame is not an extended frame. Extended frames contain
	// application data units greater than 15 bytes.
	Control1StdFrame ControlField1 = 1 << 7

	// Control1NoRepeat causes a repeated frame not to be sent on the medium. If you send two
	// identical frames, than one of them will not be sent on the medium when this flag is present.
	Control1NoRepeat ControlField1 = 1 << 5

	// Control1NoSysBroadcast causes the frame to be transmitted in normal broadcast mode, instead
	// of system broadcast mode.
	Control1NoSysBroadcast ControlField1 = 1 << 4

	// Control1WantAck requests an acknowledgement. Only works for L_Data.req.
	Control1WantAck ControlField1 = 1 << 1

	// Control1HasError indicates an error. Only relevant in L_Data.con.
	Control1HasError ControlField1 = 1
)

// Control1Prio generates the control field 1 flag for the given priority.
func Control1Prio(prio Priority) ControlField1 {
	return ControlField1(prio&3) << 2
}

// ControlField2 contains various control information.
type ControlField2 uint8

// IsGroupAddr determines if the destination address is a group address.
func (ctrl2 ControlField2) IsGroupAddr() bool {
	return ctrl2&Control2GroupAddr == Control2GroupAddr
}

// Hops retrieves the number of hops.
func (ctrl2 ControlField2) Hops() uint8 {
	return uint8(ctrl2>>7) & 7
}

const (
	// Control2GroupAddr determines that the destination address inside the frame is a group address,
	// instead of an individual address.
	Control2GroupAddr ControlField2 = 1 << 7

	// Control2LTEFrame indicates that the frame is a LTE-frame.
	Control2LTEFrame ControlField2 = 1 << 2
)

// Control2Hops generates the control field 2 flag for the given number of hops.
func Control2Hops(hops uint8) ControlField2 {
	if hops > 7 {
		hops = 7
	}

	return ControlField2(hops&7) << 4
}
