// Package run enables mocking the app singleton for integration tests.
package run

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/stalwartgiraffe/cmr/cmd"
	"github.com/stalwartgiraffe/cmr/internal/app"
)

type Runner struct {
	makeCtx        makeCtxFn
	makeApp        makeAppFn
	beginInterrupt beginInterruptFn
	runCmd         runCmdFn
	onShutdown     onShutdownFn
}

func NewRunner() *Runner {
	return &Runner{
		makeCtx:        makeCtx,
		makeApp:        makeApp,
		beginInterrupt: beginInterrupt,
		runCmd:         runCmd,
		onShutdown:     onShutdown,
	}
}

// Run the code application. All dependencies are inject-able.
func (r *Runner) Run() int {
	ctx, ctxCancel := r.makeCtx()
	defer ctxCancel()
	appErr := r.makeApp(ctx)
	if appErr.Err != nil {
		return r.onShutdown(app.StartupPhase, appErr.App, appErr.Err)
	}
	r.beginInterrupt(ctxCancel)
	runErr := r.runCmd(ctx, ctxCancel, appErr.App)
	return r.onShutdown(app.ShutdownPhase, appErr.App, runErr)
}

type makeCtxFn func() (context.Context, context.CancelFunc)

// makeCtx makes a cancel-able context
func makeCtx() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

type makeAppFn func(context.Context) app.AppErr

// makeApp initializes application singletons such as
// otel metrics, traces and logs exports
func makeApp(ctx context.Context) app.AppErr {
	const otelSchema = "cmr_cli"
	return app.NewAppErr()
		
	//return app.NewAppErr().
//		WithOtel(ctx, otelSchema)
}

type beginInterruptFn func(ctxCancel context.CancelFunc)

// beginInterrupt begin a go routine that handles system interrupt signals.
func beginInterrupt(ctxCancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan   // block on interrupt
		ctxCancel() // cancel the main context
	}()
}

type runCmdFn func(ctx context.Context, ctxCancel context.CancelFunc, app *app.App) error

// runCmd creates and runs the cmd
func runCmd(ctx context.Context, ctxCancel context.CancelFunc, app *app.App) error {
	rootCmd := cmd.AddRootCommand(app, ctxCancel)
	return rootCmd.ExecuteContext(ctx)
}

type onShutdownFn func(string, *app.App, error) int

// onShutdown handles the shutdown event and runs what ever special shutdown handlers that have been added.
func onShutdown(phase string, app *app.App, fatalErr error) int {
	return app.ShutdownAll(phase, fatalErr)
}
