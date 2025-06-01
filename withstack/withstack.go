// Package withstack wraps an error with an error and stack trace.
// The package does not depend on debug but printing stack trace is not
// recommended for performance critical code called very frequently.
// The errors package is used to wrap passed in errors.
package withstack

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// Errorf formats according to a format specifier and returns the string as a
// value along with the stack string that satisfies error.
func Errorf(format string, a ...any) error {
	return fmt.Errorf("%w\n%s",
		fmt.Errorf(format, a...),
		StackTrace(),
	)
}

// New return a new error with attached stack trace.
func New(msg string) error {
	return errors.New(msg + "\n" + StackTrace())
}

// StackTrace returns the stack trace in a newline delimited string.
func StackTrace() string {
	var b bytes.Buffer
	w := io.Writer(&b)
	stk := make([]uintptr, 32)
	n := runtime.Callers(0, stk)
	printStackRecord(w, n-2, stk, false)
	return b.String()
}

// printStackRecord prints the function + source line information
// for a single stack trace.
func printStackRecord(w io.Writer, n int, stk []uintptr, allFrames bool) {
	show := allFrames
	frames := runtime.CallersFrames(stk)
	// skip withstack
	_, _ = frames.Next()
	_, _ = frames.Next()
	_, _ = frames.Next()
	for i := 0; i < n; i++ {
		frame, more := frames.Next()
		name := frame.Function
		if name == "" {
			show = true
			fmt.Fprintf(w, "%#x\n", frame.PC)
		} else if name != "runtime.goexit" && (show || !strings.HasPrefix(name, "runtime.")) {
			// Hide runtime.goexit and any runtime functions at the beginning.
			// This is useful mainly for allocation traces.
			show = true
			fmt.Fprintf(w, "%s\t%s:%d\n", name, frame.File, frame.Line)
		}
		if !more {
			break
		}
	}
	if !show {
		// We didn't print anything; do it again,
		// and this time include runtime functions.
		printStackRecord(w, n, stk, true)
		return
	}
	fmt.Fprintf(w, "\n")
}

// WithFileLine adds the filename and line number as a prefix to the msg.
func WithFileLine(msg string) string {
	_, filename, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d %s", filename, line, msg)
}
