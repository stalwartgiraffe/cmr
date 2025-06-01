package utils

import (
	"io"
	"io/fs"
)

type PathReaderFn func(path string, r io.Reader)

// WalkFileReaders runs read on each file in dir.
func WalkFileReaders(dir fs.FS, read PathReaderFn) error {
	return fs.WalkDir(dir, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		file, err := dir.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		read(path, file)
		return nil
	})
}
