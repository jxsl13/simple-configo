package parsers

import (
	"fmt"
	"regexp"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/internal"
)

// Regex allows the value of a key to be compliant to a regular expression.
func Regex(out *string, regex, errMsg string) configo.ParserFunc {
	internal.PanicIfNil(out)

	r := regexp.MustCompile(regex)
	return func(value string) error {
		if !r.MatchString(value) {
			return fmt.Errorf("%w : %s", Error, errMsg)
		}

		*out = value
		return nil
	}
}
