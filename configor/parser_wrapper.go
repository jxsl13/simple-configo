package configor

import (
	"fmt"
	"reflect"

	"github.com/jxsl13/simple-configo/parsers"
)

func refSet(dest reflect.Value, src interface{}) error {
	source := reflect.ValueOf(src)
	if dest.Kind() != source.Kind() {
		return fmt.Errorf("trying to assign src '%s' to dest '%s'", source.Type().Name(), dest.Type().Name())
	}
	if !dest.CanSet() {
		if dest.Addr().CanSet() {
			dest.Addr().Set(source.Addr())
			return nil
		}

		return fmt.Errorf("cannot set dest of type: %s", dest.Type().Name())
	}

	dest.Set(source)
	return nil
}

func ParserFuncWrapper(out reflect.Value) func(value string) error {

	return func(value string) error {
		var err error
		switch out.Kind() {
		case reflect.Bool:
			var result bool
			err = parsers.Bool(&result)(value)
			if err != nil {
				break
			}
			err = refSet(out, result)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
			var result int
			err = parsers.Int(&result)(value)
			if err != nil {
				break
			}
			err = refSet(out, result)
		case reflect.Float32:
			var result float64
			err = parsers.Float(&result, 32)(value)
			if err != nil {
				break
			}
			err = refSet(out, result)
		case reflect.Float64:
			var result float64
			err = parsers.Float(&result, 64)(value)
			if err != nil {
				break
			}
			err = refSet(out, result)
		case reflect.Map:
			var result map[string]string
			err = parsers.Map(&result, &DefaultPairDelimiter, &DefaultKeyValueDelimiter)(value)
			if err != nil {
				break
			}
			err = refSet(out, result)
		case reflect.Slice:
			var result []string
			err = parsers.List(&result, &DefaultListDelimiter)(value)
			if err != nil {
				break
			}
			err = refSet(out, result)
		case reflect.String:
			var result string
			err = parsers.String(&result)(value)
			if err != nil {
				break
			}
			err = refSet(out, result)
		default:
			return fmt.Errorf("unsupported type: %s", out.Type().Name())
		}
		return err
	}
}
