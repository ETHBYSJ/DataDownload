package file

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/models"
	"go-file-manager/pkg/acl"
	"go-file-manager/pkg/e"
	"go-file-manager/pkg/filesystem"
	"go-file-manager/pkg/serializer"
	"path/filepath"
	"strconv"
	"strings"
)

type FileService struct {
	Path string `json:"path" binding:"required,min=1,max=65535"`
}

type ShareService struct {
	Path 	string 	`form:"path" json:"path"`
	Name 	string 	`form:"name" json:"name"`
	Share 	bool 	`form:"share" json:"share"`
}

type ListByKeywordService struct {
	By 	 	string 	`form:"by" json:"by"`
	Asc  	bool 	`form:"asc" json:"asc"`
	Path 	string 	`form:"path" json:"path"`
	Keyword string 	`form:"keyword" json:"keyword"`
}

type ListDirectoryService struct {
	By 	 string `form:"by" json:"by"`
	Asc  bool 	`form:"asc" json:"asc"`
	Path string `form:"path" json:"path"`
}

type DeleteService struct {
	Path string `form:"path" json:"path"`
	Name string `form:"name" json:"name"`
}

type CreateDirectoryService struct {
	Path string `form:"path" json:"path"`
	Name string `form:"name" json:"name"`
}

type RenameService struct {
	Path string `form:"path" json:"path"`
	OldName string `form:"oldName" json:"oldName"`
	NewName string `form:"newName" json:"newName"`
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


// 设置分享
func (service *ShareService) SetShare(c *gin.Context) serializer.Response {
	err := filesystem.GlobalFs.SetShare(service.Name, service.Path, service.Share)
	if err != nil {
		return serializer.Err(e.CodeErrSetShare, err.Error(), err)
	}
	return serializer.Response{
		Code: 0,
	}
}

// 删除文件
func (service *DeleteService) Delete(c *gin.Context) serializer.Response {
	userStore, ok := c.Get("user")
	if !ok {
		return serializer.Err(e.CodeCheckLogin, e.ErrGetUser.Error(), e.ErrGetUser)
	}
	user, _ := userStore.(*models.User)

	err := filesystem.GlobalFs.Delete(service.Name, service.Path)
	if err != nil {
		return serializer.Err(e.CodeErrDelete, err.Error(), err)
	}
	// 从权限表中移除
	acl.RemovePolicy(acl.Enforcer, strconv.Itoa(int(user.ID)), strings.ReplaceAll(filepath.Join(service.Path, service.Name), "\\", "/"), "ALL")
	return serializer.Response{
		Code: 0,
	}
}

// 文件重命名
func (service *RenameService) Rename(c *gin.Context) serializer.Response {
	userStore, ok := c.Get("user")
	if !ok {
		return serializer.Err(e.CodeCheckLogin, e.ErrGetUser.Error(), e.ErrGetUser)
	}
	user, _ := userStore.(*models.User)

	err := filesystem.GlobalFs.RenameAtomic(service.OldName, service.NewName, service.Path)
	if err != nil {
		return serializer.Err(e.CodeErrRename, err.Error(), err)
	}
	// 从权限表中移除
	acl.RemovePolicy(acl.Enforcer, strconv.Itoa(int(user.ID)), strings.ReplaceAll(filepath.Join(service.Path, service.OldName), "\\", "/"), "ALL")
	acl.AddPolicy(acl.Enforcer, strconv.Itoa(int(user.ID)), strings.ReplaceAll(filepath.Join(service.Path, service.NewName), "\\", "/"), "ALL")
	return serializer.Response{
		Code: 0,
	}
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
	// 添加进权限表中
	acl.AddPolicy(acl.Enforcer, strconv.Itoa(int(user.ID)), strings.ReplaceAll(filepath.Join(service.RelativePath, service.FileName), "\\", "/"), "ALL")
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

// 根据关键字查询
func (service *ListByKeywordService) ListByKeyword(c *gin.Context) serializer.Response {
	list, err := filesystem.GlobalFs.ListByKeyword(filesystem.Sorting{By: service.By, Asc: service.Asc}, service.Path, service.Keyword)
	if err != nil {
		return serializer.Err(e.CodeNotSet, err.Error(), err)
	}
	return serializer.Response{
		Code: 0,
		Data: list,
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
	// 添加进权限表中
	acl.AddPolicy(acl.Enforcer, strconv.Itoa(int(user.ID)), strings.ReplaceAll(filepath.Join(service.Path, service.Name), "\\", "/"), "ALL")
	return serializer.Response{
		Code: 0,
		Data: fileInfo,
	}
}



