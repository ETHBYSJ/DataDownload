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
	adapter := gormadapter.NewAdapter(conf.DatabaseConfig.Type, fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		conf.DatabaseConfig.User,
		conf.DatabaseConfig.Password,
		conf.DatabaseConfig.Host,
		conf.DatabaseConfig.Port))
	enforcer, err := casbin.NewEnforcerSafe(path, adapter)
	if err != nil {
		util.Log().Panic("初始化权限配置失败, %s", err)
	}
	/*
		enforcer.AddPolicy("alice", "/foo", "GET")
		ok, err := enforcer.EnforceSafe("alice", "/foo", "GET")
		if !ok {
			util.Log().Warning("error: %s", err)
		} else {
			util.Log().Info("success")
		}
		enforcer.RemovePolicy("alice", "/foo", "GET")
		ok, err = enforcer.EnforceSafe("alice", "/foo", "GET")
		if !ok {
			util.Log().Warning("error: %s", err)
		} else {
			util.Log().Info("success")
		}
		err = enforcer.SavePolicy()
		if err != nil {
			util.Log().Warning("save policy error, %s", err)
		}
	*/
	Enforcer = enforcer
}
