package src

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// WorkDir represents a directory.
// methods defined on WorkDir always returns a path
// as returned by filepath.ToSlash to make them non-os dependant.
type WorkDir struct {
	abs  string
	path string
}

func (wd *WorkDir) String() string {
	return fmt.Sprintf("%s|%s", wd.path, wd.abs)
}

// NewWorkDir returns new WorkDir.
func NewWorkDir(path string) (*WorkDir, error) {
	if path == "" {
		return nil, errors.New("empty path")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	var wd = &WorkDir{
		abs:  filepath.ToSlash(abs),
		path: path,
	}
	return wd, nil
}

// New returns new WorkDir relative to wd.
func NewRelativeWorkDir(wd *WorkDir, path ...string) (*WorkDir, error) {
	return NewWorkDir(wd.JoinAbs(path...))
}

// IsDir() returns bool indicating if WorkDir is a directory.
func (wd *WorkDir) IsDir() (bool, error) {
	info, err := wd.Stat()
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// Stat wraps os.Stat
func (wd *WorkDir) Stat() (os.FileInfo, error) {
	return os.Stat(wd.Abs())
}

// WalkDir wraps filepath.WalkDir.
func (wd *WorkDir) WalkDir(fn fs.WalkDirFunc) error {
	return filepath.WalkDir(wd.Abs(), fn)
}

// RemoveAll removes path in WorkDir.
func (wd *WorkDir) RemoveAll(path string) error {
	err := os.RemoveAll(wd.JoinAbs(path))
	if err != nil {
		return err
	}
	return nil
}

// Open wraps os.Open.
func (wd *WorkDir) Open(path string) (*os.File, error) {
	return os.Open(wd.JoinAbs(path))
}

// OpenFile wraps os.OpenFile
func (wd *WorkDir) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(wd.JoinAbs(name), flag, perm)
}

// MakeDir creates new directory in Workdir.
func (wd *WorkDir) MakeDir(name string, perm os.FileMode) error {
	err := os.Mkdir(wd.JoinAbs(name), perm)
	if err != nil {
		return err
	}
	return nil
}

// Abs returns absolute path for WorkDir.
func (wd *WorkDir) Abs() string {
	return filepath.ToSlash(wd.abs)
}

// JoinAbs is a filepath.Join but returns an absolute path.
func (wd *WorkDir) JoinAbs(path ...string) string {
	dir := filepath.Join(append([]string{wd.abs}, path...)...)
	return filepath.ToSlash(dir)
}
