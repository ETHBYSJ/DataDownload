package main

import (
	"fmt"
	"go-file-manager/bootstrap"
	"os"
	"path"
	"path/filepath"
)

func init() {
	bootstrap.Init("conf/app.ini")
}

func main() {
	/*
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
	*/
	/*
	api := routers.InitRouter()
	util.Log().Info("开始监听 %s", conf.SystemConfig.Listen)
	if err := api.Run(conf.SystemConfig.Listen); err != nil {
		util.Log().Error("无法监听[%s], %s", conf.SystemConfig.Listen, err)
	}

	 */
	fmt.Println(path.Clean("../a/./b"))
	fmt.Println(filepath.Dir("/a/b/c"))
	fmt.Println(filepath.IsAbs("./a/b/c"))
	fmt.Println(filepath.IsAbs("C://a"))
	fmt.Println(os.TempDir())
}
