package files

import "time"

type FileInfo struct {
	PATH     string     `json:"path,omitempty"`
	Name     string     `json:"name"`
	Size     int64      `json:"size"`
	FileType string     `json:"fileType,omitempty"`
	Files    []FileInfo `json:"files,omitempty"`
	IsDir    bool       `json:"isDir"`
	ModTime  time.Time  `json:"modTime"`
}
