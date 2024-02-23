package keystore

import "io/fs"

const (
	DirMode  fs.FileMode = 0700 // 0770
	FileMode fs.FileMode = 0600 // 0640
)
