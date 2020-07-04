package models

import (
	// "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserByID(t *testing.T) {

}

func TestGetUserByEmail(t *testing.T) {

}
func TestUserSetPassword(t *testing.T) {
	asserts := assert.New(t)
	user := User{}
	err := user.SetPassword("123456")
	asserts.NoError(err)
	asserts.NotEmpty(user.Password)
	// fmt.Print(user.Password)
}

func TestUserCheckPassword(t *testing.T) {
	asserts := assert.New(t)
	user := User{}
	err := user.SetPassword("test password")
	asserts.NoError(err)
	// 密码正确
	res, err := user.CheckPassword("test password")
	asserts.NoError(err)
	asserts.True(res)
	// 密码错误
	res, err = user.CheckPassword("wrong password")
	asserts.NoError(err)
	asserts.False(res)
	// 密码为空
	user = User{}
	res, err = user.CheckPassword("test password")
	asserts.Error(err)
	asserts.False(res)
	// 未知密码类型
	user = User{}
	user.Password = "1:2:3"
	res, err = user.CheckPassword("test password")
	asserts.Error(err)
	asserts.False(res)
	// 密码正确
	user = User{}
	user.Password = "9PmwyHJhbjcbLOG7:f14272ab12b6a886a4e3e8f7429bcfc94c6b0e78"
	res, err = user.CheckPassword("123456")
	asserts.NoError(err)
	asserts.True(res)
	// 密码错误
	user = User{}
	user.Password = "9PmwyHJhbjcbLOG7:f14272ab12b6a886a4e3e8f7429bcfc94c6b0e78"
	res, err = user.CheckPassword("12345")
	asserts.NoError(err)
	asserts.False(res)
}

func TestNewUser(t *testing.T) {
	asserts := assert.New(t)
	newUser := NewUser()
	asserts.IsType(User{}, newUser)
	asserts.Empty(newUser.Email)
	asserts.Equal(newUser.Email, "")
}