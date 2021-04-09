package parsers

import (
	"fmt"
	"strconv"

	configo "github.com/jxsl13/simple-configo"
)

var (
	// map of valid bool values that can be used in configs
	boolValues = map[string]bool{
		"0":        false,
		"1":        true,
		"true":     true,
		"TRUE":     true,
		"false":    false,
		"FALSE":    false,
		"enabled":  true,
		"ENABLED":  true,
		"disabled": false,
		"DISABLED": false,
		"yes":      true,
		"YES":      true,
		"no":       false,
		"NO":       false,
		"enable":   true,
		"ENABLE":   true,
		"disable":  false,
		"DISABLE":  false,
	}
)

// String is the default function that returns a function
// that sets the parsed value to the passed referenced variable
// out is a pointer to the variable that gets the parsed value assigned to.
func String(out *string) configo.ParserFunc {
	return func(value string) error {
		*out = value
		return nil
	}
}

// Bool is the default function that returns a function
// that sets the 'out' referenced variable to the parsed value.
func Bool(out *bool) configo.ParserFunc {
	return func(value string) error {
		b, ok := boolValues[value]
		if !ok {
			return fmt.Errorf("invalid value of type 'bool': %s", value)
		}
		*out = b
		return nil
	}
}

// Int parses a passed value and sets the passed out reference to the resulting value,
// returns an error otherwise.
func Int(out *int) configo.ParserFunc {
	return func(value string) error {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid value of type 'integer': %s : %w", value, err)
		}
		*out = i
		return nil
	}
}

// Float parses a passed value and sets the passed out reference to the resulting value,
// returns an error otherwise.
func Float(out *float64, bitSize int) configo.ParserFunc {
	return func(value string) error {
		f, err := strconv.ParseFloat(value, bitSize)
		if err != nil {
			return fmt.Errorf("invalid value of type 'float': %s : %w", value, err)
		}
		*out = f
		return nil
	}
}
