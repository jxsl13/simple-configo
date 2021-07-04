package configo

import (
	"fmt"
)

// Config is an interface that implements only two methods.
// The first method simply returns the name of the configuration.
// The second method returns a list of Option objects that
// are everything that is needed to fill the struct fields of your
// custom interface implementation of Config.
type Config interface {
	Name() string
	Options() (options Options)
}

// ParseEnv parse the environment variables and fills all of the definied options on the
// configuration.
func ParseEnv(cfg Config) error {
	return Parse(cfg, GetEnv())
}

// Parse the passed envoronment map into the config struct.
// Every Config defines, how its Options look like and how those are parsed.
func Parse(cfg Config, env map[string]string) error {
	return ParseOptions(cfg.Options(), env)
}

// for internal usage in order not to call cfg.Options() multiple times.
func ParseOptions(options Options, env map[string]string) error {
	for _, opt := range options {

		// Initially the config values are set to the default value, if the default value is valid
		// pseudo options are not checked for valid keys or descriptions, nor whether their defaultvalues
		// can be successfully parsed with the provided ParseFunc.
		if err := opt.IsValid(); err != nil {
			return fmt.Errorf("the option definition for '%s' is invalid: %w", opt.Key, err)
		}

		value, ok := env[opt.Key]
		if opt.Mandatory {
			if opt.DefaultValue == "" && !ok {
				// no default value and no value in environment
				return fmt.Errorf("error: missing mandatory key: %s", opt.Key)
			}
		}

		// if we do get a valid value from the passed map, the default value is
		// overwritten then
		// pseudo options do not evaluate the value, but get the value from somewhere else other than the passed
		// string map. They might prompt the user via the shell, read some file etc.
		if ok || opt.IsPseudoOption {
			if err := opt.ParseFunction(value); err != nil {
				return fmt.Errorf("error in value of option '%s': %w", opt.Key, err)
			}
		}
	}

	return nil
}

// Unparse is the reverse operation of Parse. It retrieves the values from the configuration and
// serializes them to their respective string values in order to be able to writ ethem back to either
// the environment or to a file.
// This is usually necessary when you want to
func Unparse(cfg Config) (map[string]string, error) {
	// TODO: check if cfg implements the Lock/Unlock interface in order to lock and unlock it in
	// case we unparse values that can be accessed in amultithreaded context
	return UnparseOptions(cfg.Options())
}

// UnparseOptions returns a key value map from the parsed options.
// This is the reverse operation of ParseOptions.
// Only options that define a UnparserFunction are serialized into their string values.
// Options that do not differ from their default values are ignored in order to keep the returned map
// as small as possible.
func UnparseOptions(options Options) (map[string]string, error) {
	env := make(map[string]string, len(options))
	for _, opt := range options {

		// also validate options in this function.
		if err := opt.IsValid(); err != nil {
			return nil, fmt.Errorf("the option definition for '%s' is invalid: %w", opt.Key, err)
		}

		// skip options that do not have an UnprserFunction
		if opt.UnparseFunction == nil {
			continue
		}
		// Unparse (serialize) option values
		value, err := opt.UnparseFunction()
		if err != nil {
			return nil, fmt.Errorf("error while unparsing the option '%s': %w", opt.Key, err)
		}

		// skip default values in order to keep the config file/env variables map small.
		if value == opt.DefaultValue {
			continue
		}

		env[opt.Key] = value
	}
	return env, nil
}
