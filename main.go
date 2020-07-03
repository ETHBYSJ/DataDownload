package main

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/bootstrap"
)

func init() {
	bootstrap.Init("conf/app.ini")
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
