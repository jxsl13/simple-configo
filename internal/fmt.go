package internal

import (
	"errors"
	"fmt"
)

// error formatting
func FmtErr(prefix string, errs []error) error {
	// format errors
	cErr := errors.New(prefix)
	for _, err := range errs {
		cErr = fmt.Errorf("%w\n - %v", cErr, err)
	}
	return cErr
}
