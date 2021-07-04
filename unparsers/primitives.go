package unparsers

import (
	"strconv"

	configo "github.com/jxsl13/simple-configo"
)

// String returns a function that returns the string representation of the in value
func String(in *string) configo.UnparserFunc {
	if in == nil {
		panic("String: nil pointer passed")
	}
	return func() (string, error) {
		return *in, nil
	}
}

// Bool returns a function that returns the string representation of the in value
func Bool(in *bool) configo.UnparserFunc {
	if in == nil {
		panic("Bool: nil pointer passed")
	}
	return func() (string, error) {
		if *in {
			return "true", nil
		} else {
			return "false", nil
		}
	}
}

// Int returns a function that returns the string representation of the in value
func Int(in *int) configo.UnparserFunc {
	if in == nil {
		panic("Int: nil pointer passed")
	}
	return func() (string, error) {
		return strconv.Itoa(*in), nil
	}
}

// Float returns a function that returns the string representation of the in value
func Float(in *float64, bitSize int) configo.UnparserFunc {
	if in == nil {
		panic("Float: nil pointer passed")
	}
	return func() (string, error) {
		return strconv.FormatFloat(*in, 'f', -1, bitSize), nil
	}
}
