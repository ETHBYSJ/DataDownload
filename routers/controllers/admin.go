package controllers

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/service/admin"
)

// 创建新用户
func CreateUser(c *gin.Context) {
	var service admin.CreateUserService
	if err := c.ShouldBind(&service); err == nil {
		res := service.CreateUserService(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// 更新用户状态
func UpdateUserStatus(c *gin.Context) {
	var service admin.UpdateUserStatusService
	if err := c.ShouldBind(&service); err == nil {
		res := service.UpdateUserStatus(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// 列出用户
func ListUser(c *gin.Context) {
	var service admin.ListUserService
	if err := c.ShouldBind(&service); err == nil {
		res := service.ListUser(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}
