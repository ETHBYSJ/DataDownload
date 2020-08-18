package filesystem

import "strings"

// 用于进行验证

// 文件/路径名保留字符
var reservedCharacter = []string{"\\", "?", "*", "<", "\"", ":", ">", "/", "|"}

func (fs *FileSystem) validateLegalName(name string) bool {
	// 是否包含保留字符
	for _, value := range reservedCharacter {
		if strings.Contains(name, value) {
			return false
		}
	}
	// 是否超出长度限制
	if len(name) >= 256 {
		return false
	}
	// 是否为空
	if len(name) == 0 {
		return false
	}
	// 结尾不能是空格
	if strings.HasSuffix(name, " ") {
		return false
	}
	return true
}

// TODO 验证文件扩展名
func (fs *FileSystem) ValidateExtension(fileName string) bool {
	return true
}
