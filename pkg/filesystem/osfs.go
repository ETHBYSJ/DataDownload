package filesystem

import (
	"os"
	"time"
)

type OsFs struct{}

func NewOsFs() Fs {
	return &OsFs{}
}

func (o *OsFs) RealPath(name string) (string, error) {
	return name, nil
}

func (o *OsFs) Name() string { return "OsFs" }

func (o *OsFs) Create(name string) (File, error) {
	f, e := os.Create(name)
	if f == nil {
		return nil, e
	}
	return f, e
}

func (o *OsFs) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (o *OsFs) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (o *OsFs) Open(name string) (File, error) {
	f, e := os.Open(name)
	if f == nil {
		return nil, e
	}
	return f, e
}

func (o *OsFs) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	f, e := os.OpenFile(name, flag, perm)
	if f == nil {
		return nil, e
	}
	return f, e
}

func (o *OsFs) Remove(name string) error {
	return os.Remove(name)
}

func (o *OsFs) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (o *OsFs) Rename(oldName, newName string) error {
	return os.Rename(oldName, newName)
}

func (o *OsFs) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (o *OsFs) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

func (o *OsFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(name, atime, mtime)
}

func (o *OsFs) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	fi, err := os.Lstat(name)
	return fi, true, err
}

func (o *OsFs) SymlinkIfPossible(oldname, newname string) error {
	return os.Symlink(oldname, newname)
}

func (o *OsFs) ReadlinkIfPossible(name string) (string, error) {
	return os.Readlink(name)
}
