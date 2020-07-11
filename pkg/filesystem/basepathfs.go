package filesystem

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type BasePathFs struct {
	source Fs
	path string
}

type BasePathFile struct {
	File
	path string
}

func (f *BasePathFile) Name() string {
	sourceName := f.File.Name()
	return strings.TrimPrefix(sourceName, filepath.Clean(f.path))
}

func NewBasePathFs(source Fs, path string) Fs {
	return &BasePathFs{source: source, path: path}
}


func (b *BasePathFs) RealPath(name string) (path string, err error) {
	if err := validateBasePathName(name); err != nil {
		return name, err
	}
	bpath := filepath.Clean(b.path)
	path = filepath.Clean(filepath.Join(bpath, name))
	if !strings.HasPrefix(path, bpath) {
		return name, os.ErrNotExist
	}
	return path, nil
}

func validateBasePathName(name string) error {
	if runtime.GOOS != "windows" {
		return nil
	}
	if filepath.IsAbs(name) {
		return os.ErrNotExist
	}
	return nil
}

func (b *BasePathFs) Chtimes(name string, atime, mtime time.Time) (err error) {
	if name, err = b.RealPath(name); err != nil {
		return &os.PathError{Op: "chtimes", Path: name, Err: err}
	}
	return b.source.Chtimes(name, atime, mtime)
}

func (b *BasePathFs) Chmod(name string, mode os.FileMode) (err error) {
	if name, err = b.RealPath(name); err != nil {
		return &os.PathError{Op: "chmod", Path: name, Err: err}
	}
	return b.source.Chmod(name, mode)
}

func (b *BasePathFs) Name() string {
	return "BasePathFs"
}

func (b *BasePathFs) Stat(name string) (fs os.FileInfo, err error) {
	if name, err = b.RealPath(name); err != nil {
		return nil, &os.PathError{Op: "stat", Path: name, Err: err}
	}
	return b.source.Stat(name)
}

func (b *BasePathFs) Rename(oldName, newName string) (err error) {
	if oldName, err = b.RealPath(oldName); err != nil {
		return &os.PathError{Op: "rename", Path: oldName, Err: err}
	}
	if newName, err = b.RealPath(newName); err != nil {
		return &os.PathError{Op: "rename", Path: newName, Err: err}
	}
	return b.source.Rename(oldName, newName)
}

func (b *BasePathFs) RemoveAll(name string) (err error) {
	if name, err = b.RealPath(name); err != nil {
		return &os.PathError{Op: "remove_all", Path: name, Err: err}
	}
	return b.source.RemoveAll(name)
}

func (b *BasePathFs) Remove(name string) (err error) {
	if name, err = b.RealPath(name); err != nil {
		return &os.PathError{Op: "remove", Path: name, Err: err}
	}
	return b.source.Remove(name)
}

func (b *BasePathFs) OpenFile(name string, flag int, mode os.FileMode) (f File, err error) {
	if name, err = b.RealPath(name); err != nil {
		return nil, &os.PathError{Op: "openfile", Path: name, Err: err}
	}
	sourcef, err := b.source.OpenFile(name, flag, mode)
	if err != nil {
		return nil, err
	}
	return &BasePathFile{sourcef, b.path}, nil
}

func (b *BasePathFs) Open(name string) (f File, err error) {
	if name, err = b.RealPath(name); err != nil {
		return nil, &os.PathError{Op: "open", Path: name, Err: err}
	}
	sourcef, err := b.source.Open(name)
	if err != nil {
		return nil, err
	}
	return &BasePathFile{File: sourcef, path: b.path}, nil
}

func (b *BasePathFs) Mkdir(name string, mode os.FileMode) (err error) {
	if name, err = b.RealPath(name); err != nil {
		return &os.PathError{Op: "mkdir", Path: name, Err: err}
	}
	return b.source.Mkdir(name, mode)
}

func (b *BasePathFs) MkdirAll(name string, mode os.FileMode) (err error) {
	if name, err = b.RealPath(name); err != nil {
		return &os.PathError{Op: "mkdir_all", Path: name, Err: err}
	}
	return b.source.MkdirAll(name, mode)
}

func (b *BasePathFs) Create(name string) (f File, err error) {
	if name, err = b.RealPath(name); err != nil {
		return nil, &os.PathError{Op: "create", Path: name, Err: err}
	}
	sourcef, err := b.source.Create(name)
	if err != nil {
		return nil, err
	}
	return &BasePathFile{File: sourcef, path: b.path}, nil
}

func (b *BasePathFs) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	name, err := b.RealPath(name)
	if err != nil {
		return nil, false, &os.PathError{Op: "lstat", Path: name, Err: err}
	}
	if lstater, ok := b.source.(Lstater); ok {
		return lstater.LstatIfPossible(name)
	}
	fi, err := b.source.Stat(name)
	return fi, false, err
}

func (b *BasePathFs) symlinkIfPossible(oldname, newname string) error {
	oldname, err := b.RealPath(oldname)
	if err != nil {
		return &os.LinkError{Op: "symlink", Old: oldname, New: newname, Err: err}
	}
	newname, err = b.RealPath(newname)
	if err != nil {
		return &os.LinkError{Op: "symlink", Old: oldname, New: newname, Err: err}
	}
	if linker, ok := b.source.(Linker); ok {
		return linker.SymlinkIfPossible(oldname, newname)
	}
	return &os.LinkError{Op: "symlink", Old: oldname, New: newname, Err: ErrNoSymlink}
}

func (b *BasePathFs) ReadlinkIfPossible(name string) (string, error) {
	name, err := b.RealPath(name)
	if err != nil {
		return "", &os.PathError{Op: "readlink", Path: name, Err: err}
	}
	if reader, ok := b.source.(LinkReader); ok {
		return reader.ReadlinkIfPossible(name)
	}
	return "", &os.PathError{Op: "readlink", Path: name, Err: ErrNoReadlink}
}

