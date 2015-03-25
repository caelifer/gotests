package main

import (
	"os"
	"testing"
)

var (
	in  chan<- struct{}
	out <-chan struct{}
)

func TestMain(m *testing.M) {
	in, out = createRing(1000)
	os.Exit(m.Run())
}

func BenchmarkRing_____1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		send(in, 1)
		recv(out, 1)
	}
}
func BenchmarkRing____10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		send(in, 10)
		recv(out, 10)
	}
}
func BenchmarkRing___100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		send(in, 100)
		recv(out, 100)
	}
}
func BenchmarkRing__1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		send(in, 1000)
		recv(out, 1000)
	}
}
func BenchmarkRing_10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		send(in, 10000)
		recv(out, 10000)
	}
}
