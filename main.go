package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/stalwartgiraffe/cmr/cmd"
)

func main() {
	// Create a context that can be cancelled
	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Run a goroutine. for the long running case, cancel the context when an interrupt signal is received
	go func() {
		<-sigChan
		cancel()
	}()

	rootCmd := cmd.AddRootCommand(cancel)
	if err := rootCmd.ExecuteContext(mainCtx); err != nil {
		os.Exit(1) // failing command is expected to emit to std err if needed.
	}
}
