// Copyright 2022 Martin Müller.
// Licensed under the MIT license which can be found in the LICENSE file.

package cemi

import (
	"testing"
)

// Test Individual Addresses
func Test_IndividualAddresses(t *testing.T) {
	type Addr struct {
		Src     string
		Valid   bool
		Printed string
	}

	var addrs = []Addr{
		{"1.2.3", true, "1.2.3"},
		{"1.3.255", true, "1.3.255"},
		{"1.3.0", true, "1.3.0"},
		{"75.235", true, "4.11.235"},
		{"65535", true, "15.15.255"},
		{"15.15.255", true, "15.15.255"},
		{"15.15.0", true, "15.15.0"},
		{"13057", true, "3.3.1"},
		{"16.17.255", false, ""},
		{"1..0", false, ""},
		{"15.15.", false, ""},
		{" . .15", false, ""},
		{"18.15.240", false, ""},
		{"1.3.450", false, ""},
		{"1.450", false, ""},
		{"-2", false, ""},
		{".400", false, ""},
		{"-11.0.0", false, ""},
		{"0.0.0", false, ""},
		{"0.0", false, ""},
		{"0", false, ""},
	}

	for _, a := range addrs {
		ia, err := NewIndividualAddrString(a.Src)
		if a.Valid {
			if err != nil {
				t.Errorf("%#v has error %s.", a.Src, err)
			} else if ia.String() != a.Printed {
				t.Errorf("%#v wrongly parsed.", a.Src)
			}
		} else if err == nil {
			t.Errorf("%#v invalid parsed.", a.Src)
		}
	}
}

// Test Group Addresses
func Test_GroupAddresses(t *testing.T) {
	type Addr struct {
		Src     string
		Valid   bool
		Printed string
	}

	var addrs = []Addr{
		{"1/2/3", true, "1/2/3"},
		{"31/7/255", true, "31/7/255"},
		{"31/2040", true, "3/7/248"},
		{"65535", true, "31/7/255"},
		{"82/8/260", false, ""},
		{"84/230", false, ""},
		{"31/2060", false, ""},
		{"0/0/0", false, ""},
		{"0/0", false, ""},
		{"0", false, ""},
		{"123/foobar", false, ""},
		{"1000/2000/3000", false, ""},
	}

	for _, a := range addrs {
		ga, err := NewGroupAddrString(a.Src)
		if a.Valid {
			if err != nil {
				t.Errorf("%#v has error %s.", a.Src, err)
			} else if ga.String() != a.Printed {
				t.Errorf("%#v wrongly parsed.", a.Src)
			}
		} else if err == nil {
			t.Errorf("%#v invalid parsed.", a.Src)
		}
	}
}
