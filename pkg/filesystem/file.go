package filesystem

import "time"

type FileInfo struct {
	*Listing
	Name string
	Path string
	Size int64
	IsDir bool
	ModTime time.Time

}
