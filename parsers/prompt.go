package parsers

import (
	"strings"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/internal"
	"github.com/manifoldco/promptui"
)

// PromptPassword prompts the user for the password in case the to be parsed map value does not contain
// any string data, meaning the user is only prompted when the e.g. environment variable doe snot exist or is empty.
func PromptPassword(out, promptPrefix *string, validateFunc ...func(string) error) configo.ParserFunc {
	internal.PanicIfNil(out, promptPrefix)

	return func(value string) error {
		if value != "" {
			*out = value
			return nil
		}

		prompt := promptui.Prompt{
			Label:       internal.ValueOrDefaultString(promptPrefix),
			Mask:        '*',
			HideEntered: true,
		}

		if len(validateFunc) > 0 {
			prompt.Validate = validateFunc[0]
		}

		text, err := prompt.Run()
		if err != nil {
			return err
		}
		*out = text
		return nil
	}
}

// PromptText prompts the user to enter a text. This only prompts the user in the case that
// the corresponding environment variable does not contain any string data.
func PromptText(out, promptPrefix *string, validateFunc ...func(string) error) configo.ParserFunc {
	return func(value string) error {
		if value != "" {
			*out = value
			return nil
		}

		prompt := promptui.Prompt{
			Label: internal.ValueOrDefaultString(promptPrefix),
		}

		if len(validateFunc) > 0 {
			prompt.Validate = validateFunc[0]
		}

		text, err := prompt.Run()
		if err != nil {
			return err
		}
		if out != nil {
			*out = text
		}
		return nil
	}
}

// PromptInt prompts the user to enter a string. This only prompts the user in the case that
// the corresponding environment variable does not contain any string data.
func PromptInt(out *int, promptPrefix *string) configo.ParserFunc {
	internal.PanicIfNil(out)
	return func(value string) error {
		return PromptText(nil, promptPrefix, Int(out))(value)
	}
}

// PromptBool prompts the user to enter a boolean. This only prompts the user in the case that
// the corresponding environment variable does not contain any string data.
func PromptBool(out *bool, promptPrefix *string) configo.ParserFunc {
	internal.PanicIfNil(out)
	return func(value string) error {
		return PromptText(nil, promptPrefix, Bool(out))(value)
	}
}

// PromptFloat prompts the user to enter a float. This only prompts the user in the case that
// the corresponding environment variable does not contain any string data.
func PromptFloat(out *float64, bitSize int, promptPrefix *string) configo.ParserFunc {
	internal.PanicIfNil(out)
	return func(value string) error {
		return PromptText(nil, promptPrefix, Float(out, bitSize))(value)
	}
}

// PromptChoiceText prompts the user to select one of the provided strings. This only prompts the user in the case that
// the corresponding environment variable does not contain any string data.
func PromptChoiceText(out *string, promptPrefix *string, allowed ...string) configo.ParserFunc {
	return promptChoice(out, promptPrefix, allowed)
}

func promptChoice(out *string, promptPrefix *string, allowed []string, parseFunc ...func(string) error) configo.ParserFunc {
	internal.PanicIfNil(promptPrefix)
	internal.PanicIfEmptyString(allowed)

	allowedList := listToSortedUniqueListString(allowed)

	return func(value string) error {
		if value != "" {
			*out = value
			return nil
		}

		sel := promptui.Select{
			Label: internal.ValueOrDefaultString(promptPrefix),
			Items: allowedList,
		}

		_, text, err := sel.Run()
		if err != nil {
			return err
		}

		if out != nil {
			*out = text
		}
		if len(parseFunc) > 0 {
			return parseFunc[0](value)
		}
		return nil
	}
}

// PromptChoiceInt prompts the user to select one of the provided integers. This only prompts the user in the case that
// the corresponding environment variable does not contain any string data.
func PromptChoiceInt(out *int, promptPrefix *string, allowed []int) configo.ParserFunc {
	internal.PanicIfNil(out, promptPrefix)
	internal.PanicIfEmptyInt(allowed)

	allowedList := intListToStringList(listToSortedUniqueListInt(allowed))

	return func(value string) error {
		if value != "" {
			return Int(out)(value)
		}
		return promptChoice(nil, promptPrefix, allowedList, Int(out))(value)
	}
}

// PromptChoiceFloat prompts the user to select one of the provided floats. This only prompts the user in the case that
// the corresponding environment variable does not contain any string data.
func PromptChoiceFloat(out *float64, bitSize int, promptPrefix *string, allowed []float64) configo.ParserFunc {
	internal.PanicIfNil(out, promptPrefix)
	internal.PanicIfEmptyFloat(allowed)

	allowedList := floatListToStringList(listToSortedUniqueListFloat(allowed), bitSize)

	return func(value string) error {
		if value != "" {
			return Float(out, bitSize)(value)
		}

		return promptChoice(nil, promptPrefix, allowedList, Float(out, bitSize))(value)
	}
}

// PromptChoiceBool prompts the user to select one of the provided "y/n". This only prompts the user in the case that
// the corresponding environment variable does not contain any string data.
func PromptChoiceBool(out *bool, promptPrefix *string, defaultSelection ...bool) configo.ParserFunc {
	internal.PanicIfNil(out, promptPrefix)

	return func(value string) error {
		if value != "" {
			return Bool(out)(value)
		}

		t := "y"
		f := "n"
		pos := 1
		if len(defaultSelection) > 0 {
			capitalize := defaultSelection[0]
			if capitalize {
				t = strings.ToUpper(t)
				pos = 0
			} else {
				f = strings.ToUpper(f)
				pos = 1
			}
		}
		sel := promptui.Select{
			Label:     internal.ValueOrDefaultString(promptPrefix),
			Items:     []string{t, f},
			CursorPos: pos,
		}

		_, text, err := sel.Run()
		if err != nil {
			return err
		}
		return Bool(out)(text)
	}
}
