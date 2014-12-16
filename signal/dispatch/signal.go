package dispatch

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

// dispatcher
var dispatcher = struct {
	*sync.Mutex
	signals map[os.Signal]chan os.Signal
}{
	new(sync.Mutex),
	make(map[os.Signal]chan os.Signal),
}

// SignalHandler is a custom function that handles os.Signal
type SignalHandler func(os.Signal)

// HandleSignal installs custom SignalHandler handler for a particular os.Signal
// provided by sig argument.
func HandleSignal(sig os.Signal, handler SignalHandler) {
	var ch chan os.Signal

	// Take exclusive lock
	dispatcher.Lock()

	// Check if we already have registered handler
	if ch, ok := dispatcher.signals[sig]; ok {
		// Signal handler already exists - do clean-up
		log.Printf("clean-up existing [%s] handler", sig)

		// Stop receiving signlas
		signal.Stop(ch)
		// Close signal channel so gorutine can safely exit
		close(ch)
	}

	log.Printf("registering new [%s] handler", sig)

	// Create buffered channel of os.Signal values
	ch = make(chan os.Signal, 1)
	// Store our new channel
	dispatcher.signals[sig] = ch

	// Unlock
	dispatcher.Unlock()

	// Set notification
	signal.Notify(ch, sig)

	// Install custom handler in the separate gorutine
	go func(c <-chan os.Signal) {
		for s := range c {
			handler(s)
		}
	}(ch)
}
