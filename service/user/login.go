package user

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/models"
	"go-file-manager/pkg/serializer"
	"go-file-manager/pkg/util"
)

type UserLoginService struct {
	UserName string `form:"userName" json:"userName" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=16"`
}

func (service *UserLoginService) Login(c *gin.Context) serializer.Response {
	expectedUser, err := models.GetUserByEmail(service.UserName)
	// 验证
	if err != nil {
		return serializer.Err(401, "用户邮箱或密码错误", err)
	}
	if authOK, _ := expectedUser.CheckPassword(service.Password); !authOK {
		return serializer.Err(401, "用户邮箱或密码错误", nil)
	}
	if expectedUser.Status == models.Banned {
		return serializer.Err(403, "账号封禁中", nil)
	}
	util.SetSession(c, map[string]interface{}{
		"user_id": expectedUser.ID,
	})
	return serializer.Response{}
}

