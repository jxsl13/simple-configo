package configo_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/parsers"
)

type ErrorConfig struct {
	SomeField int
	sync.Mutex
}

func (ec *ErrorConfig) Name() string {
	return "ErrorConfig"
}

func (ec *ErrorConfig) Options() configo.Options {
	optionsList := configo.Options{
		{
			Key:           "SOME_FIELD",
			Description:   "This is some description text.",
			DefaultValue:  "SOME FIELD",
			ParseFunction: parsers.Int(&ec.SomeField),
		},
	}

	return optionsList
}

func TestParseError(t *testing.T) {

	err := configo.Parse(map[string]string{
		"SOME_FIELD": "12345",
	}, &ErrorConfig{})

	if err == nil {
		t.Errorf("Parse() EXPECTING ERROR, BUT GOT NONE!")
	}

}

type ErrorDefaultValuConfig struct {
	SomeField bool
	mu        sync.Mutex
}

func (e *ErrorDefaultValuConfig) Lock() {
	e.mu.Lock()
}

func (e *ErrorDefaultValuConfig) Unlock() {
	e.mu.Unlock()
}

func (ec *ErrorDefaultValuConfig) Name() string {
	return "ErrorConfig"
}

func (ec *ErrorDefaultValuConfig) Options() configo.Options {
	optionsList := configo.Options{
		{
			Key:           "SOME_FIELD",
			Description:   "This is some description text.",
			DefaultValue:  "2",
			ParseFunction: parsers.Bool(&ec.SomeField),
		},
	}

	return optionsList
}

func TestParseDefaultValueError(t *testing.T) {

	err := configo.Parse(map[string]string{}, &ErrorDefaultValuConfig{})

	if err == nil {
		t.Errorf("Parse() EXPECTING ERROR, BUT GOT NONE!")
	}

}

type MyConfig struct {
	SomeBool         bool
	SomeInt          int
	SomeFloat        float64
	SomeDelimiter    string
	SomeDuration     time.Duration
	SomeList         []string
	SomeStringSet    map[string]bool
	SomeChoiceInt    int
	SomeChoiceFloat  float64
	SomeChoiceString string
	SomeRangeInt     int
	sync.Mutex
}

