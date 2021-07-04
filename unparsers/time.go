package unparsers

import (
	"time"

	configo "github.com/jxsl13/simple-configo"
)

// Duration returns the string representation of the currently present value in the in pointer.
// Example output: 13h55m33s, 5m, 1h, 5m10s, 5m3s, etc.
func Duration(in *time.Duration) configo.UnparserFunc {
	return func() (string, error) {
		return in.String(), nil
	}
}
