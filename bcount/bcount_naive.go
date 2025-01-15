package main

// naiveCountBitsInUint8 counts set bits byte represented as uint8
func naiveCountBitsInUint8(b byte) int {
	var n int

	for b > 0 {
		n += int(b & 0x1)
		b >>= 1
	}

	return n
}

// naiveCountSetBitsInUint32 counts bits in provided uint32 by
// summing bit counts of 4 bytes that comprises uint32
func naiveCountSetBitsInUint32(n uint32) int {
	return naiveCountBitsInUint8(byte(n>>24)) +
		naiveCountBitsInUint8(byte(n>>16)) +
		naiveCountBitsInUint8(byte(n>>8)) +
		naiveCountBitsInUint8(byte(n))
}
