package actions

import configo "github.com/jxsl13/simple-configo"

func panicIfEmptyActionFunc(list []configo.ActionFunc) {
	if len(list) == 0 {
		panic("action function list must not be empty")
	}
}
