package configo

import (
	"testing"
)

type MapConfig struct {
	UniqueList []string
	Mapping    map[string]string
}

func (m *MapConfig) Equal(other MapConfig) bool {
	if len(m.UniqueList) != len(other.UniqueList) {
		return false
	}

	for idx, value := range m.UniqueList {
		if value != other.UniqueList[idx] {
			return false
		}
	}

	if len(m.Mapping) != len(other.Mapping) {
		return false
	}

	for key, value := range m.Mapping {
		if other.Mapping[key] != value {
			return false
		}
	}

	return true
}

func (m *MapConfig) Name() string {
	return "Map Config"
}

func (m *MapConfig) Options() Options {
	delimiter := " "
	return Options{
		{
			Key:           "SOURCE_LIST",
			Mandatory:     true,
			Description:   "This is some description text.",
			ParseFunction: DefaultParserUniqueList(&m.UniqueList, &delimiter),
		},
		{
			Key:           "TARGET_LIST",
			Mandatory:     true,
			Description:   "This is some description text.",
			ParseFunction: DefaultParserMap(&m.Mapping, &m.UniqueList, &delimiter),
		},
	}
}

func TestMapParsing(t *testing.T) {
	type args struct {
		cfg    *MapConfig
		env    map[string]string
		result *MapConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"#1",
			args{&MapConfig{},
				map[string]string{
					"SOURCE_LIST": "1 2 3 4 5",
					"TARGET_LIST": "10 20 30 40 40",
				},
				&MapConfig{
					UniqueList: []string{"1", "2", "3", "4", "5"},
					Mapping: map[string]string{
						"1": "10",
						"2": "20",
						"3": "30",
						"4": "40",
						"5": "40",
					},
				},
			},
			false,
		},
		{
			"#2",
			args{&MapConfig{},
				map[string]string{
					"SOURCE_LIST": "1 2 3 4",
					"TARGET_LIST": "1 2 3 4 5",
				},
				&MapConfig{},
			},
			true,
		},
		{
			"#3",
			args{&MapConfig{},
				map[string]string{
					"SOURCE_LIST": "1 2 3 4 5",
					"TARGET_LIST": "1 2 3 4",
				},
				&MapConfig{},
			},
			true,
		},
		{
			"#4",
			args{&MapConfig{},
				map[string]string{
					"SOURCE_LIST": "1 2 3 4 4",
					"TARGET_LIST": "1 2 3 4 5",
				},
				&MapConfig{},
			},
			true,
		},
		{
			"#3",
			args{&MapConfig{},
				map[string]string{
					"SOURCE_LIST": "1 2 3 4 5",
					"TARGET_LIST": "1 2 3 4",
				},
				&MapConfig{},
			},
			true,
		},
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
