package reload

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestWatchDirsRecursively(t *testing.T) {
	fs := fstest.MapFS{
		"./skip.txt":      &fstest.MapFile{Data: []byte("git file")},
		".git/skip.txt":   &fstest.MapFile{Data: []byte("git file")},
		"ignore/skip.txt": &fstest.MapFile{Data: []byte("skip file")},
		"a/b/c/file.txt":  &fstest.MapFile{Data: []byte("deep file")},
		"a/b/file2.txt":   &fstest.MapFile{Data: []byte("file2")},
		"a/file3.txt":     &fstest.MapFile{Data: []byte("file3")},
		"x/y/file4.txt":   &fstest.MapFile{Data: []byte("file4")},
		"x/y/file5.txt":   &fstest.MapFile{Data: []byte("file6")},
		"root.txt":        &fstest.MapFile{Data: []byte("root")},
	}
	watcher := &mockWatcher{}
	err := watchDirsRecursively(fs, watcher)
	require.NoError(t, err)
	expectedDirs := []string{"a", "a/b", "a/b/c", "x", "x/y"}
	require.Equal(t, expectedDirs, watcher.dirs)
}

type mockWatcher struct {
	dirs []string
}

func (w *mockWatcher) Add(p string) error {
	w.dirs = append(w.dirs, p)
	return nil
}
