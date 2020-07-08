package models

import (
	"github.com/jinzhu/gorm"
	"go-file-manager/pkg/util"
)

type File struct {
	gorm.Model
	Name	string 	`gorm:"type:varchar(255) not null"`
	IsDir 	bool 	`gorm:"type:tinyint(1) not null;default:0"`
	Path 	string 	`gorm:"type:text"`
	OwnerID	uint
	Owner 	User	`gorm:"foreignkey:OwnerID"`
	Privacy	bool 	`gorm:"type:tinyint(1) not null;default:0"`
}

// 创建新文件
func (file *File) Create() (uint, error) {
	if err := DB.Create(file).Error; err != nil {
		util.Log().Warning("无法创建文件记录, %s", err)
		return 0, err
	}
	return file.ID, nil
}

// 更新文件名
func (file *File) UpdateSourceName(value string) error {
	return DB.Model(&file).Update("name", value).Error
}

func (file *File) GetName() string {
	return file.Name
}

func (file *File) GetIsDir() bool {
	return file.IsDir
}

func (file *File) GetPath() string {
	return file.Path
}

func (file *File) GetOwner() User {
	DB.Model(file).Where("id=?", file.ID).Related(&file.Owner)
	return file.Owner
}

func (file *File) GetPrivacy() bool {
	return file.Privacy
}


