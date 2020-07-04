package controllers

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/pkg/serializer"
	"go-file-manager/service/user"
)

// 用户登录
func UserLogin(c *gin.Context) {
	var service user.UserLoginService
	if err := c.ShouldBindJSON(&service); err == nil {
		res := service.Login(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}
// 用户注册
func UserRegister(c *gin.Context) {
	var service user.UserRegisterService
	if err := c.ShouldBindJSON(&service); err == nil {
		res := service.Register(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// 当前登录用户
func UserMe(c *gin.Context) {
	currentUser := CurrentUser(c)
	res := serializer.BuildUserResponse(*currentUser)
	c.JSON(200, res)
}



