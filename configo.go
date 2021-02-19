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

	options := cfg.Options()
	for _, opt := range options {

		// Initially the config values are set to the default value, if the default value is valid
		if err := opt.IsValid(); err != nil {
			return fmt.Errorf("The option definition for '%s' is invalid: %w", opt.Key, err)
		}

		value, ok := env[opt.Key]
		if opt.Mandatory {
			if opt.DefaultValue == "" && !ok {
				// no default value and no value in environment
				return fmt.Errorf("Error: missing mandatory key: %s", opt.Key)
			}
		}

		// if we do get a valid value from the passed map, the default value is
		// overwritten then
		if ok {
			if err := opt.ParseFunction(value); err != nil {
				return fmt.Errorf("Error in value of option '%s': %w", opt.Key, err)
			}
		}
	}

	return nil
}
