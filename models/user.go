package models

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/jinzhu/gorm"
	"go-file-manager/pkg/util"
	"strings"
)

const (
	// 账户状态
	Active = iota
	// 被封禁
	Banned
)

type User struct {
	gorm.Model
	Email 		string `gorm:"type:varchar(100);unique_index"`
	Password	string `json:"-"`
	Status 		int
	// CreateTime 	timestamp.Timestamp
	// UpdateTime  timestamp.Timestamp
}

func NewUser() User {
	return User{}
}

func GetUserByID(ID interface{}) (User, error) {
	var user User
	result := DB.Set("gorm:auto_preload", true).First(&user, ID)
	return user, result.Error
}

func GetActiveUserByID(ID interface{}) (User, error) {
	var user User
	result := DB.Set("gorm:auto_preload", true).Where("status = ?", Active).First(&user, ID)
	return user, result.Error
}

func GetUserByEmail(email string) (User, error) {
	var user User
	result := DB.Set("gorm:auto_preload", true).Where("status = ? and email = ?", Active, email).First(&user)
	return user, result.Error
}

// 设定用户状态
func (user *User) SetStatus(status int) {
	DB.Model(&user).Update("status", status)
}

// 根据密码明文设定User的Password字段
func (user *User) SetPassword(password string) error {
	salt := util.RandStringRunes(16)
	// salt和密码组合的sha1摘要
	hash := sha1.New()
	_, err := hash.Write([]byte(password + salt))
	bs := hex.EncodeToString(hash.Sum(nil))
	if err != nil {
		return err
	}
	// salt:摘要
	user.Password = salt + ":" + string(bs)
	return nil
}

//根据明文校验密码
func (user *User) CheckPassword(password string) (bool, error) {
	// 根据存储密码拆分为salt和摘要
	passwordStore := strings.Split(user.Password, ":")
	if len(passwordStore) != 2 {
		return false, errors.New("Unknown password type")
	}
	// 计算
	hash := sha1.New()
	_, err := hash.Write([]byte(password + passwordStore[0]))
	bs := hex.EncodeToString(hash.Sum(nil))
	if err != nil {
		return false, err
	}
	return bs == passwordStore[1], nil
}









