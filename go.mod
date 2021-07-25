module github.com/jxsl13/simple-configo

go 1.13

require (
	github.com/fatih/structtag v1.2.0
	github.com/iancoleman/strcase v0.2.0
	github.com/joho/godotenv v1.3.0
	golang.org/x/term v0.0.0-20210406210042-72f3dc4e9b72
)

replace (
	github.com/jxsl13/simple-configo => ./
	github.com/jxsl13/simple-configo/internal => ./internal/
	github.com/jxsl13/simple-configo/parsers => ./parsers/
	github.com/jxsl13/simple-configo/unparsers => ./unparsers/
)
