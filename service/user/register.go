package user

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/models"
	"go-file-manager/pkg/e"
	"go-file-manager/pkg/mail"
	"go-file-manager/pkg/serializer"
	"go-file-manager/pkg/util"
	"time"
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

type UserValidateService struct {
	UserId uint	`form:"userid" json:"userid" binding:"required"`
	Code string `form:"code" json:"code" binding:"required"`
}

type UserRequestValidateService struct {
	Email string `form:"email" json:"email" binding:"required"`
}

func (service *UserRequestValidateService) RequestValidate(c *gin.Context) serializer.Response {
	expectedUser, err := models.GetUserByEmail(service.Email)
	// 验证
	if err != nil {
		return serializer.Err(e.CodeCheckLogin, "用户邮箱或密码错误", e.ErrLogin)
	}
	// 生成6位随机验证码
	code := util.RandStringRunes(6)
	// 发送邮件
	mail.SendMail([]string{service.Email}, "validate", code)
	// 在缓存中生存期设置为10分钟
	models.Cache.Add(expectedUser.ID, time.Minute * 60, code)
	return serializer.Response{
		Data: map[string]interface{}{
			"userId": expectedUser.ID,
			"code": code,
		},
	}
}

func (service *UserValidateService) Validate(c *gin.Context) serializer.Response {
	item, _ := models.Cache.Value(service.UserId)
	if item != nil {
		code := item.Data().(string)
		if code == service.Code {
			user, _ := models.GetUserByID(service.UserId)
			user.SetValidate()
			// 验证通过,从缓存中删除
			models.Cache.Delete(service.UserId)
			util.SetSession(c, map[string]interface{}{
				"user_id": service.UserId,
			})
			return serializer.Response{
				Code: 0,
			}
		} else {
			return serializer.Response{
				Code: e.CodeValidateError,
			}
		}
	} else {
		return serializer.Response{
			Code: e.CodeValidateError,
		}
	}
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
	user.Validate = false
	// 创建用户
	if err := models.DB.Create(user).Error; err != nil {
		return serializer.DBErr("邮箱已被使用", e.ErrRegister)
	}
	util.Log().Info("新用户id: %v", user.ID)
	/*
	util.SetSession(c, map[string]interface{}{
		"user_id": user.ID,
	})
	*/
	// 生成6位随机验证码
	code := util.RandStringRunes(6)
	// 发送邮件
	mail.SendMail([]string{user.Email}, "validate", code)
	// 在缓存中生存期设置为10分钟
	models.Cache.Add(user.ID, time.Minute * 60, code)
	return serializer.Response{
		Msg: "注册成功",
		Data: map[string]interface{}{
			"userId": user.ID,
			"code": code,
		},
	}
}
