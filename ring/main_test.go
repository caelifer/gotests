package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func BenchmarkRing_____1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ringer(503, 1)
	}
}
func BenchmarkRing____10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ringer(503, 1)
	}
}
func BenchmarkRing___100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ringer(503, 100)
	}
}
func BenchmarkRing__1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ringer(503, 1000)
	}
}
func BenchmarkRing_10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ringer(503, 10000)
	}
}
