package user

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/models"
	"go-file-manager/pkg/e"
	"go-file-manager/pkg/serializer"
	"go-file-manager/pkg/util"
)


func CurrentUser(c *gin.Context) *models.User {
	if user, _ := c.Get("user"); user != nil {
		if u, ok := user.(*models.User); ok {
			return u
		}
	}
	return nil
}

type UserLanguageService struct {
	Language string `form:"language" json:"language" binding:"required"`
}

type UserLogoutService struct {

}

type MyFilesService struct {
	UID 		uint 	`form:"uid" json:"uid"`
	Page 		int 	`form:"page" json:"page"`
	PageSize	int		`form:"pageSize" json:"pageSize"`
	Keyword 	string 	`form:"keyword" json:"keyword"`
	Category 	int 	`form:"category" json:"category"`
}

// 我的文件
func (service *MyFilesService) GetMyFiles(c *gin.Context) serializer.Response {
	files, count, err := models.GetMyFiles(service.UID, service.Page, service.PageSize, service.Keyword, service.Category)
	if err != nil {
		return serializer.Err(e.CodeErrGetMyFiles, err.Error(), err)
	}
	return serializer.Response{
		Code: 0,
		Data: map[string]interface{}{
			"files": files,
			"count": count,
		},
	}
}

func (service *UserLanguageService) LanguageSet(c *gin.Context) serializer.Response {
	user := CurrentUser(c)
	if user == nil {
		return serializer.Err(e.CodeCheckLogin, "请登录", e.ErrGetUser)
	}
	language := service.Language
	err := models.LanguageSet(user, language)
	if err != nil {
		return serializer.Err(e.CodeLanguageSet, "语言设置失败", err)
	}
	return serializer.Response{Msg: "语言设置成功"}
}

func (service *UserLogoutService) Logout(c *gin.Context) serializer.Response {
	user := CurrentUser(c)
	if user == nil {
		return serializer.Err(e.CodeCheckLogin, "请登录", e.ErrGetUser)
	}
	// 清除session
	util.DeleteSession(c, "user_id")
	// 清除context中存储的用户
	c.Set("user", nil)
	return serializer.Response{Msg: "注销成功"}
}