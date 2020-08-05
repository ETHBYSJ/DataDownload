package middleware

import (
	"github.com/gin-gonic/gin"
	"go-file-manager/models"
	"go-file-manager/pkg/acl"
	"go-file-manager/pkg/serializer"
	"go-file-manager/pkg/util"
	"path/filepath"
	"strconv"
	"strings"
)

// 访问权限控制中间件
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
		util.Log().Info("intercept request. uid = %v, path = %v, method = %v\n", uid, p, m)
		// 检查请求权限
		// list和list_by_keyword不需要所有者权限，可以直接放行，Admin用户放行
		// 此外，Admin用户请求下载时无视文件审核情况
		if user.UserType == "Admin" /*|| strings.Contains(p, "list") || strings.Contains(p, "list_by_keyword")*/ {
			c.Next()
		} else {
			// 根据各方法分别判断
			/*
			if strings.Contains(p, "download") {
				// download
				var service file.DownloadService
			 	data, err := c.GetRawData()
			 	err = json.Unmarshal(data, &service)
			 	if err != nil {
			 		c.JSON(200, controllers.ErrorResponse(err))
			 		c.Abort()
			 		return
				}
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
				// 在非共享状态下，禁止下载
				fm, err := models.GetFileByNameAndPath(service.Name, service.Path)
				if err != nil {
					// 404
					c.JSON(e.CodeNotFound, serializer.NotFound())
					c.Abort()
					return
				}
				if !fm.Review {
					c.JSON(e.CodeNoPermissionErr, serializer.PermissionDenied())
					c.Abort()
					return
				}
				// util.Log().Info("owner: %v, share: %v", fm.OwnerID, fm.Share)
				if fm.OwnerID == uid || fm.Share {
					// 如果是所有者或者已经状态设为共享，则允许访问
					c.Next()
					return
				}
				// 否则返回403
				c.JSON(e.CodeNoPermissionErr, serializer.PermissionDenied())
				c.Abort()
				return
			}
			*/if strings.Contains(p, "set_share") {
				// set_share
				path, exists := c.GetQuery("path")
				if !exists {
					c.Next()
					return
				}
				name, exists := c.GetQuery("name")
				if !exists {
					c.Next()
					return
				}
				// 必须是文件的所有者才能修改共享状态
				auth, err := acl.Enforce(acl.Enforcer, strconv.Itoa(int(uid)), strings.ReplaceAll(filepath.Join(path, name), "\\", "/"), "ALL")
				if !auth || err != nil {
					c.JSON(200, serializer.PermissionDenied())
					c.Abort()
					return
				}

			} else if strings.Contains(p, "list") {
				// list
				path, exists := c.GetQuery("path")
				if !exists {
					c.Next()
					return
				}
				if path == "/" {
					// 允许各用户查看根目录下的文件
					c.Next()
					return
				}
				// 在非共享状态下，无法查看
				fm, err := models.GetFileByPath(path)
				if err != nil {
					// 404
					c.JSON(200, serializer.NotFound())
					c.Abort()
					return
				}
				if fm.OwnerID == uid || fm.Share {
					// 如果是所有者或者已经状态设为共享，则允许访问
					c.Next()
					return
				}
				c.JSON(200, serializer.PermissionDenied())
				c.Abort()
				return

			} else if strings.Contains(p, "list_by_keyword") {
				// list_by_keyword
				path, exists := c.GetQuery("path")
				if !exists {
					c.Next()
					return
				}
				if path == "/" {
					// 允许各用户查看根目录下的文件
					c.Next()
					return
				}
				// 在非共享状态下，无法查看

			} else if strings.Contains(p, "create") {
				// create
				path, exists := c.GetQuery("path")
				if !exists {
					c.Next()
					return
				}
				if path == "/" {
					// 允许各用户在根目录下创建文件夹
					c.Next()
					return
				}
				// 必须是path对应目录的所有者才能在path下创建新文件/文件夹
				auth, err := acl.Enforce(acl.Enforcer, strconv.Itoa(int(uid)), path, "ALL")
				if !auth || err != nil {
					c.JSON(200, serializer.PermissionDenied())
					c.Abort()
					return
				}
			} else if strings.Contains(p, "rename") {
				// rename
				path, exists := c.GetQuery("path")
				if !exists {
					c.Next()
					return
				}
				oldName, exists := c.GetQuery("oldName")
				if !exists {
					c.Next()
					return
				}
				// 必须是文件的所有者才能重命名文件
				auth, err := acl.Enforce(acl.Enforcer, strconv.Itoa(int(uid)), strings.ReplaceAll(filepath.Join(path, oldName), "\\", "/"), "ALL")
				if !auth || err != nil {
					c.JSON(200, serializer.PermissionDenied())
					c.Abort()
					return
				}
			} else if strings.Contains(p, "delete") {
				// delete
				path, exists := c.GetQuery("path")
				if !exists {
					c.Next()
					return
				}
				name, exists := c.GetQuery("name")
				if !exists {
					c.Next()
					return
				}
				// 必须是文件的所有者才能删除文件
				auth, err := acl.Enforce(acl.Enforcer, strconv.Itoa(int(uid)), strings.ReplaceAll(filepath.Join(path, name), "\\", "/"), "ALL")
				if !auth || err != nil {
					c.JSON(200, serializer.PermissionDenied())
					c.Abort()
					return
				}
			} else if strings.Contains(p, "chunk") {
				// chunk
				if m == "GET" {
					path, exists := c.GetQuery("relativePath")
					if !exists {
						c.Next()
						return
					}
					// 目录的所有者才能上传文件
					auth, err := acl.Enforce(acl.Enforcer, strconv.Itoa(int(uid)), path, "ALL")
					if !auth || err != nil {
						c.JSON(200, serializer.PermissionDenied())
						c.Abort()
						return
					}
				} else if m == "POST" {
					path, exists := c.GetPostForm("relativePath")
					if !exists {
						c.Next()
						return
					}
					// 目录的所有者才能上传文件
					auth, err := acl.Enforce(acl.Enforcer, strconv.Itoa(int(uid)), path, "ALL")
					if !auth || err != nil {
						c.JSON(200, serializer.PermissionDenied())
						c.Abort()
						return
					}
				} else {
					c.Next()
					return
				}
			} else if strings.Contains(p, "merge") {
				// merge
				path, exists := c.GetQuery("relativePath")
				if !exists {
					c.Next()
					return
				}
				// 目录的所有者才能上传文件
				auth, err := acl.Enforce(acl.Enforcer, strconv.Itoa(int(uid)), path, "ALL")
				if !auth || err != nil {
					c.JSON(200, serializer.PermissionDenied())
					c.Abort()
					return
				}
			} else {
				c.Next()
				return
			}

		}


		c.Next()
	}
}