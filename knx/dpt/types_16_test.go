// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"fmt"
	"testing"
)

// Test DPT 16.000 with special input
func TestDPT_16000(t *testing.T) {
	var dst DPT_16000

	// Test packing with simple string
	src := DPT_16000("KNX is OK")
	buf := src.Pack()
	r := fmt.Sprintf("%x", buf)
	if r != "004b4e58206973204f4b0000000000" {
		t.Errorf("Error packing [%s] => [%x]", src, r)
	}

	err := dst.Unpack(buf)
	if err != nil {
		t.Errorf("Error unpacking [%s]", src)
	}

	// Test with a string that will get truncated
	src = DPT_16000("This is a valid string longer than 14 chars, it will be truncated")
	buf = src.Pack()
	err = dst.Unpack(buf)

	if err != nil {
		t.Errorf("Error unpacking [%s] => [%s]", src, dst)
	}

	if dst.String() != "This is a vali" {
		t.Errorf("Unexpected unpacking value [%s] => [%s]", src, dst)
	}

	// Test with empty string
	src = DPT_16000("")
	buf = src.Pack()
	err = dst.Unpack(buf)

	if err != nil {
		t.Errorf("Unpacking error %s [%s] => [%s]", err, src, dst)
	}

	if dst.String() != "" {
		t.Errorf("Unexpected unpacking value [%s] => [%s]", src, dst)
	}

	if len(dst.String()) != len("") {
		t.Errorf("Silent drop error [%s] => [%s]", src, dst)
	}

	// Test Latin1 chars, replacing invalid chars with space
	src = DPT_16000("çäüLatin1:öàéè")
	t_src := "   Latin1:    "
	buf = src.Pack()
	err = dst.Unpack(buf)

	if err != nil {
		t.Errorf("Unpacking error %s [%s] => [%s]", err, src, dst)
	}

	if dst.String() != t_src {
		t.Errorf("Unexpected unpacking value [%s] => [%s]", src, dst)
	}

	if len(dst.String()) != len(t_src) {
		t.Errorf("Silent drop not working [%s] => [%s]", src, dst)
	}

	// Test with unsupported chars that should be silently converted to 0x20
	src = DPT_16000("hello你好")

	if !src.IsValid() {
		buf = src.Pack()
		err = dst.Unpack(buf)

		if err != nil {
			t.Errorf("%s [%s] => [%s]", err, src, dst)
		}
	}
}

// Test DPT 16.000 with special input
func TestDPT_16001(t *testing.T) {
	var dst DPT_16001

	// Test with simple string
	src := DPT_16001("KNX is OK")
	buf := src.Pack()
	r := fmt.Sprintf("%x", buf)
	if r != "004b4e58206973204f4b0000000000" {
		t.Errorf("Error packing [%s] => [%x]", src, r)
	}

	err := dst.Unpack(buf)
	if err != nil {
		t.Errorf("Error unpacking [%s]", src)
	}

	// Test with a string that will get truncated
	src = DPT_16001("This is a valid string longer than 14 chars, it will be truncated")
	e_dst := "This is a vali"
	buf = src.Pack()
	err = dst.Unpack(buf)

	if err != nil {
		t.Errorf("Error unpacking [%s] => [%s]", src, dst)
	}

	if dst.String() != e_dst {
		t.Errorf("Unexpected unpacking value [%s] => [%s]", src, dst)
	}

	// Test with empty string
	src = DPT_16001("")
	buf = src.Pack()
	err = dst.Unpack(buf)

	if err != nil {
		t.Errorf("Unpacking error %s [%s] => [%s]", err, src, dst)
	}

	if dst.String() != "" {
		t.Errorf("Unexpected unpacking value [%s] => [%s]", src, dst)
	}

	if len(dst.String()) != len("") {
		t.Errorf("Silent drop error [%s] => [%s]", src, dst)
	}

	// Test With Latin1 chars
	src = DPT_16001("çäüLatin1:öàéè")
	buf = src.Pack()
	err = dst.Unpack(buf)

	if err != nil {
		t.Errorf("Unpacking error %s [%s] => [%s]", err, src, dst)
	}

	if dst.String() != src.String() {
		t.Errorf("Unexpected unpacking value [%s] => [%s]", src, dst)
	}

	if len(dst.String()) != len(src.String()) {
		t.Errorf("Silent drop not working [%s] => [%s]", src, dst)
	}

	// Test with unsupported chars that should be silently converted to 0x20
	src = DPT_16001("hello你好")

	if !src.IsValid() {
		buf = src.Pack()
		err = dst.Unpack(buf)

		if err != nil {
			t.Errorf("%s [%s] => [%s]", err, src, dst)
		}
	}
}