func (m *MyConfig) String() string {
	b, err := json.MarshalIndent(m, " ", " ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (m *MyConfig) Equal(other *MyConfig) bool {
	eq := m.SomeBool == other.SomeBool &&
		m.SomeInt == other.SomeInt &&
		m.SomeFloat == other.SomeFloat &&
		m.SomeDelimiter == other.SomeDelimiter &&
		m.SomeDuration == other.SomeDuration &&
		m.SomeChoiceInt == other.SomeChoiceInt &&
		m.SomeChoiceFloat == other.SomeChoiceFloat &&
		m.SomeChoiceString == other.SomeChoiceString &&
		m.SomeRangeInt == other.SomeRangeInt
	if !eq {
		return false
	}
	eq = eq && (len(m.SomeList) == len(other.SomeList))
	eq = eq && (len(m.SomeStringSet) == len(other.SomeStringSet))
	if !eq {
		return false
	}
	for idx, v := range m.SomeList {
		if v != other.SomeList[idx] {
			return false
		}
	}

	for k, v := range m.SomeStringSet {
		if v != other.SomeStringSet[k] {
			return false
		}
	}

	return eq
}

// Name is the name of the configuration Cache
func (m *MyConfig) Name() (name string) {
	return "MY_CONFIG"
}

// Options returns a list of available options that can be configured for this
// config object
func (m *MyConfig) Options() (options configo.Options) {

	// NOTE: delimiter is parsed before the other values, this order is important,
	// as the delimiter is used afterwards.
	optionsList := configo.Options{
		{
			Key:           "SOME_BOOL",
			Mandatory:     true,
			Description:   "This is some description text.",
			DefaultValue:  "no",
			ParseFunction: parsers.Bool(&m.SomeBool),
		},
		{
			Key:           "SOME_INT",
			Description:   "This is some description text.",
			DefaultValue:  "42",
			ParseFunction: parsers.Int(&m.SomeInt),
		},
		{
			Key:           "SOME_FLOAT",
			Description:   "This is some description text.",
			DefaultValue:  "99.99",
			ParseFunction: parsers.Float(&m.SomeFloat, 64),
		},
		{
			Key:           "SOME_DELIMITER",
			Description:   "delimiter to split the lists below.",
			DefaultValue:  " ",
			ParseFunction: parsers.String(&m.SomeDelimiter),
		},
		{
			Key:           "SOME_DURATION",
			Description:   "This is some description text.",
			DefaultValue:  "24h12m44s",
			ParseFunction: parsers.Duration(&m.SomeDuration),
		},
		{
			Key:           "SOME_LIST",
			Description:   "Some IP list",
			DefaultValue:  "127.0.0.1 127.0.0.2 127.0.0.3",
			ParseFunction: parsers.List(&m.SomeList, &m.SomeDelimiter),
		},
		{
			Key:           "SOME_SET",
			Description:   "This is some description text.",
			DefaultValue:  "127.0.0.1 127.0.0.2 127.0.0.3 127.0.0.1",
			ParseFunction: parsers.ListToSet(&m.SomeStringSet, &m.SomeDelimiter),
		},
		{
			Key:           "SOME_CHOICE_INT",
			Description:   "This is some description text.",
			DefaultValue:  "4",
			ParseFunction: parsers.ChoiceInt(&m.SomeChoiceInt, 1, 2, 3, 4, 5, 6),
		},
		{
			Key:           "SOME_CHOICE_FLOAT",
			Description:   "This is some description text.",
			DefaultValue:  "5.5",
			ParseFunction: parsers.ChoiceFloat(&m.SomeChoiceFloat, 64, 1.1, 2.2, 3.3, 4.4, 5.5),
		},
		{
			Key:           "SOME_CHOICE_STRING",
			Description:   "This is some description text.",
			DefaultValue:  "empty",
			ParseFunction: parsers.ChoiceString(&m.SomeChoiceString, "empty", "full", "half empty"),
		},
		{
			Key:           "SOME_RANGE_INT",
			Description:   "This is some description text.",
			DefaultValue:  "42",
			ParseFunction: parsers.RangesInt(&m.SomeRangeInt, 0, 99),
		},
	}

	// add prefix
	for idx := range optionsList {
		optionsList[idx].Key = "MY_" + optionsList[idx].Key
	}

	return optionsList
}

func Test_Parse(t *testing.T) {
	type args struct {
		cfg    *MyConfig
		env    map[string]string
		result *MyConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"#1", args{&MyConfig{}, map[string]string{
			"MY_SOME_BOOL":          "true",
			"MY_SOME_INT":           "1234567",
			"MY_SOME_FLOAT":         "22.22",
			"MY_SOME_DELIMITER":     ",",
			"MY_SOME_DURATION":      "24h",
			"MY_SOME_LIST":          "1,2,3,4,5,6,7,8,9,0",
			"MY_SOME_SET":           "1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7",
			"MY_SOME_CHOICE_INT":    "5",
			"MY_SOME_CHOICE_FLOAT":  "1.1",
			"MY_SOME_CHOICE_STRING": "full",
			"MY_SOME_RANGE_INT":     "90",
		}, &MyConfig{
			SomeBool:      true,
			SomeInt:       1234567,
			SomeFloat:     22.22,
			SomeDelimiter: ",",
			SomeDuration:  time.Hour * 24,
			SomeList:      []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"},
			SomeStringSet: map[string]bool{
				"1": true,
				"2": true,
				"3": true,
				"4": true,
				"5": true,
				"6": true,
				"7": true,
				"8": true,
				"9": true,
				"0": true,
			},
			SomeChoiceInt:    5,
			SomeChoiceFloat:  1.1,
			SomeChoiceString: "full",
			SomeRangeInt:     90,
		}}, false},
		{"#2", args{&MyConfig{}, map[string]string{
			"MY_SOME_BOOL":      "false",
			"MY_SOME_INT":       "-123",
			"MY_SOME_FLOAT":     "-100",
			"MY_SOME_DELIMITER": ";",
			"MY_SOME_DURATION":  "1440m",
			"MY_SOME_LIST":      "1;2;3;",
			"MY_SOME_SET":       "1;2;3;",
		}, &MyConfig{
			SomeBool:      false,
			SomeInt:       -123,
			SomeFloat:     -100.0,
			SomeDelimiter: ";",
			SomeDuration:  time.Minute * 1440,
			SomeList:      []string{"1", "2", "3"},
			SomeStringSet: map[string]bool{
				"1": true,
				"2": true,
				"3": true,
			},
			SomeChoiceInt:    4,
			SomeChoiceFloat:  5.5,
			SomeChoiceString: "empty",
			SomeRangeInt:     42,
		}}, false},
		{"#3", args{&MyConfig{}, map[string]string{
			"MY_SOME_INT": "42",
		}, &MyConfig{
			SomeBool:      false,
			SomeInt:       42,
			SomeFloat:     99.99,
			SomeDelimiter: " ",
			SomeDuration:  24*time.Hour + 12*time.Minute + 44*time.Second,
			SomeList:      []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"},
			SomeStringSet: map[string]bool{
				"127.0.0.1": true,
				"127.0.0.2": true,
				"127.0.0.3": true,
			},
			SomeChoiceInt:    4,
			SomeChoiceFloat:  5.5,
			SomeChoiceString: "empty",
			SomeRangeInt:     42,
		}}, false},
		{"#4", args{&MyConfig{}, map[string]string{
			"MY_SOME_BOOL": "false",
		}, &MyConfig{
			SomeBool:      false,
			SomeInt:       42,
			SomeFloat:     99.99,
			SomeDelimiter: " ",
			SomeDuration:  24*time.Hour + 12*time.Minute + 44*time.Second,
			SomeList:      []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"},
			SomeStringSet: map[string]bool{
				"127.0.0.1": true,
				"127.0.0.2": true,
				"127.0.0.3": true,
			},
			SomeChoiceInt:    4,
			SomeChoiceFloat:  5.5,
			SomeChoiceString: "empty",
			SomeRangeInt:     42,
		}}, false},
		{"#5", args{&MyConfig{}, map[string]string{
			"MY_SOME_BOOL": "günni",
		}, &MyConfig{}}, true},
		{"#6", args{&MyConfig{}, map[string]string{
			"MY_SOME_BOOL": "false",
			"MY_SOME_INT":  "-123.99",
		}, &MyConfig{}}, true},
		{"#7", args{&MyConfig{}, map[string]string{
			"MY_SOME_BOOL":  "false",
			"MY_SOME_FLOAT": "123;99",
		}, &MyConfig{}}, true},
		{"#8", args{&MyConfig{}, map[string]string{
			"MY_SOME_BOOL":     "false",
			"MY_SOME_DURATION": "99hs",
		}, &MyConfig{}}, true},
		{"#9", args{&MyConfig{}, map[string]string{
			"MY_SOME_CHOICE_INT": "-5",
		}, &MyConfig{}}, true},
		{"#10", args{&MyConfig{}, map[string]string{
			"MY_SOME_CHOICE_FLOAT": "9.9",
		}, &MyConfig{}}, true},
		{"#11", args{&MyConfig{}, map[string]string{
			"MY_SOME_CHOICE_STRING": "not allowed",
		}, &MyConfig{}}, true},
		{"#12", args{&MyConfig{}, map[string]string{
			"MY_SOME_RANGE_INT": "200",
		}, &MyConfig{}}, true},

		{"#13", args{&MyConfig{}, map[string]string{
			"MY_SOME_CHOICE_INT": "-5",
		}, &MyConfig{SomeChoiceInt: 4}}, true},
		{"#14", args{&MyConfig{}, map[string]string{
			"MY_SOME_CHOICE_FLOAT": "9.9",
		}, &MyConfig{SomeChoiceFloat: 5.5}}, true},
		{"#15", args{&MyConfig{}, map[string]string{
			"MY_SOME_CHOICE_STRING": "not allowed",
		}, &MyConfig{SomeChoiceString: "empty"}}, true},
		{"#16", args{&MyConfig{}, map[string]string{
			"MY_SOME_RANGE_INT": "200",
		}, &MyConfig{SomeRangeInt: 42}}, true},
		{"#17", args{&MyConfig{}, map[string]string{
			"MY_SOME_RANGE_INT": "200",
		}, &MyConfig{SomeRangeInt: 42}}, true},
	}

	for idx, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := configo.Parse(tt.args.env, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && !tt.args.cfg.Equal(tt.args.result) {
				t.Fatalf("#%d : Parse() error = UNEXPECTED RESULT\nWANT:\n%s\nGOT:\n%s\n", idx+1, tt.args.result, tt.args.cfg)
			}
		})
	}
}

