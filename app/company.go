package app

type Company struct {
	DisplayName    string   `json:"displayName"`
	FolderPathName string   `json:"folderPathName"`
	FullPath       string   `json:"fullPath"`
	SubFolders     []string `json:"subFolders"`
}
