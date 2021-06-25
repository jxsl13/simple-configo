package unparsers

import (
	"encoding/json"
	"io/fs"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/internal"
	"gopkg.in/yaml.v3"
)

// WriteJSON marshals the passed structure or value to a json string and writes it to the
// parsed destination
func WriteJSON(in interface{}, perm ...fs.FileMode) configo.UnparserFunc {
	return func(key, value string) error {
		data, err := json.MarshalIndent(in, "", " ")
		if err != nil {
			return err
		}
		return internal.Save(string(data), value, perm...)
	}
}

// WriteFile will write the content of the 'in' string into the correstonding file defined in the environment
// variable that is the key of the option that defines this function as UnparserFunc
func WriteFile(in *string, perm ...fs.FileMode) configo.UnparserFunc {
	return func(key, value string) error {
		return internal.Save(*in, value, perm...)
	}
}

// WriteYAML writes the 'in' interface as a yaml formated file to the location thecified in the parent
// option's key's value. Basically -> option -> key -> value := env[key], whatever env is.
func WriteYAML(in interface{}, perm ...fs.FileMode) configo.UnparserFunc {
	return func(key, value string) error {

		data, err := yaml.Marshal(in)
		if err != nil {
			return err
		}
		return internal.Save(string(data), value, perm...)
	}
}
