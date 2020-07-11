package conf

import (
	"github.com/go-ini/ini"
	"github.com/go-playground/validator/v10"
	"go-file-manager/pkg/util"
)

type database struct {
	Type 		string
	User 		string
	Password 	string
	Host 		string
	Name 		string
	TablePrefix string
	Port 		int
}

type system struct {
	Debug 			bool
	Listen			string `validate:"required"`
	SessionSecret 	string
	HashIDSalt		string
	StorageRoot		string
}

var cfg *ini.File

func Init(path string) {
	var err error
	cfg, err = ini.Load(path)
	if err != nil {
		util.Log().Panic("无法解析配置文件 '%s': %s", path, err)
	}

	sections := map[string]interface{}{
		"database": DatabaseConfig,
		"system":	SystemConfig,
	}
	for sectionName, sectionStruct := range sections {
		err = mapSection(sectionName, sectionStruct)
		if err != nil {
			util.Log().Panic("配置文件 %s 分区解析失败: %s", sectionName, err)
		}
	}
}

func mapSection(sectionName string, sectionStruct interface{}) error {
	err := cfg.Section(sectionName).MapTo(sectionStruct)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(sectionStruct)
	if err != nil {
		return err
	}
	return nil
}
