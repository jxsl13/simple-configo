package configo

// ActionFunc is a function that does something
// but contrary to ParseFunc doe snot evaluate the provided option value.
// An option that only consists of actions is in fact an action that is not
// expected to be found in the string map, e.g. env map, env file, etc..
type ActionFunc func() error

func tryExecAction(f func() error) error {
	if f == nil {
		return nil
	}
	return f()
}
