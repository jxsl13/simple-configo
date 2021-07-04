package configo

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
)

var (
	// Signals are the default signals used in the UnparseOnSignal functions.
	// You may change these variables in order to add or remove signals.
	Signals = []os.Signal{os.Interrupt, syscall.SIGTERM}
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
func ParseEnv(cfg Config) error {
	return Parse(cfg, GetEnv())
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
func ParseEnvFile(cfg Config, filePathOrEnvKey string) error {
	env := GetEnv()
	filePath := getFilePathOrKey(env, filePathOrEnvKey)

	env, err := godotenv.Read(filePath)
	if err != nil {
		return err
	}
	return Parse(cfg, env)
}

// UnparseEnvFile writes the result either to the provided file path or in case
// the path is actually an environment key, then it is written to the path that is found in the environment
// under that key.
func UnparseEnvFile(filePathOrEnvKey string) UnparserHook {
	return func(envMap map[string]string, err error) {
		if err != nil {
			log.Println(err)
			return
		}

		env := GetEnv()
		filePath := getFilePathOrKey(env, filePathOrEnvKey)

		err = godotenv.Write(envMap, filePath)
		if err != nil {
			log.Println(err)
		}
	}
}

// ParseEnvFileOrEnv tries to parse the
func ParseEnvFileOrEnv(cfg Config, filePathOrEnvKey string) error {

	env := GetEnv()
	filePath := getFilePathOrKey(env, filePathOrEnvKey)

	fileMap, err := godotenv.Read(filePath)
	if err != nil {
		// parse environment
		return Parse(cfg, env)
	}
	// parse fileMap
	return Parse(cfg, fileMap)
}

// Parse the passed envoronment map into the config struct.
// Every Config defines, how its Options look like and how those are parsed.
// INFO: In case Config implements the sync.Locker inteface by either embedding an anonymous sync.Mutex or by
// implementing the methods func (cfg *Config) Lock() and func (cfg *Config) Unlock(), those methods are called before
// attempting to parse the option values.
// This allows
func Parse(cfg Config, env map[string]string) error {
	locker, ok := cfg.(sync.Locker)
	if ok {
		locker.Lock()
		defer locker.Unlock()
	}
	return ParseOptions(cfg.Options(), env)
}

// for internal usage in order not to call cfg.Options() multiple times.
// INFO: ParseOptions is not goroutine safe.
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

// UnparseOnShutdown writes the content of the config into provided file path or the filepath provided
// in the filePathOrEnvKey environment variable.
func UnparseOnShutdown(cfg Config, filePathOrEnvKey string) {
	UnparseOnSignal(cfg, UnparseEnvFile(filePathOrEnvKey))
}

// UnparseValidateOnShutdown is the same as UnparseOnShutdwon, but it also validates
// the content of the resultung key value map by parsing that map again.
func UnparseValidateOnShutdown(cfg Config, filePathOrEnvKey string) {
	UnparseValidateOnSignal(cfg, UnparseEnvFile(filePathOrEnvKey))
}

// UnparseOnSignal starts a goroutine that waits until any of Signals defined in the Signals package variable
// are received. Before just before the application terminates due to receiving a signal (not os.Exit(..))
// then the unparsing of the configuration is executed, the resulting map can be then processed in the
// callback function.
// the default callback function that can be used is UnparseEnvFile in order to write the resulting map
// to an .env file.
func UnparseOnSignal(cfg Config, callback UnparserHook) {
	go func(c Config, sigs ...os.Signal) {
		notify := make(chan os.Signal, 1)
		signal.Notify(notify, sigs...)
		<-notify
		callback(Unparse(c))
	}(cfg, Signals...)
}

// UnparseValidateOnSignal is the same as UnparseOnSignal but after the unparsing it parses the resulting
// map again in order to validate its content according to the ParseFunctions.
// This might be useful in order to validate integer ranges that have been changed while the app has been running.
func UnparseValidateOnSignal(cfg Config, callback UnparserHook) {
	go func(c Config, sigs ...os.Signal) {
		notify := make(chan os.Signal, 1)
		signal.Notify(notify, sigs...)
		<-notify
		callback(UnparseValidate(c))
	}(cfg, Signals...)
}

// Unparse is the reverse operation of Parse. It retrieves the values from the configuration and
// serializes them to their respective string values in order to be able to writ ethem back to either
// the environment or to a file.
// INFO: In case cfg implements the sync.Locker interface by either embedding an anonymous sync.Mutex or
// implementing the Lock() and Unlock() methods, then those methods are called in order to guard the configuration values.
func Unparse(cfg Config) (map[string]string, error) {
	locker, ok := cfg.(sync.Locker)
	if ok {
		locker.Lock()
		defer locker.Unlock()
	}
	return UnparseOptions(cfg.Options())
}

// UnparseValidate unparses the values and tries to parse the values again in order to validate their values
// this allows to have a complex ParserFunction but a simple UnparserFunction, as all of the validation logic is
// provided via the ParserFunction.
// INFO: UnparseValidate is goroutine safe in case cfg implements the sync.Locker interface by either embedding
// the sync.Mutex struct anonymously or by implementing the Lock() and Unlock() methods.
func UnparseValidate(cfg Config) (map[string]string, error) {
	locker, ok := cfg.(sync.Locker)
	if ok {
		locker.Lock()
		defer locker.Unlock()
	}
	options := cfg.Options()
	env, err := UnparseOptions(options)
	if err != nil {
		return nil, err
	}

	// validate through parse functions
	err = ParseOptions(options, env)
	if err != nil {
		return nil, fmt.Errorf("failed to validate unparse options: %w", err)
	}
	// return map
	return env, nil
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
			if errors.Is(err, ErrSkipUnparse) {
				// skip unparsing in case the function returns the skip error.
				continue
			}
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
