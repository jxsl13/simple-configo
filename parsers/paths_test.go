package parsers_test

import (
	"testing"

	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/parsers"
	"github.com/jxsl13/simple-configo/unparsers"
)

type pathConfig struct {
	directory    string
	file         string
	absDirectory string
	absFile      string
}

func (pc *pathConfig) Options() configo.Options {
	return configo.Options{
		{
			Key:       "DIRECTORY",
			Mandatory: true,
			ParseFunction: parsers.And(
				parsers.PathDirectory(&pc.directory),
				parsers.PathAbsDirectory(&pc.absDirectory),
			),
			UnparseFunction: unparsers.Path(&pc.directory),
		},
		{
			Key:       "FILE",
			Mandatory: true,
			ParseFunction: parsers.And(
				parsers.PathFile(&pc.file),
				parsers.PathAbsFile(&pc.absFile),
			),
			UnparseFunction: unparsers.Path(&pc.file),
		},
		{
			Key:       "FAIL_DIRECTORY",
			Mandatory: true,
			ParseFunction: parsers.And(
				parsers.PathDirectory(&pc.directory),
				parsers.PathAbsDirectory(&pc.absDirectory),
			),
		},
		{
			Key:       "FAIL_FILE",
			Mandatory: true,
			ParseFunction: parsers.And(
				parsers.PathFile(&pc.file),
				parsers.PathAbsFile(&pc.absFile),
			),
		},
	}
}

func TestPathDirectory(t *testing.T) {
	env := map[string]string{
		"DIRECTORY":      "./../",
		"FILE":           "./paths.go",
		"FAIL_DIRECTORY": "./abcdefghijklmnop.xyz",
		"FAIL_FILE":      "./abcdefghijklmnop.xyz",
	}

	pc := &pathConfig{}
	options := pc.Options()

	dirOpt := options[0]
	fileOpt := options[1]
	failDirOpt := options[2]
	failFileOpt := options[3]

	err := dirOpt.Parse(env)
	if err != nil {
		t.Fatal(err)
	}
	value, err := dirOpt.Unparse()
	if err != nil {
		t.Fatalf("failed to unparse dirOpt: %s", value)
	}

	err = fileOpt.Parse(env)
	if err != nil {
		t.Fatal(err)
	}
	value, err = fileOpt.Unparse()
	if err != nil {
		t.Fatalf("failed to unparse fileOpt: %s", value)
	}

	err = failDirOpt.Parse(env)
	if err == nil {
		t.Fatal("expecting failDirOpt to fail, as it does not exist")
	}

	err = failFileOpt.Parse(env)
	if err == nil {
		t.Fatal("expecting failFileOpt to fail, as it does not exist")
	}

}
