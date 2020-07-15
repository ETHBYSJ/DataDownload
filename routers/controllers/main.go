package controllers

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"go-file-manager/pkg/serializer"
)

func ErrorResponse(err error) serializer.Response {
	if err == nil {
		err = errors.New("")
	}
	if _, ok := err.(validator.ValidationErrors); ok {
		return serializer.ParamErr("字段验证错误", err)
	}
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return serializer.ParamErr("JSON类型不匹配", err)
	}
	return serializer.ParamErr("参数错误", err)
}


