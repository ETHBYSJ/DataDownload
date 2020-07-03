package serializer

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
