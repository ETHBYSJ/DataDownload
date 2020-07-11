package filesystem

var (
	// 字典序比较
	ByNameDictionaryStyle = "dic"
	// 自然风格
	ByNameNaturalStyle = "natural"
	BySize = "size"
	ByModified = "modified"
)

type Sorting struct {
	By 	string	`json:"by"`
	Asc bool	`json:"asc"`
}