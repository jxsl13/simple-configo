package parsers

import (
	"fmt"
	"strconv"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/internal"
)

// ChoiceString restricts the string value to a given set of values
// that are passed with the 'allowed' parameter.
func ChoiceString(out *string, allowed ...string) configo.ParserFunc {
	internal.PanicIfNil(out)
	internal.PanicIfEmptyString(allowed)

	// create set only once in order to have a fast access later on
	// in order not to waste RAM, we waste a few CPU cycles instead, if allowed contains
	// redundant string values.
	allowedSet := make(map[string]bool, len(allowed)/2)
	for _, choice := range allowed {
		allowedSet[choice] = true
	}

	return func(value string) error {

		// value not allowed
		if !allowedSet[value] {
			allowedList := setToSortedListString(allowedSet)
			return fmt.Errorf("invalid value of type 'string' got: '%s', allowed: %v", value, allowedList)
		}

		*out = value
		return nil
	}
}

// ChoiceInt restricts the integer value to a given set of values
// that are passed with the 'allowed' parameter.
func ChoiceInt(out *int, allowed ...int) configo.ParserFunc {
	internal.PanicIfNil(out)
	internal.PanicIfEmptyInt(allowed)

	// create set only once in order to have a fast access later on
	// in order not to waste RAM, we waste a few CPU cycles instead, if allowed contains
	// redundant string values.
	allowedSet := make(map[int]bool, len(allowed)/2)
	for _, choice := range allowed {
		allowedSet[choice] = true
	}

	return func(value string) error {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid value of type 'integer': %s : %w", value, err)
		}

		// value not allowed
		if !allowedSet[i] {
			allowedList := setToSortedListInt(allowedSet)
			return fmt.Errorf("invalid value of type 'integer' got: '%s', allowed: %v", value, allowedList)
		}

		*out = i
		return nil
	}
}

// ChoiceFloat restricts the float value to a given set of values
// that are passed with the 'allowed' parameter.
func ChoiceFloat(out *float64, bitSize int, allowed ...float64) configo.ParserFunc {
	internal.PanicIfNil(out)
	internal.PanicIfEmptyFloat(allowed)

	// create set only once in order to have a fast access later on
	// in order not to waste RAM, we waste a few CPU cycles instead, if allowed contains
	// redundant string values.
	allowedSet := make(map[float64]bool, len(allowed)/2)
	for _, choice := range allowed {
		allowedSet[choice] = true
	}

	return func(value string) error {
		f, err := strconv.ParseFloat(value, bitSize)
		if err != nil {
			return fmt.Errorf("invalid value of type 'float': %s : %w", value, err)
		}

		// value not allowed
		if !allowedSet[f] {
			allowedList := setToSortedListFloat(allowedSet)
			return fmt.Errorf("invalid value of type 'float' got: '%s', allowed: %v", value, allowedList)
		}

		*out = f
		return nil
	}
}