type MandatoryConfig struct {
	MandatoryRegex string
	sync.Mutex
}

func (m *MandatoryConfig) Equal(other *MandatoryConfig) bool {
	return m.MandatoryRegex == other.MandatoryRegex
}

func (m *MandatoryConfig) Name() string {
	return "Mandatory Config"
}

func (m *MandatoryConfig) Options() configo.Options {
	return configo.Options{
		{
			Key:           "MANDATORY_REGEX",
			Mandatory:     true,
			Description:   "This is some description text.",
			DefaultValue:  "mandatory",
			ParseFunction: parsers.Regex(&m.MandatoryRegex, "[a-z]+", "must only contain a-z"),
		},
	}
}

func Test_MandatoryParse(t *testing.T) {
	type args struct {
		cfg    *MandatoryConfig
		env    map[string]string
		result *MandatoryConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"#1", args{&MandatoryConfig{}, map[string]string{
			"MANDATORY_REGEX": "",
		}, &MandatoryConfig{
			MandatoryRegex: "mandatory",
		}}, true},
		{"#2", args{&MandatoryConfig{}, map[string]string{},
			&MandatoryConfig{
				MandatoryRegex: "mandatory",
			}}, false},
	}

	for idx, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := configo.Parse(tt.args.env, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && !tt.args.cfg.Equal(tt.args.result) {
				t.Fatalf("#%d : Parse() error = UNEXPECTED RESULT\nWANT:\n%s\nGOT:\n%s\n",
					idx+1,
					tt.args.result.MandatoryRegex,
					tt.args.cfg.MandatoryRegex,
				)
			}
		})
	}
}

