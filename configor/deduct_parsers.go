package configor

import (
	"reflect"

	"github.com/fatih/structtag"
	configo "github.com/jxsl13/simple-configo"
)

func deductParsers(tag *structtag.Tag, objValue reflect.Value) (configo.ParserFunc, error) {
	//options := tag.Options // extract requested parser from struct tag

	return ParserFuncWrapper(objValue), nil
}
