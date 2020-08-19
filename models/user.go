package models

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/jinzhu/gorm"
	"go-file-manager/pkg/e"
	"go-file-manager/pkg/util"
	"strings"
)

type User struct {
	gorm.Model
	Email       string `gorm:"type:varchar(100);unique_index"`
	FirstName   string `gorm:"type:varchar(100) not null"`
	LastName    string `gorm:"type:varchar(100) not null"`
	PhoneNumber string `gorm:"type:varchar(100) not null"`
	School      string `gorm:"type:varchar(100) not null"`
	Role        string `gorm:"type:varchar(100) not null"`
	Password    string `json:"-"`
	Status      bool   `gorm:"type:tinyint(1) not null;default:1"`
	UserType    string `gorm:"type:varchar(8) not null;default:'User'"`
	Language    string `gorm:"type:varchar(8) not null;default:'zh-CN'"`
	Validate 	bool   `gorm:"type:tinyint(1) not null;default:0"`
}

func NewUser() *User {
	return &User{}
}

// 我的文件
func GetMyFiles(uid uint, page int, pageSize int, keyword string, category int) ([]*File, uint64, error) {
	files := make([]*File, 0)
	var err error
	if category == 0 {
		// 未审核
		err = DB.Model(&File{}).Where("name LIKE ? AND review = ? AND owner_id = ? AND is_dir = ?", "%"+keyword+"%", 0, uid, 0).Order("updated_at desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&files).Error
	} else if category == 1 {
		// 已审核
		err = DB.Model(&File{}).Where("name LIKE ? AND review = ? AND owner_id = ? AND is_dir = ?", "%"+keyword+"%", 1, uid, 0).Order("updated_at desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&files).Error
	} else if category == 2 {
		// 全部
		err = DB.Model(&File{}).Where("name LIKE ? AND owner_id = ? AND is_dir = ?", "%"+keyword+"%", uid, 0).Order("updated_at desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&files).Error
	}

	if err != nil {
		return []*File{}, 0, err
	}
	var count uint64
	if category == 0 {
		// 未审核
		err = DB.Model(&File{}).Where("name LIKE ? AND review = ? AND owner_id = ? AND is_dir = ?", "%"+keyword+"%", 0, uid, 0).Count(&count).Error
	} else if category == 1 {
		// 已审核
		err = DB.Model(&File{}).Where("name LIKE ? AND review = ? AND owner_id = ? AND is_dir = ?", "%"+keyword+"%", 1, uid, 0).Count(&count).Error
	} else if category == 2 {
		// 全部
		err = DB.Model(&File{}).Where("name LIKE ? AND owner_id = ? AND is_dir = ?", "%"+keyword+"%", uid, 0).Count(&count).Error
	}
	if err != nil {
		return []*File{}, 0, err
	}
	return files, count, nil
}

// 分页列出用户
func ListUser(pageSize int, page int) ([]*User, uint64, error) {
	users := make([]*User, 0)
	err := DB.Order("updated_at desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&users).Error
	if err != nil {
		return []*User{}, 0, err
	}
	var count uint64
	err = DB.Table("users").Count(&count).Error
	if err != nil {
		return []*User{}, 0, err
	}
	return users, count, nil
}

func checkLanguage(language string) bool {
	if language == "zh-CN" || language == "en-US" {
		return true
	}
	return false
}

func LanguageSet(user *User, language string) error {
	if !checkLanguage(language) {
		return e.ErrLanguageInvalid
	}
	result := DB.Model(&user).Update("language", language)
	return result.Error
}

func GetUserByID(ID interface{}) (*User, error) {
	var user User
	result := DB.Set("gorm:auto_preload", true).First(&user, ID)
	return &user, result.Error
}

func GetActiveUserByID(ID interface{}) (*User, error) {
	var user User
	result := DB.Set("gorm:auto_preload", true).Where("status = ?", 1).First(&user, ID)
	return &user, result.Error
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	result := DB.Set("gorm:auto_preload", true).Where("status = ? and email = ?", 1, email).First(&user)
	return &user, result.Error
}

// 设定激活字段
func (user *User) SetValidate() error {
	return DB.Model(&user).Update("validate", true).Error
}

// 设定用户状态
func (user *User) SetStatus(status bool) error {
	if status {
		return DB.Model(&user).Update("status", true).Error
	}
	return DB.Model(&user).Update("status", false).Error

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
