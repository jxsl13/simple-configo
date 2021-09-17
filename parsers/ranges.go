package parsers

import (
	"fmt"
	"strconv"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/internal"
)

// RangesInt restricts the integer value to a distinct list of min-max ranges.
// If the passed value to the returned function is not in any of these ranges, an error is returned.
// Example minMaxRange:
// #1: 0,1024  		# range from 0 through 1024
// #2: 0,3,9,10 	# range from 0 through 3 and from 9 through 10
// #3: 0,10,2,12 	# rage from 0 through 12
func RangesInt(out *int, minMaxRanges ...int) configo.ParserFunc {
	internal.PanicIfNil(out)
	internal.PanicIfEmptyInt(minMaxRanges)

	// heap allocated and always there for lookup, does not need to be recreated on every function call
	// that is returned below
	distinctRanges := NewDistinctRangeListInt(minMaxRanges...)

	return func(value string) error {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid value of type 'integer': %s : %w", value, err)
		}

		// value not allowed
		if !distinctRanges.Contains(i) {
			return fmt.Errorf("invalid value of type 'integer' got: '%s', allowed ranges: %s", value, distinctRanges.String())
		}

		*out = i
		return nil
	}
}
