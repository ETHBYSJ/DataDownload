package bootstrap

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/models"
	"go-file-manager/pkg/acl"
	"go-file-manager/pkg/conf"
	"go-file-manager/pkg/filesystem"
)

func Init(confPath string, aclPath string) {
	// 初始化配置文件
	conf.Init(confPath)
	if !conf.SystemConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	// 初始化数据库配置
	models.Init()
	// 初始化全局虚拟文件系统
	filesystem.Init()
	// 初始化权限访问控制
	acl.Init(aclPath)
}
