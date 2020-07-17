package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-file-manager/models"
	"go-file-manager/pkg/serializer"
)

func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		uid := session.Get("user_id")
		if uid != nil {
			user, err := models.GetActiveUserByID(uid)
			if err == nil {
				// util.Log().Info("user ID = %v", user.ID)
				c.Set("user", user)
			}
		} else {
			// 此处避免上一次的残留
			c.Set("user", nil)
		}
		c.Next()
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get("user"); user != nil {
			if _, ok := user.(*models.User); ok {
				c.Next()
				return
			}
		}
		c.JSON(200, serializer.CheckLogin())
		c.Abort()
	}
}

