package acl

import (
	"fmt"
	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	"go-file-manager/pkg/conf"
	"go-file-manager/pkg/util"
)

var Enforcer *casbin.Enforcer

// 全局初始化
func Init(path string) {
	var enforcer *casbin.Enforcer
	/*
	/bytes, err := util.ReadAll(path)
	if err != nil {
		util.Log().Panic("读取权限配置失败, %s", err)
	}
	*/
	adapter := gormadapter.NewAdapter(conf.DatabaseConfig.Type, fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		conf.DatabaseConfig.User,
		conf.DatabaseConfig.Password,
		conf.DatabaseConfig.Host,
		conf.DatabaseConfig.Port))
	/*
	casbinModel := string(bytes)
	enforcer, err = casbin.NewEnforcerSafe(casbin.NewModel(casbinModel))
	*/
	enforcer, err := casbin.NewEnforcerSafe(path, adapter)
	if err != nil {
		util.Log().Panic("初始化权限配置失败, %s", err)
	}
	enforcer.AddPolicy("alice", "data1", "read")
	ok, err := enforcer.EnforceSafe("alice", "data1", "read")
	if !ok {
		util.Log().Warning("error: %s", err)
	} else {
		util.Log().Info("success")
	}
	ok, err = enforcer.EnforceSafe("bob", "data1", "read")
	util.Log().Warning("ok = %v, err = %s", ok, err)
	// util.Log().Info("enforce %v", ok)
	// enforcer.Enforce("alice", "data1", "read")
	err = enforcer.SavePolicy()
	if err != nil {
		util.Log().Warning("save policy error, %s", err)
	}
	Enforcer = enforcer
}
