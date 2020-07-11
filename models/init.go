package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-file-manager/pkg/conf"
	"go-file-manager/pkg/util"
	"time"
)
// 全局变量
var DB *gorm.DB

func Init() {
	util.Log().Info("初始化数据库连接")
	var (
		db *gorm.DB
		err error
	)
	db, err = gorm.Open(conf.DatabaseConfig.Type, fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.DatabaseConfig.User,
		conf.DatabaseConfig.Password,
		conf.DatabaseConfig.Host,
		conf.DatabaseConfig.Port,
		conf.DatabaseConfig.Name))
	/*
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return conf.DatabaseConfig.TablePrefix + defaultTableName
	}
	*/
	if conf.SystemConfig.Debug {
		db.LogMode(true)
	} else {
		db.LogMode(false)
	}

	if err != nil {
		util.Log().Panic("连接数据库失败, %s", err)
	}

	db.DB().SetMaxIdleConns(50)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Second * 30)

	DB = db

	// 执行迁移
	migration()

}
