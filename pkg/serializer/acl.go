package serializer

import (
	"go-file-manager/pkg/e"
)

func PermissionDenied() Response {
	return Response{
		Code: e.CodeNoPermissionErr,
		Msg:  "对资源的访问被拒绝",
	}
}
