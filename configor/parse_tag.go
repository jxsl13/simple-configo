package configor

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/fatih/structtag"
	configo "github.com/jxsl13/simple-configo"
)

var extractOptionRegexp = regexp.MustCompile(`^(\S+)='(.+)'|(\S+)=(.+)$`)
var errNoMatch = errors.New("no match")

// getOption extracts the key=value or key='value' of an option
func getOption(optionStr string) (key, value string, err error) {
	matches := extractOptionRegexp.FindStringSubmatch(optionStr)
	if len(matches) != 5 {
		return "", "", errNoMatch
	}

	if matches[1] != "" && matches[2] != "" {
		return matches[1], matches[2], nil
	} else if matches[3] != "" && matches[4] != "" {
		return matches[3], matches[4], nil
	}
	return "", "", errNoMatch
}

func parseTag(tag *structtag.Tag, objValue reflect.Value) (*configo.Option, error) {

	option := &configo.Option{
		Key:            tag.Name,
		Mandatory:      tag.HasOption(ConfigorStructTagOptionMandatory),
		IsPseudoOption: tag.HasOption(ConfigorStructTagOptionPseudo),
	}

	for _, to := range tag.Options {
		key, value, err := getOption(to) // extract key=value/key='value'
		if err != nil {
			fmt.Printf("unknown tag option: '%s'\n", to)
			continue
		}
		switch key {
		case ConfigorStructTagOptionDefault:
			// extracted default value
			option.DefaultValue = value
		case ConfigorStructTagOptionDescription:
			// extracted description
			option.Description = value
		}
	}

	parser, err := deductParsers(tag, objValue)
	if err != nil {
		return nil, err
	}
	option.ParseFunction = parser

	// set default description value
	if option.Description == "" {
		option.Description = option.Key
	}

	return option, nil
}
