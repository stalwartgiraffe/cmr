package xr

import (
	"context"
	"fmt"
)

type Linter struct {
	Name          string
	AllowedStatus int
	CmdArgs       []string
}

// CheckEach run each Lint. Stop on failure.
func CheckEach(ctx context.Context, linters []Linter, fn Funcs) error {
	for _, a := range linters {
		out, err := Run(ctx, a.Name, a.AllowedStatus, fn, a.CmdArgs...)
		if err != nil {
			return err
		}
		if len(out) < 1 {
			return fmt.Errorf(out)
		}
	}
	return nil
}
