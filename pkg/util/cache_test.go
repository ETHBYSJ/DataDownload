package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExpire(t *testing.T) {
	asserts := assert.New(t)
	var table *Table = NewTable()
	table.Add("key1", time.Second*5, "1")
	time.Sleep(3 * time.Second)
	table.Add("key2", time.Second*5, "2")
	time.Sleep(3 * time.Second)
	val, err := table.Value("key1")
	asserts.NoError(err)
	asserts.Nil(val)
	val, err = table.Value("key2")
	asserts.NoError(err)
	asserts.NotNil(val)
}
