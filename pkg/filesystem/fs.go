package filesystem

import (
	MD5 "crypto/md5"
	"encoding/hex"
	"go-file-manager/models"
	"go-file-manager/pkg/e"
	"go-file-manager/pkg/util"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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

type CheckInfo struct {
	SkipUpload  	bool   `json:"skipUpload"`
	Uploaded 		[]int  `json:"uploaded"`
	NeedMerge 		bool   `json:"needMerge"`
	Identifier 		string `json:"identifier"`
	FileName 		string `json:"filename"`
	RelativePath 	string `json:"relativePath"`
	TotalChunks 	int    `json:"totalChunks"`
}

// 实际使用的文件系统结构
type FileSystem struct {
	// 互斥锁
	locker *util.Table
	// 文件系统
	// Fs *FsWrapper
	Fs Fs
	// 钩子函数
	Hooks map[string][]Hook
}

// 设置文件分享状态
func (fs *FileSystem) SetShare(name string, path string, share bool) error {
	fm, err := models.GetFileByNameAndPath(name, path)
	if err != nil {
		return e.ErrUploadPathNotExists
	}
	err = fm.UpdateShare(share)
	if err != nil {
		return err
	}
	return nil
}

// 获取锁，立即返回
func (fs *FileSystem) Lock(key string) error {
	cacheItem, err := fs.locker.Value(key)
	if err != nil {
		util.Log().Error("获取锁错误 %s", err)
		return err
	}
	if cacheItem != nil && cacheItem.Data() != nil {
		return e.ErrLock
	}
	duration := 12 * time.Hour
	fs.locker.Add(key, duration, true)
	return nil
}

func (fs *FileSystem) Unlock(key string) error {
	exist := fs.locker.Exists(key)
	if exist {
		_, err := fs.locker.Delete(key)
		return err
	} else {
		util.Log().Error("解锁错误")
		return e.ErrUnlock
	}
}

// 保存文件，此函数取自Gin
func (fs *FileSystem) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	out, err := fs.Fs.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return err
}

