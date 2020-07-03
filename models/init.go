package models

import (
	// "fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-file-manager/pkg/util"
)
var DB *gorm.DB

func Init() {
	util.Log().Info("初始化数据库连接")
	var (
		db *gorm.DB
		err error
	)


}
