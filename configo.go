package configo

import (
	"errors"
	"fmt"
	"sync"

	"github.com/joho/godotenv"
)

// Config is an interface that implements only two methods.
// The first method simply returns the name of the configuration.
// The second method returns a list of Option objects that
// are everything that is needed to fill the struct fields of your
// custom interface implementation of Config.
// WARNING: In case your configuration also implements the sync.Locker interface,
// you MUST NOT Lock()/RLock() your mutex in the Options() method.
type Config interface {
	Name() string
	Options() (options Options)
}

// ParseEnv parse the environment variables and fills all of the definied options on the
// configuration.
// INFO: goroutine safe if config implements the sync.Locker interface.
func ParseEnv(cfgs ...Config) error {
	return Parse(GetEnv(), cfgs...)
}

func getFilePathOrKey(env map[string]string, filePathOrEnvKey string) string {
	filePath := filePathOrEnvKey
	value, found := env[filePathOrEnvKey]
	if found && value != "" {
		filePath = value
	}
	return filePath
}

// ParseEnVFile parses either the file at the provided location filePathOrEnvKey
// or checks if the povided location is actually an environment variable pointing to a
// file location.
func ParseEnvFile(filePathOrEnvKey string, cfgs ...Config) error {
	env := GetEnv()
	filePath := getFilePathOrKey(env, filePathOrEnvKey)

	env, err := godotenv.Read(filePath)
	if err != nil {
		return err
	}
	return Parse(env, cfgs...)
}

// UnparseEnvFile is the opposite of ParseEnvFile. It serializes the map back into
// the file.
func UnparseEnvFile(filePathOrEnvKey string, cfgs ...Config) error {
	env, err := Unparse(cfgs...)
	if err != nil {
		return err
	}
	return UpdateEnvFile(env, filePathOrEnvKey)
}

// ParseEnvFileOrEnv tries to parse the env file first and then the environment in case the file
// parsing fails.
func ParseEnvFileOrEnv(filePathOrEnvKey string, cfgs ...Config) error {

	env := GetEnv()
	filePath := getFilePathOrKey(env, filePathOrEnvKey)

	fileMap, err := godotenv.Read(filePath)
	if err != nil {
		// parse environment
		return Parse(env, cfgs...)
	}
	// parse fileMap
	return Parse(fileMap, cfgs...)
}

// Parse the passed envoronment map into the config struct.
// Every Config defines, how its Options look like and how those are parsed.
// INFO: In case Config implements the sync.Locker inteface by either embedding an anonymous sync.Mutex or by
// implementing the methods func (cfg *Config) Lock() and func (cfg *Config) Unlock(), those methods are called before
// attempting to parse the option values.
// This allows
func Parse(env map[string]string, cfgs ...Config) error {
	for _, cfg := range cfgs {
		err := func(c Config) error {
			locker, ok := cfg.(sync.Locker)
			if ok {
				locker.Lock()
				defer locker.Unlock()
			}
			err := ParseOptions(cfg.Options(), env)
			if err != nil {
				return err
			}
			return nil
		}(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

// for internal usage in order not to call cfg.Options() multiple times.
// INFO: ParseOptions is not goroutine safe.
func ParseOptions(options Options, env map[string]string) error {
	for _, opt := range options {
		err := opt.Parse(env)
		if err != nil {
			return err
		}
	}

	return nil
}

// Unparse is the reverse operation of Parse. It retrieves the values from the configuration and
// serializes them to their respective string values in order to be able to writ ethem back to either
// the environment or to a file.
// INFO: In case cfg implements the sync.Locker interface by either embedding an anonymous sync.Mutex or
// implementing the Lock() and Unlock() methods, then those methods are called in order to guard the configuration values.
func Unparse(cfgs ...Config) (map[string]string, error) {
	resultMap := make(map[string]string)
	for _, cfg := range cfgs {
		// wrapped in a function call in order to directly unlock
		// the mutex after the parsing is done.
		err := func(c Config) error {
			locker, ok := c.(sync.Locker)
			if ok {
				locker.Lock()
				defer locker.Unlock() // unlocked on return
			}
			env, err := UnparseOptions(c.Options())
			if err != nil {
				return err
			}
			for k, v := range env {
				resultMap[k] = v
			}
			return nil
		}(cfg)
		if err != nil {
			return nil, err
		}

	}
	return resultMap, nil
}

// UnparseValidate unparses the values and tries to parse the values again in order to validate their values
// this allows to have a complex ParserFunction but a simple UnparserFunction, as all of the validation logic is
// provided via the ParserFunction.
// INFO: UnparseValidate is goroutine safe in case cfg implements the sync.Locker interface by either embedding
// the sync.Mutex struct anonymously or by implementing the Lock() and Unlock() methods.
func UnparseValidate(cfgs ...Config) (map[string]string, error) {
	resultEnv := make(map[string]string)
	for _, cfg := range cfgs {

		err := func(c Config) error {
			locker, ok := c.(sync.Locker)
			if ok {
				locker.Lock()
				defer locker.Unlock()
			}
			options := c.Options()
			env, err := UnparseOptions(options)
			if err != nil {
				return err
			}

			// validate through parse functions
			err = ParseOptions(options, env)
			if err != nil {
				return fmt.Errorf("failed to validate unparse options: %w", err)
			}
			// add to result map
			for k, v := range env {
				resultEnv[k] = v
			}
			return nil
		}(cfg)
		if err != nil {
			return nil, err
		}

	}
	return resultEnv, nil
}

// UnparseOptions returns a key value map from the parsed options.
// This is the reverse operation of ParseOptions.
// Only options that define a UnparserFunction are serialized into their string values.
// Options that do not differ from their default values are ignored in order to keep the returned map
// as small as possible.
// INFO: Not goroutine safe
func UnparseOptions(options Options) (map[string]string, error) {
	env := make(map[string]string, len(options))
	for _, opt := range options {

		value, err := opt.UnparseFunction()
		if err != nil {
			if errors.Is(err, ErrSkipUnparse) {
				continue
			}
			// unknown error
			return nil, err
		}
		// set map value
		env[opt.Key] = value
	}
	return env, nil
}
