// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"fmt"
)

// DPT_8001 represents DPT 8.001 / Counter.
type DPT_8001 int16

func (d DPT_8001) Pack() []byte {
	return packV16(int16(d))
}

func (d *DPT_8001) Unpack(data []byte) error {
	return unpackV16(data, (*int16)(d))
}

func (d DPT_8001) Unit() string {
	return ""
}

func (d DPT_8001) String() string {
	return fmt.Sprintf("%d °", int16(d))
}

// DPT_8002 represents DPT 8.002 / delta time ms.
type DPT_8002 int16

func (d DPT_8002) Pack() []byte {
	return packV16(int16(d))
}

func (d *DPT_8002) Unpack(data []byte) error {
	return unpackV16(data, (*int16)(d))
}

func (d DPT_8002) Unit() string {
	return "ms"
}

func (d DPT_8002) String() string {
	return fmt.Sprintf("%d ms", int16(d))
}

// DPT_8003 represents DPT 8.003 / delta time centiseconds
type DPT_8003 int16

func (d DPT_8003) Pack() []byte {
	return packV16(int16(d))
}

func (d *DPT_8003) Unpack(data []byte) error {
	return unpackV16(data, (*int16)(d))
}

func (d DPT_8003) Unit() string {
	return "cs"
}

func (d DPT_8003) String() string {
	return fmt.Sprintf("%d cs", int16(d))
}

// DPT_8004 represents DPT 8.004 / delta time deciseconds
type DPT_8004 int16

func (d DPT_8004) Pack() []byte {
	return packV16(int16(d))
}

func (d *DPT_8004) Unpack(data []byte) error {
	return unpackV16(data, (*int16)(d))
}

func (d DPT_8004) Unit() string {
	return "ds"
}

func (d DPT_8004) String() string {
	return fmt.Sprintf("%d ds", int16(d))
}

// DPT_8005 represents DPT 8.005 / delta time seconds
type DPT_8005 int16

func (d DPT_8005) Pack() []byte {
	return packV16(int16(d))
}

func (d *DPT_8005) Unpack(data []byte) error {
	return unpackV16(data, (*int16)(d))
}

func (d DPT_8005) Unit() string {
	return "s"
}

func (d DPT_8005) String() string {
	return fmt.Sprintf("%d s", int16(d))
}

// DPT_8006 represents DPT 8.005 / delta time minutes
type DPT_8006 int16

func (d DPT_8006) Pack() []byte {
	return packV16(int16(d))
}

func (d *DPT_8006) Unpack(data []byte) error {
	return unpackV16(data, (*int16)(d))
}

func (d DPT_8006) Unit() string {
	return "min"
}

func (d DPT_8006) String() string {
	return fmt.Sprintf("%d min", int16(d))
}

// DPT_8007 represents DPT 8.007 / delta time hours
type DPT_8007 int16

func (d DPT_8007) Pack() []byte {
	return packV16(int16(d))
}

func (d *DPT_8007) Unpack(data []byte) error {
	return unpackV16(data, (*int16)(d))
}

func (d DPT_8007) Unit() string {
	return "h"
}

func (d DPT_8007) String() string {
	return fmt.Sprintf("%d h", int16(d))
}

// DPT_8010 represents DPT 8.010 / percentage
type DPT_8010 int16

func (d DPT_8010) Pack() []byte {
	return packV16(int16(d))
}

func (d *DPT_8010) Unpack(data []byte) error {
	return unpackV16(data, (*int16)(d))
}

func (d DPT_8010) Unit() string {
	return "%"
}

func (d DPT_8010) String() string {
	return fmt.Sprintf("%d %%", int16(d))
}

// DPT_8011 represents DPT 8.011 / Rotation angle °.
type DPT_8011 int16

func (d DPT_8011) Pack() []byte {
	return packV16(int16(d))
}

func (d *DPT_8011) Unpack(data []byte) error {
	return unpackV16(data, (*int16)(d))
}

func (d DPT_8011) Unit() string {
	return "°"
}

func (d DPT_8011) String() string {
	return fmt.Sprintf("%d °", int16(d))
}
