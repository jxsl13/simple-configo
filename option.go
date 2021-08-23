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
// The UnparseFunction allows to do the opposite of the ParseFunvction. It is called once your application shuts down.
// This way you may serialize previously deserialized values back into a file, set environment variables,
// close client connections, and so on.
// IsPseudoOption is an option that does not necessary relate to any actual Key value in any configuration map but does actually just do
// some operation that relies on previously computed config values e.g. the construction of a file path that
// needs a previously configured and evaluated directory path and some filename in order to construct that path.
// INFO: A pseudo option enforces the execution of the parsing function, even if the corresponding key does not exist in e.g. the environment.
type Option struct {
	Key          string
	Description  string
	Mandatory    bool
	DefaultValue string

	PreParseAction  ActionFunc
	ParseFunction   ParserFunc
	PostParseAction ActionFunc

	PreUnparseAction ActionFunc
	UnparseFunction  UnparserFunc
}

func (o *Option) IsAction() bool {
	return o.ParseFunction == nil &&
		o.UnparseFunction == nil &&
		(o.PreParseAction != nil ||
			o.PostParseAction != nil ||
			o.PreUnparseAction != nil)
}

func (o *Option) IsOption() bool {
	return o.ParseFunction != nil || o.UnparseFunction != nil
}

func (o *Option) Parse(m map[string]string) error {

	if err := tryExecAction(o.PreParseAction); err != nil {
		return fmt.Errorf("pre parse action: %w", err)
	}

	// evaluation
	value, ok := m[o.Key]
	if !ok {
		// value not found in env map
		if o.Mandatory && o.DefaultValue == "" {
			// no default value and no value in environment
			return fmt.Errorf("error: missing mandatory key: %s", o.Key)
		}
	}

	// mandatory values may be empty but only of the env value exists,
	if !o.Mandatory || o.DefaultValue != "" {
		if err := tryParse(o.DefaultValue, o.ParseFunction); err != nil {
			return fmt.Errorf("error in default value of option '%s': %w", o.Key, err)
		}
	}

	// if we do get a valid value from the passed map, the default value is
	// overwritten then
	// pseudo options do not evaluate the value, but get the value from somewhere else other than the passed
	// string map. They might prompt the user via the shell, read some file etc.
	if ok {
		if err := tryParse(value, o.ParseFunction); err != nil {
			return fmt.Errorf("error in value of option '%s': %w", o.Key, err)
		}
	}

	if err := tryExecAction(o.PostParseAction); err != nil {
		return fmt.Errorf("post parse action: %w", err)
	}
	return nil
}

func (o *Option) Unparse() (string, error) {

	err := tryExecAction(o.PreUnparseAction)
	if err != nil {
		if errors.Is(err, ErrSkipUnparse) {
			return "", ErrSkipUnparse
		}
		return "", fmt.Errorf("pre unparse action: %w", err)
	}

	// Unparse (serialize) option values
	value, err := tryUnparse(o.UnparseFunction)
	if err != nil {
		if errors.Is(err, ErrSkipUnparse) {
			return "", ErrSkipUnparse
		}
		return "", fmt.Errorf("error while unparsing the option '%s': %w", o.Key, err)
	}

	// skip default values in order to keep the config file/env variables map small.
	if value == o.DefaultValue {
		// allow user to manually decide whether to use or not to us ethe value
		return value, ErrSkipUnparse
	}
	return value, nil
}

// Options are usually unique, so one MUST NOT use redundant Option parameters
type Options []Option

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
