package parsers

import (
	"encoding/json"
	"io/ioutil"

	"github.com/joho/godotenv"
	configo "github.com/jxsl13/simple-configo"
	"gopkg.in/yaml.v3"
)

//DotEnvTuple is used as a bridging struct that connects Keys found in a .env file and the corresponding destination values
// where the values of these keys are written into.
type DotEnvTuple struct {
	// Key .env -> USERNAME=xxx, username is the key
	Key string
	// KeyPtr is the dynamic version of the key that might be used to change the Key at runtime.
	KeyPtr *string
	//OutValuePtr is the target string variable that receives the value parsed from the .env file.
	OutValuePtr *string
}

// ReadDotEnvFileMulti reads the provided file and sets the DotEnvTuple.Value variable to the value found in the .env file
// under the DotEnvTuple.Key.
func ReadDotEnvFileMulti(outTuples ...DotEnvTuple) configo.ParserFunc {
	return func(filePath string) error {

		env, err := godotenv.Read(filePath)
		if err != nil {
			return err
		}
		for _, tuple := range outTuples {
			keyPtr := tuple.KeyPtr
			key := tuple.Key
			outPtr := tuple.OutValuePtr

			if keyPtr != nil {
				*outPtr = env[*keyPtr]
			} else {
				// set out to expected value, will be empty in case the value was not found.
				*outPtr = env[key]
			}
		}

		return nil
	}
}

// ReadDotEnvFileMulti reads the provided file and sets the DotEnvTuple.Value variable to the value found in the .env file
// under the DotEnvTuple.Key.
func ReadDotEnvFileMap(outMap map[string]*string) configo.ParserFunc {
	return func(filePath string) error {
		tuples := make([]DotEnvTuple, 0, len(outMap))
		for key, outPtr := range outMap {
			tuples = append(tuples, DotEnvTuple{
				Key:         key,
				OutValuePtr: outPtr,
			})
		}
		f := ReadDotEnvFileMulti(tuples...)
		return f(filePath)
	}
}

// ReadDotEnvFile reads a single variable from the .env file provided in the surrounding option's DefaultValue or from the environment.
// This functions allows to fetch a single value from the file.
// If you want to fetch multiple values at once, you may want to try either the easy one: ReadDotEnvFileMap
// or the slightly more complex version: ReadDotEnvFileMulti
func ReadDotEnvFile(out *string, key string) configo.ParserFunc {
	return func(filePath string) error {
		f := ReadDotEnvFileMulti(DotEnvTuple{
			Key:         key,
			OutValuePtr: out,
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
