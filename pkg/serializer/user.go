package serializer

import (
	"go-file-manager/models"
	"go-file-manager/pkg/e"
)

type User struct {
	ID 			uint   	`json:"id"`
	Email 		string 	`json:"email"`
	FirstName 	string 	`json:"firstName"`
	LastName 	string 	`json:"lastName"`
	PhoneNumber string 	`json:"phoneNumber"`
	School 		string 	`json:"school"`
	Role 		string 	`json:"role"`
	Status 		bool  	`json:"status"`
	CreatedAt 	int64 	`json:"createdAt"`
	UpdatedAt 	int64 	`json:"updatedAt"`
	UserType 	string 	`json:"userType"`
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
		ID: user.ID,
		Email: user.Email,
		FirstName: user.FirstName,
		LastName: user.LastName,
		PhoneNumber: user.PhoneNumber,
		School: user.School,
		Role: user.Role,
		Status: user.Status,
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
		UserType: user.UserType,
	}
}
