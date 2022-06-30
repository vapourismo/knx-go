// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"fmt"
	"testing"
)

// Test DPT 1.xxx (B₁)
func TestDPT_1(t *testing.T) {
	type DPT1 struct {
		Dpv     DatapointValue
		OnFalse string
		OnTrue  string
	}

	var types_1 = []DPT1{
		{new(DPT_1001), "Off", "On"},
		{new(DPT_1002), "False", "True"},
		{new(DPT_1003), "Disable", "Enable"},
		{new(DPT_1004), "No ramp", "Ramp"},
		{new(DPT_1005), "No alarm", "Alarm"},
		{new(DPT_1006), "Low", "High"},
		{new(DPT_1007), "Decrease", "Increase"},
		{new(DPT_1008), "Up", "Down"},
		{new(DPT_1009), "Open", "Close"},
		{new(DPT_1010), "Stop", "Start"},
		{new(DPT_1011), "Inactive", "Active"},
		{new(DPT_1012), "Not inverted", "Inverted"},
		{new(DPT_1013), "Start/stop", "Cyclically"},
		{new(DPT_1014), "Fixed", "Calculated"},
		{new(DPT_1015), "no action", "reset command"},
		{new(DPT_1016), "no action", "acknowledge command"},
		{new(DPT_1017), "trigger", "trigger"},
		{new(DPT_1018), "not occupied", "occupied"},
		{new(DPT_1019), "closed", "open"},
		{new(DPT_1021), "OR", "AND"},
		{new(DPT_1022), "scene A", "scene B"},
		{new(DPT_1023), "only move Up/Down mode", "move Up/Down + StepStop mode"},
		{new(DPT_1024), "Day", "Night"},
		{new(DPT_1100), "cooling", "heating"}}

	for _, e := range types_1 {
		src := e.Dpv
		if fmt.Sprintf("%s", src) != e.OnFalse {
			t.Errorf("%#v has wrong default value [%v]. Should be [%s].", e.Dpv, e.Dpv, e.OnFalse)
		}

		e.Dpv.Unpack(packB1(false))
		if fmt.Sprintf("%s", e.Dpv) != e.OnFalse {
			t.Errorf("%#v has wrong false value [%v]. Should be [%s].", e.Dpv, e.Dpv, e.OnFalse)
		}

		e.Dpv.Unpack(packB1(true))
		if fmt.Sprintf("%s", e.Dpv) != e.OnTrue {
			t.Errorf("%#v has wrong true value [%v]. Should be [%s].", e.Dpv, e.Dpv, e.OnTrue)
		}
	}
}
