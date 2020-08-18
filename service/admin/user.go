package admin

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/models"
	"go-file-manager/pkg/e"
	"go-file-manager/pkg/serializer"
)

type ListUserService struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"pageSize" json:"pageSize"`
}

type UpdateUserStatusService struct {
	ID     uint `form:"id" json:"id"`
	Status bool `form:"status" json:"status"`
}

type CreateUserService struct {
	Email       string `form:"email" json:"email" binding:"required,email"`
	FirstName   string `form:"firstName" json:"firstName" binding:"required"`
	LastName    string `form:"lastName" json:"lastName" binding:"required"`
	PhoneNumber string `form:"phoneNumber" json:"phoneNumber" binding:"required"`
	School      string `form:"school" json:"school" binding:"required"`
	Role        string `form:"role" json:"role" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required,min=6,max=16"`
	UserType    string `form:"userType" json:"userType" binding:"required"`
}

// 创建新用户
func (service *CreateUserService) CreateUserService(c *gin.Context) serializer.Response {
	user := models.NewUser()
	user.Email = service.Email
	user.FirstName = service.FirstName
	user.LastName = service.LastName
	user.PhoneNumber = service.PhoneNumber
	user.School = service.School
	user.Role = service.Role
	user.SetPassword(service.Password)
	user.Status = true
	user.UserType = service.UserType
	// 创建用户
	if err := models.DB.Create(user).Error; err != nil {
		return serializer.DBErr("邮箱已被使用", e.ErrRegister)
	}
	return serializer.Response{Msg: "创建新用户成功"}
}

// 更新用户状态
func (service *UpdateUserStatusService) UpdateUserStatus(c *gin.Context) serializer.Response {
	// util.Log().Info("user status %v", service.Status)
	user, err := models.GetUserByID(service.ID)
	if err != nil {
		return serializer.Err(e.CodeSetStatusErr, err.Error(), err)
	}
	err = user.SetStatus(service.Status)
	if err != nil {
		return serializer.Err(e.CodeSetStatusErr, err.Error(), err)
	}
	return serializer.Response{
		Code: 0,
	}

}

func (service *ListUserService) ListUser(c *gin.Context) serializer.Response {
	users, count, err := models.ListUser(service.PageSize, service.Page)
	if err != nil {
		return serializer.Err(e.CodeListUserErr, err.Error(), err)
	} else {
		return serializer.Response{
			Code: 0,
			Data: map[string]interface{}{
				"users": users,
				"count": count,
			},
		}
	}
}
