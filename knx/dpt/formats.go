// Copyright 2017 Ole Krüger.

package dpt

import (
	"encoding/binary"
	"unsafe"
)

const sizeB1 = 1

func packB1(buffer []byte, b bool) {
	if b {
		buffer[0] = 1
	} else {
		buffer[0] = 0
	}
}

func unpackB1(data []byte, b *bool) {
	*b = data[0]&1 == 1
}

const sizeB2 = 1

func packB2(buffer []byte, b1, b2 bool) {
	buffer[0] = 0

	if b1 {
		buffer[0] |= 2
	}

	if b2 {
		buffer[0] |= 1
	}
}

func unpackB2(data []byte, b1, b2 *bool) {
	*b1 = data[0]&2 == 2
	*b2 = data[0]&1 == 1
}

const sizeB1U3 = 1

func packB1U3(buffer []byte, b bool, u uint8) {
	buffer[0] = uint8(u & 7)

	if b {
		buffer[0] |= 1 << 3
	}
}

func unpackB1U3(data []byte, b *bool, u *uint8) {
	*b = data[0]&(1<<3) == 1<<3
	*u = data[0] & 7
}

const sizeA8 = 2

func packA8(buffer []byte, a byte) {
	buffer[1] = a
}

func unpackA8(data []byte, a *byte) {
	*a = data[1]
}

const sizeU8 = 2

func packU8(buffer []byte, u uint8) {
	buffer[1] = u
}

func unpackU8(data []byte, u *uint8) {
	*u = data[1]
}

const sizeV8 = 2

func packV8(buffer []byte, v int8) {
	buffer[1] = byte(v)
}

func unpackV8(data []byte, v *int8) {
	*v = int8(data[1])
}

const sizeB5N3 = 2

func packB5N3(buffer []byte, b1, b2, b3, b4, b5 bool, n uint8) {
	if b1 {
		buffer[1] |= 1 << 7
	}

	if b2 {
		buffer[1] |= 1 << 6
	}

	if b3 {
		buffer[1] |= 1 << 5
	}

	if b4 {
		buffer[1] |= 1 << 4
	}

	if b5 {
		buffer[1] |= 1 << 3
	}

	buffer[1] |= n & 7
}

func unpackB5N3(data []byte, b1, b2, b3, b4, b5 *bool, n *uint8) {
	*b1 = data[1]&(1<<7) == 1<<7
	*b2 = data[1]&(1<<6) == 1<<6
	*b3 = data[1]&(1<<5) == 1<<5
	*b4 = data[1]&(1<<4) == 1<<4
	*b5 = data[1]&(1<<3) == 1<<3
	*n = data[1] & 7
}

const sizeU16 = 3

func packU16(buffer []byte, u uint16) {
	binary.BigEndian.PutUint16(buffer[1:], u)
}

func unpackU16(data []byte, u *uint16) {
	*u = binary.BigEndian.Uint16(data[1:])
}

const sizeV16 = 3

func packV16(buffer []byte, v int16) {
	binary.BigEndian.PutUint16(buffer[1:], uint16(v))
}

func unpackV16(data []byte, v *int16) {
	*v = int16(binary.BigEndian.Uint16(data[1:]))
}

const sizeF16 = 3

func packF16(buffer []byte, f float32) {
	buffer[2] = 0
	buffer[1] = 0
	buffer[0] = 0

	if f > 670760.96 {
		f = 670760.96
	} else if f < -671088.64 {
		f = -671088.64
	}

	signedMantissa := int(f * 100)
	exp := 0

	for signedMantissa > 2047 || signedMantissa < -2048 {
		signedMantissa /= 2
		exp++
	}

	buffer[1] |= uint8(exp&15) << 3

	if signedMantissa < 0 {
		signedMantissa += 2048
		buffer[1] |= 1 << 7
	}

	mantissa := uint(signedMantissa)

	buffer[1] |= uint8(mantissa>>8) & 7
	buffer[2] |= uint8(mantissa)
}

