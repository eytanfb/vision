package app

import "time"

type FileInfo struct {
	Name      string
	Content   string
	UpdatedAt time.Time
	FullPath  string
	IsDir     bool
}

func (f *FileInfo) FileNameWithoutExtension() string {
	if f.IsDir {
		return f.Name
	}

	extension := ".md"
	return f.Name[0 : len(f.Name)-len(extension)]
}
