package xr

import (
	"bytes"
	"fmt"
	"os/exec"
)

type Linter struct {
	Name          string
	AllowedStatus int
	CmdArgs       []string
}

// CheckLints run each Lint. Stop on failure.
func CheckEach(linters []Linter, fn Funcs) error {
	for _, a := range linters {
		out, err := Run(a.Name, a.AllowedStatus, fn, a.CmdArgs...)
		if err != nil {
			return err
		}
		if len(out) < 1 {
			return fmt.Errorf(out)
		}
	}
	return nil
}

// Run will run program name at cwd with args.
// Will return the output written fully into a string or error.
func Run(name string, allowedStatus int, fn Funcs, args ...string) (string, error) {
	if fn == nil {
		fn = newFuncs()
	}
	if dir, err := fn.Getwd(); err != nil {
		return "", err
	} else {
		return RunAt(dir, allowedStatus, name, fn, args...)
	}
}

// RunAt will run program name at dir with args.
// Will return the output written fully into a string or error.
func RunAt(dir string, allowedStatus int, name string, fn Funcs, args ...string) (string, error) {
	if _, err := fn.LookPath(name); err != nil {
		return "", err
	}

	env := fn.Environ()
	args = expandArgs(env, args)

	var out bytes.Buffer
	var serr bytes.Buffer

	runner := fn.MakeRunner(dir, env, &out, &serr, name, args...)

	if err := runner.Run(); err != nil {
		exiterr, ok := err.(*exec.ExitError)
		if !ok {
			return "", fmt.Errorf("%s had unexpected exit: %w", name, err)
		}
		code := exiterr.ExitCode()
		if code != 0 && code != allowedStatus {
			errMsg := serr.String()
			if len(errMsg) < 1 {
				errMsg = exiterr.Error()
			}
			errMsg += "\nout\n" + out.String() + "\nerr\n" + serr.String()
			return "", fmt.Errorf("%s had status %s", name, errMsg)
		}
	}

	return out.String(), nil
}
