package unparsers_test

import (
	"reflect"
	"testing"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/internal"
	"github.com/jxsl13/simple-configo/parsers"
	"github.com/jxsl13/simple-configo/unparsers"
)

var (
	fileName = "./test.txt"
	env      = map[string]string{}
)

type TestStruct struct {
	Integer   int
	IntPtr    *int
	Float     float64
	FloatPtr  *float64
	String    string
	StringPtr *string
	SubStruct *TestStruct
}

func newTestStruct() TestStruct {
	sub := &TestStruct{
		Integer: 63,
		Float:   63.63,
		String:  "string-1",
	}
	sub.FloatPtr = &sub.Float
	sub.IntPtr = &sub.Integer
	sub.StringPtr = &sub.String

	main := TestStruct{
		Integer:   64,
		Float:     64.64,
		String:    "string",
		SubStruct: sub,
	}
	main.IntPtr = &main.Integer
	main.FloatPtr = &main.Float
	main.StringPtr = &main.String
	return main
}

type TestConfig struct {
	Test  TestStruct
	Index int
}

func (c *TestConfig) Name() string {
	return "TestConfig"
}

func (c *TestConfig) Options() configo.Options {
	return configo.Options{
		{
			Key:             "JSON",
			Description:     "test key that parses the file located at the Location specified in DefaultValue",
			DefaultValue:    fileName,
			ParseFunction:   parsers.ReadJSON(&c.Test),
			UnparseFunction: unparsers.WriteJSON(&c.Test),
		},
		{
			Key:             "YAML",
			Description:     "test key that parses the file located at the Location specified in DefaultValue",
			DefaultValue:    fileName,
			ParseFunction:   parsers.ReadYAML(&c.Test),
			UnparseFunction: unparsers.WriteYAML(&c.Test),
		},
		{
			Key:             "TEXT",
			Description:     "test key that parses the file located at the Location specified in DefaultValue",
			DefaultValue:    fileName,
			ParseFunction:   parsers.ReadFile(&c.Test.String),
			UnparseFunction: unparsers.WriteFile(&c.Test.String),
		},
	}[c.Index : c.Index+1]
}

func TestWriteYAML(t *testing.T) {

	test := newTestStruct()
	err := unparsers.WriteYAML(test)("YAML", fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer internal.Delete(fileName)

	cfg := &TestConfig{
		Index: 1,
	}

	cfg2 := &TestConfig{
		Index: 1,
	}
	err = configo.Parse(cfg, env)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg.Test, test) {
		t.Fatalf("got: %v want: %v\n", cfg.Test, test)
	}

	unparse := configo.Unparse(cfg, env)
	unparse()

	err = configo.Parse(cfg2, env)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(cfg, cfg2) {
		t.Fatalf("got: %v want: %v\n", cfg, cfg2)
	}

}

func TestWriteJSON(t *testing.T) {

	test := newTestStruct()
	err := unparsers.WriteJSON(test)("JSON", fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer internal.Delete(fileName)

	cfg := &TestConfig{
		Index: 0,
	}

	cfg2 := &TestConfig{
		Index: 0,
	}
	err = configo.Parse(cfg, env)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg.Test, test) {
		t.Fatalf("got: %v want: %v\n", cfg.Test, test)
	}

	unparse := configo.Unparse(cfg, env)
	unparse()

	err = configo.Parse(cfg2, env)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(cfg, cfg2) {
		t.Fatalf("got: %v want: %v\n", cfg, cfg2)
	}
}
