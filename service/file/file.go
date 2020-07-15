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

type ChunkService struct {
	ChunkNumber 		int 	`form:"chunkNumber"`
	ChunkSize 			int		`form:"chunkSize"`
	CurrentChunkSize	int 	`form:"currentChunkSize"`
	TotalSize			int64 	`form:"totalSize"`
	Identifier 			string 	`form:"identifier"`
	FileName 			string 	`form:"filename"`
	RelativePath 		string 	`form:"relativePath"`
	TotalChunks 		int		`form:"totalChunks"`
}

type MergeService struct {
	FileName 		string 	`form:"filename"`
	RelativePath 	string 	`form:"relativePath"`
	Identifier 		string 	`form:"identifier"`
	TotalChunks 	int 	`form:"totalChunks"`
}

// 合并分块
func (service *MergeService) MergeChunk(c *gin.Context) serializer.Response {
	userStore, ok := c.Get("user")
	if !ok {
		return serializer.Err(e.CodeCheckLogin, e.ErrGetUser.Error(), e.ErrGetUser)
	}
	user, _ := userStore.(*models.User)
	checkInfo, err := filesystem.GlobalFs.MergeChunk(service.FileName, service.RelativePath, user, service.Identifier, service.TotalChunks)
	if err != nil {
		return serializer.Err(e.CodeErrMerge, err.Error(), err)
	}
	return serializer.Response{
		Code: 0,
		Data: checkInfo,
	}
}

// 检查分块
func (service *ChunkService) CheckChunk(c *gin.Context) serializer.Response {
	userStore, ok := c.Get("user")
	if !ok {
		return serializer.Err(e.CodeCheckLogin, e.ErrGetUser.Error(), e.ErrGetUser)
	}
	user, _ := userStore.(*models.User)
	checkInfo, err := filesystem.GlobalFs.CheckChunk(service.FileName, service.RelativePath, user, service.Identifier, service.TotalSize, service.TotalChunks)
	if err != nil {
		return serializer.Err(e.CodeCheckChunk, err.Error(), err)
	}
	return serializer.Response{
		Code: 0,
		Data: checkInfo,
	}
}

// 上传分块
func (service *ChunkService) UploadChunk(c *gin.Context) serializer.Response {
	userStore, ok := c.Get("user")
	if !ok {
		return serializer.Err(e.CodeCheckLogin, e.ErrGetUser.Error(), e.ErrGetUser)
	}
	user, _ := userStore.(*models.User)
	file, err := c.FormFile("file")
	if err != nil {
		return serializer.Err(e.CodeErrGetUploadChunk, e.ErrGetUploadChunk.Error(), e.ErrGetUploadChunk)
	}
	checkInfo, err := filesystem.GlobalFs.UploadChunk(service.FileName, service.RelativePath, user, service.Identifier, service.TotalSize, service.ChunkNumber, service.TotalChunks, file)
	if err != nil {
		return serializer.Err(e.CodeUploadChunk, err.Error(), err)
	}
	return serializer.Response{
		Code: 0,
		Data: checkInfo,
	}
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



