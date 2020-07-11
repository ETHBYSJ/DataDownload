package e

import (
	"errors"
)

// 文件相关
var (
	ErrIllegalName = errors.New("目标名称非法")
	ErrFileExisted = errors.New("同名文件已存在")
	ErrFolderExisted = errors.New("同名目录已存在")
	ErrFileNotExist = errors.New("文件不存在")
	ErrPathNotExist = errors.New("路径不存在")
	ErrClientCancelled = errors.New("客户端取消操作")
	ErrGetContext = errors.New("获取上下文错误")
	ErrType = errors.New("类型错误")
)
// 用户相关
var (
	ErrGetUser = errors.New("获取用户错误")
	ErrLogin = errors.New("登录失败")
	ErrRegister = errors.New("注册失败")
	ErrUserStatus = errors.New("账号被封禁")
	ErrLanguageInvalid = errors.New("不支持的语言类型")
	ErrLanguageSet = errors.New("语言设置失败")
)