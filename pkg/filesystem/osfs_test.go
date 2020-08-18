package filesystem

import (
	"github.com/stretchr/testify/assert"
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

/*
func TestOsFs_Mkdir(t *testing.T) {
	asserts := assert.New(t)
	osFs := &OsFs{}
	e := osFs.Mkdir(baseDir + string(filepath.Separator) + "test", os.ModePerm)
	asserts.NoError(e)
	e = osFs.Mkdir(baseDir + string(filepath.Separator) + "test", os.ModePerm)
	asserts.Error(e)
}
*/
func TestOsFs_Rename(t *testing.T) {
	asserts := assert.New(t)
	osFs := &OsFs{}
	e := osFs.Rename("D:\\storage\\a0\\test1.txt", "D:\\storage\\a0\\test.txt")
	asserts.NoError(e)
}
