package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-file-manager/models"
	"go-file-manager/pkg/serializer"
)

func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userStore, ok := c.Get("user")
		if !ok {
			// 登录状态错误
			c.JSON(200, serializer.CheckLogin())
			c.Abort()
			return
		}
		user, _ := userStore.(*models.User)
		// 用户ID
		uid := user.ID
		p := c.Request.URL.Path
		m := c.Request.Method
		fmt.Printf("intercept request. uid = %v, path = %v, method = %v\n", uid, p, m)
		c.Next()
	}
}