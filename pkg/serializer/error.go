package serializer

import (
	"go-file-manager/pkg/e"
)

func NotFound() Response {
	return Response {
		Code: e.CodeNotFound,
		Msg: "资源不存在",
	}
}

// 错误处理
type AppError struct {
	Code		int
	Msg 		string
	RawError	error
}

func NewError(code int, msg string, err error) AppError {
	return AppError{
		Code: 		code,
		Msg: 		msg,
		RawError: 	err,
	}
}

func (err AppError) Error() string {
	return err.Msg
}

func DBErr(msg string, err error) Response {
	if msg == "" {
		msg = "数据库操作失败"
	}
	return Err(e.CodeDBError, msg, err)
}

func ParamErr(msg string, err error) Response {
	if msg == "" {
		msg = "参数错误"
	}
	return Err(e.CodeParamError, msg, err)
}

func Err(errCode int, msg string, err error) Response {
	if appError, ok := err.(AppError); ok {
		errCode = appError.Code
		err = appError.RawError
		msg = appError.Msg
	}

	res := Response{
		Code: 	errCode,
		Msg: 	msg,
	}
	res.Error = err.Error()
	return res
}
