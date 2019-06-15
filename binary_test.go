package binutils

import (
	"math"
	"testing"
	"gotest.tools/assert"
)

func b(bytes ...byte) []byte {
	slice := make([]byte, len(bytes))
	for idx, val := range bytes {
		slice[idx] = val
	}
	return slice
}

var knownEncodings map[int32][]byte = map[int32][]byte {
	// Values alternate positive/negative as we ascend binary values:
	0:   b(0x00),
	-1:  b(0x01),
	1:   b(0x02),
	-2:  b(0x03),
	2:   b(0x04),

	// -64 is the 'biggest' signed number that can be represented in a single byte
	63:   b(0x7e),
	-64:  b(0x7f),
	64:   b(0x80, 0x01),
	-65:  b(0x81, 0x01),
	65:   b(0x82, 0x01),


	// In this encoding, the largest integer is immediately followed by the smallest integer (in the binary view), because we know the *absolute* value 
	// of the smallest integer is one greater than the largest integer
	// After decoding this will be unravelled into an integer with all but the most significant bit set, which matches the largest possible integer.
	math.MaxInt32: b(0xfe, 0xff, 0xff, 0xff, 0x0f),
	// For the minimum possible integer, the bits will all be the same as above, except for the least significant, which is set. This means that after
	// decoding it will be the 32 bit integer with all bits set.
	math.MinInt32: b(0xff, 0xff, 0xff, 0xff, 0x0f),
}

func TestEncodeVarInt(t *testing.T) {
	for value, encoded := range knownEncodings {
		buffer := &[]byte{}
		WriteVarInt(buffer, value)
		assert.DeepEqual(t, *buffer, encoded)
	}
}

func TestDecodeVarInt(t *testing.T) {
	for value, encoded := range knownEncodings {
		off := 0
		read := ReadVarInt(&encoded, &off)
		assert.Equal(t, value, read)
	}
}