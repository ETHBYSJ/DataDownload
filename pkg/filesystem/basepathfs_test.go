package filesystem

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRealPath(t *testing.T) {
	asserts := assert.New(t)
	baseFs := &BasePathFs{source: NewOsFs(), path: "D:\\storage"}
	_, err := baseFs.RealPath("abc")
	asserts.NoError(err)
}

func TestBasePath(t *testing.T) {
	asserts := assert.New(t)
	baseFs := &BasePathFs{source: NewOsFs(), path: "D:\\storage"}
	_, err := baseFs.Create("test1.txt")
	asserts.NoError(err)
}
