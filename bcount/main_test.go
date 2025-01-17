package main

import (
	"math/bits"
	"testing"
)

var tests = []struct {
	test uint32
	want int
}{
	{0x0, 0},
	{0x1, 1},
	{0x2, 1},
	{0x3, 2},
	{0xF, 4},
	{0xFF, 8},
	{0xFF00, 8},
	{0xFF0000, 8},
	{0xFF000000, 8},
}

func TestAlgos(t *testing.T) {
	for i := uint32(0xFFFFFFFF); i > 0; i >>= 1 {
		naive := naiveCountSetBitsInUint32(i)
		fast := fastCountSetBitsInUint32(i)
		fastest := fastestCountSetBitsInUint32(i)

		if naive != fast || fast != fastest || fastest != naive {
			t.Errorf("Expected naive (%d) to be equal to fast (%d) and to fastest (%d)",
				naive, fast, fastest)
		}
	}
}

func TestFastestCount(t *testing.T) {
	for _, tst := range tests {
		got := fastestCountSetBitsInUint32(tst.test)
		if tst.want != got {
			t.Errorf("Got: %d, expected: %d for %d", got, tst.want, tst.test)
		}
	}
}

func TestFastCount(t *testing.T) {
	for _, tst := range tests {
		got := fastCountSetBitsInUint32(tst.test)
		if tst.want != got {
			t.Errorf("Got: %d, expected: %d for %d", got, tst.want, tst.test)
		}
	}
}

func TestNaiveCount(t *testing.T) {
	for _, tst := range tests {
		got := naiveCountSetBitsInUint32(tst.test)
		if tst.want != got {
			t.Errorf("Got: %d, expected: %d for %d", got, tst.want, tst.test)
		}
	}
}

func BenchmarkFastestCount(b *testing.B) {
	var res int
	for i := 0; i < b.N; i++ {
		res = fastestCountSetBitsInUint32(uint32(i))
	}
	_ = res
}

func BenchmarkFastCount(b *testing.B) {
	var res int
	for i := 0; i < b.N; i++ {
		res = fastCountSetBitsInUint32(uint32(i))
	}
	_ = res
}

func BenchmarkNaiveCount(b *testing.B) {
	var res int
	for i := 0; i < b.N; i++ {
		res = naiveCountSetBitsInUint32(uint32(i))
	}
	_ = res
}

func BenchmarkStdLibOnesCount32(b *testing.B) {
	var res int
	for i := 0; i < b.N; i++ {
		res = bits.OnesCount32(uint32(i))
	}
	_ = res
}