func unpackF16(data []byte, f *float32) {
	m := int(data[1]&7)<<8 | int(data[2])
	if data[1]&128 == 128 {
		m -= 2048
	}

	e := (data[1] >> 3) & 15

	*f = 0.01 * float32(m) * float32(uint(1)<<e)
}

const sizeN3N5R2N6R2N6 = 4

func packN3N5R2N6R2N6(buffeR []byte, n1, n2, n3, n4 uint8) {
	buffeR[0] = 0
	buffeR[1] = (n1&7)<<5 | (n2 & 31)
	buffeR[2] = n3 & 63
	buffeR[3] = n4 & 63
}

func unpackN3N5R2N6R2N6(data []byte, n1, n2, n3, n4 *uint8) {
	*n1 = (data[1] >> 5) & 7
	*n2 = data[1] & 31
	*n3 = data[2] & 63
	*n4 = data[3] & 63
}

const sizeR3N5R4N4R1U7 = 4

func packR3U5R4U4R1U7(buffer []byte, n1, n2, u uint8) {
	buffer[0] = 0
	buffer[1] = n1 & 31
	buffer[2] = n2 & 15
	buffer[3] = u & 127
}

func unpackR3U5R4U4R1U7(buffer []byte, n1, n2, u *uint8) {
	*n1 = buffer[1] & 31
	*n2 = buffer[2] & 15
	*u = buffer[3] & 127
}

const sizeU32 = 5

func packU32(buffer []byte, u uint32) {
	binary.BigEndian.PutUint32(buffer[1:], u)
}

func unpackU32(buffer []byte, u *uint32) {
	*u = binary.BigEndian.Uint32(buffer[1:])
}

const sizeV32 = 5

func packV32(buffer []byte, v int32) {
	binary.BigEndian.PutUint32(buffer[1:], uint32(v))
}

func unpackV32(buffer []byte, v *int32) {
	*v = int32(binary.BigEndian.Uint32(buffer[1:]))
}

const sizeF32 = 5

func packF32(buffer []byte, f float32) {
	binary.BigEndian.PutUint32(buffer[1:], *(*uint32)(unsafe.Pointer(&f)))
}

func unpackF32(data []byte, f *float32) {
	u32 := binary.BigEndian.Uint32(data[1:])
	*f = *(*float32)(unsafe.Pointer(&u32))
}

const sizeU4U4U4U4U4U4B4N4 = 5

func packU4U4U4U4U4U4B4N4(buffer []byte, u1, u2, u3, u4, u5, u6 uint8, b1, b2, b3, b4 bool, n uint8) {
	buffer[0] = 0
	buffer[1] = (u1&15)<<4 | (u2 & 15)
	buffer[2] = (u3&15)<<4 | (u4 & 15)
	buffer[3] = (u5&15)<<4 | (u6 & 15)
	buffer[4] = n & 15

	if b1 {
		buffer[4] |= 1 << 7
	}

	if b2 {
		buffer[4] |= 1 << 6
	}

	if b3 {
		buffer[4] |= 1 << 5
	}

	if b4 {
		buffer[4] |= 1 << 4
	}
}

func unpackU4U4U4U4U4U4B4N4(buffer []byte, u1, u2, u3, u4, u5, u6 *uint8, b1, b2, b3, b4 *bool, n *uint8) {
	*u1 = (buffer[1] >> 4) & 15
	*u2 = buffer[1] & 15
	*u3 = (buffer[2] >> 4) & 15
	*u4 = buffer[2] & 15
	*u5 = (buffer[3] >> 4) & 15
	*u6 = buffer[3] & 15
	*n = buffer[4] & 15
	*b1 = buffer[4]&(1<<7) == 1<<7
	*b2 = buffer[4]&(1<<6) == 1<<6
	*b3 = buffer[4]&(1<<5) == 1<<5
	*b4 = buffer[4]&(1<<4) == 1<<4
}

const sizeA112 = 15

func packA112(buffer []byte, a [14]byte) {
	copy(buffer[1:], a[:])
}

func unpackA112(data []byte, a *[14]byte) {
	copy(a[:], data[1:])
}
