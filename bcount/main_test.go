package main

import (
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
	for i := 0; i < b.N; i++ {
		_ = fastestCountSetBitsInUint32(uint32(i))
	}
}

func BenchmarkFastCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fastCountSetBitsInUint32(uint32(i))
	}
}

func BenchmarkNaiveCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = naiveCountSetBitsInUint32(uint32(i))
	}
}
