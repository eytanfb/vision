package app

import "time"

type ProgressInfo struct {
	Completed       int
	Total           int
	Percentage      int
	IsInactive      bool
	IsComplete      bool
	ActiveCount     int
	IncompleteCount int
	LastCalculated  time.Time
}

type FileInfo struct {
	Name            string
	Content         string
	UpdatedAt       time.Time
	FullPath        string
	DisplayName     string
	CachedProgress  *ProgressInfo
}

func (f *FileInfo) FileNameWithoutExtension(extension string) string {
	if f.DisplayName != "" {
		return f.DisplayName
	}
	if len(f.Name) < len(extension) {
		return f.Name
	}
	return f.Name[0 : len(f.Name)-len(extension)]
}
