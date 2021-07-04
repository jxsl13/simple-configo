package configo

import (
	"os"
	"strings"
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
