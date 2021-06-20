package parsers_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/parsers"
	"gopkg.in/yaml.v3"
)

func TestReadYAML(t *testing.T) {
	fileName := "test.yaml"
	defer os.Remove(fileName)
	cfg := yamlCfg{
		FilePath: fileName,
		Initial: yamlStruct{
			time.Now().Truncate(time.Millisecond),
			10,
			"test_string",
		},
		Result: yamlStruct{},
	}

	data, err := yaml.Marshal(cfg.Initial)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile(fileName, data, 0770)
	if err != nil {
		t.Error(err)
	}

	env := map[string]string{}
	err = configo.Parse(&cfg, env)
	if err != nil {
		t.Error(err)
	}
	i := cfg.Initial
	j := cfg.Result

	if !reflect.DeepEqual(i, j) {
		t.Errorf("wanted: %v got: %v", cfg.Initial, cfg.Result)
	}

}

type yamlStruct struct {
	Date   time.Time
	Num    int
	String string
}

type yamlCfg struct {
	FilePath string
	Initial  yamlStruct
	Result   yamlStruct
}

func (c *yamlCfg) Name() string {
	return ""
}

func (c *yamlCfg) Options() configo.Options {
	return configo.Options{
		{
			Key:           "FILE_NAME",
			Description:   "file path to the struct",
			DefaultValue:  c.FilePath,
			ParseFunction: parsers.ReadYAML(&c.Result),
		},
	}
}

func TestReadJSON(t *testing.T) {
	fileName := "test.json"
	defer os.Remove(fileName)
	cfg := jsonCfg{
		FilePath: fileName,
		Initial: jsonStruct{
			time.Now().Round(time.Millisecond),
			10,
			"test_string",
		},
		Result: jsonStruct{},
	}

	data, err := json.Marshal(cfg.Initial)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile(fileName, data, 0770)
	if err != nil {
		t.Error(err)
	}

	env := map[string]string{}
	err = configo.Parse(&cfg, env)
	if err != nil {
		t.Error(err)
	}
	i := &cfg.Initial
	j := &cfg.Result

	if !reflect.DeepEqual(i, j) {
		t.Errorf("wanted: %v got: %v", cfg.Initial, cfg.Result)
	}

}

type jsonStruct struct {
	Date   time.Time
	Num    int
	String string
}

type jsonCfg struct {
	FilePath string
	Initial  jsonStruct
	Result   jsonStruct
}

func (c *jsonCfg) Name() string {
	return ""
}

func (c *jsonCfg) Options() configo.Options {
	return configo.Options{
		{
			Key:           "FILE_NAME",
			Description:   "file path to the struct",
			DefaultValue:  c.FilePath,
			ParseFunction: parsers.ReadJSON(&c.Result),
		},
	}
}
