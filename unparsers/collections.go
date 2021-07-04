package unparsers

import (
	"sort"
	"strings"

	configo "github.com/jxsl13/simple-configo"
)

// List returns a function that returns the string representation of the 'in' list
func List(in *[]string, delimiter *string) configo.UnparserFunc {
	return func() (string, error) {
		return strings.Join(*in, *delimiter), nil

	}
}

// SetToList returns a function that returns the set as a 'delimiter' separated string
// the returned list is sorted, as it does not matter in a set context, whether the list s sorted or not,
// we do want the list to be sorted for readability reasons.
func SetToList(in *map[string]bool, delimiter *string) configo.UnparserFunc {
	return func() (string, error) {
		list := make([]string, 0, len(*in))
		for k := range *in {
			list = append(list, k)
		}
		sort.Strings(list)
		return strings.Join(list, *delimiter), nil
	}
}

// Map returns a function that returns a sorted list of key value pairs delimited by the pairDelimiter.
// Each key value is delimited by the keyValueDelimiter
func Map(in *map[string]string, pairDelimiter, keyValueDelimiter *string) configo.UnparserFunc {
	return func() (string, error) {
		pairs := make([]string, 0, len(*in))
		for key, value := range *in {
			pairs = append(pairs, key+(*keyValueDelimiter)+value)
		}
		sort.Strings(pairs)
		return strings.Join(pairs, *pairDelimiter), nil
	}
}
