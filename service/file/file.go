package file

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/models"
	"go-file-manager/pkg/e"
	"go-file-manager/pkg/filesystem"
	"go-file-manager/pkg/serializer"
)

type FileService struct {
	Path string `json:"path" binding:"required,min=1,max=65535"`
}

type ListDirectoryService struct {
	By 	 string `form:"by" json:"by"`
	Asc  bool 	`form:"asc" json:"asc"`
	Path string `form:"path" json:"path"`
}

type CreateDirectoryService struct {
	Path string `form:"path" json:"path"`
	Name string `form:"name" json:"name"`
}

// 列出目录内容
func (service *ListDirectoryService) ListDirectory(c *gin.Context) serializer.Response {
	list, err := filesystem.GlobalFs.List(filesystem.Sorting{By: service.By, Asc: service.Asc}, service.Path)
	if err != nil {
		return serializer.Err(e.CodeNotSet, err.Error(), err)
	}
	return serializer.Response{
		Code: 0,
		Data: list,
	}
}

// 创建新目录
func (service *CreateDirectoryService) CreateDirectory(c *gin.Context) serializer.Response {
	userStore, ok := c.Get("user")
	if !ok {
		return serializer.Err(e.CodeCheckLogin, e.ErrGetUser.Error(), e.ErrGetUser)
	}
	user, _ := userStore.(*models.User)
	fileInfo, err := filesystem.GlobalFs.CreateDirectory(user, service.Name, service.Path)
	if err != nil {
		return serializer.Err(e.CodeCreateFolderFailed, err.Error(), err)
	}
	return serializer.Response{
		Code: 0,
		Data: fileInfo,
	}
}



