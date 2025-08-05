package run

import (
	"context"

	"github.com/stalwartgiraffe/cmr/internal/app"
)

type FakeRunner struct {
	Runner

	app *app.App

	cancelCalled bool

	shutdownCalled bool
	shutdownPhase  *string
	shutdownErr    *error

	onFatalCalled bool
	fatalPhase    *string
	fatalErr      *error
}

func NewFakeRunner() *FakeRunner {
	fake := &FakeRunner{
		app: &app.App{},
	}

	fake.makeCtx = func() (context.Context, context.CancelFunc) {
		return context.Background(), func() {
			fake.cancelCalled = true
		}
	}
	fake.makeApp = func(ctx context.Context) app.AppErr {
		return app.AppErr{
			App: fake.app,
		}
	}
	fake.beginInterrupt = func(ctxCancel context.CancelFunc) {
		// no-op for test
	}
	fake.runCmd = func(ctx context.Context, ctxCancel context.CancelFunc, app *app.App) error {
		return nil
	}
	fake.onShutdown = func(phase string, _ *app.App, fatalErr error) int {
		fake.shutdownCalled = true
		fake.shutdownPhase = &phase
		fake.shutdownErr = &fatalErr
		if fatalErr == nil {
			return app.OkCode
		}
		return app.ErrorCode
	}

	fake.app.OnShutdownError = func(phase string, fatalErr error) {
		fake.onFatalCalled = true
		fake.fatalPhase = &phase
		fake.fatalErr = &fatalErr
	}

	return fake
}
