package configor_test

import (
	"testing"

	"github.com/jxsl13/simple-configo/configor"
)

type TestStruct struct {
	Empty       string
	Key         string  `configo:",mandatory,default='peter'"`
	Name        string  `configo:"Name,mandatory,default='peter'"`
	SomeInt     int     `configo:"SoME_INT,mandatory,default='19'"`
	SomeBool    bool    `configo:"SOMe_BOOL,default='false'"`
	SomeFloat64 float64 `configo:"SoMe_FLOAT64,default='64',parser='regex:^[A-Z]+$:invalid key must be capital case:'"`
}

func TestParse(t *testing.T) {
	env := map[string]string{
		"EMPTY":        "empty text",
		"KEY":          "key string",
		"Name":         "some name",
		"SoME_INT":     "77",
		"SOMe_BOOL":    "true",
		"SoMe_FLOAT64": "69",
	}
	config := TestStruct{}
	err := configor.Parse(env, config)
	if err != nil {
		t.Fatal(err)
	}
	t.Error("expected")
}
