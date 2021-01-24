package configo

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (

	// ErrOptionMissingKey The option is missing the 'Key' field
	ErrOptionMissingKey = errors.New("The option is missing the 'Key' field")
	// ErrOptionMissingDescription The option is missing the 'Description' field
	ErrOptionMissingDescription = errors.New("The option is missing the 'Description' field")
	// ErrOptionInvalidDefaultValue The option has an invalid 'DefaultValue' field, please check its 'type' field
	ErrOptionInvalidDefaultValue = errors.New("The option has an invalid 'DefaultValue' field")
	// ErrOptionMissingParseFunction The option is missing its 'ParseFunc' field
	ErrOptionMissingParseFunction = errors.New("The option is missing its 'ParseFunc' field")
)

// Option is a value that can be set to configure
type Option struct {
	Key           string
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

	if err := o.ParseFunction(o.DefaultValue); !o.Mandatory && err != nil {
		return fmt.Errorf("%w : %v", ErrOptionInvalidDefaultValue, err)
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

// Options are usually unique, so one MUST NOT use redundant Option parameters
type Options []Option

// MustValid panics if any of the option definitions is not valid.
func (o *Options) MustValid() {
	for _, option := range *o {
		option.MustValid()
	}
}

func (o *Option) String() string {
	type SubOption struct {
		Key          string
		Description  string
		Mandatory    bool
		DefaultValue string
	}

	so := SubOption{
		Key:          o.Key,
		Description:  o.Description,
		Mandatory:    o.Mandatory,
		DefaultValue: o.DefaultValue,
	}

	b, err := json.MarshalIndent(&so, " ", " ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
