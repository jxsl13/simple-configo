package parsers

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/internal"
)

var (
	// ErrPathNotFile is returned when a given path is not a file
	ErrPathNotFile = errors.New("not a file")
	// ErrPathNotDirectory is returned when an expected path is not a directory
	ErrPathNotDirectory = errors.New("not a directory")
	// ErrPathExistsButNotDir is returned when we expect a directory but get a file at a specific file path
	ErrPathExistsButNotDir = errors.New("path exists but is not a directory")
	// ErrPathExistsButNotFile is returned when we expect a file but get a directory at a specific file path
	ErrPathExistsButNotFile = errors.New("path exists but is not a file")
	// ErrPathFailedToCreateDirs is returned when the creation of a subfolder structure failed
	ErrPathFailedToCreateDirs = errors.New("failed to create directories")
	// ErrPathFailedToConstructAbs is returned when we fail to construct an absolute file path from a relative one.
	ErrPathFailedToConstructAbs = errors.New("failed to construct absolute path")
)

// PathDirectory checks whether the value (e.g. environment value, map string value) is a valid
// directory path. Returns an error in case that the path location does not exist or is not a directory.
func PathDirectory(outPath *string) configo.ParserFunc {
	internal.PanicIfNil(outPath)
	return func(value string) error {
		normalizedPath := filepath.ToSlash(value)
		fi, err := os.Stat(normalizedPath)
		if err != nil {
			return err
		}
		if !fi.IsDir() {
			return fmt.Errorf("%w: %s", ErrPathNotDirectory, value)
		}
		*outPath = normalizedPath
		return nil
	}
}

// PathDirectoryCreate checks whether the value (e.g. environment value, map string value) is a valid
// directory path. Returns an error in case that the path location does not exist or is not a directory.
// dirPerms is set to a single default value. You may provide a single octal 0700 directory permission that
// is applied to all newly created sub directories.
func PathDirectoryCreate(outPath *string, dirPerms ...fs.FileMode) configo.ParserFunc {
	internal.PanicIfNil(outPath)
	return func(value string) error {
		normalizedPath := filepath.ToSlash(value)
		fi, err := os.Stat(normalizedPath)
		if err == nil {
			// dir already exists
			if fi.IsDir() {
				return nil
			}
			return fmt.Errorf("%w: %s", ErrPathExistsButNotDir, normalizedPath)
		}

		// does not exist, gotta create
		err = internal.MakeAllDir(normalizedPath, dirPerms...)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrPathFailedToCreateDirs, err)
		}
		*outPath = normalizedPath
		return nil
	}
}

// PathAbsDirectory checks whether the value (e.g. environment value, map string value) is a valid
// directory path. Returns an error in case that the path location does not exist or is not a directory.
// The returned outPath is an absolute path.
func PathAbsDirectory(outPath *string) configo.ParserFunc {
	internal.PanicIfNil(outPath)
	return func(value string) error {
		normalizedPath := filepath.ToSlash(value)
		absolutePath, err := filepath.Abs(normalizedPath)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrPathFailedToConstructAbs, err)
		}

		fi, err := os.Stat(absolutePath)
		if err != nil {
			return err
		}
		if !fi.IsDir() {
			return fmt.Errorf("%w: %s", ErrPathNotDirectory, absolutePath)
		}
		*outPath = absolutePath
		return nil
	}
}

// PathAbsDirectoryCreate checks whether the value (e.g. environment value, map string value) is a valid
// directory path. Returns an error in case that the path location does not exist or is not a directory.
// dirPerms is set to a single default value. You may provide a single octal 0700 directory permission that
// is applied to all newly created sub directories.
// The returned directory path is an absolute path.
func PathAbsDirectoryCreate(outPath *string, dirPerms ...fs.FileMode) configo.ParserFunc {
	internal.PanicIfNil(outPath)
	return func(value string) error {
		normalizedPath := filepath.ToSlash(value)
		absolutePath, err := filepath.Abs(normalizedPath)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrPathFailedToConstructAbs, err)
		}

		fi, err := os.Stat(absolutePath)
		if err == nil {
			// dir already exists
			if fi.IsDir() {
				return nil
			}
			return fmt.Errorf("%w: %s", ErrPathExistsButNotDir, absolutePath)
		}

		// does not exist, gotta create
		err = internal.MakeAllDir(absolutePath, dirPerms...)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrPathFailedToCreateDirs, err)
		}
		*outPath = absolutePath
		return nil
	}
}

