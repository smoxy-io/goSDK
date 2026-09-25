package os

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// Go signal usage
//
// Signals the runtime uses for itself:
// - SIGPROF: CPU profiling sample delivery
// - SIGURG: Async preemption
// - SIGSETXID, SIGCANCEL, SIGSYNCCALL: glibc/musl internals
// - SIGPIPE: special-case for writes to closed file descriptors
//
// Signals with default runtime behavior:
// - SIGQUIT: full goroutine traceback then exit
// - SIGABRT: full goroutine traceback then exit
//
// Signals that cannot be caught via `_SigNotify`:
// - SIGSYS
// - SIGSEGV
// - SIGBUS
// - SIGFPE
// - SIGILL
// - SIGTRAP
// - SIGSTKFLT

var (
	notifySigs = []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGALRM,
		syscall.SIGTERM,
		syscall.SIGCHLD,
		syscall.SIGCONT,
		syscall.SIGTSTP,
		syscall.SIGTTIN,
		syscall.SIGTTOU,
		syscall.SIGXCPU,
		syscall.SIGXFSZ,
		syscall.SIGVTALRM,
		syscall.SIGWINCH,
		syscall.SIGIO,
		syscall.SIGQUIT,
		syscall.SIGABRT,
	}
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
	syscall.SIGQUIT: DefaultSigQuitHandler,
	syscall.SIGABRT: DefaultSigAbrtHandler,
}

// RegisterSignalHandler registers a handler function for a signal.
// registering a handler for a signal that already has a handler will replace the
// old handler with the new handler
func RegisterSignalHandler(sig os.Signal, handler SignalHandler) {
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
	// SIGQUIT and SIGABRT also preserve the default runtime behavior, but can be overridden by custom handlers
	signal.Notify(sigs, notifySigs...)

	go func() {
		defer close(sigHandlerChan)

		// this loop is never broken so that we can continue to receive signals which enables support
		// for more advanced signal handling such as first SIGINT is graceful shutdown and second SIGINT
		// is forced process exit.
		// the main go routine will unblock when a signal handler returns true (by default this is
		// when SIGINT or SIGTERM is received)
		for {
			sig := <-sigs

			// terminals send a SIGHUP signal when the terminal is closed.
			// reset the signal handler to ensure that we can continue to receive signals in
			// the event that the process continues to run after the terminal is closed
			if sig == syscall.SIGHUP {
				signal.Reset()
				signal.Notify(sigs, notifySigs...)
			}

			// check for handler that should be called for all signals
			// if this handler returns false, then we will check for a signal specific handler
			// this can be useful for logging all signals received by the process or propagating signals
			// using an application specific method
			if hndlr, ok := sigHandlers[SIGALL]; ok {
				if hndlr(sig) {
					sigHandlerChan <- true
					continue
				}
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

func DefaultSigQuitHandler(_ os.Signal) bool {
	return DumpTraceAndExit()
}

func DefaultSigAbrtHandler(_ os.Signal) bool {
	return DumpTraceAndExit()
}

// DumpTraceAndExit writes a full stack trace to stderr and exits the process
func DumpTraceAndExit() bool {
	buf := make([]byte, 1<<22)
	n := runtime.Stack(buf, true)

	_, _ = os.Stderr.Write(buf[:n])
	_, _ = os.Stderr.Write([]byte{'\n'})

	return true
}

// SendSignal sends an os signal to the signal handler. allows complex applications a chance to exit gracefully from
// anywhere within the application
func SendSignal(sig os.Signal) {
	sigs <- sig
}
