package bootstrap

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/models"
	"go-file-manager/pkg/conf"
)

func Init(path string) {
	conf.Init(path)
	if !conf.SystemConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	models.Init()
}
