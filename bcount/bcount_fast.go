package main

// fastCountBitsInUint8 counts set bits byte represented as uint8
func fastCountBitsInUint8(b byte) int {
	switch b {
	case 0:
		return 0
	case 1, 2, 4, 8, 16, 32, 64, 128:
		return 1
	case 3, 5, 6, 9, 10, 12, 17, 18, 20, 24, 33,
		34, 36, 40, 48, 65, 66, 68, 72, 80, 96,
		129, 130, 132, 136, 144, 160, 192:
		return 2
	case 7, 11, 13, 14, 19, 21, 22, 25, 26, 28,
		35, 37, 38, 41, 42, 44, 49, 50, 52, 56,
		67, 69, 70, 73, 74, 76, 81, 82, 84, 88,
		97, 98, 100, 104, 112, 131, 133, 134, 137,
		138, 140, 145, 146, 148, 152, 161, 162,
		164, 168, 176, 193, 194, 196, 200, 208, 224:
		return 3
	case 15, 23, 27, 29, 30, 39, 43, 45, 46, 51, 53,
		54, 57, 58, 60, 71, 75, 77, 78, 83, 85, 86,
		89, 90, 92, 99, 101, 102, 105, 106, 108, 113,
		114, 116, 120, 135, 139, 141, 142, 147, 149,
		150, 153, 154, 156, 163, 165, 166, 169, 170,
		172, 177, 178, 180, 184, 195, 197, 198, 201,
		202, 204, 209, 210, 212, 216, 225, 226, 228,
		232, 240:
		return 4
	case 31, 47, 55, 59, 61, 62, 79, 87, 91, 93, 94, 103,
		107, 109, 110, 115, 117, 118, 121, 122, 124, 143,
		151, 155, 157, 158, 167, 171, 173, 174, 179, 181,
		182, 185, 186, 188, 199, 203, 205, 206, 211, 213,
		214, 217, 218, 220, 227, 229, 230, 233, 234, 236,
		241, 242, 244, 248:
		return 5
	case 63, 95, 111, 119, 123, 125, 126, 159, 175, 183, 187,
		189, 190, 207, 215, 219, 221, 222, 231, 235, 237,
		238, 243, 245, 246, 249, 250, 252:
		return 6
	case 127, 191, 223, 239, 247, 251, 253, 254:
		return 7
	case 255:
		return 8
	}

	// Never reached
	return 0
}

// fastCountSetBitsInUint32 counts bits in provided uint32 by
// summing bit counts of 4 bytes that comprises uint32
func fastCountSetBitsInUint32(n uint32) int {
	return fastCountBitsInUint8(byte(n>>24)) +
		fastCountBitsInUint8(byte(n>>16)) +
		fastCountBitsInUint8(byte(n>>8)) +
		fastCountBitsInUint8(byte(n))
}
