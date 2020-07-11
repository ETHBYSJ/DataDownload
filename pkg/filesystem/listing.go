package filesystem

import (
	"sort"
	"strconv"
	"strings"
)

type Listing struct {
	Items []*FileInfo		`json:"items"`
	NumDirs int				`json:"numDirs"`
	NumFiles int			`json:"numFiles"`
	Sorting Sorting			`json:"sorting"`
}


func (l Listing) ApplySort() {
	if !l.Sorting.Asc {
		switch l.Sorting.By {
		case ByNameDictionaryStyle:
			sort.Sort(sort.Reverse(byNameDictionary(l)))
		case ByNameNaturalStyle:
			sort.Sort(sort.Reverse(byNameNatural(l)))
		case BySize:
			sort.Sort(sort.Reverse(bySize(l)))
		case ByModified:
			sort.Sort(sort.Reverse(byModified(l)))
		default:
			sort.Sort(sort.Reverse(byNameDictionary(l)))
		}
	} else {
		switch l.Sorting.By {
		case ByNameDictionaryStyle:
			sort.Sort(byNameDictionary(l))
		case ByNameNaturalStyle:
			sort.Sort(byNameNatural(l))
		case BySize:
			sort.Sort(bySize(l))
		case ByModified:
			sort.Sort(byModified(l))
		default:
			sort.Sort(byNameDictionary(l))
			return
		}
	}
}

// 取数字位
func digits(s string) int {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return i
		}
	}
	return len(s)
}


// 寻找公共前缀，不包含数字
func commonPrefix(a, b string) int {
	m := len(a)
	if n := len(b); n < m {
		m = n
	}
	if m == 0 {
		return 0
	}
	_ = a[m-1]
	_ = b[m-1]
	for i := 0; i < m; i++ {
		ca := a[i]
		cb := b[i]
		if (ca >= '0' && ca <= '9') || (cb >= '0' && cb <= '9') || ca != cb {
			return i
		}
	}
	return m
}

func Less(a, b string) bool {
	for {
		if a == b {
			return false
		}
		if p := commonPrefix(a, b); p != 0 {
			a = a[p:]
			b = b[p:]
		}
		if ia := digits(a); ia > 0 {
			if ib := digits(b); ib > 0 {
				an, aerr := strconv.ParseUint(a[:ia], 10, 64)
				bn, berr := strconv.ParseUint(b[:ib], 10, 64)
				if aerr == nil && berr == nil {
					if an != bn {
						return an < bn
					}
					if ia != len(a) && ib != len(b) {
						a = a[ia:]
						b = b[ib:]
						continue
					}
				}
			}
		}
		return a < b
	}
}

type byNameDictionary Listing
type byNameNatural Listing
type bySize Listing
type byModified Listing


// 排序时，文件夹在前

// 通过文件名比较
func (l byNameDictionary) Len() int {
	return len(l.Items)
}

func (l byNameDictionary) Swap(i, j int) {
	l.Items[i], l.Items[j] = l.Items[j], l.Items[i]
}
// 字典序
func (l byNameDictionary) Less(i, j int) bool {
	if l.Items[i].IsDir && !l.Items[j].IsDir {
		return true
	}
	if !l.Items[i].IsDir && l.Items[j].IsDir {
		return false
	}
	return strings.ToLower(l.Items[i].Name) < strings.ToLower(l.Items[j].Name)
}

func (l byNameNatural) Len() int {
	return len(l.Items)
}

func (l byNameNatural) Swap(i, j int) {
	l.Items[i], l.Items[j] = l.Items[j], l.Items[i]
}

// 自然风格
func (l byNameNatural) Less(i, j int) bool {
	if l.Items[i].IsDir && !l.Items[j].IsDir {
		return true
	}
	if !l.Items[i].IsDir && l.Items[j].IsDir {
		return false
	}
	return Less(strings.ToLower(l.Items[i].Name), strings.ToLower(l.Items[j].Name))
	// return Less(l.Items[i].Name, l.Items[j].Name)
}


// 通过文件大小比较
func (l bySize) Len() int {
	return len(l.Items)
}

func (l bySize) Swap(i, j int) {
	l.Items[i], l.Items[j] = l.Items[j], l.Items[i]
}

const directoryOffset = -1 << 31
func (l bySize) Less(i, j int) bool {
	if l.Items[i].IsDir && !l.Items[j].IsDir {
		return true
	}
	if !l.Items[i].IsDir && l.Items[j].IsDir {
		return false
	}
	iSize, jSize := l.Items[i].Size, l.Items[j].Size
	if l.Items[i].IsDir {
		iSize = directoryOffset + iSize
	}
	if l.Items[j].IsDir {
		jSize = directoryOffset + jSize
	}
	return iSize < jSize
}


func (l byModified) Len() int {
	return len(l.Items)
}

func (l byModified) Swap(i, j int) {
	l.Items[i], l.Items[j] = l.Items[j], l.Items[i]
}

func (l byModified) Less(i, j int) bool {
	if l.Items[i].IsDir && !l.Items[j].IsDir {
		return true
	}
	if !l.Items[i].IsDir && l.Items[j].IsDir {
		return false
	}
	iModified, jModified := l.Items[i].ModTime, l.Items[j].ModTime
	return iModified.Sub(jModified) < 0
}










