package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// cors配置
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			// c.AbortWithStatus(http.StatusNoContent)
			c.JSON(http.StatusOK, "ok")
		}
		// 处理请求
		c.Next()
		// fmt.Print("cors")
	}
}
