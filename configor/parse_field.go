package configor

import (
	"reflect"

	"github.com/fatih/structtag"
	configo "github.com/jxsl13/simple-configo"
)

func getTagOrDefaultTag(tagKey, tag, defaultTag string) (*structtag.Tag, error) {
	// parse provided tags string
	tags, err := structtag.Parse(tag)
	if err != nil {
		return nil, err
	}

	// try to find our configo tag
	resultTag, err := tags.Get(tagKey)
	if err != nil {
		// failed to fetch cofigo tag
		// use default tag string
		tags, err = structtag.Parse(defaultTag)
		if err != nil {
			return nil, err
		}
		resultTag, err = tags.Get(tagKey)
	}
	// return configo tag
	return resultTag, err
}

func parseField(objValue reflect.Value, objField reflect.StructField) (*configo.Option, error) {
	tag, err := getTagOrDefaultTag(ConfigorStructTagName, string(objField.Tag), DefaultStructTag)
	if err != nil {
		return nil, err
	}
	switch tag.Name {
	case "":
		tag.Name = convertToKey(objField.Name)
	case "-":
		// skip parsing of field if field name is '-' = first element in struct tag
		return nil, nil
	}

	return parseTag(tag, objValue)
}
