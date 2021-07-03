package parsers

import (
	"encoding/json"
	"io/ioutil"

	"github.com/joho/godotenv"
	configo "github.com/jxsl13/simple-configo"
	"gopkg.in/yaml.v3"
)

// ReadDotEnvFileMulti is basically the ability to use the whole of simple configo
// in a nested way for .env files as a single ParseFunc inside of a higher level
// Option.
func ReadDotEnvFileMulti(options ...configo.Option) configo.ParserFunc {
	return func(filePath string) error {
		env, err := godotenv.Read(filePath)
		if err != nil {
			return err
		}
		return configo.ParseOptions(options, env)
	}
}

// ReadDotEnvFileMap expects all of the provided keys to exist in the .env file
// those keys are associated with their output string pointers that are filled with the key's values
// upon parsing of the file.
func ReadDotEnvFileMap(outMap map[string]*string) configo.ParserFunc {
	return func(filePath string) error {
		options := make(configo.Options, 0, len(outMap))
		for key, outPtr := range outMap {
			options = append(options, configo.Option{
				Key:           key,
				Mandatory:     true,
				ParseFunction: String(outPtr),
			})
		}
		f := ReadDotEnvFileMulti(options...)
		return f(filePath)
	}
}

// ReadDotEnvFile reads a single variable from the .env file provided in the surrounding option's DefaultValue or from the environment.
// This functions allows to fetch a single value from the file.
// If you want to fetch multiple values at once, you may want to try either the easy one: ReadDotEnvFileMap
// or the slightly more complex version: ReadDotEnvFileMulti
func ReadDotEnvFile(out *string, key string) configo.ParserFunc {
	return func(filePath string) error {
		f := ReadDotEnvFileMulti(configo.Option{
			Key:           key,
			ParseFunction: String(out),
		})
		return f(filePath)
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
