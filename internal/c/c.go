// The c package contains miscellaneous code common to all C targets.
package c

import (
	"fmt"
)

// WidthToReadType returns type that is sufficient to represent data
// of given width in the C language for read functions.
func WidthToReadType(width int64) Type {
	if width > 64 {
		return ByteArray{}
	} else if width > 32 {
		return Uint64Ptr{}
	} else if width > 16 {
		return Uint32Ptr{}
	} else if width > 8 {
		return Uint16Ptr{}
	}
	return Uint8Ptr{}
}

// WidthToWriteType returns type that is sufficient to represent data
// of given width in the C language for write functions.
func WidthToWriteType(width int64) Type {
	if width > 64 {
		return ByteArray{}
	} else if width > 32 {
		return Uint64{}
	} else if width > 16 {
		return Uint32{}
	} else if width > 8 {
		return Uint16{}
	}
	return Uint8{}
}

// WidthToWordByteShift returns byte shift required to advance to the next word in memory.
func WidthToWordByteShift(width int64) int64 {
	if width%8 != 0 {
		panic(fmt.Sprintf("unsupported word width %d", width))
	}

	return width / 8
}

// MaskToValue returns bit mask represented as value based on the masks start bit and end bit.
// The mask is always shifted to the right.
// For example, the mask for start bit 5 and end bit 8 is 0xF, not 0xF0.
// It panics if required conditions are not met.
func MaskToValue(startBit, endBit int64) uint64 {
	if startBit > 64 {
		panic(fmt.Sprintf("start bit (%d) greater than 64", startBit))
	}
	if endBit > 64 {
		panic(fmt.Sprintf("end bit (%d) greater than 64", endBit))
	}
	if startBit < 0 {
		panic(fmt.Sprintf("negative start bit (%d)", startBit))
	}
	if endBit < 0 {
		panic(fmt.Sprintf("negative end bit (%d)", endBit))
	}
	if startBit > endBit {
		panic(fmt.Sprintf("start bit (%d) is greater than end bit (%d) ", startBit, endBit))
	}

	mask := uint64(0)
	for i := startBit; i <= endBit; i++ {
		mask |= 1 << i
	}
	return (mask >> startBit)
}
