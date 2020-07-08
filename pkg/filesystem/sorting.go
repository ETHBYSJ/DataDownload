package filesystem

var (
	// 字典序比较
	ByNameDictionaryStyle = "dic"
	// win10风格
	ByNameNaturalStyle = "win10"
	BySize = "size"
	ByModified = "modified"
)

type Sorting struct {
	By 	string
	Asc bool
}