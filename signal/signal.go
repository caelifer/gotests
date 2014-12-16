package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

var sigCounter = 0

func main() {
	// Install hadler
	HandleSignal(os.Interrupt, func(signal os.Signal) {
		sigCounter++
		fmt.Printf("\nGot SIGINT [%v]\n", signal)
		if sigCounter > 2 {
			os.Exit(1)
		}
	})

	periodic := time.Tick(1 * time.Second)
	for _ = range periodic {
		fmt.Print(".")
	}
}

func HandleSignal(sig os.Signal, f func(os.Signal)) {
	// Create SIGINT buffered channel
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sig)

	go func() {
		for s := range ch {
			f(s)
		}
	}()
}
