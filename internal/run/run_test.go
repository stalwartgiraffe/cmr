package run

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stalwartgiraffe/cmr/internal/app"
)

func TestRunnerRun(t *testing.T) {
	const (
		runCmdErrorCode   = 13
		shutdownErrorCode = 42
	)
	t.Run("fake execution", func(t *testing.T) {
		t.Parallel()
		fake := NewFakeRunner()
		require.False(t, fake.cancelCalled)
		exitCode := fake.Run()
		require.Equal(t, app.OkCode, exitCode)
		require.True(t, fake.cancelCalled)
		require.True(t, fake.shutdownCalled)
		require.False(t, fake.onFatalCalled)
	})

	t.Run("RunCmd returns a fatal error", func(t *testing.T) {
		t.Parallel()
		fake := NewFakeRunner()

		runCmdFatalErr := errors.New("test RunCmd has failed fatally")
		fake.runCmd = func(ctx context.Context, ctxCancel context.CancelFunc, app *app.App) error {
			return runCmdFatalErr
		}

		baseShutdown := fake.onShutdown
		fake.onShutdown = func(phase string, app *app.App, fatalErr error) int {
			baseShutdown(phase, app, fatalErr)
			require.Equal(t, runCmdFatalErr, fatalErr)
			return runCmdErrorCode
		}
		exitCode := fake.Run()
		require.Equal(t, runCmdErrorCode, exitCode)
		require.True(t, fake.cancelCalled)
		require.True(t, fake.shutdownCalled)
		require.False(t, fake.onFatalCalled)
	})

	t.Run("makeApp returns a fatal error", func(t *testing.T) {
		t.Parallel()
		fake := NewFakeRunner()

		initFatalErr := errors.New("init failed fatally")
		fake.makeApp = func(ctx context.Context) app.AppErr {
			return app.AppErr{
				App: fake.app,
				Err: initFatalErr,
			}
		}
		exitCode := fake.Run()
		require.Equal(t, app.ErrorCode, exitCode)
		require.True(t, fake.cancelCalled)
		require.True(t, fake.shutdownCalled)
		require.False(t, fake.onFatalCalled)
	})

	t.Run("makeApp returns a fatal error, shutdown err", func(t *testing.T) {
		t.Parallel()
		fake := NewFakeRunner()

		initFatalErr := errors.New("init failed fatally")
		fake.makeApp = func(ctx context.Context) app.AppErr {
			return app.AppErr{
				App: fake.app,
				Err: initFatalErr,
			}
		}
		fake.onShutdown = func(phase string, app *app.App, fatalErr error) int {
			fake.shutdownCalled = true
			fake.shutdownPhase = &phase
			fake.shutdownErr = &fatalErr
			return shutdownErrorCode
		}
		exitCode := fake.Run()
		require.Equal(t, shutdownErrorCode, exitCode)
		require.True(t, fake.cancelCalled)
		require.True(t, fake.shutdownCalled)
		require.False(t, fake.onFatalCalled)
	})

	t.Run("runCmd returns an error", func(t *testing.T) {
		t.Parallel()
		fake := NewFakeRunner()

		runFatalErr := errors.New("run error")
		fake.runCmd = func(ctx context.Context, ctxCancel context.CancelFunc, app *app.App) error {
			return runFatalErr
		}
		exitCode := fake.Run()
		require.Equal(t, app.ErrorCode, exitCode)
		require.True(t, fake.cancelCalled)
		require.True(t, fake.shutdownCalled)
		require.False(t, fake.onFatalCalled)
	})
}
