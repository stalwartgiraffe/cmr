package elog

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

type ElapsedLogger struct {
	start time.Time
	out   io.Writer
	flush flushFn
}

func New() *ElapsedLogger {
	return &ElapsedLogger{
		start: time.Now(),
		out:   os.Stdout,
		flush: noop,
	}
}

type flushFn func()

func noop() {
	// do nothing
}

func NewWithBuffer() *ElapsedLogger {
	sb := &strings.Builder{}
	return &ElapsedLogger{
		start: time.Now(),
		out:   sb,
		flush: func() {
			fmt.Println(sb.String())
		},
	}
}

func (l *ElapsedLogger) Flush() {
	l.flush()
}

// lastName returns the last name in a slash delimited path
func lastName(path string) string {
	short := path
	for i := len(path) - 1; i > 0; i-- {
		if path[i] == '/' { // FIXME ask os path delimiter
			short = path[i+1:]
			break
		}
	}
	return short
}

// prefix returns a customer log prefix.
// milliseconds_elapsed file:func:line
func (l *ElapsedLogger) prefix() string {
	elapsed := time.Since(l.start)
	ms := elapsed.Milliseconds()
	us := elapsed.Microseconds() % 1000
	pc, file, line, ok := runtime.Caller(2)
	funcName := "?"
	if ok {
		funcName = lastName(runtime.FuncForPC(pc).Name())
	}
	file = lastName(file)
	return fmt.Sprintf("%5d.%03d %16s:%s:%d ", ms, us, file, funcName, line)
}

// Implement class standard log package prints interface.

func (l *ElapsedLogger) Println(v ...any) {
	fmt.Fprint(l.out, l.prefix())
	fmt.Fprintln(l.out, v...)
}

func (l *ElapsedLogger) Printf(format string, v ...any) {
	fmt.Fprint(l.out, l.prefix())
	fmt.Fprintf(l.out, format, v...)
}

func (l *ElapsedLogger) Print(v ...any) {
	fmt.Fprint(l.out, l.prefix())
	fmt.Fprint(l.out, v...)
}

type NoopLogger struct {
}

func (l *NoopLogger) Println(v ...any) {
}
func (l *NoopLogger) Printf(format string, v ...any) {
}
func (l *NoopLogger) Print(v ...any) {
}
func (l *NoopLogger) Flush() {
}
