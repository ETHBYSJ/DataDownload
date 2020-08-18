package e

import (
	"errors"
)

// 文件相关
var (
	ErrIllegalName          = errors.New("目标名称非法")
	ErrFileExisted          = errors.New("同名文件已存在")
	ErrFolderExisted        = errors.New("同名目录已存在")
	ErrFileNotExist         = errors.New("文件不存在")
	ErrPathNotExist         = errors.New("路径不存在")
	ErrClientCancelled      = errors.New("客户端取消操作")
	ErrGetContext           = errors.New("获取上下文错误")
	ErrType                 = errors.New("类型错误")
	ErrFileCover            = errors.New("已存在同名文件")
	ErrCreateFileRecord     = errors.New("创建文件记录失败")
	ErrGetUploadChunk       = errors.New("获取上传文件分块失败")
	ErrUploadPathNotExists  = errors.New("上传路径不存在")
	ErrGetUploadRecord      = errors.New("获取文件记录失败")
	ErrCheckChunks          = errors.New("统计文件分块出错")
	ErrAlreadyMerged        = errors.New("文件分块已经合并")
	ErrMergeFile            = errors.New("合并文件出错")
	ErrMergeFileMD5         = errors.New("MD5与预期不一致")
	ErrLock                 = errors.New("获取锁失败")
	ErrUnlock               = errors.New("解锁失败")
	ErrDelete               = errors.New("删除失败")
	ErrTryToGetRoot         = errors.New("试图获取根目录在数据库中的记录")
	ErrDownloadFileNotExist = errors.New("下载链接过期或文件不存在")
)

// 用户相关
var (
	ErrGetUser         = errors.New("获取用户错误")
	ErrLogin           = errors.New("登录失败")
	ErrRegister        = errors.New("注册失败")
	ErrUserStatus      = errors.New("账号被封禁")
	ErrLanguageInvalid = errors.New("不支持的语言类型")
	ErrLanguageSet     = errors.New("语言设置失败")
)


