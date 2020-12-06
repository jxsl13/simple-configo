# Simple-Configo

Simple-Configo streamlines the creation of multiple independent configuration structs to the implementation of a single interface, the `Config` interface, that has only two methods.

```go
type Config interface {
    Name() string
    Options() (options Options)
}
```

The `Name` method simply returns a string of your choice.
The `Options` method returns a list of `Option` objects that contain all of the logic needed to parse a config file into your custom struct fields.

I usually fetch key-value pairs from a `.env` file as well the environment variables of your current user session.
That's why the `configo.Parse`function looks the way it does, you pass an `env map[string]string` as a parameter to the funcion as well as a type implementing the `Config` interface.

## Example

In order to create your own custom configuration struct that is supposed to fetch values from your environment or a `.env` file, use a third party package or the `os` package to fetch a map of your envirnonment variables.

Go Playground example: [CLICK ME](https://play.golang.org/p/MRJxvSyzc0d)

```go
package main

import (
    "encoding/json"
    "fmt"
    "time"

    configo "github.com/jxsl13/simple-configo"
)

// MyConfig is a custom configuration that I want to use.
type MyConfig struct {
    SomeBool      bool
    SomeInt       int
    SomeFloat     float64
    SomeDelimiter string
    SomeDuration  time.Duration
    SomeList      []string
    SomeStringSet map[string]bool
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
            Type:          "bool",
            Mandatory:     true,
            Description:   "This is some description text.",
            DefaultValue:  "no",
            ParseFunction: configo.DefaultParserBool(&m.SomeBool),
        },
        {
            Key:           "SOME_INT",
            Type:          "int",
            Description:   "This is some description text.",
            DefaultValue:  "42",
            ParseFunction: configo.DefaultParserInt(&m.SomeInt),
        },
        {
            Key:           "SOME_FLOAT",
            Type:          "float",
            Description:   "This is some description text.",
            DefaultValue:  "99.99",
            ParseFunction: configo.DefaultParserFloat(&m.SomeFloat, 64),
        },
        {
            Key:           "SOME_DELIMITER",
            Type:          "string",
            Description:   "delimiter to split the lists below.",
            DefaultValue:  " ",
            ParseFunction: configo.DefaultParserString(&m.SomeDelimiter),
        },
        {
            Key:           "SOME_DURATION",
            Type:          "duration",
            Description:   "This is some description text.",
            DefaultValue:  "24h12m44s",
            ParseFunction: configo.DefaultParserDuration(&m.SomeDuration),
        },
        {
            Key:           "SOME_LIST",
            Type:          "list",
            Description:   "Some IP list",
            DefaultValue:  "127.0.0.1 127.0.0.2 127.0.0.3",
            ParseFunction: configo.DefaultParserList(m.SomeDelimiter, &m.SomeList),
        },
        {
            Key:           "SOME_SET",
            Type:          "",
            Description:   "This is some description text.",
            DefaultValue:  "127.0.0.1 127.0.0.2 127.0.0.3 127.0.0.1",
            ParseFunction: configo.DefaultParserListToSet(m.SomeDelimiter, &m.SomeStringSet),
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
    if err := configo.Parse(myCfg, env); err != nil {
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
}
```

This is everything you need to write in order to parse a configuration file with key value pairs into a struct of your choise.

## TODO

- Remove 'Type' field of 'Option' struct, as all that's needed is actually the ParseFunction, ut gotta look further into it.
