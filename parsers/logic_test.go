package parsers_test

import (
	"errors"
	"testing"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/parsers"
)

func errParser(value string) error {
	return errors.New("expected error")
}

func nilParser(value string) error {
	return nil
}

type testCfg struct {
	Start int
	End   int
}

func (s *testCfg) Name() string {
	return "test cfg"
}

func (s *testCfg) Options() configo.Options {
	return configo.Options{
		{
			Key:            "#0",
			IsPseudoOption: true,
			ParseFunction:  parsers.Xor(errParser, errParser, nilParser, errParser),
		},
		{
			Key:            "#1",
			IsPseudoOption: true,
			ParseFunction:  parsers.Or(errParser, errParser, nilParser, errParser),
		},
		{
			Key:            "#2",
			IsPseudoOption: true,
			ParseFunction:  parsers.Or(errParser, nilParser),
		},
		{
			Key:            "#3",
			IsPseudoOption: true,
			ParseFunction:  parsers.And(nilParser, nilParser),
		},
		{
			Key:            "#4",
			IsPseudoOption: true,
			ParseFunction:  parsers.And(errParser, nilParser),
		},
		{
			Key:            "#5",
			IsPseudoOption: true,
			ParseFunction:  parsers.Or(errParser, errParser),
		},
		{
			Key:            "#6",
			IsPseudoOption: true,
			ParseFunction:  parsers.Xor(errParser, errParser, nilParser, errParser, nilParser, errParser, errParser),
		},
	}[s.Start:s.End]
}

func TestLogic(t *testing.T) {
	// pseudo values
	var env map[string]string = make(map[string]string)
	// for i := 0; i < 1000; i++ {
	// 	env[fmt.Sprintf("#%d", i)] = "#value"
	// }

	tests := []struct {
		name    string
		cfg     configo.Config
		wantErr bool
	}{
		{"#1", &testCfg{0, 4}, false},
		{"#2", &testCfg{4, 5}, true},
		{"#3", &testCfg{5, 6}, true},
		{"#4", &testCfg{6, 7}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := configo.Parse(env, tt.cfg)
			if tt.wantErr && err == nil || !tt.wantErr && err != nil {
				t.Errorf("Want Error: %v, got: %v\n", tt.wantErr, err)
			}
		})
	}
}
