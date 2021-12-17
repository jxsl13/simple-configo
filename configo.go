package configo

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config is an interface that implements only two methods.
// The first method simply returns the name of the configuration.
// The second method returns a list of Option objects that
// are everything that is needed to fill the struct fields of your
// custom interface implementation of Config.
type Config interface {
	Options() (options Options)
}

// ParseEnv parse the environment variables and fills all of the definied options on the
// configuration.
func ParseEnv(cfgs ...Config) error {
	return Parse(GetEnv(), cfgs...)
}

// OptionDefaults returns a map of option keys and option default values
// for options that define at least a ParseFunction or UnparseFunction
func OptionDefaults(cfgs ...Config) map[string]string {
	m := make(map[string]string, len(cfgs)*2)
	for _, c := range cfgs {
		for _, opt := range c.Options() {
			if opt.IsOption() {
				m[opt.Key] = opt.DefaultValue
			}
		}
	}
	return m
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
// filePathOrEnvKey may be a file path or an environment key containing a file path
// In case a variable is not found in theenv file the next level is tried which is the environment.
func ParseEnvFileOrEnv(filePathOrEnvKey string, cfgs ...Config) error {
	env := GetEnv()
	filePath := getFilePathOrKey(env, filePathOrEnvKey)
	fileMap, err := godotenv.Read(filePath)
	if err != nil {
		// no file, parse env only
		return Parse(env, cfgs...)
	}
	// environment extends and overrides env file values
	env = update(fileMap, env)

	return Parse(env, cfgs...)
}

// ParseFlags parses the flags provided to the application based on the
// provided option definitions in every passed Config
func ParseFlags(cfgs ...Config) error {
	return parseFlags(os.Args[1:], cfgs...)
}

// allows to pass custom args for testing
func parseFlags(args []string, cfgs ...Config) error {
	flagMap, err := GetFlagMap(args, cfgs...)
	if err != nil {
		return err
	}
	return Parse(flagMap, cfgs...)
}

// ParseEnvOrFlags fetches config values from the .env file, the environment
// and from the flags and parses the configurations with those values provided as key value map.
func ParseEnvOrFlags(cfgs ...Config) error {
	return parseEnvOrFlags(os.Args[1:], cfgs...)
}

// parseEnvOrFlags allows passing of custom args for testing
func parseEnvOrFlags(args []string, cfgs ...Config) error {
	// override & extend env values with flag values
	flagMap, err := GetFlagMap(args, cfgs...)
	if err != nil {
		return err
	}
	env := update(GetEnv(), flagMap)

	// parse the combined map
	return Parse(env, cfgs...)
}

// ParseEnvFileOrEnvOrFlags fetches config values from the .env file, the environment
// and from the flags and parses the configurations with those values provided as key value map.
// Warning: do not call this function multiple times with the same configurations, as redefiition of flag names
// may cause a panic.
func ParseEnvFileOrEnvOrFlags(filePathOrEnvKey string, cfgs ...Config) error {
	return parseEnvFileOrEnvOrFlags(filePathOrEnvKey, os.Args[1:], cfgs...)
}

// parseEnvFileOrEnvOrFlags allows to pass custom os.Args[1:] for testing
func parseEnvFileOrEnvOrFlags(filePathOrEnvKey string, args []string, cfgs ...Config) error {
	// must always be parsed in order to fetch the potential file path
	env := GetEnv()
	flags, err := GetFlagMap(args, cfgs...)
	if err != nil {
		return err
	}

	filePath := getFilePathOrKey(env, filePathOrEnvKey)
	fileMap, err := godotenv.Read(filePath)
	if err != nil {
		// parse environment variables extended and overwritten by flag parameters
		return Parse(update(env, flags), cfgs...)
	}

	// override and update .env file with environment variables
	// override and update .env file and environment variables with flag values
	env = update(update(fileMap, env), flags)

	// parse the combined map
	return Parse(env, cfgs...)
}

// Parse the passed envoronment map into the config struct.
// Every Config defines, how its Options look like and how those are parsed.
func Parse(env map[string]string, cfgs ...Config) error {
	for _, cfg := range cfgs {
		err := ParseOptions(cfg.Options(), env)
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
func Unparse(cfgs ...Config) (map[string]string, error) {
	resultMap := make(map[string]string)
	for _, cfg := range cfgs {
		env, err := UnparseOptions(cfg.Options())
		if err != nil {
			return nil, err
		}
		for k, v := range env {
			resultMap[k] = v
		}
	}
	return resultMap, nil
}

// UnparseValidate unparses the values and tries to parse the values again in order to validate their values
// this allows to have a complex ParserFunction but a simple UnparserFunction, as all of the validation logic is
// provided via the ParserFunction.
func UnparseValidate(cfgs ...Config) (map[string]string, error) {
	resultEnv := make(map[string]string)
	for _, cfg := range cfgs {

		err := func(c Config) error {
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

		value, err := opt.Unparse()
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
