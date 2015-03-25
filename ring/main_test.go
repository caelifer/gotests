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

func BenchmarkRing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		send(in, 1000)
		recv(out, 1000)
	}
}