// PathFile checks whether the value (e.g. environment value, map string value) is a valid
// file path. Returns an error in case that the path location does not exist or is not a file.
func PathFile(outPath *string) configo.ParserFunc {
	internal.PanicIfNil(outPath)
	return func(value string) error {
		normalizedPath := filepath.ToSlash(value)
		fi, err := os.Stat(normalizedPath)
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return fmt.Errorf("%w: %s", ErrPathNotFile, normalizedPath)
		}
		*outPath = normalizedPath
		return nil
	}
}

// PathAbsFile checks whether the value (e.g. environment value, map string value) is a valid
// file path. Returns an error in case that the path location does not exist or is not a file.
// The returned outPath is an absolute path.
func PathAbsFile(outPath *string) configo.ParserFunc {
	internal.PanicIfNil(outPath)
	return func(value string) error {
		normalizedPath := filepath.ToSlash(value)
		absolutePath, err := filepath.Abs(normalizedPath)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrPathFailedToConstructAbs, err)
		}

		fi, err := os.Stat(absolutePath)
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return fmt.Errorf("%w: %s", ErrPathNotFile, absolutePath)
		}
		*outPath = absolutePath
		return nil
	}
}

// PathFileDirectoryCreate checks whether the value (e.g. environment value, map string value) is a valid
// file path. In case that the file path does not exist, all directories of that file path are created
// In case that the final file already exists, we do not creat eanything.
// Beware, the file itself is never created with this, only the folder structure that is supposed to
// contain the file.
// dirPerms takes the first passed octal formated permission or defaults to 0700.
func PathFileDirectoryCreate(outFilePath *string, dirPerms ...fs.FileMode) configo.ParserFunc {
	internal.PanicIfNil(outFilePath)
	return func(value string) error {
		normalizedPath := filepath.ToSlash(value) // file path
		fi, err := os.Stat(normalizedPath)
		if err == nil {
			// file already exists
			if !fi.IsDir() {
				// is file, we don't do anything
				return nil
			}
			return fmt.Errorf("%w: %s", ErrPathExistsButNotDir, normalizedPath)
		}

		// does not exist, gotta create the folder structure for that file
		err = internal.MkdirAll(normalizedPath, dirPerms...)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrPathFailedToCreateDirs, err)
		}
		*outFilePath = normalizedPath
		return nil
	}
}

// PathFileDirectoryCreate checks whether the value (e.g. environment value, map string value) is a valid
// file path. In case that the file path does not exist, all directories of that file path are created
// In case that the final file already exists, we do not creat eanything.
// Beware, the file itself is never created with this, only the folder structure that is supposed to
// contain the file.
// dirPerms takes the first passed octal formated permission or defaults to 0700.
func PathAbsFileDirectoryCreate(outAbsFilePath *string, dirPerms ...fs.FileMode) configo.ParserFunc {
	internal.PanicIfNil(outAbsFilePath)
	return func(value string) error {
		normalizedPath := filepath.ToSlash(value) // file path
		absolutePath, err := filepath.Abs(normalizedPath)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrPathFailedToConstructAbs, err)
		}
		fi, err := os.Stat(absolutePath)
		if err == nil {
			// file already exists
			if !fi.IsDir() {
				// is file, we don't do anything
				return nil
			}
			return fmt.Errorf("%w: %s", ErrPathExistsButNotFile, absolutePath)
		}

		// does not exist, gotta create the folder structure for that file
		err = internal.MkdirAll(absolutePath, dirPerms...)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrPathFailedToCreateDirs, err)
		}
		*outAbsFilePath = absolutePath
		return nil
	}
}
