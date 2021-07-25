package configor

import (
	"errors"
	"reflect"

	configo "github.com/jxsl13/simple-configo"
)

func Parse(env map[string]string, i interface{}) error {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Struct {
		return errors.New("i must be a struct")
	}
	options, err := iterateOverFields(v, reflect.StructField{})
	if err != nil {
		return err
	}

	return configo.ParseOptions(options, env)

}

func iterateOverFields(obj reflect.Value, objField reflect.StructField) ([]configo.Option, error) {
	options := make([]configo.Option, 0, 1)
	if reflect.Invalid < obj.Kind() && obj.Kind() < reflect.Struct {
		// primitive types, parse tags for every primitive type
		option, err := parseField(obj, objField)
		if err != nil {
			return options, err
		}
		if option == nil {
			return nil, nil
		}
		return append(options, *option), nil
	}

	// struct -> parse each individual struct field
	for i := 0; i < obj.NumField(); i++ {
		fieldValue := obj.Field(i)
		fieldType := obj.Type().Field(i)

		// parse struct fields
		subOptions, err := iterateOverFields(fieldValue, fieldType)
		if err != nil {
			return nil, err
		}
		if len(subOptions) == 0 {
			// skip empty options
			continue
		}
		options = append(options, subOptions...)
	}
	return options, nil
}
