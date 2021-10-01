package unparsers

import (
	"errors"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/internal"
)

// Not negates the result of a given unparser
func Not(unparser configo.UnparserFunc) configo.UnparserFunc {
	internal.PanicIfNil(unparser)

	return func() (string, error) {
		value, err := unparser()
		if err != nil {
			return value, nil
		}
		return value, errors.New("expected parsing failure")
	}
}

// Or succeeds if any of the provided functions result in a successful state
// Basically trying to unparse the value with all provided unparsing function,
// until the first of those functions does not return an error.
// In case all of these parsers return an error, the returned unparser function
// also returns an error constructed of all other errors
// This basically expects ANY of the parsers to succeed
func Or(unparsers ...configo.UnparserFunc) configo.UnparserFunc {
	panicIfEmptyUnparseFunc(unparsers)

	return func() (string, error) {
		errs := make([]error, 0, len(unparsers))

		for _, f := range unparsers {
			value, err := f()
			if err != nil {
				errs = append(errs, err)
			} else {
				// return on first success
				return value, nil
			}
		}
		return "", internal.FmtErr("could not unparseparse: ", errs)
	}
}

// If conditional allows to use different unparsers based on the passed condition.
func If(condition *bool, trueCase configo.UnparserFunc, falseCase configo.UnparserFunc) configo.UnparserFunc {
	internal.PanicIfNil(condition, trueCase, falseCase)

	return func() (string, error) {
		if *condition {
			return trueCase()
		}
		return falseCase()
	}
}

// OnlyIf executes the trueCase action only in the case that the condition is true
func OnlyIf(condition *bool, trueCase configo.UnparserFunc) configo.UnparserFunc {
	internal.PanicIfNil(condition, trueCase)

	return func() (string, error) {
		if *condition {
			return trueCase()
		}
		return "", configo.ErrSkipUnparse
	}
}

// OnlyIfNot executes the falseCase action only in the case that the condition is false
func OnlyIfNot(condition *bool, falseCase configo.UnparserFunc) configo.UnparserFunc {
	internal.PanicIfNil(condition, falseCase)

	return func() (string, error) {
		if *condition {
			return "", configo.ErrSkipUnparse
		}
		return falseCase()
	}
}

// OnlyIfNotNil executes the trueCase function when the condition is NOT nil at evaluation time
func OnlyIfNotNil(condition interface{}, trueCase configo.UnparserFunc) configo.UnparserFunc {
	internal.PanicIfNil(trueCase)

	return func() (string, error) {
		if condition != nil {
			return trueCase()
		}
		return "", configo.ErrSkipUnparse
	}
}

// OnlyIfNotNil executes the trueCase function when the condition is nil at evaluation time
func OnlyIfNil(condition interface{}, trueCase configo.UnparserFunc) configo.UnparserFunc {
	internal.PanicIfNil(trueCase)

	return func() (string, error) {
		if condition == nil {
			return trueCase()
		}
		return "", configo.ErrSkipUnparse
	}
}
