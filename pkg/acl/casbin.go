package acl

import (
	"github.com/casbin/casbin"
)

// 简单封装casbin api

// 授权
func AddPolicy(enforcer *casbin.Enforcer, uid string, resource string, method string) bool {
	return enforcer.AddPolicy(uid, resource, method)
}

// 取消授权
func RemovePolicy(enforcer *casbin.Enforcer, uid string, resource string, method string) bool {
	return enforcer.RemovePolicy(uid, resource, method)
}

// 判断是否具有权限
func Enforce(enforcer *casbin.Enforcer, uid string, resource string, method string) (bool, error) {
	return enforcer.EnforceSafe(uid, resource, method)
}
