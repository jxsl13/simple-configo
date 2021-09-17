package parsers

import (
	"fmt"
	"time"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/internal"
)

// Duration parses a string containing a duration value in the format of:
// 13h55m33s, 5m, 1h, 5m10s, 5m3s, etc.
// and returns the duration to the passed out pointer
func Duration(out *time.Duration) configo.ParserFunc {
	internal.PanicIfNil(out)

	return func(value string) error {
		d, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("invalid value of type 'duration': %s : %w", value, err)
		}
		*out = d
		return nil
	}
}
