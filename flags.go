package configo

import (
	"flag"
	"os"
	"regexp"
	"strings"
)

var (
	// KeyToFlagNameTransformer is the function that takes an Option.Key and transforms it into a
	// cli flag name. e.g. FLAG_NAME to --flag-name for any Option that is NOT an .IsAction()
	KeyToFlagNameTransformer = DefaultKeyToFlagNameTransformer

	trimAndSpecialTransformer = regexp.MustCompile(`^-+|-+$|,|\.+|;+|:+`)
	flagNameTransformer       = regexp.MustCompile(`\s+|_+`)
)

func DefaultKeyToFlagNameTransformer(key string) string {
	lcKey := strings.TrimSpace(strings.ToLower(key))
	// replace spaces & underscores with single dash
	flagName := flagNameTransformer.ReplaceAllString(lcKey, "-")
	// remove special characters
	flagName = trimAndSpecialTransformer.ReplaceAllString(flagName, "")
	return flagName
}

// GetFlags parses your config options and extracts allof the option's keys, constructs
// flag names from the option keys, parses the os.Args and then returns a map with all
// non null flag values
// This function may exit the application in case the flag parsing fails somehow.
func GetFlags(cfgs ...Config) map[string]string {
	flagMap, _ := getFlagMapWithErrorHandling(os.Args[1:], flag.ExitOnError, cfgs...)
	return flagMap
}

// GetFlagMap returns a map of flags that consists of flag values passed via osArgs that can be found in
// the cfg Options' keys.
func GetFlagMap(osArgs []string, cfgs ...Config) (map[string]string, error) {
	return getFlagMapWithErrorHandling(osArgs, flag.ContinueOnError, cfgs...)
}

// getFlagMapWithErrorHandling parses the provided args according to your configo definitions.
func getFlagMapWithErrorHandling(osArgs []string, errHandling flag.ErrorHandling, cfgs ...Config) (map[string]string, error) {
	options := Options{}
	for _, cfg := range cfgs {
		options = append(options, cfg.Options()...)
	}

	// single key -> last description of that key
	flagDefinitions := make(map[string]string, len(options))
	envP := make(map[string]*string, len(options))

	for _, opt := range options {
		if opt.IsAction() {
			// skip actions that do not have value parsing logic
			continue
		}

		_, found := flagDefinitions[opt.Key]
		if !found || (found && opt.Description != "") {
			flagDefinitions[opt.Key] = opt.Description
		}
	}

	flags := flag.NewFlagSet("", errHandling)
	for key, description := range flagDefinitions {
		flagName := KeyToFlagNameTransformer(key)
		envP[key] = flags.String(flagName, "", description)
	}

	err := flags.Parse(osArgs)
	if err != nil {
		return nil, err
	}

	env := make(map[string]string, len(envP))
	for k, v := range envP {
		if v != nil && *v != "" {
			env[k] = *v
		}
	}
	return env, nil
}
