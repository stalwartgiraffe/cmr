// Package xr wraps the go exec package with helper functions.
package xr

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

// Run will run program name at cwd with args.
// Will return the output written fully into a string or error.
func Run(
	ctx context.Context,
	name string,
	allowedStatus int,
	fn Funcs,
	args ...string,
) (string, error) {
	if fn == nil {
		fn = newFuncs()
	}
	if dir, err := fn.Getwd(); err != nil {
		return "", err
	} else {
		return RunAt(ctx, dir, allowedStatus, name, fn, args...)
	}
}

// RunAt will run program name at dir with args.
// Will return the output written fully into a string or error.
func RunAt(
	ctx context.Context,
	dir string,
	allowedStatus int,
	name string,
	fn Funcs,
	args ...string,
) (string, error) {
	if _, err := fn.LookPath(name); err != nil {
		return "", err
	}

	env := fn.Environ()
	args = expandArgs(env, args)

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	runner := fn.MakeRunner(ctx, dir, env, &stdOut, &stdErr, name, args...)

	if err := runner.Run(); err != nil {
		exiterr, ok := err.(*exec.ExitError)
		if !ok {
			return "", fmt.Errorf("%s had unexpected exit: %w", name, err)
		}
		code := exiterr.ExitCode()
		if code != 0 && code != allowedStatus {
			errMsg := stdErr.String()
			if len(errMsg) < 1 {
				errMsg = exiterr.Error()
			}
			errMsg += "\nout\n" + stdOut.String() + "\nerr\n" + stdErr.String()
			return "", fmt.Errorf("%s had status %s", name, errMsg)
		}
	}

	return stdOut.String(), nil
}
