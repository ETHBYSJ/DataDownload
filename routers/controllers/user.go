package controllers

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/pkg/serializer"
	"go-file-manager/service/user"
)

// 我的文件
func MyFiles(c *gin.Context) {
	var service user.MyFilesService
	if err := c.ShouldBind(&service); err == nil {
		res := service.GetMyFiles(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// 用户登录
func UserLogin(c *gin.Context) {
	var service user.UserLoginService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Login(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// 用户请求验证
func UserRequestActivate(c *gin.Context) {
	var service user.UserRequestValidateService
	if err := c.ShouldBind(&service); err == nil {
		res := service.RequestActivate(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// 用户邮箱验证
func UserActivate(c *gin.Context) {
	var service user.UserValidateService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Activate(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// 用户注册
func UserRegister(c *gin.Context) {
	var service user.UserRegisterService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Register(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// 用户注销
func UserLogout(c *gin.Context) {
	var service user.UserLogoutService
	res := service.Logout(c)
	c.JSON(200, res)
}

// 当前登录用户
func UserMe(c *gin.Context) {
	currentUser := user.CurrentUser(c)
	res := serializer.BuildUserResponse(*currentUser)
	c.JSON(200, res)
}

// 设置用户偏好语言

func LanguageSet(c *gin.Context) {
	var service user.UserLanguageService
	if err := c.ShouldBind(&service); err == nil {
		res := service.LanguageSet(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}
