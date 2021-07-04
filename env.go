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

// ReadEnvFile allows to read the env map from a key value file
// File content: key=value
func ReadEnvFile(filePaths ...string) (map[string]string, error) {
	return godotenv.Read(filePaths...)
}

func WriteEnvFile(env map[string]string, filePath string) error {
	return godotenv.Write(env, filePath)
}