type InvalidMandatoryConfig struct {
	MandatoryRegex string
}

func (m *InvalidMandatoryConfig) Equal(other InvalidMandatoryConfig) bool {
	return m.MandatoryRegex == other.MandatoryRegex
}

func (m *InvalidMandatoryConfig) Name() string {
	return "Mandatory Config"
}

func (m *InvalidMandatoryConfig) Options() configo.Options {
	return configo.Options{
		{
			Key:           "MANDATORY_REGEX",
			Mandatory:     true,
			Description:   "This is some description text.",
			DefaultValue:  "15", // Configuration definition is invalid at this point
			ParseFunction: parsers.Regex(&m.MandatoryRegex, "[a-z]+", "must only contain a-z"),
		},
	}
}

func Test_InvalidMandatoryParse(t *testing.T) {
	type args struct {
		cfg    *InvalidMandatoryConfig
		env    map[string]string
		result *InvalidMandatoryConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"#1", args{&InvalidMandatoryConfig{}, map[string]string{
			"MANDATORY_REGEX": "",
		}, &InvalidMandatoryConfig{}}, true},
		{"#2", args{&InvalidMandatoryConfig{}, map[string]string{},
			&InvalidMandatoryConfig{}}, true},
	}

	for idx, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := configo.Parse(tt.args.env, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && !tt.args.cfg.Equal(*tt.args.result) {
				t.Fatalf("#%d : Parse() error = UNEXPECTED RESULT\nWANT:\n%s\nGOT:\n%s\n", idx+1, tt.args.result, tt.args.cfg)
			}
		})
	}
}

