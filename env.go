package configo

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
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
	return godotenv.Write(env, filePath)
}
