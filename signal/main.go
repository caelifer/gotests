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
	for _ = range periodic {
		fmt.Print(".")
	}
}

func handleSIGINT_ONE(signal os.Signal) {
	sigCounter++
	fmt.Printf("\nHandler ONE: Got signal [%v]\n", signal)

	// Install different signal in mid-flight
	dispatch.HandleSignal(os.Interrupt, handleSIGINT_TWO)
}

func handleSIGINT_TWO(signal os.Signal) {
	sigCounter++
	fmt.Printf("\nHandler TWO: Got signal [%v]\n", signal)
	if sigCounter > 1 {
		os.Exit(1)
	}
}
