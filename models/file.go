package models

import (
	"github.com/jinzhu/gorm"
	"go-file-manager/pkg/e"
	"go-file-manager/pkg/util"
	"strings"
)

type File struct {
	gorm.Model
	Name		string 	`gorm:"type:varchar(255) not null;unique_index:idx"`
	IsDir 		bool 	`gorm:"type:tinyint(1) not null;default:0"`
	Path 		string 	`gorm:"type:varchar(255);unique_index:idx"`
	OwnerID		uint
	Owner 		User	`gorm:"foreignkey:OwnerID"`
	// 0代表私有 1代表公开
	Share		bool 	`gorm:"type:tinyint(1) not null;default:0"`
	Size 		int64
	// 是否通过审核
	Review  	bool 	`gorm:"type:tinyint(1) not null;default:0"`
	MD5 		string 	`gorm:"type:varchar(255);"`
	Uploaded 	bool 	`gorm:"type:tinyint(1) not null;default:0"`
	Merge 		bool 	`gorm:"type:tinyint(1) not null;default:0"`
}

func NewFile() *File {
	return &File{}
}

// 根据路径名获取文件
func GetFileByPath(path string) (*File, error) {
	if path == "/" {
		util.Log().Panic("试图获取根目录在数据库中的记录")
		return nil, e.ErrTryToGetRoot
	}
	if path[len(path) - 1] == '/' {
		path = path[:len(path) - 1]
	}
	n := ""
	p := ""
	index := strings.LastIndex(path, "/")
	if index == 0 {
		// ex: /test
		p = "/"
	} else {
		// ex: /test/test1
		p = path[:index]
	}
	n = path[index + 1:]
	util.Log().Info("path = %v, name = %v", p, n)

	return GetFileByNameAndPath(n, p)
}

// 根据文件名和路径名获取文件
func GetFileByNameAndPath(name string, path string) (*File, error) {
	var file File
	result := DB.Where("name = ? AND path = ?", name, path).Find(&file)
	return &file, result.Error
}

// 根据ID获取文件
func GetFileByID(ID interface{}) (*File, error) {
	var file File
	result := DB.First(&file, ID)
	return &file, result.Error
}

// 创建新文件
func (file *File) Create() (uint, error) {
	if err := DB.Create(file).Error; err != nil {
		util.Log().Warning("无法创建文件记录, %s", err)
		return 0, err
	}
	return file.ID, nil
}

// 通过ID删除文件
func DeleteFileByID(id uint) error {
	result := DB.Unscoped().Where("id = ?", id).Delete(&File{})
	return result.Error
}

// 通过文件名和路径删除文件
func DeleteFileByNameAndPath(name string, path string) error {
	result := DB.Unscoped().Where("name = ? AND path = ?", name, path).Delete(&File{})
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

// 更新是否已上传字段
func (file *File) UpdateUploaded(uploaded bool) error {
	return DB.Model(&file).Update("uploaded", uploaded).Error
}

// 更新共享字段
func (file *File) UpdateShare(share bool) error {
	return DB.Model(&file).Update("share", share).Error
}

func (file *File) UpdateRename(oldName, newName, path string) error {
	return DB.Model(&file).Update("name", newName).Error
}

// 更新是否已合并字段
func (file *File) UpdateMerge(merge bool) error {
	return DB.Model(&file).Update("merge", merge).Error
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

func (file *File) GetOwner() *User {
	DB.Model(file).Where("id=?", file.ID).Related(&file.Owner)
	return &file.Owner
}

func (file *File) GetShare() bool {
	return file.Share
}


