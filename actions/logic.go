package actions

import (
	"errors"
	"fmt"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/internal"
)

// error formatting
func fmtErr(prefix string, errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	// format errors
	cErr := errors.New(prefix)
	for _, err := range errs {
		cErr = fmt.Errorf("%w\n - %v", cErr, err)
	}
	return cErr
}

// Not negates the result of a given action
func Not(action configo.ActionFunc) configo.ActionFunc {
	internal.PanicIfNil(action)

	return func() error {
		err := action()
		if err != nil {
			return nil
		}
		return errors.New("expected action failure")
	}
}

// Or succeeds if any of the provided functions result in a successful state
// Basically trying to parse the value with all prvided parsing function,
// until the first of those functions does not return an error.
// In case all of these actions return an error, the returned parser function
// also returns an error constructed of all other errors
// This basically expects ANY of the actions to succeed
func Or(actions ...configo.ActionFunc) configo.ActionFunc {
	panicIfEmptyActionFunc(actions)

	return func() error {
		errs := make([]error, 0, len(actions))
		for _, f := range actions {
			err := f()
			if err != nil {
				errs = append(errs, err)
			} else {
				// return on first successful parser
				return nil
			}
		}
		return fmtErr("could not execute action: ", errs)
	}
}

// Xor enforces that only one of the functions results in a successful result.
// all of the other results are expected to yield an error, otherwise this function returns an error.
// This basically expects only ONE of the actions to succeed.
func Xor(actions ...configo.ActionFunc) configo.ActionFunc {
	panicIfEmptyActionFunc(actions)

	return func() error {
		errs := make([]error, 0, len(actions))
		successIndexes := make([]int, 0, 2)
		for idx, f := range actions {
			err := f()
			if err != nil {
				errs = append(errs, err)
			} else {
				successIndexes = append(successIndexes, idx)
			}
		}

		diff := len(actions) - len(errs)
		if diff == 0 {
			// no success
			return fmtErr("could not execute action: ", errs)
		} else if diff != 1 {
			// more than one success
			return fmt.Errorf("multiple actions succeeded, but only one was allowed to succeed: %v", successIndexes)
		}
		return nil
	}
}

// And expects that all provided actions succeed
// This expects ALL actions to succeed.
func And(actions ...configo.ActionFunc) configo.ActionFunc {
	panicIfEmptyActionFunc(actions)

	return func() error {
		for _, f := range actions {
			err := f()
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// If conditional allows to use different actions based on the passed condition.
func If(condition *bool, trueCase configo.ActionFunc, falseCase configo.ActionFunc) configo.ActionFunc {
	internal.PanicIfNil(condition, trueCase, falseCase)

	return func() error {
		if *condition {
			return trueCase()
		}
		return falseCase()
	}
}

// IfAction conditional allows to use different actions based on the passed condition.
func IfAction(condition configo.ActionFunc, trueCase configo.ActionFunc, falseCase configo.ActionFunc) configo.ActionFunc {
	internal.PanicIfNil(condition, trueCase, falseCase)

	return func() error {
		if condition() == nil {
			return trueCase()
		}
		return falseCase()
	}
}

// OnlyIf executes the trueCase action only in the case that the condition is true at the time when the
// parent option is parsed. Not at the time of option definition
func OnlyIf(condition *bool, trueCase configo.ActionFunc) configo.ActionFunc {
	internal.PanicIfNil(condition, trueCase)

	return func() error {
		if *condition {
			return trueCase()
		}
		return nil
	}
}

// OnlyIfNot executes the falseCase action only in the case that the condition is false at the time when the
// parent option is parsed. Not at the time of option definition.
func OnlyIfNot(condition *bool, falseCase configo.ActionFunc) configo.ActionFunc {
	internal.PanicIfNil(condition, falseCase)

	return func() error {
		if *condition {
			return nil
		}
		return falseCase()
	}
}

// OnlyIfAction executes the true case action at the time of option parsing only
// when the condition does return nil at the time of option parsing
func OnlyIfAction(condition configo.ActionFunc, trueCase configo.ActionFunc) configo.ActionFunc {
	internal.PanicIfNil(condition, trueCase)

	return func() error {
		if condition() == nil {
			return trueCase()
		}
		return nil
	}
}

// OnlyIfNotNil executes the true case action at the time of option parsing only
// when the condition is NOT nil at the time of option parsing
func OnlyIfNotNil(condition interface{}, trueCase configo.ActionFunc) configo.ActionFunc {
	internal.PanicIfNil(trueCase)

	return func() error {
		if condition != nil {
			return trueCase()
		}
		return nil
	}
}

// OnlyIfNil executes the true case action at the time of option parsing only
// when the condition is nil at the time of option parsing
func OnlyIfNil(condition interface{}, trueCase configo.ActionFunc) configo.ActionFunc {
	internal.PanicIfNil(trueCase)

	return func() error {
		if condition == nil {
			return trueCase()
		}
		return nil
	}
}
