package main

import "fmt"

func main() {
	// First 15 elems
	for i := uint32(0); i < 0x10; i++ {
		j := uint32(i + 0x0)
		fmt.Printf("%10d %032b %d\n", j, j, fastestCountSetBitsInUint32(j))
	}

	// Second set
	for i := uint32(0); i < 0x10; i++ {
		j := uint32(i + 0xFF0)
		fmt.Printf("%10d %032b %d\n", j, j, fastestCountSetBitsInUint32(j))
	}

	// Third set
	for i := uint32(0); i < 0x10; i++ {
		j := uint32(i + 0xFFFF0)
		fmt.Printf("%10d %032b %d\n", j, j, fastestCountSetBitsInUint32(j))
	}

	// Last 15 elems
	for i := uint32(0); i < 0x10; i++ {
		j := uint32(i + 0xFFFFFFF0)
		fmt.Printf("%10d %032b %d\n", j, j, fastestCountSetBitsInUint32(j))
	}
}
