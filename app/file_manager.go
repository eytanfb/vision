package app

import (
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type FileManager struct {
	FilesCursor  int
	Files        []FileInfo
	FileCache    map[string][]FileInfo
	TaskCache    map[string]map[string][]Task
	SelectedFile FileInfo
}

func (fm *FileManager) FetchFiles(dm *DirectoryManager, tm *TaskManager) []FileInfo {
	log.Info("Fetching files")
	var files []FileInfo
	companyFolderPath := dm.CurrentFolderPath()
	categoryPath := strings.ToLower(dm.SelectedCategory)

	path := notesPath() + "/" + companyFolderPath + "/" + categoryPath
	log.Info("Path: " + path)

	sorting := "default"

	if categoryPath == "tasks" {
		sorting = "active"
	}

	files = readFilesInDirecory(path, sorting, tm)

	filenames := []string{}
	for _, file := range files {
		filenames = append(filenames, file.Name)
	}
	log.Info("Files: \n" + strings.Join(filenames, "\n"))

	if categoryPath == "standups" {
		lastStandup := files[0] // The first one is the most recent
		todayInFormat := time.Now().Format("2006-01-02")
		todayInFormat += ".md"

		if lastStandup.Name != todayInFormat && isWorkingDay() {
			fm.CreateStandup(companyFolderPath)
			files = readFilesInDirecory(path, sorting, tm)
		}
	}

	fm.FetchTasks(dm, tm)

	return files
}

func (fm *FileManager) FetchTasks(dm *DirectoryManager, tm *TaskManager) []Task {
	var tasks []Task
	log.Info("Fetching tasks")
	companyFolderPath := dm.CurrentFolderPath()

	path := notesPath() + "/" + companyFolderPath + "/tasks"
	log.Info("Path: " + path)

	files := readFilesInDirecory(path, "updatedAt", tm)
	for _, file := range files {
		tasks := tm.ExtractTasks(file.Name, file.Content)
		tm.TaskCollection.Add(file.Name, tasks)
	}

	tasks = tm.TaskCollection.GetAll()
	return tasks
}

func (fm *FileManager) GetCurrentFilePath(companyName string, categoryName string) string {
	return notesPath() + "/" + companyName + "/" + categoryName + "/" + fm.currentFileName()
}

func (fm FileManager) CurrentFile() FileInfo {
	return fm.SelectedFile
}

func (fm FileManager) CurrentFileContent() string {
	return fm.CurrentFile().Content
}

func (fm FileManager) CreateStandup(company string) {
	todayInFormat := time.Now().Format("2006-01-02")

	filePath := notesPath() + "/" + company + "/standups/" + todayInFormat + ".md"
	templatePath := notesPath() + "/obsidian/templates/" + company + "_standup.md"

	err := copyFile(templatePath, filePath)
	if err != nil {
		log.Fatal(err)
	}
}

func (fm FileManager) CreateTask(company string, taskName string) {
	filePath := notesPath() + "/" + company + "/tasks/" + taskName + ".md"
	templatePath := notesPath() + "/obsidian/templates/" + company + "_task.md"

	err := copyFile(templatePath, filePath)
	if err != nil {
		log.Fatal(err)
	}
}

func (fm *FileManager) ResetCache() {
	fm.FileCache = make(map[string][]FileInfo)
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
	return fm.SelectedFile
}

func (fm FileManager) currentFileName() string {
	return fm.CurrentFile().Name
}

func readFilesInDirecory(path string, sortBy string, tm *TaskManager) []FileInfo {
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
			FullPath:  fullPath,
		}

		fileInfos = append(fileInfos, newFileInfo)
	}

	if sortBy == "updatedAt" {
		slices.SortFunc(fileInfos, updatedAtCmp)
	} else if sortBy == "active" {
		fileInfos = sortedFiles(fileInfos, tm)
	} else {
		slices.SortFunc(fileInfos, nameCmp)
	}

	return fileInfos
}

func sortedFiles(fileInfos []FileInfo, tm *TaskManager) []FileInfo {
	filenames := []string{}
	for _, file := range fileInfos {
		filenames = append(filenames, file.Name)
	}

	activeCmp(filenames, tm)

	sortedFiles := []FileInfo{}
	for _, filename := range filenames {
		for _, file := range fileInfos {
			if file.Name == filename {
				sortedFiles = append(sortedFiles, file)
			}
		}
	}

	return sortedFiles
}

func activeCmp(filenames []string, tm *TaskManager) {
	sort.Slice(filenames, func(i, j int) bool {
		iFilename := filenames[i]
		jFilename := filenames[j]

		log.Info("Comparing " + iFilename + " with " + jFilename)

		iCompleted := tm.TaskCollection.IsCompleted(iFilename)
		jCompleted := tm.TaskCollection.IsCompleted(jFilename)

		iInactive := tm.TaskCollection.IsInactive(iFilename)
		jInactive := tm.TaskCollection.IsInactive(jFilename)

		if iCompleted {
			log.Info(iFilename + " is completed returning false")
			return false
		}

		if jCompleted {
			log.Info(jFilename + " is completed returning true")
			return true
		}

		if iInactive {
			log.Info(iFilename + " is inactive returning false")
			return false
		}

		if jInactive {
			log.Info(jFilename + " is inactive returning true")
			return true
		}

		iCompletedTasks, iTotalTasks := tm.TaskCollection.Progress(iFilename)
		jCompletedTasks, jTotalTasks := tm.TaskCollection.Progress(jFilename)

		iPercentage := float64(iCompletedTasks) / float64(iTotalTasks)
		jPercentage := float64(jCompletedTasks) / float64(jTotalTasks)

		iRoundedUp := int(iPercentage*10) * 10
		jRoundedUp := int(jPercentage*10) * 10

		if iRoundedUp == 100 {
			return false
		}

		if jRoundedUp == 100 {
			return true
		}

		if iRoundedUp == jRoundedUp {
			return iFilename < jFilename
		}

		return iRoundedUp > jRoundedUp
	})
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

func isWorkingDay() bool {
	return time.Now().Weekday() != time.Saturday && time.Now().Weekday() != time.Sunday
}

func notesPath() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir + "/Notes"
}
