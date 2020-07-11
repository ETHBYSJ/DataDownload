package filesystem

import (
	"context"
)

// 钩子函数
type Hook func(ctx context.Context, fs *FileSystem) error

// 注入钩子函数
func (fs *FileSystem) Use(name string, hook Hook) {
	if fs.Hooks == nil {
		fs.Hooks = make(map[string][]Hook)
	}
	if _, ok := fs.Hooks[name]; ok {
		fs.Hooks[name] = append(fs.Hooks[name], hook)
		return
	}
	fs.Hooks[name] = []Hook{hook}
}

// 清空钩子，name为空表示全部清空
func (fs *FileSystem) CleanHooks(name string) {
	if name == "" {
		fs.Hooks = nil
	} else {
		delete(fs.Hooks, name)
	}
}