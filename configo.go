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
		value, ok := env[opt.Key]
		if opt.Mandatory {
			if !ok {
				return fmt.Errorf("Error: missing mandatory key: %s", opt.Key)
			}
		}

		if !ok {
			if err := opt.ParseFunction(opt.DefaultValue); err != nil {
				return err
			}
		} else {
			if err := opt.ParseFunction(value); err != nil {
				return err
			}
		}
	}

	return nil
}
