package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go-file-manager/models"
	"go-file-manager/pkg/serializer"
)

func ErrorResponse(err error) serializer.Response {
	if _, ok := err.(validator.ValidationErrors); ok {
		return serializer.ParamErr("字段验证错误", err)
	}
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return serializer.ParamErr("JSON类型不匹配", err)
	}
	return serializer.ParamErr("", err)
}

func CurrentUser(c *gin.Context) *models.User {
	if user, _ := c.Get("user"); user != nil {
		if u, ok := user.(*models.User); ok {
			return u
		}
	}
	return nil
}
