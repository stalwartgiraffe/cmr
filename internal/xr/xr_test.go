package xr

import (
	"fmt"
	"io"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type errRunner struct {
	err error
}

func (r *errRunner) Run() error {
	return r.err
}

func mockFuncs(runnerErr error) *funcs {
	return &funcs{
		environ: func() []string {
			return nil
		},
		getwd: func() (dir string, err error) {
			return "/", nil
		},

		lookPath: func(file string) (string, error) {
			return "/user/bin/noopcmd", nil
		},
		makeRunner: func(
			dir string,
			env []string,
			sout io.Writer,
			serr io.Writer,
			name string,
			arg ...string) Runner {
			return &errRunner{
				err: runnerErr,
			}
		},
	}
}

var _ = Describe("RunAt", func() {
	DescribeTable("command lines ",
		func(want string, isErr bool, dir string, allowedStatus int, name string, fn Funcs, args ...string) {
			have, err := RunAt(dir, allowedStatus, name, fn, args...)
			Expect(have).To(Equal(want))
			Expect(err != nil).To(Equal(isErr))
		},
		Entry("run noop", "", false, "/", 1, "noop", mockFuncs(nil), "dummy.txt"),

		Entry("run err", "", true, "/", 1, "run error",
			func() Funcs {
				f := mockFuncs(nil)
				f.lookPath = func(file string) (string, error) {
					return "", fmt.Errorf("exec.LookPath err")
				}
				return f
			}(),
			"dummy.txt"),

		Entry("run err", "", true, "/", 1, "run error",
			mockFuncs(
				fmt.Errorf("unexpected exit error"),
			),
			"dummy.txt"),
	)
})
