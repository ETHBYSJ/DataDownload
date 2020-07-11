package user

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/models"
	"go-file-manager/pkg/e"
	"go-file-manager/pkg/serializer"
	"go-file-manager/pkg/util"
)

type UserRegisterService struct {
	UserName string `form:"userName" json:"userName" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=16"`
}

func (service *UserRegisterService) Register(c *gin.Context) serializer.Response {
	user := models.NewUser()
	user.Email = service.UserName
	user.SetPassword(service.Password)
	user.Status = true
	// 创建用户
	if err := models.DB.Create(&user).Error; err != nil {
		return serializer.DBErr("邮箱已被使用", e.ErrRegister)
	}
	util.Log().Info("新用户id: %v", user.ID)
	util.SetSession(c, map[string]interface{}{
		"user_id": user.ID,
	})
	return serializer.Response{Msg: "注册成功"}
}