package models

import (
	"github.com/jinzhu/gorm"
	"go-file-manager/pkg/conf"
	"go-file-manager/pkg/util"
)

// 执行数据迁移
func migration() {
	util.Log().Info("开始进行数据库初始化...")
	if conf.DatabaseConfig.Type == "mysql" {
		DB = DB.Set("gorm:table_options", "ENGINE=InnoDB")
	}
	DB.AutoMigrate(&User{}, &File{})
	addDefaultUser()

	util.Log().Info("数据库初始化结束")
}

func addDefaultUser() {
	_, err := GetUserByID(1)
	if gorm.IsRecordNotFoundError(err) {
		defaultUser := NewUser()
		defaultUser.Email = "2296176046@qq.com"
		defaultUser.FirstName = "FirstName"
		defaultUser.LastName = "LastName"
		defaultUser.PhoneNumber = "13262285856"
		defaultUser.School = "SJTU"
		defaultUser.Role = "Student"
		defaultUser.Status = true
		defaultUser.UserType = "Admin"
		err := defaultUser.SetPassword("123456")
		if err != nil {
			util.Log().Panic("设置管理员密码失败, %s", err)
		}
		if err := DB.Create(&defaultUser).Error; err != nil {
			util.Log().Panic("无法创建初始用户, %s", err)
		}
	}
}