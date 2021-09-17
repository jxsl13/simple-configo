# Simple-Configo

[![Test](https://github.com/jxsl13/simple-configo/actions/workflows/build.yaml/badge.svg)](https://github.com/jxsl13/simple-configo/actions/workflows/build.yaml) [![Go Report Card](https://goreportcard.com/badge/github.com/jxsl13/simple-configo)](https://goreportcard.com/report/github.com/jxsl13/simple-configo) [![codecov](https://codecov.io/gh/jxsl13/simple-configo/branch/master/graph/badge.svg?token=noNR6ork0u)](https://codecov.io/gh/jxsl13/simple-configo) [![Total alerts](https://img.shields.io/lgtm/alerts/g/jxsl13/simple-configo.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/jxsl13/simple-configo/alerts/) [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Simple-Configo streamlines the creation of multiple independent configuration structs to the implementation of a single interface, the `Config` interface, that has only two methods.

```go
type Config interface {
    Options() (options Options)
}
```

The `Options` method returns a list of `Option` objects that contain all of the logic needed to parse a config file into your custom struct fields.
The `Option` function must not lock a mutex in case your implementing struct embeds an anonymous `sync.Mutex`.

I usually fetch key-value pairs from a `.env` file as well the environment variables of your current user session.
That's why the `configo.Parse`function looks the way it does, you pass an `env map[string]string` as a parameter to the funcion as well as a type implementing the `Config` interface.

## Example

In order to create your own custom configuration struct that is supposed to fetch values from your environment or a `.env` file, use a third party package or the `os` package to fetch a map of your envirnonment variables.

Go Playground example: [https://play.golang.org/p/lsyBJv9ItzV](https://play.golang.org/p/lsyBJv9ItzV)

```go
package main

import (
    "encoding/json"
    "fmt"
    "sync"
    "time"

    configo "github.com/jxsl13/simple-configo"
    "github.com/jxsl13/simple-configo/parsers"
    "github.com/jxsl13/simple-configo/unparsers"
)

// MyConfig is a custom configuration that I want to use.
type MyConfig struct {
    sync.Mutex    // optional mutex to make the config goroutine safe
    SomeBool      bool
    SomeInt       int
    SomeFloat     float64
    SomeDelimiter string
    SomeDuration  time.Duration
    SomeList      []string
    SomeStringSet map[string]bool
}

// Options returns a list of available options that can be configured for this
// config object
func (m *MyConfig) Options() (options configo.Options) {
    // WARNING: no locking in this function.
    // NOTE: delimiter is parsed before the other values, this order is important,
    // as the delimiter is used afterwards.
    optionsList := configo.Options{
        {
            Key:             "SOME_BOOL",
            Mandatory:       true,
            Description:     "This is some description text.",
            DefaultValue:    "no",
            ParseFunction:   parsers.Bool(&m.SomeBool),
            UnparseFunction: unparsers.Bool(&m.SomeBool),
        },
        {
            Key:             "SOME_INT",
            Description:     "This is some description text.",
            DefaultValue:    "42",
            ParseFunction:   parsers.Int(&m.SomeInt),
            UnparseFunction: unparsers.Int(&m.SomeInt),
        },
        {
            Key:             "SOME_FLOAT",
            Description:     "This is some description text.",
            DefaultValue:    "99.99",
            ParseFunction:   parsers.Float(&m.SomeFloat, 64),
            UnparseFunction: unparsers.Float(&m.SomeFloat, 64),
        },
        {
            Key:             "SOME_DELIMITER",
            Description:     "delimiter to split the lists below.",
            DefaultValue:    " ",
            ParseFunction:   parsers.String(&m.SomeDelimiter),
            UnparseFunction: unparsers.String(&m.SomeDelimiter),
        },
        {
            Key:             "SOME_DURATION",
            Description:     "This is some description text.",
            DefaultValue:    "24h12m44s",
            ParseFunction:   parsers.Duration(&m.SomeDuration),
            UnparseFunction: unparsers.Duration(&m.SomeDuration),
        },
        {
            Key:             "SOME_LIST",
            Description:     "Some IP list",
            DefaultValue:    "127.0.0.1 127.0.0.2 127.0.0.3",
            ParseFunction:   parsers.List(&m.SomeList, &m.SomeDelimiter),
            UnparseFunction: unparsers.List(&m.SomeList, &m.SomeDelimiter),
        },
        {
            Key:             "SOME_SET",
            Description:     "This is some description text.",
            DefaultValue:    "127.0.0.1 127.0.0.2 127.0.0.3 127.0.0.1",
            ParseFunction:   parsers.ListToSet(&m.SomeStringSet, &m.SomeDelimiter),
            UnparseFunction: unparsers.SetToList(&m.SomeStringSet, &m.SomeDelimiter),
        },
    }

    // add prefix
    for idx := range optionsList {
        optionsList[idx].Key = "MY_" + optionsList[idx].Key
    }

    return optionsList
}

func main() {

    env := map[string]string{
        "MY_SOME_BOOL":      "true",
        "MY_SOME_INT":       "10",
        "MY_SOME_FLOAT":     "12.5",
        "MY_SOME_DELIMITER": ";",
        "MY_SOME_DURATION":  "12h",
        "MY_SOME_LIST":      "99;15;13;77",
        "MY_SOME_SET":       "99;15;13;77;99",
    }

    myCfg := &MyConfig{}
    if err := configo.Parse(env, myCfg); err != nil {
        panic(err)
    }

    for _, opt := range myCfg.Options() {
        fmt.Println(opt.String())
    }

    b, err := json.MarshalIndent(&myCfg, " ", " ")
    if err != nil {
        panic(err)
    }
    fmt.Println(string(b))


    newEnvMap, err := configo.Unparse(myCfg)
    if err != nil {
        panic(err)
    }
    // write map to file, update some database, redis, etc.
    b, err = json.MarshalIndent(&newEnvMap, " ", " ")
    if err != nil {
        panic(err)
    }
    fmt.Println(string(b))
}
```

This is everything you need to write in order to parse a configuration file with key value pairs into a struct of your choise.