// 合并分块文件
func (fs *FileSystem) MergeChunk(name string, path string, user *models.User, md5 string, totalChunks int) (*CheckInfo, error) {
	target := path
	if path == "/" {
		target += name
	} else {
		// path: /test
		target += "/" + name
	}
	err := fs.Lock(target)
	if err != nil {
		return nil, err
	}
	defer fs.Unlock(target)
	// 查看上传的路径是否存在
	if path == "/" {

	} else {
		util.Log().Info(" merge path = %s\n", path)
		last := strings.LastIndex(path, "/")
		n := path[last + 1:]
		p := path[:last]
		if p == "" {
			p = "/"
		}
		_, err := models.GetFileByNameAndPath(n, p)

		if err != nil {
			return nil, e.ErrUploadPathNotExists
		}
	}
	// 查数据库中是否已有相同文件名、路径名、MD5的记录
	fm, err := models.GetFileByNameAndPath(name, path)
	if err != nil || fm.MD5 != md5 || !fm.Uploaded {
		return nil, e.ErrGetUploadRecord
	} else {
		if fm.Merge {
			// 已经合并过了
			return nil, e.ErrAlreadyMerged
		} else {
			// 执行合并
			fs.Fs.Remove(target)
			// 不能是追加模式
			f, err := fs.Fs.OpenFile(target, os.O_CREATE|os.O_RDWR, 0777)
			f.Close()
			if err != nil {
				return nil, e.ErrMergeFile
			} else {
				// writer := bufio.NewWriter(f)
				for i := 1; i <= totalChunks; i++ {
					currentPath := path
					if path == "/" {
						currentPath += md5 + "_" + name + "_" + strconv.Itoa(i)
					} else {
						// path: /test
						currentPath += "/" + md5 + "_" + name + "_" + strconv.Itoa(i)
					}
					fileBuffer, err := ReadFile(fs.Fs, currentPath)
					util.Log().Info("file buffer length = %v", len(fileBuffer))
					if err != nil {
						return nil, e.ErrMergeFile
					}
					err = AppendFile(fs.Fs, target, fileBuffer, 0777)
					if err != nil {
						return nil, err
					}
				}
				// 重新打开文件，目的是获取md5
				f, err = fs.Fs.OpenFile(target, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
				defer f.Close()
				md5Merge := MD5.New()
				writtenNum, err := io.Copy(md5Merge, f)
				if err != nil {
					return nil, e.ErrMergeFileMD5
				}
				util.Log().Info("复制文件 %v", writtenNum)
				md5Val := hex.EncodeToString(md5Merge.Sum(nil))
				if md5Val != md5 {
					return nil, e.ErrMergeFileMD5
				}
				// 校验成功，更新数据库相应字段
				fm.UpdateMerge(true)
				for i := 1; i <= totalChunks; i++ {
					currentPath := path
					if path == "/" {
						currentPath += md5 + "_" + name + "_" + strconv.Itoa(i)
					} else {
						// path: /test
						currentPath += "/" + md5 + "_" + name + "_" + strconv.Itoa(i)
					}
					fs.Fs.Remove(currentPath)
				}
				return &CheckInfo{SkipUpload: true, Uploaded: []int{}, NeedMerge: false, Identifier: md5, FileName: name, RelativePath: path, TotalChunks: totalChunks}, nil
			}
		}
	}
}

// 分块上传文件
func (fs *FileSystem) UploadChunk(name string, path string, user *models.User, md5 string, size int64, chunkNumber int, totalChunks int, file *multipart.FileHeader) (*CheckInfo, error) {
	savePath := path
	if path == "/" {
		savePath += md5 + "_" + name + "_" + strconv.Itoa(chunkNumber)
	} else {
		// path: /test
		savePath += "/" + md5 + "_" + name + "_" + strconv.Itoa(chunkNumber)
	}
	util.Log().Info("savePath: %s", savePath)
	err := fs.Lock(savePath)
	if err != nil {
		return nil, err
	}
	defer fs.Unlock(savePath)
	// 查看上传的路径是否存在
	if path == "/" {

	} else {
		util.Log().Info("upload path = %s\n", path)
		last := strings.LastIndex(path, "/")
		n := path[last + 1:]
		p := path[:last]
		if p == "" {
			p = "/"
		}
		_, err := models.GetFileByNameAndPath(n, p)

		if err != nil {
			return nil, e.ErrUploadPathNotExists
		}
	}
	// 查数据库中是否已有相同文件名、路径名、MD5的记录
	fm, err := models.GetFileByNameAndPath(name, path)
	if err != nil || fm.MD5 != md5 {
		return nil, e.ErrGetUploadRecord
	} else {
		// 是否已经上传
		if fm.Uploaded {
			return &CheckInfo{SkipUpload: true, Uploaded: []int{}, NeedMerge: !fm.Merge, Identifier: md5, FileName: name, RelativePath: path, TotalChunks: totalChunks}, nil
		} else {
			// 把上传的文件保存到提供的路径
			err := fs.SaveUploadedFile(file, savePath)
			if err != nil {
				return nil, err
			} else {
				// 查看当前文件夹下的分块情况
				files, err := ReadDir(fs.Fs, path)
				if err != nil {
					return nil, e.ErrCheckChunks
				}
				uploaded := make([]int, 0)
				for _, f := range files {
					// 分块文件命名格式:
					// md5_name_index
					fileName := f.Name()
					split := strings.Split(fileName, "_")
					identifier := split[0]
					// 必须md5和文件名一致
					if identifier == md5 && name == split[1] {
						index, err := strconv.Atoi(split[2])
						if err == nil {
							uploaded = append(uploaded, index)
						}
					}
				}
				sort.Ints(uploaded)
				util.Log().Info("uploaded = %v", uploaded)
				// 全部分块已经上传
				if len(uploaded) == totalChunks {
					// 更新数据库
					fm.UpdateUploaded(true)
					return &CheckInfo{SkipUpload: false, Uploaded: uploaded, NeedMerge: true, Identifier: md5, FileName: name, RelativePath: path, TotalChunks: totalChunks}, nil
				}
				return &CheckInfo{SkipUpload: false, Uploaded: uploaded, NeedMerge: false, Identifier: md5, FileName: name, RelativePath: path, TotalChunks: totalChunks}, nil
			}
		}
	}

}

// 检查文件分块上传情况
func (fs *FileSystem) CheckChunk(name string, path string, user *models.User, md5 string, size int64, totalChunks int) (*CheckInfo, error) {
	// 查看上传的路径是否存在
	if path == "/" {

	} else {
		last := strings.LastIndex(path, "/")
		n := path[last + 1:]
		p := path[:last]
		if p == "" {
			p = "/"
		}
		_, err := models.GetFileByNameAndPath(n, p)
		util.Log().Info("name: %s, path: %s", n, p)
		if err != nil {
			return nil, e.ErrUploadPathNotExists
		}
	}
	// 查数据库中是否已有相同文件名、路径名、MD5的记录
	file, err := models.GetFileByNameAndPath(name, path)
	if err != nil {
		// 没有相关记录
		file = models.NewFile()
		file.Name = name
		file.MD5 = md5
		file.Path = path
		file.IsDir = false
		file.OwnerID = user.ID
		file.Owner = *user
		file.Size = size
		// 文件默认是共享的
		file.Share = true
		if err := models.DB.Create(&file).Error; err != nil {
			util.Log().Warning("创建文件记录失败")
			return nil, e.ErrCreateFileRecord
		} else {
			return &CheckInfo{SkipUpload: false, Uploaded: []int{}, NeedMerge: false, Identifier: md5, FileName: name, RelativePath: path, TotalChunks: totalChunks}, nil
		}
	} else {
		// 存在相关记录，需要进一步判断
		// MD5匹配，说明是相同的文件，检查分块情况
		if file.MD5 == md5 {
			// 如果已经上传过
			if file.Uploaded == true {
				return &CheckInfo{SkipUpload: true, Uploaded: []int{}, NeedMerge: !file.Merge, Identifier: md5, FileName: name, RelativePath: path, TotalChunks: totalChunks}, nil
			} else {
				files, err := ReadDir(fs.Fs, path)
				if err != nil {
					return nil, e.ErrCheckChunks
				}
				uploaded := make([]int, 0)
				for _, f := range files {
					// 分块文件命名格式:
					// md5_name_index
					fileName := f.Name()
					split := strings.Split(fileName, "_")
					identifier := split[0]
					if identifier == md5 {
						index, err := strconv.Atoi(split[2])
						if err == nil {
							uploaded = append(uploaded, index)
						}
					}
				}
				sort.Ints(uploaded)
				return &CheckInfo{SkipUpload: false, Uploaded: uploaded, NeedMerge: true, Identifier: md5, FileName: name, RelativePath: path, TotalChunks: totalChunks}, nil
			}

		} else {
			// MD5不一致，说明是同名不同内容的文件。拒绝请求
			return nil, e.ErrFileCover
		}
	}
}

// 重命名，且只允许在同一目录下操作
// 注意，如果重命名操作的对象是文件，且目标文件已经存在的话，目标文件将被覆盖
// 因此需要先查数据库
func (fs *FileSystem) Rename(oldName, newName, dirPath string) error {
	// 验证新名字
	if !fs.validateLegalName(newName) || !fs.ValidateExtension(newName) {
		return e.ErrIllegalName
	}
	file, err := models.GetFileByNameAndPath(oldName, dirPath)
	if err != nil {
		return e.ErrFileNotExist
	}
	// 先查数据库中是否有同名文件
	_, err = models.GetFileByNameAndPath(newName, dirPath)
	if err == nil {
		return e.ErrFileExisted
	}
	// 需要先将目录与文件路径合并
	err = fs.Fs.Rename(filepath.Join(dirPath, oldName), filepath.Join(dirPath, newName))
	if err != nil {
		return err
	}
	// 修改数据库中的记录
	err = file.UpdateRename(oldName, newName, dirPath)
	return err
}

// 带有保护的重命名操作
func (fs *FileSystem) RenameAtomic(oldName, newName, dirPath string) error {
	fullPath := filepath.Join(dirPath, newName)
	err := fs.Lock(fullPath)
	defer fs.Unlock(fullPath)
	if err != nil {
		return err
	}
	return fs.Rename(oldName, newName, dirPath)

}



// 删除文件
func (fs *FileSystem) Delete(name string, path string) error {
	err := fs.Fs.Remove(filepath.Join(path, name))
	if err != nil {
		return e.ErrDelete
	}
	err = models.DeleteFileByNameAndPath(name, path)
	if err != nil {
		return e.ErrDelete
	}
	return nil
}

// 创建文件夹
func (fs *FileSystem) CreateDirectory(user *models.User, name, dirPath string) (*FileInfo, error) {
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
		Review: true,	// 文件夹默认不需要审核
		Share: true,  // 文件夹默认是共享的
	}
	id, err := newFolder.Create()
	if err != nil {
		return nil, e.ErrFolderExisted
	}
	err = fs.Fs.Mkdir(filepath.Join(dirPath, name), os.ModePerm)
	if err != nil {
		// 删除数据库插入的记录
		_ = models.DeleteFileByID(id)
		return nil, err
	}
	file := FileInfo{
		ID: id,
		Name: name,
		Path: dirPath,
		IsDir: true,
		ModTime: newFolder.CreatedAt,
		Review: true,
		OwnerID: user.ID,
		Share: true,
	}
	return &file, nil
}

// 根据关键字查询
func (fs *FileSystem) ListByKeyword(sorting Sorting, dirPath string, keyword string) (*Listing, error) {
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
		if keyword != "" && strings.Index(name, keyword) == -1 {
			continue
		}
		file := &FileInfo{
			Name: name,
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
		if !fileModel.Merge {
			continue
		}
		file.ID = fileModel.ID
		file.OwnerID = fileModel.OwnerID
		file.Share = fileModel.Share
		file.Review = fileModel.Review
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
		if !fileModel.IsDir && !fileModel.Merge {
			continue
		}
		file.ID = fileModel.ID
		file.OwnerID = fileModel.OwnerID
		file.Share = fileModel.Share
		file.Review = fileModel.Review
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


