package unparsers

import (
	"sync"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/internal"
)

// SyncWrapper allows to wrap UnparserFunc in a synchronized context.
// This allows to execute the configuration asynchronously or repeatedly without
// having a data race.
func SyncWrapper(mu sync.Locker, f configo.UnparserFunc) configo.UnparserFunc {
	internal.PanicIfNil(mu, f)
	return func() (string, error) {
		mu.Lock()
		defer mu.Unlock()
		return f()
	}
}
