package configo

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	// valid type strings
	// TODO: could be extended with a RWMutex and the ability to
	// be extended with a package function call, where new types can register
	// themselves.
	validTypes = []string{
		"bool",
		"string",
		"int",
		"float",
		"duration",
		"list",
	}
	// ErrOptionMissingKey The option is missing the 'Key' field
	ErrOptionMissingKey = errors.New("The option is missing the 'Key' field")
	// ErrOptionMissingDescription The option is missing the 'Description' field
	ErrOptionMissingDescription = errors.New("The option is missing the 'Description' field")
	// ErrOptionMissingType The option is missing its 'Type' field
	ErrOptionMissingType = errors.New("The option is missing its 'Type' field")
	// ErrOptionUnknownType The option has an unknown 'Type' field
	ErrOptionUnknownType = errors.New("The option has an unknown 'Type' field")
	// ErrOptionMissingDefaultValue The option is missing its default fallback value
	ErrOptionMissingDefaultValue = errors.New("The option is missing its default fallback value")
	// ErrOptionInvalidDefaultValue The option has an invalid 'DefaultValue' field
	ErrOptionInvalidDefaultValue = errors.New("The option has an invalid 'DefaultValue' field")
	// ErrOptionMissingParseFunction The option is missing its 'ParseFunc' field
	ErrOptionMissingParseFunction = errors.New("The option is missing its 'ParseFunc' field")
)

// Option is a value that can be set to configure
// Type can be one of the following values:
// bool, string, integer, duration, list
// list is a list of strings
type Option struct {
	Key           string
	Type          string
	Description   string
	Mandatory     bool
	DefaultValue  string
	ParseFunction ParserFunc
}

// IsValid retrurns true if there are no programming errors
func (o *Option) IsValid() error {

	if o.Key == "" {
		return ErrOptionMissingKey
	}

	if o.Description == "" {
		return ErrOptionMissingDescription
	}

	if !o.Mandatory && o.DefaultValue == "" {
		return ErrOptionMissingDefaultValue
	}

	if o.Type == "" {
		return ErrOptionMissingType
	}

	isValidType := false
	for _, t := range validTypes {
		if t == o.Type {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return ErrOptionUnknownType
	}

	if ok := isValidValue(o.Type, o.DefaultValue); !ok {
		return ErrOptionInvalidDefaultValue
	}
	return nil
}

// MustValid enforces validity of the option.
// panics if the programmer did a mistake
func (o *Option) MustValid() {
	if err := o.IsValid(); err != nil {
		panic(fmt.Sprintf("Error: %s : %v", o.Key, err))
	}
}

func isValidValue(typestr, value string) bool {
	isValid := false

	switch typestr {
	case "bool":
		_, ok := boolValues[value]
		isValid = ok
	case "string":
		return true
	case "int":
		_, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			isValid = true
		}
	case "duration":
		_, err := time.ParseDuration(value)
		if err == nil {
			isValid = true
		}
	case "list":
		return true
	case "float":
		_, err := strconv.ParseFloat(value, 64)
		if err == nil {
			isValid = true
		}
	default:
		// unknown types
		isValid = false
	}

	return isValid
}

// Options are usually unique, so one MUST NOT use redundant Option parameters
type Options []Option

// MustValid panics if any of the option definitions is not valid.
func (o *Options) MustValid() {
	for _, option := range *o {
		option.MustValid()
	}
}
