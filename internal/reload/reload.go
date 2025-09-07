// Package reload watches directories
package reload

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"

	"github.com/stalwartgiraffe/cmr/internal/utils"
)

func BeginWatchPwd(ctxCancel context.CancelFunc) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("creating a new watcher: %w", err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}
	fs := os.DirFS(pwd)
	if err := watchDirsRecursively(fs, watcher); err != nil {
		return err
	}

	// Start listening for events.
	go runWatch(ctxCancel, watcher)

	return nil
}

func runWatch(ctxCancel context.CancelFunc, watcher *fsnotify.Watcher) {
	for {
		select {
		case err, _ := <-watcher.Errors:
			panic(err) // some other error
			return
		case _, ok := <-watcher.Events:
			if ok { // Channel was closed (i.e. Watcher.Close() was called).
				watcher.Close()
				ctxCancel()
			}
			return
		}
	}
}

type Watcher interface {
	Add(string) error
}

// func watchDirsRecursively(fs fs.FS, w *fsnotify.Watcher) error {
func watchDirsRecursively(fs fs.FS, w Watcher) error {
	var errs error
	ignorePrefix := []string{".git", "ignore"}
	utils.WalkDirs(fs, func(path string) {
		for _, p := range ignorePrefix {
			if strings.HasPrefix(path, p) {
				return
			}
		}
		if err := w.Add(path); err != nil {
			errs = errors.Join(errs, err)
		}
	})

	return errs
}
