package filesystem

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
	"testing"
)

func TestSubstr(t *testing.T) {
	asserts := assert.New(t)
	s := "/test"
	last := strings.LastIndex(s, "/")
	n := s[last+1:]
	p := s[:last]
	if p == "" {
		p = "/"
	}
	fmt.Printf("n = %s, p = %s\n", n, p)
	// fmt.Print(len(s))
	asserts.True(true)
}

func TestIndex(t *testing.T) {
	asserts := assert.New(t)
	s := "LastIndex返回字符串str在字符串s中最后出现位置的索引"
	index := strings.Index(s, "字符串")
	asserts.True(index != -1)
	index = strings.Index(s, "Index")
	asserts.True(index != -1)
}

func TestFilePath(t *testing.T) {
	asserts := assert.New(t)
	path := "\\test\\"
	name := "/a.txt"
	// fmt.Println(strings.ReplaceAll(filepath.Join(path, name), "\\", "/"))
	newPath := strings.ReplaceAll(filepath.Join(path, name), "\\", "/")
	asserts.True(newPath == "/test/a.txt")
}
