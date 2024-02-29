package app

import (
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type FileManager struct {
	FilesCursor int
	Files       []FileInfo
	Cache       map[string][]FileInfo
	Errors      []string
}

func (fm *FileManager) FetchFiles(dm *DirectoryManager) []FileInfo {
	var files []FileInfo
	companyFolderPath := dm.CurrentFolderPath()
	categoryPath := strings.ToLower(dm.SelectedCategory)

	path := "/Users/eytananjel/Notes/" + companyFolderPath + "/" + categoryPath

	cacheKey := companyFolderPath + ":" + categoryPath
	cached, ok := fm.Cache[cacheKey]

	if !ok {
		sorting := "default"

		if categoryPath == "tasks" {
			sorting = "updatedAt"
		}

		files = readFilesInDirecory(path, sorting)
		if categoryPath == "standups" {
			lastStandup := files[0] // The first one is the most recent
			todayInFormat := time.Now().Format("2006-01-02")
			todayInFormat += ".md"

			if lastStandup.Name != todayInFormat {
				fm.CreateStandup(companyFolderPath)
				files = readFilesInDirecory(path, sorting)
			}
		}

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

func (fm FileManager) CurrentFileContent() string {
	return fm.CurrentFile().Content
}

func (fm FileManager) CreateStandup(company string) {
	todayInFormat := time.Now().Format("2006-01-02")

	filePath := "/Users/eytananjel/Notes/" + company + "/standups/" + todayInFormat + ".md"
	templatePath := "/Users/eytananjel/Notes/obsidian/templates/" + company + "_standup.md"

	err := copyFile(templatePath, filePath)
	if err != nil {
		log.Fatal(err)
	}
}

// If it does exist, it will be overwritten.
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}

func (fm FileManager) currentFile() FileInfo {
	return fm.Files[fm.FilesCursor]
}

func (fm FileManager) currentFileName() string {
	return fm.CurrentFile().Name
}

func readFilesInDirecory(path string, sortBy string) []FileInfo {
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

		if strings.HasPrefix(file.Name(), "sortspec") {
			continue
		}

		fullPath := filepath.Join(path, file.Name())
		content, err := os.ReadFile(fullPath)
		if err != nil {
			log.Fatal(err)
		}

		fileInfo, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}

		newFileInfo := FileInfo{
			Name:      file.Name(),
			Content:   string(content),
			UpdatedAt: fileInfo.ModTime(),
		}

		fileInfos = append(fileInfos, newFileInfo)
	}

	if sortBy == "updatedAt" {
		slices.SortFunc(fileInfos, updatedAtCmp)
	} else {
		slices.SortFunc(fileInfos, nameCmp)
	}

	return fileInfos
}

func nameCmp(a, b FileInfo) int {
	if a.Name < b.Name {
		return 1
	}
	if a.Name > b.Name {
		return -1
	}
	return 0
}

func updatedAtCmp(a, b FileInfo) int {
	if a.UpdatedAt.Before(b.UpdatedAt) {
		return 1
	}
	if a.UpdatedAt.After(b.UpdatedAt) {
		return -1
	}
	return 0
}