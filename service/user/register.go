package user

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/models"
	"go-file-manager/pkg/e"
	"go-file-manager/pkg/serializer"
	"go-file-manager/pkg/util"
)

type UserRegisterService struct {
	Email       string `form:"email" json:"email" binding:"required,email"`
	FirstName   string `form:"firstName" json:"firstName" binding:"required"`
	LastName    string `form:"lastName" json:"lastName" binding:"required"`
	PhoneNumber string `form:"phoneNumber" json:"phoneNumber" binding:"required"`
	School      string `form:"school" json:"school" binding:"required"`
	Role        string `form:"role" json:"role" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required,min=6,max=16"`
}

func (service *UserRegisterService) Register(c *gin.Context) serializer.Response {
	user := models.NewUser()
	user.Email = service.Email
	user.FirstName = service.FirstName
	user.LastName = service.LastName
	user.PhoneNumber = service.PhoneNumber
	user.School = service.School
	user.Role = service.Role
	user.SetPassword(service.Password)
	user.Status = true
	// 创建用户
	if err := models.DB.Create(user).Error; err != nil {
		return serializer.DBErr("邮箱已被使用", e.ErrRegister)
	}
	util.Log().Info("新用户id: %v", user.ID)
	util.SetSession(c, map[string]interface{}{
		"user_id": user.ID,
	})
	return serializer.Response{Msg: "注册成功"}
}
