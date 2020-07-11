package filesystem

// 上下文相关
type key int
const (
	// Gin的上下文
	GinCtx key = iota
)
