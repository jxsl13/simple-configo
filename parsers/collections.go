package parsers

import (
	"fmt"
	"strings"

	configo "github.com/jxsl13/simple-configo"
)

// List parses a string containing a 'delimiter'(space, comma, semicolon, etc.) delimited list
// into the string list 'out'
func List(out *[]string, delimiter *string) configo.ParserFunc {
	return func(value string) error {
		list := strings.Split(value, *delimiter)

		if len(list) > 0 && list[len(list)-1] == "" {
			list = list[:len(list)-1]
		}
		*out = list
		return nil
	}
}

// ListToSet parses a string containing a 'delimiter'(space, comma, semicolon, etc.) delimited list
// into set 'out'
func ListToSet(out *map[string]bool, delimiter *string) configo.ParserFunc {
	return func(value string) error {
		list := strings.Split(value, *delimiter)

		if len(list) > 0 && list[len(list)-1] == "" {
			list = list[:len(list)-1]
		}

		*out = make(map[string]bool, len(list))
		for _, s := range list {
			(*out)[s] = true
		}
		return nil
	}
}

// UniqueList enforces that the passed list contains only unique values
func UniqueList(out *[]string, delimiter *string) configo.ParserFunc {
	return func(value string) error {
		list := strings.Split(value, *delimiter)

		testMap := make(map[string]bool, len(list))
		for _, key := range list {
			testMap[key] = true
		}

		if len(list) != len(testMap) {
			return fmt.Errorf("the list must contain only unique values, but has %d redundant values", len(list)-len(testMap))
		}

		*out = list
		return nil
	}
}

// Map allows to define key->value associations directly inside of a single parameter
func Map(out *map[string]string, pairDelimiter, keyValueDelimiter *string) configo.ParserFunc {
	return func(value string) error {
		if *pairDelimiter == *keyValueDelimiter {
			return fmt.Errorf("pairDelimiter and keyValueDelimiter must not be equal: '%s'", *pairDelimiter)
		}
		pairs := strings.Split(value, *pairDelimiter)
		m := make(map[string]string, len(pairs))
		for _, pair := range pairs {
			list := strings.Split(pair, *keyValueDelimiter)
			if len(list) != 2 {
				return fmt.Errorf("'%s' is not a key value pair with delimiter '%s'", pair, *keyValueDelimiter)
			}
			key := list[0]
			value := list[1]

			if _, ok := m[key]; ok {
				return fmt.Errorf("duplicate key detected: %s", key)
			}

			m[key] = value
		}

		if *out == nil {
			*out = make(map[string]string, len(pairs))
		}

		for key, value := range m {
			(*out)[key] = value
		}
		return nil
	}
}

// MapReverse allows to define a value->key associations directly inside of a single parameter
// This might be useful if you want to construct a two way mapping that is able to associate a key to a value as well as
// the association form a value back to its key.
func MapReverse(out *map[string]string, pairDelimiter, keyValueDelimiter *string) configo.ParserFunc {
	return func(value string) error {
		if *pairDelimiter == *keyValueDelimiter {
			return fmt.Errorf("pairDelimiter and keyValueDelimiter must not be equal: '%s'", *pairDelimiter)
		}
		pairs := strings.Split(value, *pairDelimiter)
		m := make(map[string]string, len(pairs))
		for _, pair := range pairs {
			list := strings.Split(pair, *keyValueDelimiter)
			if len(list) != 2 {
				return fmt.Errorf("'%s' is not a key value pair with delimiter '%s'", pair, *keyValueDelimiter)
			}

			// the only difference between Map & MapReverse
			key := list[1]
			value := list[0]

			if _, ok := m[key]; ok {
				return fmt.Errorf("duplicate key detected: %s", key)
			}

			m[key] = value
		}

		if *out == nil {
			*out = make(map[string]string, len(pairs))
		}

		for key, value := range m {
			(*out)[key] = value
		}
		return nil
	}
}

// MapFromKeysSlice fills the 'out' map with the keys and for each key the corresponding
// value at the exact same position of the passed "value" string that is split into another
// list is creates. keys{0, 1, 2, 3} -> values{0, 1, 2, 3}
func MapFromKeysSlice(out *map[string]string, keys *[]string, delimiter *string) configo.ParserFunc {
	return func(value string) error {
		values := strings.Split(value, *delimiter)

		if len(values) != len(*keys) {
			return fmt.Errorf("passed key slice(len=%d) has a different length than the parsed values list(len=%d)", len(*keys), len(values))
		}

		keySet := make(map[string]bool, len(*keys))
		for _, key := range *keys {
			keySet[key] = true
		}

		if len(keySet) != len(*keys) {
			return fmt.Errorf("the passed key list must contain only unique keys, %d of %d keys are unique", len(keySet), len(*keys))
		}

		if *out == nil {
			*out = make(map[string]string, len(*keys))
		}

		for idx, source := range *keys {
			(*out)[source] = values[idx]
		}
		return nil
	}
}
