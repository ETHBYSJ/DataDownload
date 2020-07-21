package controllers

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/service/file"
)

// 下载文件
func Download(c *gin.Context) {
	var service file.DownloadService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Download(c)
		if res.Code != 0 {
			c.JSON(200, res)
		}
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// 设置分享
func SetShare(c *gin.Context) {
	var service file.ShareService
	if err := c.ShouldBind(&service); err == nil {
		res := service.SetShare(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}


// 重命名
func Rename(c *gin.Context) {
	var service file.RenameService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Rename(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// 合并分块
func MergeChunk(c *gin.Context) {
	var service file.MergeService
	if err := c.ShouldBind(&service); err == nil {
		res := service.MergeChunk(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// 检查分块
func CheckChunk(c *gin.Context) {
	var service file.ChunkService
	if err := c.BindQuery(&service); err == nil {
		res := service.CheckChunk(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// 上传分块
func UploadChunk(c *gin.Context) {
	var service file.ChunkService
	if err := c.ShouldBind(&service); err == nil {
		res := service.UploadChunk(c)
		c.JSON(200, res)

	} else {
		// util.Log().Info("bind error, err = %s", err)
		c.JSON(200, ErrorResponse(err))
	}
}



// 删除文件
func Delete(c *gin.Context) {
	var service file.DeleteService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Delete(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// 创建新目录
func CreateDirectory(c *gin.Context) {
	var service file.CreateDirectoryService
	if err := c.BindQuery(&service); err == nil {
		res := service.CreateDirectory(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// 根据关键字查询
func ListByKeyword(c *gin.Context) {
	var service file.ListByKeywordService
	if err := c.ShouldBind(&service); err == nil {
		res := service.ListByKeyword(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// 列出目录内容
func ListDirectory(c *gin.Context) {
	var service file.ListDirectoryService
	if err := c.BindQuery(&service); err == nil {
		res := service.ListDirectory(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

func TestController(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 200,
		"success": true,
	})
}