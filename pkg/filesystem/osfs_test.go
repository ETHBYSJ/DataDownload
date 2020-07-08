package filesystem

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

var (
	baseDir = "D:\\storage"
)

func TestCreate(t *testing.T) {
	asserts := assert.New(t)
	osFs := &OsFs{}
	_, e := osFs.Create(baseDir + string(filepath.Separator) + "test.txt")
	asserts.NoError(e)
}

func TestMkdir(t *testing.T) {
	asserts := assert.New(t)
	osFs := &OsFs{}
	e := osFs.Mkdir(baseDir + string(filepath.Separator) + "test", os.ModePerm)
	asserts.NoError(e)
	e = osFs.Mkdir(baseDir + string(filepath.Separator) + "test", os.ModePerm)
	asserts.Error(e)
}