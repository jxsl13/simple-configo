package configo

import "fmt"

// Config is an interface that implements only two methods.
// The first method simply returns the name of the configuration.
// The second method returns a list of Option objects that
// are everything that is needed to fill the struct fields of your
// custom interface implementation of Config.
type Config interface {
	Name() string
	Options() (options Options)
}

// Parse the passed envoronment map into the config struct.
// Every Config defines, how its Options look like and how those are parsed.
func Parse(cfg Config, env map[string]string) error {
	return parse(cfg.Options(), env)
}

// for internal usage in order not to call cfg.Options() multiple times.
func parse(options Options, env map[string]string) error {
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

// ParseWithUnparse parses the configuration the same way as the 'Parse' function and
// initializes the shutdown hooks that are defined in the configuration's options.
// The returned function is a blocking function that can be used in a 'defer unparse()' context where
// the function blocks until the application received any of the defined shutdown signals.
// once the signals are received, the function stops blocking and starts to execute the
// UnparseFunctions of all of the defined Options in REVERSE order of their occurrance in the Options array.
// One may also handle any errors that the returned function returns upon invokation.
func ParseWithUnparse(cfg Config, env map[string]string) (func() error, error) {
	// only call this function once in order not to cause any side effects wheh calling it again.
	options := cfg.Options()

	if err := parse(options, env); err != nil {
		return nil, err
	}
	return unparse(options, env), nil
}
