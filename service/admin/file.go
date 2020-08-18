package admin

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/models"
	"go-file-manager/pkg/e"
	"go-file-manager/pkg/serializer"
)

type AdminGetFilesService struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"pageSize" json:"pageSize"`
	Keyword  string `form:"keyword" json:"keyword"`
	Category int    `form:"category" json:"category"`
}

type AdminReviewService struct {
	ID     uint `form:"id" json:"id"`
	Review bool `form:"review" json:"review"`
}

// 修改审核状态
func (service *AdminReviewService) SetReview(c *gin.Context) serializer.Response {
	file, err := models.GetFileByID(service.ID)
	if err != nil {
		return serializer.Err(e.CodeReviewErr, err.Error(), err)
	}
	err = file.UpdateReview(service.Review)
	if err != nil {
		return serializer.Err(e.CodeReviewErr, err.Error(), err)
	}
	return serializer.Response{
		Code: 0,
	}
}

// 获取文件列表
func (service *AdminGetFilesService) AdminGetFiles(c *gin.Context) serializer.Response {
	files, count, err := models.AdminGetFiles(service.Page, service.PageSize, service.Keyword, service.Category)
	if err != nil {
		return serializer.Err(e.CodeGetFilesErr, err.Error(), err)
	}
	return serializer.Response{
		Code: 0,
		Data: map[string]interface{}{
			"files": files,
			"count": count,
		},
	}
}
