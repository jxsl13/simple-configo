module github.com/jxsl13/simple-configo

go 1.13

require (
	golang.org/x/term v0.0.0-20210406210042-72f3dc4e9b72
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace (
	github.com/jxsl13/simple-configo => ./
	github.com/jxsl13/simple-configo/parsers => ./parsers/
)
