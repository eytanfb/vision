package app

import "time"

type FileInfo struct {
	Name      string
	Content   string
	UpdatedAt time.Time
	FullPath  string
}
