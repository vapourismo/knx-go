// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"fmt"
	"strings"
	"testing"
)

// Test DPT 28.001 with some special inputs
func TestDPT_28001(t *testing.T) {
	var dst DPT_28001

	// Test with ASCII string
	src := DPT_28001("KNX is OK")
	buf := src.Pack()

	if fmt.Sprintf("%x", buf) != "004b4e58206973204f4b00" {
		t.Errorf("Error packing [%s] => [%x]", src, buf)
	}

	err := dst.Unpack(buf)
	if err != nil {
		t.Errorf("Error unpacking [%x]", buf)
	}

	// Test with an empty string
	src = DPT_28001("")
	buf = src.Pack()
	err = dst.Unpack(buf)

	if err != nil {
		t.Errorf("Error unpacking [%x] => [%s]", buf, dst)
	}

	if dst.String() != src.String() {
		t.Errorf("Error comparing [%s]/%d => [%s]/%d", src, len(src), dst, len(dst))
	}

	// Test with a very long string
	src = DPT_28001(strings.Repeat("|01234567890123456789", 50))
	buf = src.Pack()
	err = dst.Unpack(buf)

	if err != nil {
		t.Errorf("Error unpacking [%x] => [%s]", buf, dst)
	}

	if dst.String() != string(src) {
		t.Errorf("Error comparing [%s]/%d => [%s]/%d", src, len(src), dst, len(dst))
	}

	// Test with Latin1 chars
	src = DPT_28001("Latin1:çäüöàéè")
	buf = src.Pack()
	err = dst.Unpack(buf)

	if err != nil {
		t.Errorf("Unpacking error [%x] => [%s]", buf, dst)
	}

	if dst.String() != src.String() {
		t.Errorf("Unexpected ISO-8859-1 chars [%s] => [%s]", src, dst)
	}

	// Test with mixed chars, 1, 2 and 3 bytes
	// Length = 17 = 3 * 2 + 5 * 1 + 2 * 3
	src = DPT_28001("çäühello你好")
	buf = src.Pack()
	err = dst.Unpack(buf)

	if err != nil {
		t.Errorf("Error unpacking [%s] => [%s]", src, dst)
	}

	if dst.String() != src.String() {
		t.Errorf("Error when comparing [%s] => [%s]", src, dst)
	}

	if len(dst) != 17 {
		t.Errorf("Lost or gained something [%s] => [%s]", src, dst)
	}
}
