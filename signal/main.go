package main

import (
	"fmt"
	"os"
	"time"
	"sync/atomic"

	"github.com/caelifer/gotests/signal/dispatch"
)

// Signal event counter - global shared value
var sigCounter int32 = 0

func main() {
	// Install hadler
	dispatch.HandleSignal(os.Interrupt, handleSIGINT_ONE)

	periodic := time.Tick(1 * time.Second)
	start := time.Now()

	for _ = range periodic {
		fmt.Print(".")
		if time.Since(start) > 5*time.Second {
			// Stop handling INT signal
			dispatch.StopSignalHandler(os.Interrupt)

			// Give time to user to press CTRL-C
			time.Sleep(3 * time.Second)

			// Exit program
			return
		}
	}
}

// SIGINT custom handler one
func handleSIGINT_ONE(signal os.Signal) {
	// Atomically adjust counter
	atomic.AddInt32(&sigCounter, 1)

	fmt.Printf("\nHandler ONE: Got signal [%v]\n", signal)

	// Atomically compare
	if atomic.LoadInt32(&sigCounter) > 3 {
		os.Exit(1)
	}

	// Install different signal in mid-flight
	dispatch.HandleSignal(os.Interrupt, handleSIGINT_TWO)
}

// SIGINT custom handler two
func handleSIGINT_TWO(signal os.Signal) {
	// Atomically adjust counter
	atomic.AddInt32(&sigCounter, 1)

	fmt.Printf("\nHandler TWO: Got signal [%v]\n", signal)

	// Atomically compare
	if atomic.LoadInt32(&sigCounter) > 3 {
		os.Exit(1)
	}

	// Install different signal in mid-flight
	dispatch.HandleSignal(os.Interrupt, handleSIGINT_ONE)
}
