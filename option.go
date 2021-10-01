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

	PreUnparseAction  ActionFunc   // used to prepare values for serialization, may return an ErrSkipUnparse to skip the unparsing step.
	UnparseFunction   UnparserFunc // execute string serialization, may return an ErrSkipUnparse
	PostUnparseAction ActionFunc   // may be used for closing handles after parameter serialization, cannot invoke unparse skipping
}

// IsAction returns true in case the provided option has no ParseFunction nor UnparseFunction
// defined and at least one Action defined.
func (o *Option) IsAction() bool {
	return o.ParseFunction == nil &&
		o.UnparseFunction == nil &&
		(o.PreParseAction != nil ||
			o.PostParseAction != nil ||
			o.PreUnparseAction != nil ||
			o.PostUnparseAction != nil)
}

// IsOption retursns true in case either ParseFunction or UnparseFunction is not nil.
func (o *Option) IsOption() bool {
	return o.ParseFunction != nil || o.UnparseFunction != nil
}

// Parse evaluates the passed key/value map.
// Initially the PreparseAction is executed in case it exists, then the ParseFunction is executed with the map value
// that may or may not exist. Depending on whether the option is mandatory and not default value is set in that option,
// it requires that the option value is provided by the key/value map.
// After parsing the PostParseAction is executed.
// In case that the expected value is not provided by the env map, the default value is parsed instead.
// The default value is always parsed in order to check whether it is valid according to the ParseFunction.
// In case a custom value is defined the default value is overwritten by the custom value.
func (o *Option) Parse(m map[string]string) error {

	if err := tryExecAction(o.PreParseAction); err != nil {
		return fmt.Errorf("pre parse action of option '%s': %w", o.Key, err)
	}

	// mandatory values may be empty but only if the env value exists
	// parse default value in case the option ir not mandatory or in
	// the case that the option has a non-empty default value
	if !o.Mandatory || o.DefaultValue != "" {
		if err := tryParse(o.DefaultValue, o.ParseFunction); err != nil {
			return fmt.Errorf("error in default value of option '%s': %w", o.Key, err)
		}
	}

	// evaluation of environment map
	value, ok := m[o.Key]
	if !ok {
		// value not found in env map
		if o.Mandatory && o.DefaultValue == "" {
			// no default value and no value in environment
			return fmt.Errorf("error: missing mandatory key: %s", o.Key)
		}
	} else {
		// if we do get a valid value from the passed map, the default value is
		// overwritten then
		// pseudo options do not evaluate the value, but get the value from somewhere else other than the passed
		// string map. They might prompt the user via the shell, read some file etc.
		if err := tryParse(value, o.ParseFunction); err != nil {
			return fmt.Errorf("error in value of option '%s': %w", o.Key, err)
		}
	}

	if err := tryExecAction(o.PostParseAction); err != nil {
		return fmt.Errorf("post parse action of option '%s': %w", o.Key, err)
	}
	return nil
}

// Unparse executes the PreUnparseAction, UnparseFunction and the PostUnparseAction
// Depending on the UnparseFunction we get a string representation of a value back.
// This is the inverse operation of parsing, usually a serialization operation
// The returned value represents the value at the key o.Key in the resulting env map
// that is of type map[string]string.
// This function may return a configo.ErrSkipUnparse error which indicates that
// that the returned value is not added to any map or that we do not want to unparse (serialize)
// any values of this option struct.
func (o *Option) Unparse() (string, error) {

	err := tryExecAction(o.PreUnparseAction)
	if err != nil {
		if errors.Is(err, ErrSkipUnparse) {
			return "", ErrSkipUnparse
		}
		return "", fmt.Errorf("pre unparse action of the option '%s': %w", o.Key, err)
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
		// allow user to manually decide whether to use or not to use the value
		// this is important in the case that we do define a lot of sane default values
		// in our application that do not necessarily need to be written to the config map
		return value, ErrSkipUnparse
	}

	// PostUnparseAction may be used to close connections, file handles, etc.
	err = tryExecAction(o.PostUnparseAction)
	if err != nil {
		if errors.Is(err, ErrSkipUnparse) {
			return "", ErrSkipUnparse
		}
		// at this point we cannot skip the unparsing(serialization),
		// as it has already happened.
		return "", fmt.Errorf("post unparse action of the option '%s': %w", o.Key, err)
	}

	return value, nil
}

// Options are usually unique, so one MUST NOT use redundant Option parameters
type Options []Option

// String returns a string representation of the option. w/o any function pointers
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
