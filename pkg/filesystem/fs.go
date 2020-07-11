package filesystem

import (
	"go-file-manager/models"
	"go-file-manager/pkg/e"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// 通用接口，与具体实现分离
type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Writer
	io.WriterAt

	Name() string
	Readdir(count int) ([]os.FileInfo, error)
	Readdirnames(n int) ([] string, error)
	Stat() (os.FileInfo, error)
	Sync() error
	Truncate(size int64) error
	WriteString(s string) (ret int, err error)
}

// 通用接口，与具体实现分离
type Fs interface {
	Create(name string) (File, error)
	Mkdir(name string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Open(name string) (File, error)
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldName, newName string) error
	Stat(name string) (os.FileInfo, error)
	Name() string
	Chmod(name string, mode os.FileMode) error
	Chtimes(name string, atime time.Time, mtime time.Time) error
	// 实际路径
	RealPath(name string) (string, error)
}
/*
type FsWrapper struct {
	Fs
}
*/
// 实际使用的文件系统结构
type FileSystem struct {
	// 互斥锁
	Lock sync.Mutex
	// 文件系统
	// Fs *FsWrapper
	Fs Fs
	// 钩子函数
	Hooks map[string][]Hook
}

// 重命名，且只允许在同一目录下操作
// 注意，如果重命名操作的对象是文件，且目标文件已经存在的话，目标文件将被覆盖
// 因此需要先查数据库
func (fs *FileSystem) Rename(oldName, newName, dirPath string) error {
	// 验证新名字
	if !fs.validateLegalName(newName) || !fs.ValidateExtension(newName) {
		return e.ErrIllegalName
	}
	// 先查数据库中是否有同名文件
	_, err := models.GetFileByNameAndPath(newName, dirPath)
	if err != nil {
		return e.ErrFileExisted
	}
	// 需要先将目录与文件路径合并
	err = fs.Fs.Rename(filepath.Join(dirPath, oldName), filepath.Join(dirPath, newName))
	return err
}

// func (fs *FileSystem) Delete()

// 创建文件夹
func (fs *FileSystem) CreateDirectory(user *models.User, name, dirPath string) (*FileInfo, error) {
	/*
	ginCtx, ok := ctx.Value(GinCtx).(*gin.Context)
	if !ok {
		return nil, ErrGetContext
	}

	userStore, ok := ginCtx.Get("user")
	if !ok {
		return nil, ErrGetUser
	}
	user, _ := userStore.(*models.User)
	*/
	// 检查目录名是否合法
	if !fs.validateLegalName(name) {
		return nil, e.ErrIllegalName
	}
	// 数据库操作
	newFolder := models.File{
		Name: name,
		Path: dirPath,
		IsDir: true,
		Owner: *user,
		OwnerID: user.ID,
		Review: true,	// 文件夹默认是审核通过的
	}
	id, err := newFolder.Create()
	if err != nil {
		return nil, e.ErrFolderExisted
	}
	err = fs.Fs.Mkdir(filepath.Join(dirPath, name), os.ModePerm)
	if err != nil {
		// 删除数据库插入的记录
		_ = models.DeleteFileByID(id)
		return nil, e.ErrFolderExisted
	}
	file := FileInfo{
		ID: id,
		Name: name,
		Path: dirPath,
		IsDir: true,
		ModTime: newFolder.CreatedAt,
		Review: true,
	}
	return &file, nil
}

// 列出路径下的内容
func (fs *FileSystem) List(sorting Sorting, dirPath string) (*Listing, error) {
	dir, err := ReadDir(fs.Fs, dirPath)
	if err != nil {
		return nil, err
	}
	listing := &Listing{
		Items: []*FileInfo{},
		NumDirs: 0,
		NumFiles: 0,
		Sorting: sorting,
	}
	for _, f := range dir {
		name := f.Name()

		file := &FileInfo{
			Name: name,
			// Path: fPath,
			Path: dirPath,
			Size: f.Size(),
			IsDir: f.IsDir(),
			ModTime: f.ModTime(),
		}
		// 查数据库获取唯一ID
		fileModel, err := models.GetFileByNameAndPath(name, dirPath)
		if err != nil {
			continue
		}
		file.ID = fileModel.ID
		if file.IsDir {
			listing.NumDirs++
		} else {
			listing.NumFiles++
		}
		listing.Items = append(listing.Items, file)
	}
	listing.ApplySort()
	return listing, nil
}


