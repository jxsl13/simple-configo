package parsers

import (
	"sync"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/internal"
)

// SyncWrapper allows to wrap ParseFunc in a synchronized context.
// This allows to execute the parsing in an asynchronous context.
// which prevents data races even in the case that the configuration
// periodically kept up to date.
func SyncWrapper(mu sync.Locker, f configo.ParserFunc) configo.ParserFunc {
	internal.PanicIfNil(mu, f)
	return func(value string) error {
		mu.Lock()
		defer mu.Unlock()
		return f(value)
	}
}
