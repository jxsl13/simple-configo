package unparsers

import configo "github.com/jxsl13/simple-configo"

func panicIfEmptyUnparseFunc(list []configo.UnparserFunc) {
	if len(list) == 0 {
		panic("unparse function list must not be empty")
	}
}
