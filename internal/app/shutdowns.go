package app

import (
	"context"
	"errors"
	"fmt"
	"os"
)

type Shutdowns struct {
	shutdowns []func(context.Context) error

	OnShutdownError FatalFn
}

func newShutdowns() Shutdowns {
	return Shutdowns{
		OnShutdownError: onShutdownError,
	}
}

// AddShutdown adds f to the list of shutdowns
func (s *Shutdowns) AddShutdown(f func(context.Context) error) {
	s.shutdowns = append(s.shutdowns, f)
}

const StartupPhase = "Startup"
const ShutdownPhase = "Shutdown"

const (
	OkCode            = 0
	ErrorCode         = 1
	ShutdownErrorCode = 2
)

func (s *Shutdowns) ShutdownAll(phase string, fatalErr error) int {
	exitCode := OkCode
	if fatalErr != nil {
		exitCode = ErrorCode
	}

	err := fatalErr
	ctx := context.Background()
	for _, shutdown := range s.shutdowns {
		if shutdownErr := shutdown(ctx); shutdownErr != nil {
			exitCode = ShutdownErrorCode
			err = errors.Join(err, shutdownErr)
		}
	}
	if exitCode == ShutdownErrorCode {
		s.OnShutdownError(phase, err)
	}
	return exitCode
}

type FatalFn func(string, error)

// onShutdownError is the default fatal error handler
func onShutdownError(phase string, fatalErr error) {
	// default fatal error handler
	fmt.Fprintf(os.Stderr, "%s Fatal Error: %v\n", phase, fatalErr)
}
