package app

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

type FileManager struct {
	FilesCursor int
	Files       []FileInfo
	Cache       map[string][]FileInfo
}

func (fm *FileManager) FetchFiles(dm *DirectoryManager) []FileInfo {
	var files []FileInfo
	companyFolderPath := dm.CurrentFolderPath()
	categoryPath := strings.ToLower(dm.SelectedCategory)

	path := "/Users/eytananjel/Notes/" + companyFolderPath + "/" + categoryPath

	cacheKey := companyFolderPath + ":" + categoryPath
	cached, ok := fm.Cache[cacheKey]

	if !ok {
		files = readFilesInDirecory(path)
		fm.Cache[cacheKey] = files
	} else {
		log.Info("Read from cache")
		files = cached
	}

	return files
}

func (fm *FileManager) GetCurrentFilePath(companyName string, categoryName string) string {
	return "/Users/eytananjel/Notes/" + companyName + "/" + categoryName + "/" + fm.currentFileName()
}

func (fm FileManager) CurrentFile() FileInfo {
	return fm.Files[fm.FilesCursor]
}

func (fm FileManager) currentFileName() string {
	return fm.CurrentFile().Name
}

func (fm FileManager) CurrentFileContent() string {
	return fm.CurrentFile().Content
}

func (fm FileManager) currentFile() FileInfo {
	return fm.Files[fm.FilesCursor]
}

func readFilesInDirecory(path string) []FileInfo {
	var fileInfos []FileInfo

	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		fullPath := filepath.Join(path, file.Name())
		content, err := os.ReadFile(fullPath)
		if err != nil {
			log.Fatal(err)
		}

		newFileInfo := FileInfo{
			Name:    file.Name(),
			Content: string(content),
		}

		fileInfos = append(fileInfos, newFileInfo)
	}

	return fileInfos
}
