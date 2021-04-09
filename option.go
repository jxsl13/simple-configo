package configo

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (

	// ErrOptionMissingKey The option is missing the 'Key' field
	ErrOptionMissingKey = errors.New("the option is missing the 'Key' field")
	// ErrOptionMissingDescription The option is missing the 'Description' field
	ErrOptionMissingDescription = errors.New("the option is missing the 'Description' field")
	// ErrOptionInvalidDefaultValue the option has an invalid 'DefaultValue' field, please check its 'type' field
	ErrOptionInvalidDefaultValue = errors.New("the option has an invalid 'DefaultValue' field")
	// ErrOptionMissingParseFunction the option is missing its 'ParseFunc' field
	ErrOptionMissingParseFunction = errors.New("the option is missing its 'ParseFunc' field")
)

// Option is a value that can be set to configure
// The key is any kind of key that defines the option's values in like an environment variable or and kind
// of map structure that is passed to the configo.Parse(config, map) function.
// The description describes this option with a long text
// The Mandatory parameter enforces this value to be present, either by having a non-empty DefaultValue string
//	 or by being present in the map that is used to fill the resulting struct.
// The DefaultValue can be any non-empty string value that can for example be configured, but must not be configured, as
// the default value is good enough without being changed. At specific constellations and with specific parsing functions this value is
// also checked for validity with the following ParseFunction.
// The ParseFunction heavily relies on side effects, as it does only return the error in case the parsing failed.
// Usually the pattern is followed where another function gets parameters or struct property references passed and returns a ParseFunc
// that modifies the parameters that were passed to the parent function which returns the ParseFunc.
// IsPseudoOption is an option that does not necessary relate to any actual Key value in any configuration map but does actually just do
// some operation that relies on previously computed config values e.g. the construction of a file path that
// needs a previously configured and evaluated directory path and some filename in order to construct that path.
type Option struct {
	Key            string
	Description    string
	Mandatory      bool
	DefaultValue   string
	ParseFunction  ParserFunc
	IsPseudoOption bool
}

// IsValid retrurns true if there are no programming errors
func (o *Option) IsValid() error {

	if o.Key == "" && !o.IsPseudoOption {
		return ErrOptionMissingKey
	}

	if o.Description == "" && !o.IsPseudoOption {
		return ErrOptionMissingDescription
	}

	if (o.DefaultValue != "" || !o.Mandatory) && !o.IsPseudoOption {
		err := o.ParseFunction(o.DefaultValue)
		if err != nil {
			return fmt.Errorf("%w : %v", ErrOptionInvalidDefaultValue, err)
		}
	}

	return nil
}

// MustValid enforces validity of the option.
// panics if the programmer did a mistake
func (o *Option) MustValid() {
	if err := o.IsValid(); err != nil {
		panic(fmt.Sprintf("error: %s : %v", o.Key, err))
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
		Key            string
		Description    string
		Mandatory      bool
		DefaultValue   string
		IsPseudoOption bool
	}

	so := SubOption{
		Key:            o.Key,
		Description:    o.Description,
		Mandatory:      o.Mandatory,
		DefaultValue:   o.DefaultValue,
		IsPseudoOption: o.IsPseudoOption,
	}

	b, err := json.MarshalIndent(&so, " ", " ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
