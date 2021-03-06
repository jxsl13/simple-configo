package internal

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path"
)

// Exists reports whether the named file or directory exists.
func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Save allows to save the text at a given filePath
func Save(text, filePath string, perm ...fs.FileMode) error {
	var mode fs.FileMode = 0600
	if len(perm) > 0 {
		mode = perm[0]
	}
	dirPath := path.Dir(filePath)

	if !Exists(filePath) {
		if err := os.MkdirAll(dirPath, mode); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(filePath, []byte(text), mode)
}

// Load allows to load a text from a given filePath that points to a file
// which contains the text
func Load(filePath string) (text string, err error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// Deletechecks if the given filePath exists. In case is does exist,the file is deleted.
// In case it doe sexist and it is a directory, the path and all of its children are deleted.
// In case the path does not exist, this function returnsan os.ErrNotExist error.
func Delete(filePath string) (err error) {
	if !Exists(filePath) {
		return os.ErrNotExist
	}
	return os.RemoveAll(filePath)
}
