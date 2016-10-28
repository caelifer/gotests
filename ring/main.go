package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

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
	for i < nodes {
		out = makeLink(i, out, res)
		i++
	}
	// fmt.Println("Creating final thread #:", i)
	go link(i, out, in, res)
	return in
}

func link(id int, in <-chan int, out chan<- int, res chan int) {
	for {
		i := <-in
		if i == 0 {
			res <- id
			return
		}
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
