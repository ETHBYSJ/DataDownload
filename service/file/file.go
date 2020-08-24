package file

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-file-manager/models"
	"go-file-manager/pkg/acl"
	"go-file-manager/pkg/conf"
	"go-file-manager/pkg/e"
	"go-file-manager/pkg/exec"
	"go-file-manager/pkg/filesystem"
	"go-file-manager/pkg/mail"
	"go-file-manager/pkg/serializer"
	"go-file-manager/pkg/util"
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
	Path  string `form:"path" json:"path"`
	Name  string `form:"name" json:"name"`
	Share bool   `form:"share" json:"share"`
}

type ListByKeywordService struct {
	By      string `form:"by" json:"by"`
	Asc     bool   `form:"asc" json:"asc"`
	Path    string `form:"path" json:"path"`
	Keyword string `form:"keyword" json:"keyword"`
}

type ListDirectoryService struct {
	By   string `form:"by" json:"by"`
	Asc  bool   `form:"asc" json:"asc"`
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
	Path    string `form:"path" json:"path"`
	OldName string `form:"oldName" json:"oldName"`
	NewName string `form:"newName" json:"newName"`
}

type ChunkService struct {
	ChunkNumber      int    `form:"chunkNumber"`
	ChunkSize        int    `form:"chunkSize"`
	CurrentChunkSize int    `form:"currentChunkSize"`
	TotalSize        int64  `form:"totalSize"`
	Identifier       string `form:"identifier"`
	FileName         string `form:"filename"`
	RelativePath     string `form:"relativePath"`
	TotalChunks      int    `form:"totalChunks"`
}

type MergeService struct {
	FileName     string `form:"filename"`
	RelativePath string `form:"relativePath"`
	Identifier   string `form:"identifier"`
	TotalChunks  int    `form:"totalChunks"`
}

type DownloadService struct {
	Name string `form:"name" json:"name"`
	Path string `form:"path" json:"path"`
}

type DownloadStaticService struct {
	Name string `form:"name"`
	Path string `form:"path"`
}

type DownloadSessionService struct {
	// Name string `form:"name"`
	// 图片文件夹路径
	Path string `form:"path"`
	Email string `form:"email"`
}

type DownloadNoAuthService struct {
	ID string `form:"id"`
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
		// defer diskFile.Close()
		// 设置头，通知浏览器为下载而不是预览
		c.Header("Content-Disposition", "attachment; filename=\""+url.QueryEscape(item.Name)+"\"; filename*=utf-8''"+url.QueryEscape(item.Name))
		c.Header("Content-Type", "application/octet-stream")
		err = filesystem.ServeContent(c.Writer, c.Request, item.Name, time.Now(), diskFile)
		diskFile.Close()
		if err != nil {
			if filesystem.IsUnexpectedError(err) {
				// 发生了未知错误，说明数据很可能已经损坏。删除临时文件和数据库记录
				util.Log().Error("error occurred during download. %v", err)
				filesystem.GlobalFs.Fs.Remove(filepath.Join(item.Path, item.Name))
				models.DeleteRecordByID(service.ID)
				return serializer.Err(e.CodeErrDownload, e.ErrDownloadFileNotExist.Error(), e.ErrDownloadFileNotExist)
			} else {
				// 客户端主动断开连接，更新数据库信息
				record, err := models.GetRecordByID(service.ID)
				if err == nil {
					record.SetStart(true)
				}
			}
		} else {
			// 下载完成，删除临时文件和数据库记录
			// 删除临时文件
			util.Log().Info("下载成功，删除临时文件, %v, %v", item.Path, item.Name)
			filesystem.GlobalFs.Fs.Remove(filepath.Join(item.Path, item.Name))
			// 删除数据库记录
			models.DeleteRecordByID(service.ID)
		}
		return serializer.Response{
			Code: 0,
		}
	}
}

// 创建文件下载会话
func (service *DownloadSessionService) CreateDownloadSession(c *gin.Context) serializer.Response {
	// 查看是否已有正在下载的记录
	/*
	user, err := models.GetUserByEmail(service.Email)
	if err != nil {
		util.Log().Error("获取用户失败, %v", err)
		return serializer.Response{
			Code: e.CodeErrDownloadGetUser,
			Error: e.ErrGetUser.Error(),
			// Data: downloadURL,
		}
	}
	*/
	record, err := models.GetRecordByEmail(service.Email)

	if err == nil {
		// 已存在下载记录，返回旧链接
		go func() {
			base := conf.SystemConfig.Host + conf.SystemConfig.Out
			uri := fmt.Sprintf("/api/v1/file/download_noauth?id=%s", record.ID)
			url := base + uri
			mail.SendMail([]string{service.Email}, "download link", url)
		}()
		return serializer.Response{
			Code: 0,
		}
	}
	// 生成临时文件及数据库记录，并发送给目标用户
	randName := util.RandStringRunes(16)
	go func() {
		// 执行水印算法
		err := exec.ExecPython(conf.SystemConfig.Script, conf.SystemConfig.ImageDir, randName)
		if err != nil {
			util.Log().Error("execute script error, %v", err)
			return
		}
		// 生成数据库记录
		record := models.NewDownloadRecord()
		record.Email = service.Email
		record.Start = false
		record.ID = randName
		if err = models.DB.Create(&record).Error; err != nil {
			util.Log().Error("创建下载记录失败, %v", err)
			return
		}
		// 获取url
		downloadURL := filesystem.GlobalFs.GetDownloadURL(randName, service.Path, time.Hour * 120)
		mail.SendMail([]string{service.Email}, "数据分享平台数据下载","这是您的下载链接："+ downloadURL+"。\n请在5天内完成下载，如果您未申请下载，请忽略本邮件。\n系统邮件，请勿回复。")
		time.AfterFunc(time.Hour * 120, func() {
			util.Log().Info("删除临时文件")
			// 删除临时文件
			filesystem.GlobalFs.Fs.Remove(filepath.Join(service.Path, randName + ".zip"))
			// 删除数据库记录
			models.DeleteRecordByID(randName)
		})
	}()
	return serializer.Response{
		Code: 0,
		// Data: downloadURL,
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
	c.Header("Content-Disposition", "attachment; filename=\""+url.QueryEscape(service.Name)+"\"; filename*=utf-8''"+url.QueryEscape(service.Name))
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
	c.Header("Content-Disposition", "attachment; filename=\""+url.QueryEscape(service.Name)+"\"; filename*=utf-8''"+url.QueryEscape(service.Name))
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
