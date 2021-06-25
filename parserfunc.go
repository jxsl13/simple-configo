package configo

// ParserFunc is a custom parser function that can be used to parse specific option values
// A Option struct must contain a ParseFunc in order to know, how to parse a specific value and where the
// Such a function is usually created using a generator function, that specifies the output type.
// This function is kept as simple as possible, in order to be handled exactly the same way for every
// possible return value
type ParserFunc func(value string) error

// UnparserFunc through another function that provides a reference to the actual configuration value.
// Contrary to the ParserFunc it does provide the key as well as the value of the parent option.
type UnparserFunc func(key, value string) error
