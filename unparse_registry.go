package configo

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

var (
	shutdownJobs *jobList
	initOnce     sync.Once

	shutdownChan      chan os.Signal
	shutdownAwaitFunc func() error
)

func init() {
	shutdownJobs = &jobList{
		jobs: make([]job, 0),
	}
}

// Unparse returns a function that you can call 'unparse'.
// That function is supposed to be called via 'defer unparse()' in the main function
// in order to execute the shutdown hooks that were registered in the provided Config options
// with the UnparserFunc option.
func Unparse(cfg Config, env map[string]string) func() error {
	return unparse(cfg.Options(), env)
}

func unparse(options Options, env map[string]string) func() error {

	// initialization
	initOnce.Do(func() {
		// blocking function that awaits the closing of the channel
		shutdownAwaitFunc = func() error {

			// do the actual parsing
			if err := shutdownJobs.ExecuteJobs(); err != nil {
				return err
			}
			return nil
		}
	})

	// add more shutdown hooks on every invokation if they are defined
	for _, option := range options {
		if option.UnparseFunction == nil {
			continue
		}
		value, ok := env[option.Key]
		// only add to shutdown hooks if value actually exists and is not empty
		// will be added to shutdown hooks in case the value exists, but is empty
		if ok {
			shutdownJobs.Add(option.Key, value, option.UnparseFunction)
		} else {
			shutdownJobs.Add(option.Key, option.DefaultValue, option.UnparseFunction)
		}
	}

	return shutdownAwaitFunc
}

type job struct {
	Key   string
	Value string
	Func  UnparserFunc
	err   error
}

func (j *job) Execute() error {
	if err := j.Func(j.Key, j.Value); err != nil {
		j.err = err
		return err
	}
	return nil
}

func (j job) Error() string {
	if j.err != nil {
		return fmt.Sprintf("%s : %v", j.Key, j.err)
	}
	return ""
}

type jobList struct {
	sync.Mutex
	jobs []job
}

// Add a new shutdown job
func (jl *jobList) Add(key, value string, unparseFunc UnparserFunc) {
	if unparseFunc == nil {
		return
	}
	jl.Lock()
	defer jl.Unlock()

	jl.jobs = append(jl.jobs, job{
		Key:   key,
		Value: value,
		Func:  unparseFunc,
	})
}

func (jl *jobList) ExecuteJobs() error {
	jl.Lock()
	defer jl.Unlock()

	errs := make([]error, 0)
	// execute jobs in reverse
	for i := len(jl.jobs) - 1; i >= 0; i-- {
		job := jl.jobs[i]
		if err := job.Execute(); err != nil {
			// job implements the error interface
			errs = append(errs, job)
		}
	}

	// reset joblist
	jl.jobs = jl.jobs[:0]

	// multi line error
	if len(errs) > 0 {
		return constructErr(errs)
	}
	return nil
}

func (jl *jobList) Size() int {
	jl.Lock()
	defer jl.Unlock()
	return len(jl.jobs)
}

// constructErr constructs a multi line error
func constructErr(errs []error) error {
	sb := strings.Builder{}
	sb.Grow(len(errs) * 16)
	for idx, err := range errs {
		sb.WriteString(fmt.Sprintf("%-3d: %s\n", idx+1, err))
	}
	return errors.New(sb.String())
}
