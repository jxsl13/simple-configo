package unparsers

import (
	"fmt"
	"path/filepath"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/internal"
)

// Path returns the host specific path representation
func Path(inPath *string) configo.UnparserFunc {
	internal.PanicIfNil(inPath)
	return func() (string, error) {
		return filepath.FromSlash(*inPath), nil
	}
}

// PathAbs returns the host specific path representation as an absolute path
func PathAbs(inPath *string) configo.UnparserFunc {
	internal.PanicIfNil(inPath)
	return func() (string, error) {
		normalizedPath := filepath.FromSlash(*inPath)
		absolutePath, err := filepath.Abs(normalizedPath)
		if err != nil {
			return "", fmt.Errorf("failed to construct absolute path: %w", err)
		}
		return absolutePath, nil
	}
}
