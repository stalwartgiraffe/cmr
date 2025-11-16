package xr

import (
	"context"
	"io"
	"os"
	"os/exec"
)

// can be run
type Runner interface {
	Run() error
}

// Testing exec well is a pain.
// see https://abhinavg.net/2022/05/15/hijack-testmain/
//
// Function dependencies for unit testing code that calls functions.
// This enables monkey patching in test.
// This works but requires some manual work.
// Alternatives - use an
// Actual monkey patching
//
//	https://bou.ke/blog/monkey-patching-in-go/
//	either hack the go runtime
//	https://github.com/undefinedlabs/go-mpatch
type funcs struct {
	environ func() []string
	getwd   func() (dir string, err error)

	lookPath func(file string) (string, error)

	makeRunner func(
		ctx context.Context,
		dir string,
		env []string,
		stdOut io.Writer,
		stdErr io.Writer,
		name string,
		arg ...string) Runner
}

func newCmdCtxRunner(
	ctx context.Context,
	dir string,
	env []string,
	stdOut io.Writer,
	stdErr io.Writer,
	name string,
	args ...string) Runner {
	c := exec.CommandContext(ctx, name, args...)
	c.Dir = dir       // working directory
	c.Env = env       // process environment to the child process
	c.Stdout = stdOut // writer for standard out
	c.Stderr = stdErr // writer for standard err
	return c          // ful
}

// newFuncs returns default dependencies
func newFuncs() *funcs {
	return &funcs{
		os.Environ,
		os.Getwd,

		exec.LookPath,
		newCmdCtxRunner,
	}
}
func (f *funcs) Environ() []string {
	return f.environ()
}
func (f *funcs) Getwd() (string, error) {
	return f.getwd()
}

func (f *funcs) LookPath(file string) (string, error) {
	return f.lookPath(file)
}

func (f *funcs) MakeRunner(
	ctx context.Context,
	dir string,
	env []string,
	stdOut io.Writer,
	stdErr io.Writer,
	name string,
	args ...string) Runner {
	return f.makeRunner(ctx, dir, env, stdOut, stdErr, name, args...)
}

// could be mocked
type Funcs interface {
	Environ() []string
	Getwd() (dir string, err error)

	LookPath(file string) (string, error)

	MakeRunner(
		ctx context.Context,
		dir string,
		env []string,
		stdOut io.Writer,
		stdErr io.Writer,
		name string,
		args ...string) Runner
}

var _ Funcs = &funcs{} // enforce defaults support interface
