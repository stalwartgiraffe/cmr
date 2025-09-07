// Package reload watches directories
package reload

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func BeginWatchPwd(ctxCancel context.CancelFunc) error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}
	return BeginWatch(ctxCancel, pwd)
}

func BeginWatch(ctxCancel context.CancelFunc, rootDir string) error {
	paths, err := getSubdirectories(rootDir)
	if err != nil {
		return fmt.Errorf("could not find directories under %s: %w", rootDir, err)
	}
	if len(paths) < 1 {
		return errors.New("must specify at least one path to watch")
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("creating a new watcher: %w", err)
	}

	// Start listening for events.
	go runWatch(ctxCancel, w)

	// watch these paths
	for _, p := range paths {
		if err := w.Add(p); err != nil {
			return fmt.Errorf("%q: %w", p, err)
		}
	}
	return nil
}

func runWatch(ctxCancel context.CancelFunc, w *fsnotify.Watcher) {
	for {
		select {
		case err, ok := <-w.Errors:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			panic(err) // some other error
		case _, ok := <-w.Events:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			w.Close()
			ctxCancel()
			return
		}
	}
}

func getSubdirectories(dirPath string) ([]string, error) {
	if isIgnore(dirPath) {
		return nil, nil
	}
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory '%s': %w", dirPath, err)
	}
	dirs := []string{dirPath}
	for _, entry := range entries {
		if entry.IsDir() {
			subPath := path.Join(dirPath, entry.Name())
			if subs, err := getSubdirectories(subPath); err != nil {
				return nil, err
			} else {
				dirs = append(dirs, subs...)
			}
		}
	}

	return dirs, nil
}

func isIgnore(dir string) bool {
	return strings.HasSuffix(dir, ".git") || strings.HasSuffix(dir, "ignore")
}