type EmptyMandatoryConfig struct {
	MandatoryRegex string
}

func (m *EmptyMandatoryConfig) Equal(other EmptyMandatoryConfig) bool {
	return m.MandatoryRegex == other.MandatoryRegex
}

func (m *EmptyMandatoryConfig) Name() string {
	return "Mandatory Config"
}

func (m *EmptyMandatoryConfig) Options() configo.Options {
	return configo.Options{
		{
			Key:           "MANDATORY_REGEX",
			Mandatory:     true,
			Description:   "This is some description text.",
			ParseFunction: parsers.Regex(&m.MandatoryRegex, "[a-z]+", "must only contain a-z"),
		},
	}
}

func TestEmptyMandatoryParse(t *testing.T) {
	type args struct {
		cfg    *EmptyMandatoryConfig
		env    map[string]string
		result *EmptyMandatoryConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"#1", args{&EmptyMandatoryConfig{}, map[string]string{
			"MANDATORY_REGEX": "",
		}, &EmptyMandatoryConfig{}}, true},
		{"#2", args{&EmptyMandatoryConfig{}, map[string]string{},
			&EmptyMandatoryConfig{}}, true},
		{"#3", args{&EmptyMandatoryConfig{}, map[string]string{
			"MANDATORY_REGEX": "mandatory",
		}, &EmptyMandatoryConfig{
			MandatoryRegex: "mandatory",
		}}, false},
	}

	for idx, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := configo.Parse(tt.args.env, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && !tt.args.cfg.Equal(*tt.args.result) {
				t.Fatalf("#%d : Parse() error = UNEXPECTED RESULT\nWANT:\n%s\nGOT:\n%s\n", idx+1, tt.args.result, tt.args.cfg)
			}
		})
	}
}

type InvalidDefaultValueConfig struct {
	MandatoryRegex string
}

func (m *InvalidDefaultValueConfig) Equal(other InvalidDefaultValueConfig) bool {
	return m.MandatoryRegex == other.MandatoryRegex
}

func (m *InvalidDefaultValueConfig) Name() string {
	return "Mandatory Config"
}

func (m *InvalidDefaultValueConfig) Options() configo.Options {
	return configo.Options{
		{
			Key:           "MANDATORY_REGEX",
			Mandatory:     true,
			Description:   "This is some description text.",
			DefaultValue:  "15",
			ParseFunction: parsers.Regex(&m.MandatoryRegex, "[a-z]+", "must only contain a-z"),
		},
	}
}

func TestInvalidDefaultValueParse(t *testing.T) {
	type args struct {
		cfg    *InvalidDefaultValueConfig
		env    map[string]string
		result *InvalidDefaultValueConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"#1", args{&InvalidDefaultValueConfig{}, map[string]string{
			"MANDATORY_REGEX": "",
		}, &InvalidDefaultValueConfig{}}, true},
		{"#2", args{&InvalidDefaultValueConfig{}, map[string]string{},
			&InvalidDefaultValueConfig{}}, true},
		{"#3", args{&InvalidDefaultValueConfig{}, map[string]string{
			"MANDATORY_REGEX": "mandatory",
		}, &InvalidDefaultValueConfig{
			MandatoryRegex: "mandatory",
		}}, true},
	}

	for idx, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := configo.Parse(tt.args.env, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && !tt.args.cfg.Equal(*tt.args.result) {
				t.Fatalf("#%d : Parse() error = UNEXPECTED RESULT\nWANT:\n%s\nGOT:\n%s\n", idx+1, tt.args.result, tt.args.cfg)
			}
		})
	}
}

