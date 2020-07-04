package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"go-file-manager/pkg/e"
	"go-file-manager/pkg/serializer"
	"go-file-manager/pkg/util"
)

var Store memstore.Store

func Session(secret string) gin.HandlerFunc {
	Store = memstore.NewStore([]byte(secret))
	Store.Options(sessions.Options{HttpOnly: true, MaxAge: 7 * 86400, Path: "/"})
	return sessions.Sessions("session", Store)
}

func CSRFInit() gin.HandlerFunc {
	return func(c *gin.Context) {
		util.SetSession(c, map[string]interface{}{"CSRF":true})
		c.Next()
	}
}

func CSRFCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if check, ok := util.GetSession(c, "CSRF").(bool); ok && check {
			c.Next()
			return
		}
		c.JSON(200, serializer.Err(e.CodeNoPermissionErr, "来源非法", nil))
		c.Abort()
	}
}

