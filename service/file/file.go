package file

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/models"
	"go-file-manager/pkg/acl"
	"go-file-manager/pkg/e"
	"go-file-manager/pkg/filesystem"
	"go-file-manager/pkg/serializer"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
	Path 	string `form:"path" json:"path"`
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

type DownloadService struct {
	Name	string	`form:"name" json:"name"`
	Path 	string 	`form:"path" json:"path"`
}

type DownloadStaticService struct {
	Name	string `form:"name"`
	Path	string `form:"path"`
}

type DownloadNoAuthService struct {
	ID	string	`form:"id"`
}

func (service *DownloadNoAuthService) Download(c *gin.Context) serializer.Response {
	item, exist := filesystem.GlobalFs.CheckDownload(service.ID)
	// fmt.Println(item)
	if !exist {
		return serializer.Err(e.CodeDownloadFileNotExist, e.ErrDownloadFileNotExist.Error(), e.ErrDownloadFileNotExist)
	} else {
		diskFile, err := filesystem.GlobalFs.Fs.Open(filepath.Join(item.Path, item.Name))
		if err != nil {
			return serializer.Err(e.CodeErrDownload, err.Error(), err)
		}
		defer diskFile.Close()
		// 设置头，通知浏览器为下载而不是预览
		c.Header("Content-Disposition", "attachment; filename=\"" + url.QueryEscape(item.Name) + "\"; filename*=utf-8''" + url.QueryEscape(item.Name))
		c.Header("Content-Type", "application/octet-stream")
		http.ServeContent(c.Writer, c.Request, item.Name, time.Now(), diskFile)
		return serializer.Response{
			Code: 0,
		}
	}
}

// 创建文件下载会话
func (service *DownloadStaticService) CreateDownloadSession(c *gin.Context) serializer.Response {
	// 获取url
	downloadURL := filesystem.GlobalFs.GetDownloadURL(service.Name, service.Path, time.Hour)
	return serializer.Response{
		Code: 0,
		Data: downloadURL,
	}
}

// 下载文件，针对静态路径
func (service *DownloadStaticService) Download(c *gin.Context) serializer.Response {
	diskFile, err := filesystem.GlobalFs.Fs.Open(filepath.Join(service.Path, service.Name))
	if err != nil {
		return serializer.Err(e.CodeErrDownload, err.Error(), err)
	}
	defer diskFile.Close()
	// 设置头，通知浏览器为下载而不是预览
	c.Header("Content-Disposition", "attachment; filename=\"" + url.QueryEscape(service.Name) + "\"; filename*=utf-8''" + url.QueryEscape(service.Name))
	c.Header("Content-Type", "application/octet-stream")
	http.ServeContent(c.Writer, c.Request, service.Name, time.Now(), diskFile)

	return serializer.Response{
		Code: 0,
	}
}


// 下载文件
func (service *DownloadService) Download(c *gin.Context) serializer.Response {
	diskFile, err := filesystem.GlobalFs.Fs.Open(filepath.Join(service.Path, service.Name))
	if err != nil {
		return serializer.Err(e.CodeErrDownload, err.Error(), err)
	}

	fm, err := models.GetFileByNameAndPath(service.Name, service.Path)
	if err != nil {
		return serializer.Err(e.CodeErrDownload, err.Error(), err)
	}
	/*
	if !fm.Review {
		return serializer.Err(e.CodeErrUnReviewed, e.ErrFileUnReviewed.Error(), e.ErrFileUnReviewed)
	}
	*/
	defer diskFile.Close()

	// 设置头，通知浏览器为下载而不是预览
	c.Header("Content-Disposition", "attachment; filename=\"" + url.QueryEscape(service.Name) + "\"; filename*=utf-8''" + url.QueryEscape(service.Name))
	c.Header("Content-Type", "application/octet-stream")
	http.ServeContent(c.Writer, c.Request, service.Name, fm.UpdatedAt, diskFile)
	return serializer.Response{
		Code: 0,
	}
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



