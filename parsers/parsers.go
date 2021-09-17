package parsers

import (
	"errors"

	configo "github.com/jxsl13/simple-configo"
)

var (
	// Error is returned when parsing of a key fails.
	Error = errors.New("parsing error")
)

func panicIfEmptyParseFunc(list []configo.ParserFunc) {
	if len(list) == 0 {
		panic("parse function list must not be empty")
	}
}
