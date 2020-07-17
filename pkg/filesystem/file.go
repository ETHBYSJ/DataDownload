package filesystem

import "time"

type FileInfo struct {
	*Listing
	ID 		uint			`json:"id"`
	Name 	string			`json:"name"`
	Path 	string			`json:"path"`
	Size 	int64			`json:"size"`
	IsDir 	bool			`json:"isDir"`
	ModTime time.Time		`json:"modified"`
	Review  bool 			`json:"review"`
	OwnerID uint 			`json:"ownerId"`
	Share 	bool 			`json:"share"`
}
