package utils

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestGetDirs(t *testing.T) {
	tests := []struct {
		name     string
		fs       fs.FS
		expected []string
	}{
		{
			name: "empty filesystem",
			fs:   fstest.MapFS{},
			expected: []string{
				".",
			},
		},
		{
			name: "filesystem with files only",
			fs: fstest.MapFS{
				"file1.txt": &fstest.MapFile{Data: []byte("content1")},
				"file2.go":  &fstest.MapFile{Data: []byte("content2")},
			},
			expected: []string{
				".",
			},
		},
		{
			name: "filesystem with directories",
			fs: fstest.MapFS{
				"dir1/file.txt":      &fstest.MapFile{Data: []byte("content")},
				"dir2/subdir/go.mod": &fstest.MapFile{Data: []byte("module test")},
				"dir2/emptydir/":     &fstest.MapFile{Data: nil},
				"file.txt":           &fstest.MapFile{Data: []byte("root file")},
			},
			expected: []string{
				".",
				"dir1",
				"dir2",
				"dir2/subdir",
				"dir2/emptydir",
			},
		},
		{
			name: "nested directory structure",
			fs: fstest.MapFS{
				"a/b/c/file.txt": &fstest.MapFile{Data: []byte("deep file")},
				"a/b/file2.txt":  &fstest.MapFile{Data: []byte("file2")},
				"a/file3.txt":    &fstest.MapFile{Data: []byte("file3")},
				"x/y/file4.txt":  &fstest.MapFile{Data: []byte("file4")},
				"root.txt":       &fstest.MapFile{Data: []byte("root")},
			},
			expected: []string{
				".",
				"a",
				"a/b",
				"a/b/c",
				"x",
				"x/y",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual []string

			err := WalkDirs(tt.fs, func(path string) {
				actual = append(actual, path)
			})

			require.NoError(t, err)
			require.ElementsMatch(t, tt.expected, actual)
		})
	}
}

func TestGetDirs_EarlyReturn(t *testing.T) {
	fs := fstest.MapFS{
		"dir1/file.txt": &fstest.MapFile{Data: []byte("content")},
		"dir2/file.txt": &fstest.MapFile{Data: []byte("content")},
		"dir3/file.txt": &fstest.MapFile{Data: []byte("content")},
	}

	var actual []string
	ok := false
	err := WalkDirs(fs, func(path string) {
		actual = append(actual, path)
		if path == "dir2" {
			ok = true
		}
	})

	require.True(t, ok)
	require.NoError(t, err)
	require.Contains(t, actual, ".")
	require.Contains(t, actual, "dir2")
}