type ActionConfig struct {
	peter string
	start int
	end   int
}

func (m *ActionConfig) Options() configo.Options {
	return configo.Options{
		{
			Key: "some key that will be irgnored",
			PreParseAction: func() error {
				m.peter = "pre parse"
				return nil
			},
			PostParseAction: func() error {
				if m.peter != "pre parse" {
					return errors.New("m.peter is not pre parse")
				}
				m.peter = "post parse"
				return nil
			},
		},
		{
			Key: "second key but it has a parsefunc",
			PreParseAction: func() error {
				m.peter = "pre parse"
				return nil
			},
			ParseFunction: func(value string) error {
				m.peter = "parse function"

				return nil
			},
			PostParseAction: func() error {
				if m.peter != "pre parse" {
					return fmt.Errorf("m.peter is not pre parse: %s", m.peter)
				}
				m.peter = "post parse"
				return nil
			},
		},
	}[m.start:m.end]
}

func TestActionOption(t *testing.T) {
	type args struct {
		cfg         *ActionConfig
		env         map[string]string
		resultPeter string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"#1", args{&ActionConfig{start: 0, end: 1}, map[string]string{}, "post parse"}, false},
		{"#2", args{&ActionConfig{start: 1, end: 2}, map[string]string{}, "post parse"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := configo.Parse(tt.args.env, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.args.cfg.peter != tt.args.resultPeter {
				t.Errorf("want: %s != got: %s", tt.args.cfg.peter, tt.args.resultPeter)
			}

		})
	}
}

type flagCfg struct {
	test  string
	peter int
}

func (fc *flagCfg) Options() configo.Options {
	return configo.Options{
		{
			Key:           "TEST",
			DefaultValue:  "default test",
			ParseFunction: parsers.String(&fc.test),
		},
		{
			Key:           "PETER",
			DefaultValue:  "99",
			ParseFunction: parsers.Int(&fc.peter),
		},
	}
}

func TestParseFlags(t *testing.T) {
	os.Args = append(os.Args,
		"-test=value",
	)
	os.Setenv("PETER", "1")
	os.Setenv("INVALID_PATH", "./invalid.env")

	fc := &flagCfg{}
	err := configo.ParseEnvFileOrEnvOrFlags("INVALID_PATH", fc)
	if err != nil {
		t.Fatal(err)
	}

	if fc.peter != 1 {
		t.Fatalf("fc.peter: want %d, got  %d", 1, fc.peter)
	}

	if fc.test != "value" {
		t.Fatalf("fc.test: want %s, got  %s", "value", fc.test)
	}
}

type flag2Cfg struct {
	test  string
	peter int
}

func (fc *flag2Cfg) Options() configo.Options {
	return configo.Options{
		{
			Key:           "TEST2",
			DefaultValue:  "default test",
			ParseFunction: parsers.String(&fc.test),
		},
		{
			Key:           "PETER2",
			DefaultValue:  "99",
			ParseFunction: parsers.Int(&fc.peter),
		},
	}
}

func TestParseFlags2(t *testing.T) {
	os.Args = append(os.Args,
		"-test2=value2",
	)
	os.Setenv("PETER2", "2")

	fc2 := &flag2Cfg{}
	configo.ParseEnvOrFlags(fc2)
	if fc2.peter != 2 {
		t.Fatalf("fc2.peter2: want %d, got  %d", 1, fc2.peter)
	}

	if fc2.test != "value2" {
		t.Fatalf("fc2.test: want %s, got  %s", "value2", fc2.test)
	}
}
