package user

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/models"
	"go-file-manager/pkg/e"
	"go-file-manager/pkg/serializer"
	"go-file-manager/pkg/util"
)

type UserLoginService struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=16"`
}

func (service *UserLoginService) Login(c *gin.Context) serializer.Response {
	expectedUser, err := models.GetUserByEmail(service.Email)
	// 验证
	if err != nil {
		return serializer.Err(e.CodeCheckLogin, "用户邮箱或密码错误", e.ErrLogin)
	}
	if !expectedUser.Validate {
		return serializer.Err(e.CodeNoPermissionErr, "账号未激活", e.ErrUserStatus)
	}
	if authOK, _ := expectedUser.CheckPassword(service.Password); !authOK {
		return serializer.Err(e.CodeCheckLogin, "用户邮箱或密码错误", e.ErrLogin)
	}
	if !expectedUser.Status {
		return serializer.Err(e.CodeNoPermissionErr, "账号封禁中", e.ErrUserStatus)
	}
	// util.DeleteSession(c, "user_id")
	util.SetSession(c, map[string]interface{}{
		"user_id": expectedUser.ID,
	})
	return serializer.Response{Msg: "登录成功"}
}
