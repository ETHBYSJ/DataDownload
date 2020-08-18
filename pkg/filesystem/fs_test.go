package filesystem

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestList(t *testing.T) {
	asserts := assert.New(t)
	fs := FileSystem{Fs: &BasePathFs{source: NewOsFs(), path: "D:\\storage"}}
	listing, err := fs.List(Sorting{By: "dic", Asc: true}, "testcase")
	data, _ := json.Marshal(listing)
	fmt.Printf("%s\n", data)
	asserts.NoError(err)
}

func TestRename(t *testing.T) {
	asserts := assert.New(t)
	fs := FileSystem{Fs: &BasePathFs{source: NewOsFs(), path: "D:\\storage"}}
	err := fs.Rename("test1.txt", "test.txt", "a0/a1")
	asserts.NoError(err)
}
