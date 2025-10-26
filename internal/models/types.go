package models

// FileInfo represents a file or directory in the listing
type FileInfo struct {
	Name  string
	IsDir bool
	Path  string
	Size  int64
}

// Directory represents a root directory being served
type Directory struct {
	Name string
	Path string
}

// PageData represents the data passed to the directory listing template
type PageData struct {
	CurrentPath string
	Files       []FileInfo
	Directories []Directory
	IsRoot      bool
}
