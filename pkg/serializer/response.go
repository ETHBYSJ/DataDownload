package serializer

// 响应消息体
type Response struct {
	Code  int			`json:"code"`
	Data  interface{}	`json:"data,omitempty"`
	Msg   string		`json:"msg"`
	Error string		`json:"error, omitempty"`
}