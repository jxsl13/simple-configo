package parsers

import (
	"encoding/json"
	"io/ioutil"

	configo "github.com/jxsl13/simple-configo"
	"gopkg.in/yaml.v3"
)

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
