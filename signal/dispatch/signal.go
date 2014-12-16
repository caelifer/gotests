package dispatch

import (
	"log"
	"os"
	syssignal "os/signal"
	"sync"
)

func init() {
	// Set logger
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
}

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

// HandleSignal installs custom handler for a particular os.Signal provided by signal.
func HandleSignal(signal os.Signal, handler SignalHandler) {
	// Uneregister handler if it exists
	StopSignalHandler(signal)

	log.Printf("registering new [%s] handler", signal)

	// Create buffered channel of os.Signal values
	ch := make(chan os.Signal, 1)

	/////////////////////// protected section ///////////////////
	// Take exclusive lock
	dispatcher.Lock()

	// Install our new channel
	dispatcher.signals[signal] = ch

	// Fast unlock
	dispatcher.Unlock()
	/////////////////////// protected section ///////////////////

	// Set notification
	syssignal.Notify(ch, signal)

	// Install custom handler in the separate gorutine
	go func(c <-chan os.Signal, sig os.Signal) {
		for s := range c {
			handler(s)
		}
		log.Printf("exiting [%s] handler", sig)
	}(ch, signal)
}

// StopSignalHandler safely stops signal handling for signal specified by signal.
// If no hanlder exists, this function is noop.
func StopSignalHandler(signal os.Signal) {
	// Take exclusive lock
	dispatcher.Lock()
	defer dispatcher.Unlock()

	// Check if we already have registered handler
	if ch, ok := dispatcher.signals[signal]; ok {
		// Signal handler already exists - do clean-up
		log.Printf("unregistering existing [%s] handler", signal)

		// Stop receiving signlas
		syssignal.Stop(ch)
		// Close signal channel so gorutine can safely exit
		close(ch)

		// Clear our signal table
		delete(dispatcher.signals, signal)
	}
}
