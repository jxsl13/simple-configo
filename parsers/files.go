package parsers

import (
	"encoding/json"
	"io/ioutil"

	"github.com/joho/godotenv"
	configo "github.com/jxsl13/simple-configo"
	"gopkg.in/yaml.v3"
)

// ReadDotEnvFile reads the provided files and sets the out variable to the expected key provided in the surrounding option.
// The key then is used to lookup in the files.
func ReadDotEnvFile(out *string, filePaths ...*string) configo.ParserFunc {
	return func(key string) error {
		lookupPaths := make([]string, 0, len(filePaths))
		for _, p := range filePaths {
			lookupPaths = append(lookupPaths, *p)
		}
		env, err := godotenv.Read(lookupPaths...)
		if err != nil {
			return err
		}
		// set out to expected value, will be empty in case the value was not found.
		*out = env[key]
		return nil
	}
}

// ReadFile will read the file specified at the path that is provided via the defined
// environment variable. So the location of the file can be altered hoever the user wants.
// A default value allows to set a static value without the user having to fiddle around
// in their environment variables.
func ReadFile(out *string) configo.ParserFunc {
	return func(value string) error {
		data, err := ioutil.ReadFile(value)
		if err != nil {
			return err
		}

		*out = string(data)
		return nil
	}
}

// ReadYAML reads the YAML file content from
// the file provided via environment key.
func ReadYAML(out interface{}) configo.ParserFunc {
	return func(value string) error {
		data, err := ioutil.ReadFile(value)
		if err != nil {
			return err
		}
		return yaml.Unmarshal(data, out)
	}
}

// ReadJSON reads the JSON content of the file provided via the environment variable.
func ReadJSON(out interface{}) configo.ParserFunc {
	return func(value string) error {
		data, err := ioutil.ReadFile(value)
		if err != nil {
			return err
		}
		return json.Unmarshal(data, out)
	}
}
