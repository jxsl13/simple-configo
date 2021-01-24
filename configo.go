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
		if err := opt.IsValid(); err != nil {
			return fmt.Errorf("The option definition for '%s' is invalid: %w", opt.Key, err)
		}

		value, ok := env[opt.Key]
		if opt.Mandatory {
			if !ok {
				return fmt.Errorf("Error: missing mandatory key: %s", opt.Key)
			}
		}

		// always write the default value in order to check if the programmed values are actually
		// properly set to valid values.
		// Only check this when the value is not mandatory, as you may have invalid
		// default values, because you expect user input!
		if err := opt.ParseFunction(opt.DefaultValue); !opt.Mandatory && err != nil {
			return fmt.Errorf("Error in default value of option '%s': %w", opt.Key, err)
		}

		// overwrite default value with config value
		if ok {
			if err := opt.ParseFunction(value); err != nil {
				return fmt.Errorf("Error in value of option '%s': %w", opt.Key, err)
			}
		}
	}

	return nil
}
