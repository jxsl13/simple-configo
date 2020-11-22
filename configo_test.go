package configo

import (
	"encoding/json"
	"testing"
	"time"
)

type ErrorConfig struct {
	SomeField string
}

func (ec *ErrorConfig) Name() string {
	return "ErrorConfig"
}

func (ec *ErrorConfig) Options() Options {
	optionsList := Options{
		{
			Key:           "SOME_FIELD",
			Type:          "",
			Description:   "This is some description text.",
			DefaultValue:  "SOME FIELD",
			ParseFunction: DefaultParserString(&ec.SomeField),
		},
	}

	return optionsList
}

func TestParseError(t *testing.T) {

	err := Parse(&ErrorConfig{}, map[string]string{
		"SOME_FIELD": "12345",
	})

	if err == nil {
		t.Errorf("Parse() EXPECTING ERROR, BUT GOT NONE!")
	}

}

type ErrorDefaultValuConfig struct {
	SomeField bool
}

func (ec *ErrorDefaultValuConfig) Name() string {
	return "ErrorConfig"
}

func (ec *ErrorDefaultValuConfig) Options() Options {
	optionsList := Options{
		{
			Key:           "SOME_FIELD",
			Type:          "bool",
			Description:   "This is some description text.",
			DefaultValue:  "2",
			ParseFunction: DefaultParserBool(&ec.SomeField),
		},
	}

	return optionsList
}

func TestParseDefaultValueError(t *testing.T) {

	err := Parse(&ErrorDefaultValuConfig{}, map[string]string{})

	if err == nil {
		t.Errorf("Parse() EXPECTING ERROR, BUT GOT NONE!")
	}

}

type MyConfig struct {
	SomeBool      bool
	SomeInt       int
	SomeFloat     float64
	SomeDelimiter string
	SomeDuration  time.Duration
	SomeList      []string
	SomeStringSet map[string]bool
}

func (m *MyConfig) String() string {
	b, err := json.MarshalIndent(m, " ", " ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (m *MyConfig) Equal(other MyConfig) bool {
	eq := m.SomeBool == other.SomeBool &&
		m.SomeInt == other.SomeInt &&
		m.SomeFloat == other.SomeFloat &&
		m.SomeDelimiter == other.SomeDelimiter &&
		m.SomeDuration == other.SomeDuration
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
func (m *MyConfig) Options() (options Options) {

	// NOTE: delimiter is parsed before the other values, this order is important,
	// as the delimiter is used afterwards.
	optionsList := Options{
		{
			Key:           "SOME_BOOL",
			Type:          "bool",
			Mandatory:     true,
			Description:   "This is some description text.",
			DefaultValue:  "no",
			ParseFunction: DefaultParserBool(&m.SomeBool),
		},
		{
			Key:           "SOME_INT",
			Type:          "int",
			Description:   "This is some description text.",
			DefaultValue:  "42",
			ParseFunction: DefaultParserInt(&m.SomeInt),
		},
		{
			Key:           "SOME_FLOAT",
			Type:          "float",
			Description:   "This is some description text.",
			DefaultValue:  "99.99",
			ParseFunction: DefaultParserFloat(&m.SomeFloat, 64),
		},
		{
			Key:           "SOME_DELIMITER",
			Type:          "string",
			Description:   "delimiter to split the lists below.",
			DefaultValue:  " ",
			ParseFunction: DefaultParserString(&m.SomeDelimiter),
		},
		{
			Key:           "SOME_DURATION",
			Type:          "duration",
			Description:   "This is some description text.",
			DefaultValue:  "24h12m44s",
			ParseFunction: DefaultParserDuration(&m.SomeDuration),
		},
		{
			Key:           "SOME_LIST",
			Type:          "list",
			Description:   "Some IP list",
			DefaultValue:  "127.0.0.1 127.0.0.2 127.0.0.3",
			ParseFunction: DefaultParserList(&m.SomeDelimiter, &m.SomeList),
		},
		{
			Key:           "SOME_SET",
			Type:          "list",
			Description:   "This is some description text.",
			DefaultValue:  "127.0.0.1 127.0.0.2 127.0.0.3 127.0.0.1",
			ParseFunction: DefaultParserListToSet(&m.SomeDelimiter, &m.SomeStringSet),
		},
	}

	// add prefix
	for idx := range optionsList {
		optionsList[idx].Key = "MY_" + optionsList[idx].Key
	}

	optionsList.MustValid()
	return optionsList
}

func TestParse(t *testing.T) {
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
			"MY_SOME_BOOL":      "true",
			"MY_SOME_INT":       "1234567",
			"MY_SOME_FLOAT":     "22.22",
			"MY_SOME_DELIMITER": ",",
			"MY_SOME_DURATION":  "24h",
			"MY_SOME_LIST":      "1,2,3,4,5,6,7,8,9,0",
			"MY_SOME_SET":       "1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6,7",
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
		}}, false},
		{"#3", args{&MyConfig{}, map[string]string{
			"MY_SOME_INT": "-123",
		}, &MyConfig{}}, true},
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
	}

	for idx, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := Parse(tt.args.cfg, tt.args.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && !tt.args.cfg.Equal(*tt.args.result) {
				t.Fatalf("#%d : Parse() error = UNEXPECTED RESULT\nWANT:\n%s\nGOT:\n%s\n", idx+1, tt.args.result, tt.args.cfg)
			}
		})
	}
}
