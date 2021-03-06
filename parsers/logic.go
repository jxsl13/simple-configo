package parsers

import (
	"errors"
	"fmt"

	configo "github.com/jxsl13/simple-configo"
)

// error formatting
func fmtErr(prefix string, errs []error) error {
	// format errors
	cErr := errors.New(prefix)
	for _, err := range errs {
		cErr = fmt.Errorf("%w\n -%v", cErr, err)
	}
	return cErr
}

// Or succeeds if any of the provided functions result in a successful state
// Basically trying to parse the value with all prvided parsing function,
// until the first of those functions does not return an error.
// In case all of these parsers return an error, the returned parser function
// also returns an error constructed of all other errors
// This basically expects ANY of the parsers to succeed
func Or(parsers ...configo.ParserFunc) configo.ParserFunc {
	return func(value string) error {
		errs := make([]error, 0, len(parsers))
		for _, f := range parsers {
			err := f(value)
			if err != nil {
				errs = append(errs, err)
			} else {
				// return the first working parser
				return nil
			}
		}
		return fmtErr("could not parse: ", errs)
	}
}

// Xor enforces that only one of the functions results in a successful result.
// all of the other results are expected to yield an error, otherwise this function returns an error.
// This basically expects only ONE of the parsers to succeed.
func Xor(parsers ...configo.ParserFunc) configo.ParserFunc {
	return func(value string) error {
		errs := make([]error, 0, len(parsers))
		successIndexes := make([]int, 0, 2)
		for idx, f := range parsers {
			err := f(value)
			if err != nil {
				errs = append(errs, err)
			} else {
				successIndexes = append(successIndexes, idx)
			}
		}

		diff := len(parsers) - len(errs)
		if diff == 0 {
			// no success
			return fmtErr("could not parse: ", errs)
		} else if diff != 1 {
			// more than one success
			return fmt.Errorf("multiple parsers succeeded, but only one was allowed to succeed: %v", successIndexes)
		}
		return nil
	}
}

// And expects that all provided parsers succeed
// This expects ALL parsers to succeed.
func And(parsers ...configo.ParserFunc) configo.ParserFunc {
	return func(value string) error {
		for _, f := range parsers {
			err := f(value)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// If conditional allows to use different parsers base don the passed condition.
func If(condition bool, trueCase configo.ParserFunc, falseCase configo.ParserFunc) configo.ParserFunc {
	if condition {
		return trueCase
	}
	return falseCase
}
