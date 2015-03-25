package main

import (
	"fmt"
	"time"
)

func ringer(nodes, nmsgs int) time.Duration {
	// Create ring
	in, out := createRing(nodes)

	// Start timer
	t0 := time.Now()

	// Send and receive nmsgs messages
	send(in, nmsgs)
	recv(out, nmsgs)

	// Calculate time
	return time.Since(t0)
}

func createRing(nodes int) (chan<- struct{}, <-chan struct{}) {
	in := make(chan struct{})
	out := in

	for i := 0; i < nodes; i++ {
		out = makeLink(out)
	}

	return in, out
}

func send(in chan<- struct{}, n int) {
	go func() {
		for i := 0; i < n; i++ {
			// Send message around the ring
			in <- struct{}{}
		}
	}()
}

func recv(out <-chan struct{}, n int) {
	for i := 0; i < n; i++ {
		<-out
	}
}

func makeLink(out <-chan struct{}) chan struct{} {
	in := make(chan struct{})
	go func() {
		for {
			in <- <-out
		}
	}()
	return in
}

func main() {
	fmt.Println(ringer(1000, 1000))
}
