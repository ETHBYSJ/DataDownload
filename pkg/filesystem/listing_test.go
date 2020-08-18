package filesystem

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDicOrder(t *testing.T) {
	asserts := assert.New(t)
	s1 := "a00b00"
	s2 := "a0b00"
	asserts.True(s1 < s2)
}

func TestInit(t *testing.T) {
	l := Listing{}
	fmt.Print(l)
}

/*
func TestWin10Sort(t *testing.T) {
	asserts := assert.New(t)
	var items = []*FileInfo{
		{Name: "a00b00", IsDir: true, Size: 100},
		{Name: "a0b00", IsDir: true, Size: 52},
		{Name: "a0B", IsDir: true, Size: 68},
	}
	var sorting = Sorting{By: ByNameWin10Style, Asc: true}
	var listing = Listing{Items: items, Sorting: sorting}

	listing.ApplySort()
	asserts.Equal(listing.Items[0].Name, "a00b00")
	asserts.Equal(listing.Items[1].Name, "a0b00")
	asserts.Equal(listing.Items[2].Name, "a0B")
}
*/
/*
func TestStrComp(t *testing.T) {
	asserts := assert.New(t)
	s1 := "a0b00"
	s2 := "a0B"
	asserts.True(Less(s1, s2))
}
*/
