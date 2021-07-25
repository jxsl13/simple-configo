package configor

import "github.com/iancoleman/strcase"

func convertToKey(fieldName string) string {
	switch DefaultKeyCase {
	case KeyCaseEnv:
		return strcase.ToScreamingSnake(fieldName)
	case KeyCaseSnake:
		return strcase.ToSnake(fieldName)

	}
	return fieldName
}
