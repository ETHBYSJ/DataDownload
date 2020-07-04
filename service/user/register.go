package user

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/models"
	"go-file-manager/pkg/serializer"
)

type UserRegisterService struct {
	UserName string `form:"userName" json:"userName" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=16"`
}

func (service *UserRegisterService) Register(c *gin.Context) serializer.Response {
	user := models.NewUser()
	user.Email = service.UserName
	user.SetPassword(service.Password)
	user.Status = models.Active
	// 创建用户
	if err := models.DB.Create(&user).Error; err != nil {
		return serializer.DBErr("邮箱已被使用", err)
	}
	return serializer.Response{}
}