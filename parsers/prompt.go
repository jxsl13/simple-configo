package parsers

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"

	configo "github.com/jxsl13/simple-configo"
	"golang.org/x/term"
)

var (
	stdinFd = int(os.Stdin.Fd())
)

func promptPassword(linePrefix string) (string, error) {

	fmt.Print(linePrefix)
	password, err := term.ReadPassword(stdinFd)

	if err != nil {
		return "", err
	}

	return string(password), nil
}

// loads or prompts and saves password
func loadOrPromptPassword(promptPrefix, filePath string, perm ...fs.FileMode) (string, error) {
	text, err := load(filePath)
	if err == nil {
		return text, nil
	}
	// could not load
	text, err = promptPassword(promptPrefix)
	if err != nil {
		return "", err
	}
	err = save(text, filePath, perm...)
	if err != nil {
		return "", err
	}
	return text, nil
}

// PromptPassword prompts the user for the password in case the to be parsed map value does not contain
// any string data, meaning the user is only prompted when the e.g. environment variable doe snot exist or is empty.
func PromptPassword(out *string, promptPrefix string) configo.ParserFunc {
	return func(value string) error {
		if value != "" {
			*out = value
			return nil
		}

		password, err := promptPassword(promptPrefix)
		fmt.Print("\n")
		if err != nil {
			return err
		}
		*out = password
		return nil
	}
}

// LoadOrPromptPassword tries to load the passed file and extract its string content or prompts the user in the shell
// for entering their password (invisible) and then saves it to the passed file.
func LoadOrPromptPassword(out *string, promptPrefix string, filePath *string, perm ...fs.FileMode) configo.ParserFunc {
	return func(value string) error {
		text, err := loadOrPromptPassword(promptPrefix, *filePath, perm...)
		if err != nil {
			return err
		}
		*out = text
		return nil
	}
}

// loads or prompts and saves password
func loadOrPromptText(promptPrefix, filePath string, perm ...fs.FileMode) (string, error) {
	text, err := load(filePath)
	if err == nil {
		return text, nil
	}
	// could not load
	text = promptText(promptPrefix)
	err = save(text, filePath, perm...)
	if err != nil {
		return "", err
	}
	return text, nil
}

func promptText(linePrefix string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(linePrefix)
	scanner.Scan()
	return scanner.Text()
}

// PromptText prompts the user to enter a text. This only prompts the user in the case that
// the corresponding environment variable doe snot contain any string data.
func PromptText(out *string, promptPrefix string) configo.ParserFunc {
	return func(value string) error {
		if value != "" {
			*out = value
			return nil
		}

		text := promptText(promptPrefix)
		*out = text
		return nil
	}
}

// LoadOrPromptText either loads the content of the filePath and sets the string value to the file's content or die prompt the user to enter
// the data and then saves the result in the specified file.
func LoadOrPromptText(out *string, promptPrefix string, filePath *string, perm ...fs.FileMode) configo.ParserFunc {
	return func(value string) error {
		text, err := loadOrPromptText(promptPrefix, *filePath, perm...)
		if err != nil {
			return err
		}
		*out = text
		return nil
	}
}
