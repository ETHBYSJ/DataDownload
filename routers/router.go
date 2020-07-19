package routers

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/middleware"
	"go-file-manager/pkg/conf"
	"go-file-manager/routers/controllers"
)



func InitRouter() *gin.Engine {
	r := gin.Default()
	r.MaxMultipartMemory = 1024 << 20
	v1 := r.Group("/api/v1")
	v1.Use(middleware.Session(conf.SystemConfig.SessionSecret))
	v1.Use(middleware.CurrentUser())
	// v1.Use(middleware.Cors())
	{
		user := v1.Group("user")
		{
			user.POST("register", controllers.UserRegister)
			user.POST("login", controllers.UserLogin)
			user.GET("logout", controllers.UserLogout)
		}
		// 需要管理员权限才能访问
		admin := v1.Group("admin")
		admin.Use(middleware.AdminRequired())
		{
			admin.GET("list_user", controllers.ListUser)
			admin.GET("update_status", controllers.UpdateUserStatus)
			admin.POST("create_user", controllers.CreateUser)
		}
		// 需要登录才能访问
		auth := v1.Group("")
		auth.Use(middleware.AuthRequired())
		{
			user := auth.Group("user")
			{
				user.GET("me", controllers.UserMe)
				user.GET("profile", controllers.UserMe)
				user.GET("language", controllers.LanguageSet)
			}
			// file := auth.Group("file")
			file := auth.Group("file")
			file.Use(middleware.CasbinMiddleware())
			{
				file.GET("test", controllers.TestController)
				file.GET("list", controllers.ListDirectory)
				file.PUT("create", controllers.CreateDirectory)
				file.GET("rename", controllers.Rename)
				file.GET("delete", controllers.Delete)
				file.GET("list_by_keyword", controllers.ListByKeyword)
				file.GET("set_share", controllers.SetShare)
				// 分块上传相关
				file.GET("chunk", controllers.CheckChunk)
				file.POST("chunk", controllers.UploadChunk)
				file.GET("merge", controllers.MergeChunk)
			}
		}

	}
	return r
}
