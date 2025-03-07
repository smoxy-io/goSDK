package os

import (
	"os"
	"os/signal"
	"syscall"
)

// SignalHandler defines the function signature for a signal handler function
// sig is the signal that triggered the handler
// return true to exit signal handling (implies that the main process will exit)
type SignalHandler func(sig os.Signal) (exit bool)

var sigs chan os.Signal
var sigHandlerChan chan bool
var sigHandlers map[os.Signal]SignalHandler = map[os.Signal]SignalHandler{
	syscall.SIGINT:  DefaultSigIntHandler,
	syscall.SIGTERM: DefaultSigTermHandler,
}

// RegisterSignalHandler registers a handler function for a signal.
// registering a handler for a signal that already has a handler will replace the
// old handler with the new handler
func RegisterSignalHandler(sig syscall.Signal, handler SignalHandler) {
	sigHandlers[sig] = handler
}

// WaitForExitSignal blocks until a signal handler indicates that the process should exit
// StartSignalHandler MUST be called BEFORE this function
func WaitForExitSignal() {
	// wait for handler to signal that the process should exit
	// ignore all messages and errors (they all mean the process should exit)
	_, _ = <-sigHandlerChan
}

// StartSignalHandler starts the signal handler go routine
// all signal handlers MUST be registered with RegisterSignalHandler BEFORE to calling this function
func StartSignalHandler() {
	sigs = make(chan os.Signal, 128) // large buffer because we are listening for all signals
	sigHandlerChan = make(chan bool, 1)

	// by default, only SIGINT or SIGTERM will cause the process to exit
	// this can be overridden by registering a custom handler for SIGINT or SIGTERM
	signal.Notify(sigs)

	go func() {
		defer close(sigHandlerChan)

		// this loop is never broken so that we can continue to receive signals which enables support
		// for more advanced signal handling such as first SIGINT is graceful shutdown and second SIGINT
		// is forced process exit.
		// the main go routine will unblock when a signal handler returns true (by default this is
		// when SIGINT or SIGTERM is received)
		for {
			sig := <-sigs

			if sig == syscall.SIGURG {
				// ignore this signal as it doesn't mean anything and is used internally by go
				// so it shouldn't be used in the application
				continue
			}

			// check for a registered signal handler
			handler, ok := sigHandlers[sig]

			if !ok {
				// no handler for this signal. ignore it
				continue
			}

			// call the signal handler
			if handler(sig) {
				// unblock the main go routine. by default SIGINT or SIGTERM will trigger this
				sigHandlerChan <- true
			}
		}
	}()
}

func DefaultSigIntHandler(_ os.Signal) bool {
	return true
}

func DefaultSigTermHandler(_ os.Signal) bool {
	return true
}

// SendSignal sends an os signal to the signal handler. allows complex applications a chance to exit gracefully from
// anywhere within the application
func SendSignal(sig os.Signal) {
	sigs <- sig
}
