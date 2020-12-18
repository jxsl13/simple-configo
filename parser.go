package configo

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
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

	// ErrParsing is returned when parsing of a key fails.
	ErrParsing = errors.New("Parsing Error")
)

// ParserFunc is a custom parser function that can be used to parse specific option values
// A Option struct must contain a ParseFunc in order to know, how to parse a specific value and where the
// Such a function is usually created using a generator function, that specifies the output type.
// This function is kept as simple as possible, in order to be handled exactly the same way for every
// possible return value
type ParserFunc func(value string) error

// DefaultParserString is the default function that returns a function
// that sets the parsed value to the passed referenced variable
// out is a pointer to the variable that gets the parsed value assigned to.
func DefaultParserString(out *string) ParserFunc {
	return func(value string) error {
		*out = value
		return nil
	}
}

// DefaultParserBool is the default function that returns a function
// that sets the 'out' referenced variable to the parsed value.
func DefaultParserBool(out *bool) ParserFunc {
	return func(value string) error {
		b, ok := boolValues[value]
		if !ok {
			return fmt.Errorf("Invalid value of type 'bool': %s", value)
		}
		*out = b
		return nil
	}
}

// DefaultParserInt parses a passed value and sets the passed out reference to the resulting value,
// returns an error otherwise.
func DefaultParserInt(out *int) ParserFunc {
	return func(value string) error {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("Invalid value of type 'integer': %s : %w", value, err)
		}
		*out = i
		return nil
	}
}

// DefaultParserFloat parses a passed value and sets the passed out reference to the resulting value,
// returns an error otherwise.
func DefaultParserFloat(out *float64, bitSize int) ParserFunc {
	return func(value string) error {
		f, err := strconv.ParseFloat(value, bitSize)
		if err != nil {
			return fmt.Errorf("Invalid value of type 'float': %s : %w", value, err)
		}
		*out = f
		return nil
	}
}

// DefaultParserDuration parses a string containing a duration value in the format of:
// 13h55m33s, 5m, 1h, 5m10s, 5m3s, etc.
// and returns the duration to the passed out pointer
func DefaultParserDuration(out *time.Duration) ParserFunc {
	return func(value string) error {
		d, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("Invalid value of type 'duration': %s : %w", value, err)
		}
		*out = d
		return nil
	}
}

// DefaultParserList parses a string containing a 'delimiter'(space, comma, semicolon, etc.) delimited list
// into the string list 'out'
func DefaultParserList(delimiter *string, out *[]string) ParserFunc {
	return func(value string) error {
		list := strings.Split(value, *delimiter)

		if len(list) > 0 && list[len(list)-1] == "" {
			list = list[:len(list)-1]
		}
		*out = list
		return nil
	}
}

// DefaultParserListToSet parses a string containing a 'delimiter'(space, comma, semicolon, etc.) delimited list
// into set 'out'
func DefaultParserListToSet(delimiter *string, out *map[string]bool) ParserFunc {
	return func(value string) error {
		list := strings.Split(value, *delimiter)

		if len(list) > 0 && list[len(list)-1] == "" {
			list = list[:len(list)-1]
		}

		*out = make(map[string]bool, len(list))
		for _, s := range list {
			(*out)[s] = true
		}
		return nil
	}
}

// DefaultParserChoiceString restricts the string value to a given set of values
// that are passed with the 'allowed' parameter.
func DefaultParserChoiceString(out *string, allowed ...string) ParserFunc {

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
			return fmt.Errorf("Invalid value of type 'string' got: '%s', allowed: %v", value, allowedList)
		}

		*out = value
		return nil
	}
}

// DefaultParserChoiceInt restricts the integer value to a given set of values
// that are passed with the 'allowed' parameter.
func DefaultParserChoiceInt(out *int, allowed ...int) ParserFunc {

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
			return fmt.Errorf("Invalid value of type 'integer': %s : %w", value, err)
		}

		// value not allowed
		if !allowedSet[i] {
			allowedList := setToSortedListInt(allowedSet)
			return fmt.Errorf("Invalid value of type 'integer' got: '%s', allowed: %v", value, allowedList)
		}

		*out = i
		return nil
	}
}

// DefaultParserChoiceFloat restricts the float value to a given set of values
// that are passed with the 'allowed' parameter.
func DefaultParserChoiceFloat(out *float64, bitSize int, allowed ...float64) ParserFunc {

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
			return fmt.Errorf("Invalid value of type 'float': %s : %w", value, err)
		}

		// value not allowed
		if !allowedSet[f] {
			allowedList := setToSortedListFloat(allowedSet)
			return fmt.Errorf("Invalid value of type 'float' got: '%s', allowed: %v", value, allowedList)
		}

		*out = f
		return nil
	}
}

// DefaultParserRangesInt restricts the integer value to a distinct list of min-max ranges.
// If the passed value to the returned function is not in any of these ranges, an error is returned.
// Example minMaxRange:
// #1: 0,1024  		# range from 0 through 1024
// #2: 0,3,9,10 	# range from 0 through 3 and from 9 through 10
// #3: 0,10,2,12 	# rage from 0 through 12
func DefaultParserRangesInt(out *int, minMaxRanges ...int) ParserFunc {

	distinctRanges := newDistinctRangeListInt(minMaxRanges...)

	return func(value string) error {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("Invalid value of type 'integer': %s : %w", value, err)
		}

		// value not allowed
		if !distinctRanges.Contains(i) {
			return fmt.Errorf("Invalid value of type 'integer' got: '%s', allowed ranges: %s", value, distinctRanges.String())
		}

		*out = i
		return nil
	}
}

// DefaultParserRegex allows the value of a key to be compliant to a regular expression.
func DefaultParserRegex(out *string, regex, errMsg string) ParserFunc {
	r := regexp.MustCompile(regex)
	return func(value string) error {
		if !r.MatchString(value) {
			return fmt.Errorf("%w : %s", ErrParsing, errMsg)
		}

		*out = value
		return nil
	}
}
