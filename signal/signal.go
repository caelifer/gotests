package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

var sigCounter = 0

// SignalHandler is a custom function that handles os.Signal
type SignalHandler func(os.Signal)

func main() {
	// Install hadler
	HandleSignal(os.Interrupt, func(signal os.Signal) {
		sigCounter++
		fmt.Printf("\nGot signal [%v]\n", signal)
		if sigCounter > 1 {
			os.Exit(1)
		}
	})

	periodic := time.Tick(1 * time.Second)
	for _ = range periodic {
		fmt.Print(".")
	}
}

// HandleSignal installs custom SignalHandler handler for a particular os.Signal
// provided by sig argument.
func HandleSignal(sig os.Signal, handler SignalHandler) {
	// Create SIGINT buffered channel
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sig)

	go func() {
		for s := range ch {
			handler(s)
		}
	}()
}
