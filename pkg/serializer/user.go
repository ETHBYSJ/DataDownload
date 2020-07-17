package serializer

import (
	"go-file-manager/models"
	"go-file-manager/pkg/e"
)

type User struct {
	ID 			uint `json:"id"`
	Email 		string `json:"userName"`
	Status 		bool `json:"status"`
	CreatedAt 	int64 `json:"createdAt"`
	UpdatedAt 	int64 `json:"updatedAt"`
}

func CheckLogin() Response {
	return Response{
		Code: e.CodeCheckLogin,
		Msg: "未登录",
	}
}

func BuildUserResponse(user models.User) Response {
	return Response{
		Data: BuildUser(user),
	}
}

func BuildUser(user models.User) User {
	return User {
		// ID: hashid.HashID(user.ID, hashid.UserID),
		ID: user.ID,
		Email: user.Email,
		Status: user.Status,
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}
}
