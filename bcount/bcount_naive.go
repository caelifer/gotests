package main

// naiveCountBitsInUint8 counts set bits byte represented as uint8
func naiveCountBitsInUint8(b uint8) int {
	n := uint8(0)
	for i := uint8(0); i < 8; i++ {
		n += 0x1 & (b >> i)
	}

	// Never reached
	return int(n)
}

// naiveCountSetBitsInUint32 counts bits in provided uint32 by
// summing bit counts of 4 bytes that comprises uint32
func naiveCountSetBitsInUint32(ui uint32) int {
	return naiveCountBitsInUint8(uint8(ui>>24)) +
		naiveCountBitsInUint8(uint8(ui>>16)) +
		naiveCountBitsInUint8(uint8(ui>>8)) +
		naiveCountBitsInUint8(uint8(ui))
}