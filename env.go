package configo

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/jxsl13/simple-configo/internal"
)

// GetEnv returns a map of OS environment variables
func GetEnv() map[string]string {
	pairs := os.Environ()
	env := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		keyPairs := strings.SplitN(pair, "=", 2)
		if len(keyPairs) == 2 {
			env[keyPairs[0]] = keyPairs[1]
		}
	}
	return env
}

// SetEnv sets the environment variables found in the env map.
func SetEnv(env map[string]string) {
	for k, v := range env {
		os.Setenv(k, v)
	}
}

// ReadEnvFile allows to read the env map from a key value file
// File content: key=value
func ReadEnvFile(filePathOrEnvKey string) (map[string]string, error) {
	filePath := getFilePathOrKey(GetEnv(), filePathOrEnvKey)
	return godotenv.Read(filePath)
}

// WriteEnvFile writes the map content into an env file
func WriteEnvFile(env map[string]string, filePathOrEnvKey string) error {
	filePath := getFilePathOrKey(GetEnv(), filePathOrEnvKey)

	// try creating folder incase it's needed
	err := internal.MkdirAll(filePath)
	if err != nil {
		return err
	}
	return godotenv.Write(env, filePath)
}

// UpdateEnvFile reads the file and update sits content to the new values.
func UpdateEnvFile(env map[string]string, filePathOrEnvKey string) error {
	filePath := getFilePathOrKey(GetEnv(), filePathOrEnvKey)
	var err error
	old := map[string]string{}

	// try creating folder incase it's needed
	err = internal.MkdirAll(filePath)
	if err != nil {
		return err
	}

	// check if .env file exists
	if internal.Exists(filePath) {
		// update old values in case we can read the env file
		old, err = godotenv.Read(filePath)
		if err != nil {
			return err
		}
	}
	// update map and write back to filePath location
	return godotenv.Write(update(old, env), filePath)
}

func update(old, new map[string]string) map[string]string {
	m := old
	for k, v := range new {
		m[k] = v
	}
	return m
}
