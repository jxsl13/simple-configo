package parsers

import (
	"errors"
	"testing"

	configo "github.com/jxsl13/simple-configo"
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
			ParseFunction:  Xor(errParser, errParser, nilParser, errParser),
		},
		{
			Key:            "#1",
			IsPseudoOption: true,
			ParseFunction:  Or(errParser, errParser, nilParser, errParser),
		},
		{
			Key:            "#2",
			IsPseudoOption: true,
			ParseFunction:  Or(errParser, nilParser),
		},
		{
			Key:            "#3",
			IsPseudoOption: true,
			ParseFunction:  And(nilParser, nilParser),
		},
		{
			Key:            "#4",
			IsPseudoOption: true,
			ParseFunction:  And(errParser, nilParser), // TODO: this should actually return an error
		},
	}[s.Start:s.End]
}

func TestLogic(t *testing.T) {

	tests := []struct {
		name    string
		cfg     configo.Config
		wantErr bool
	}{
		{"#1", &testCfg{0, 3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := configo.Parse(tt.cfg, map[string]string{})
			if tt.wantErr && err == nil || !tt.wantErr && err != nil {
				t.Errorf("Want Error: %v, got: %v\n", tt.wantErr, err)
			}
		})
	}
}
