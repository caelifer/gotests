package main

import (
	"fmt"
	"os"
	"time"

	"github.com/caelifer/gotests/signal/dispatch"
)

var sigCounter = 0

func main() {
	// Install hadler
	dispatch.HandleSignal(os.Interrupt, handleSIGINT_ONE)

	periodic := time.Tick(1 * time.Second)
	start := time.Now()

	for _ = range periodic {
		fmt.Print(".")
		if time.Since(start) > 5*time.Second {
			// No longer handle SIGINT
			dispatch.StopSignalHandler(os.Interrupt)

			// Give time to user to press CTRL-C
			time.Sleep(3 * time.Second)

			// Exit program
			return
		}
	}
}

func handleSIGINT_ONE(signal os.Signal) {
	sigCounter++
	fmt.Printf("\nHandler ONE: Got signal [%v]\n", signal)

	if sigCounter > 3 {
		os.Exit(1)
	}

	// Install different signal in mid-flight
	dispatch.HandleSignal(os.Interrupt, handleSIGINT_TWO)
}

func handleSIGINT_TWO(signal os.Signal) {
	sigCounter++
	fmt.Printf("\nHandler TWO: Got signal [%v]\n", signal)

	if sigCounter > 3 {
		os.Exit(1)
	}

	// Install different signal in mid-flight
	dispatch.HandleSignal(os.Interrupt, handleSIGINT_ONE)
}
