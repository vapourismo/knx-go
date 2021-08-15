// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"fmt"
)

// DPT_7001 represents DPT 7.001 / Value 2 Ucount.
type DPT_7001 uint16

func (d DPT_7001) Pack() []byte {
	return packU16(uint16(d))
}

func (d *DPT_7001) Unpack(data []byte) error {
	return unpackU16(data, (*uint16)(d))
}

func (d DPT_7001) Unit() string {
	return "pulses"
}

func (d DPT_7001) String() string {
	return fmt.Sprintf("%d pulses", uint16(d))
}

// DPT_7002 represents DPT 7.002 / Time Period MSec.
type DPT_7002 uint16

func (d DPT_7002) Pack() []byte {
	return packU16(uint16(d))
}

func (d *DPT_7002) Unpack(data []byte) error {
	return unpackU16(data, (*uint16)(d))
}

func (d DPT_7002) Unit() string {
	return "ms"
}

func (d DPT_7002) String() string {
	return fmt.Sprintf("%d ms", uint16(d))
}

// DPT_7003 represents DPT 7.003 / Time Period 10 MSec.
type DPT_7003 uint16

func (d DPT_7003) Pack() []byte {
	return packU16(uint16(d))
}

func (d *DPT_7003) Unpack(data []byte) error {
	return unpackU16(data, (*uint16)(d))
}

func (d DPT_7003) Unit() string {
	return "s"
}

func (d DPT_7003) String() string {
	return fmt.Sprintf("%d s", uint16(d))
}

// DPT_7004 represents DPT 7.004 / Time Period 100 MSec.
type DPT_7004 uint16

func (d DPT_7004) Pack() []byte {
	return packU16(uint16(d))
}

func (d *DPT_7004) Unpack(data []byte) error {
	return unpackU16(data, (*uint16)(d))
}

func (d DPT_7004) Unit() string {
	return "s"
}

func (d DPT_7004) String() string {
	return fmt.Sprintf("%d s", uint16(d))
}

// DPT_7005 represents DPT 7.005 / Time Period Sec.
type DPT_7005 uint16

func (d DPT_7005) Pack() []byte {
	return packU16(uint16(d))
}

func (d *DPT_7005) Unpack(data []byte) error {
	return unpackU16(data, (*uint16)(d))
}

func (d DPT_7005) Unit() string {
	return "s"
}

func (d DPT_7005) String() string {
	return fmt.Sprintf("%d s", uint16(d))
}

// DPT_7006 represents DPT 7.006 / Time Period Min.
type DPT_7006 uint16

func (d DPT_7006) Pack() []byte {
	return packU16(uint16(d))
}

func (d *DPT_7006) Unpack(data []byte) error {
	return unpackU16(data, (*uint16)(d))
}

func (d DPT_7006) Unit() string {
	return "m"
}

func (d DPT_7006) String() string {
	return fmt.Sprintf("%d m", uint16(d))
}

// DPT_7007 represents DPT 7.007 / Time Period Hrs.
type DPT_7007 uint16

func (d DPT_7007) Pack() []byte {
	return packU16(uint16(d))
}

func (d *DPT_7007) Unpack(data []byte) error {
	return unpackU16(data, (*uint16)(d))
}

func (d DPT_7007) Unit() string {
	return "h"
}

func (d DPT_7007) String() string {
	return fmt.Sprintf("%d h", uint16(d))
}

// DPT_7010 represents DPT 7.010 / Property DataType.
type DPT_7010 uint16

func (d DPT_7010) Pack() []byte {
	return packU16(uint16(d))
}

func (d *DPT_7010) Unpack(data []byte) error {
	return unpackU16(data, (*uint16)(d))
}

func (d DPT_7010) Unit() string {
	return ""
}

func (d DPT_7010) String() string {
	return fmt.Sprintf("%d", uint16(d))
}

// DPT_7011 represents DPT 7.011 / Length mm.
type DPT_7011 uint16

func (d DPT_7011) Pack() []byte {
	return packU16(uint16(d))
}

func (d *DPT_7011) Unpack(data []byte) error {
	return unpackU16(data, (*uint16)(d))
}

func (d DPT_7011) Unit() string {
	return "mm"
}

func (d DPT_7011) String() string {
	return fmt.Sprintf("%d mm", uint16(d))
}

// DPT_7012 represents DPT 7.012 / Current mA.
type DPT_7012 uint16

func (d DPT_7012) Pack() []byte {
	return packU16(uint16(d))
}

func (d *DPT_7012) Unpack(data []byte) error {
	return unpackU16(data, (*uint16)(d))
}

func (d DPT_7012) Unit() string {
	return "mA"
}

func (d DPT_7012) String() string {
	return fmt.Sprintf("%d mA", uint16(d))
}

// DPT_7013 represents DPT 7.013 / Brightness lux.
type DPT_7013 uint16

func (d DPT_7013) Pack() []byte {
	return packU16(uint16(d))
}

func (d *DPT_7013) Unpack(data []byte) error {
	return unpackU16(data, (*uint16)(d))
}

func (d DPT_7013) Unit() string {
	return "lux"
}

func (d DPT_7013) String() string {
	return fmt.Sprintf("%d lux", uint16(d))
}
