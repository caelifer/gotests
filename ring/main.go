package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

// ringer is a ring of nodes connected by one-directional channels. Each node will substract one from
// recieved integer and send the result to the following neighbor. Once count is zero, stop.
func ringer(nodes, nmsgs int) int {
	// Create ring
	res := make(chan int)
	in := createRing(nodes, res)

	// Send token
	in <- nmsgs
	runtime.Gosched()

	// Report id
	return <-res
}

func createRing(nodes int, res chan int) chan int {
	in := make(chan int)
	out := in
	i := 1
	for ; i < nodes; i++ {
		out = makeLink(i, out, res)
	}
	// fmt.Println("Creating final thread #:", i)
	go link(i, out, in, res)
	return in
}

func link(id int, in <-chan int, out chan<- int, res chan int) {
	for {
		i := <-in

		if i == 0 {
			// Send id of a node that received a zero.
			res <- id
			return
		}
		// Opearation: send one less to the neighbor
		out <- i - 1
	}
}

func makeLink(id int, in <-chan int, res chan int) chan int {
	out := make(chan int)
	go link(id, in, out, res)
	return out
}

func main() {
	n := 1000
	if len(os.Args) > 1 {
		n, _ = strconv.Atoi(os.Args[1])
	}
	// Run single-threaded
	runtime.GOMAXPROCS(1)
	// Start timer
	t0 := time.Now()
	r := ringer(503, n)
	fmt.Println(r, time.Since(t0))
}

// vim: :ts=4:sw=4:noexpandtab:ai
