package configo

import (
	"flag"
	"regexp"
	"strings"
)

// GetFlags parses the provided fags according to your configo definitions.
// Warning: You must not redefine flag values as the flag.Parse function does not like that.
func GetFlags(cfgs ...Config) map[string]string {

	options := Options{}
	for _, cfg := range cfgs {
		options = append(options, cfg.Options()...)
	}
	envP := make(map[string]*string, len(options))

	for _, opt := range options {
		if opt.IsAction() {
			// skip actions that do not have value parsing logic
			continue
		}

		flagName := transformKeyToFlagName(opt.Key)
		// define a flag for every option, do not set default values,
		// as the resulting map should not contain value sin case they are not set as flags
		envP[opt.Key] = flag.String(flagName, "", opt.Description)
	}
	flag.Parse()

	env := make(map[string]string, len(envP))
	for k, v := range envP {
		if v != nil {
			value := *v
			if value != "" {
				env[k] = value
			}
		}
	}
	return env
}

var (
	flagNameTransformer       = regexp.MustCompile(`\s+|_+`)
	trimAndSpecialTransformer = regexp.MustCompile(`^-+|-+$|,|\.+|;+|:+`)
)

func transformKeyToFlagName(key string) string {
	lcKey := strings.TrimSpace(strings.ToLower(key))
	flagName := flagNameTransformer.ReplaceAllString(lcKey, "-")
	flagName = trimAndSpecialTransformer.ReplaceAllString(flagName, "")
	return flagName
}
