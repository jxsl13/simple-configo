package configo

import (
	"os"
	"syscall"
)

var (
	// ShutdownSignalListFunc can be exchanged in order to listen to a different set of functions
	// that indicate a shutdown.
	ShutdownSignalListFunc = DefaultSignalListFunc()
)

// DefaultSignalListFunc returns a function that returns the signals os.Interrupt and syscall.SIGTERM
func DefaultSignalListFunc() func() []os.Signal {
	return func() []os.Signal {
		return []os.Signal{
			os.Interrupt,
			syscall.SIGTERM,
		}
	}
}
