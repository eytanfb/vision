package app

import "time"

type FileInfo struct {
	Name      string
	Content   string
	UpdatedAt time.Time
	FullPath  string
}

func (f *FileInfo) FileNameWithoutExtension(extension string) string {
	return f.Name[0 : len(f.Name)-len(extension)]
}
