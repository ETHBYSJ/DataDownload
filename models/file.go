package models

import (
	"github.com/jinzhu/gorm"
	"go-file-manager/pkg/util"
)

type File struct {
	gorm.Model
	Name	string 	`gorm:"type:varchar(255) not null;unique_index:idx"`
	IsDir 	bool 	`gorm:"type:tinyint(1) not null;default:0"`
	Path 	string 	`gorm:"type:varchar(255);unique_index:idx"`
	OwnerID	uint
	Owner 	User	`gorm:"foreignkey:OwnerID"`
	Privacy	bool 	`gorm:"type:tinyint(1) not null;default:0"`
	Size 	int64
	Review  bool 	`gorm:"type:tinyint(1) not null;default:0"`
}
// 根据文件名和路径名获取文件
func GetFileByNameAndPath(name string, path string) (File, error) {
	var file File
	result := DB.Where("name = ? AND path = ?", name, path).Find(&file)
	return file, result.Error
}

// 根据ID获取文件
func GetFileByID(ID interface{}) (File, error) {
	var file File
	result := DB.First(&file, ID)
	return file, result.Error
}

// 创建新文件
func (file *File) Create() (uint, error) {
	if err := DB.Create(file).Error; err != nil {
		util.Log().Warning("无法创建文件记录, %s", err)
		return 0, err
	}
	return file.ID, nil
}

func DeleteFileByID(id uint) error {
	result := DB.Unscoped().Where("id = ?", id).Delete(&File{})
	return result.Error
}

// 重命名文件
func (file *File) Rename(newName string) error {
	return DB.Model(&file).Update("name", newName).Error
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


