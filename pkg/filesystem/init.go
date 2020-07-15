package filesystem

import (
	"go-file-manager/pkg/conf"
	"go-file-manager/pkg/util"
)

// 全局变量
var GlobalFs FileSystem

func Init() {
	GlobalFs.Fs = &BasePathFs{source: NewOsFs(), path: conf.SystemConfig.StorageRoot}
	GlobalFs.locker = util.NewTable()
}
