package filesystem

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSubstr(t *testing.T) {
	asserts := assert.New(t)
	s := "/test"
	last := strings.LastIndex(s, "/")
	n := s[last + 1:]
	p := s[:last]
	if p == "" {
		p = "/"
	}
	fmt.Printf("n = %s, p = %s\n", n, p)
	// fmt.Print(len(s))
	asserts.True(true)
}

