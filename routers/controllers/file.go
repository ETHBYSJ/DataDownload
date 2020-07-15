package controllers

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/service/file"
)


// 文件上传
/*
func FileUploadStream(c *gin.Context) {
	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 取得文件大小
	fileSize, err := strconv.ParseUint(c.Request.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		c.JSON(200, ErrorResponse(err))
		return
	}
	// 解析信息

}
*/

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