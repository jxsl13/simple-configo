package configo

// ParserFunc is a custom parser function that can be used to parse specific option values
// A Option struct must contain a ParseFunc in order to know, how to parse a specific value and where the
// Such a function is usually created using a generator function, that specifies the output type.
// This function is kept as simple as possible, in order to be handled exactly the same way for every
// possible return value
type ParserFunc func(value string) error

// UnparserFunc is a function that receives a key and returns the key's value
// UnparseFunctions go back to creating a map[string]string from the previously parse configuration struct.
type UnparserFunc func() (string, error)
