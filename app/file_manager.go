package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"
	"vision/mindmap"

	"github.com/charmbracelet/log"
)

type FileManager struct {
	FilesCursor            int
	Files                  []FileInfo
	FileCache              map[string][]FileInfo
	TaskCache              map[string]map[string][]Task
	SelectedFile           FileInfo
	TaskSuggestions        []string
	PeopleSuggestions      []string
	SuggestionsFilterValue string
	FileExtension          string
	Updater                mindmap.MindMapUpdaterInterface
}

func NewFileManager() *FileManager {
	return &FileManager{
		FileCache:     make(map[string][]FileInfo),
		TaskCache:     make(map[string]map[string][]Task),
		FileExtension: ".md",
		Updater:       mindmap.NewNullUpdater(),
	}
}

func (fm *FileManager) FetchFiles(dm *DirectoryManager, tm *TaskManager) []FileInfo {
	var files []FileInfo
	companyFolderPath := dm.CurrentFolderPath()
	categoryPath := strings.ToLower(dm.SelectedCategory)
	log.Info("Fetching files for category: " + categoryPath)

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

	if categoryPath == "standups" {
		lastStandup := files[0] // The first one is the most recent
		todayInFormat := time.Now().Format("2006-01-02")
		todayInFormat += fm.FileExtension

		if lastStandup.Name != todayInFormat && isWorkingDay() {
			fm.CreateStandup(companyFolderPath)
			files = readFilesInDirecory(path, sorting, tm)
		}
	}

	log.Info("Files count: " + fmt.Sprintf("%d", len(files)))
	fm.FetchTasks(dm, tm)
	fm.Files = files

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
		tasks := tm.ExtractTasks(dm.SelectedCompany.DisplayName, file.Name, file.Content)
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

	filePath := notesPath() + "/" + company + "/standups/" + todayInFormat + fm.FileExtension
	templatePath := notesPath() + "/obsidian/templates/" + company + "_standup.md"

	err := copyFile(templatePath, filePath)
	if err != nil {
		log.Fatal(err)
	}
}

func (fm FileManager) CreateTask(company string, taskName string) {
	filePath := notesPath() + "/" + company + "/tasks/" + taskName + fm.FileExtension
	templatePath := notesPath() + "/obsidian/templates/" + company + "_task.md"

	err := copyFile(templatePath, filePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Calling MindMapUpdater.AppendTask", "taskName", taskName)
	fm.Updater.AppendTask(taskName, taskName)
}

func (fm FileManager) CreateSubTask(company string, file FileInfo, taskName string) {
	filePath := filepath.Join(notesPath(), "/", company, "/tasks/", file.Name)

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.Contains(line, "### Sub-tasks") {
			lines = append(lines[:i+2], append([]string{"- [ ] " + taskName}, lines[i+2:]...)...)
			break
		}
	}

	newContent := strings.Join(lines, "\n")
	err = os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	parentTaskID := strings.TrimSuffix(file.Name, fm.FileExtension)
	log.Info("Calling MindMapUpdater.AppendSubtask", "parentTaskID", parentTaskID, "taskName", taskName)
	fm.Updater.AppendSubtask(parentTaskID, taskName, taskName)
}

func (fm *FileManager) ResetCache() {
	fm.FileCache = make(map[string][]FileInfo)
}

func (fm *FileManager) UpdateTask(task Task, status string) {
	filename := task.FileName // This is the parent task's filename
	text := task.Text         // This is the subtask text

	filePath := notesPath() + "/" + task.Company + "/tasks/" + filename
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(file), "\n")
	for i, line := range lines {
		if strings.Contains(line, text) {
			if status == "scheduled" {
				line := lines[i]
				regex := regexp.MustCompile(`🛫\s+\d{4}-\d{2}-\d{2}`)
				lines[i] = regex.ReplaceAllString(line, "")
				lines[i] = strings.ReplaceAll(lines[i], "- [x]", "- [ ]")

				if !strings.Contains(line, ScheduledIcon) {
					lines[i] = lines[i] + " " + ScheduledIcon + " " + time.Now().Format("2006-01-02")
				}

				break
			} else if status == "completed" {
				line = strings.ReplaceAll(line, "- [ ]", "- [x]")
				lines[i] = line + " " + CompletedIcon + " " + time.Now().Format("2006-01-02")

				break
			} else if status == "started" {
				line := lines[i]
				regex := regexp.MustCompile(`✅\s+\d{4}-\d{2}-\d{2}`)
				lines[i] = regex.ReplaceAllString(line, "")

				if !strings.Contains(line, StartedIcon) {
					lines[i] = lines[i] + " " + StartedIcon + " " + time.Now().Format("2006-01-02")
				}

				break
			} else if status == "unscheduled" {
				line := lines[i]
				regex := regexp.MustCompile(`⏳\s+\d{4}-\d{2}-\d{2}`)
				lines[i] = regex.ReplaceAllString(line, "")

				break
			}

			if status == "priority" {
				line := lines[i]
				regex := regexp.MustCompile(`\s+\d{4}-\d{2}-\d{2}`)
				lines[i] = regex.ReplaceAllString(line, "")

				if !strings.Contains(line, PriorityIcon) {
					checkboxRegex := regexp.MustCompile(`- \[[ x]\]`)
					if loc := checkboxRegex.FindStringIndex(line); loc != nil {
						prefix := line[:loc[1]]
						suffix := line[loc[1]:]
						lines[i] = prefix + PriorityIcon + suffix
					}
				}

				break
			} else if status == "unpriority" {
				line := lines[i]
				if strings.Contains(line, PriorityIcon) {
					lines[i] = strings.ReplaceAll(lines[i], PriorityIcon, "")
				}

				break
			}
		}
	}

	newContent := strings.Join(lines, "\n")
	err = os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	if status == "scheduled" || status == "started" || status == "completed" || status == "unscheduled" {
		// Clean the task text before sending to mind map updater
		cleanText := task.textWithoutDates()
		cleanText = strings.TrimSpace(cleanText)

		// Use the filename (without extension) as the parent task ID
		parentTaskID := strings.TrimSuffix(filename, fm.FileExtension)
		log.Info("Calling MindMapUpdater.AppendSubtask and UpdateSubtaskStatus",
			"parentTaskID", parentTaskID,
			"taskName", cleanText,
			"status", status)

		// First ensure the subtask exists under the parent
		fm.Updater.AppendSubtask(parentTaskID, cleanText, cleanText)
		// Then update its status
		fm.Updater.UpdateSubtaskStatus(parentTaskID, cleanText, status)
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

func (fm *FileManager) SelectFile(filename string) {
	for _, file := range fm.Files {
		if file.Name == filename {
			fm.SelectedFile = file
		}
	}
}

func (fm *FileManager) PeopleFilenames(dm *DirectoryManager, tm *TaskManager, filterValue string) []string {
	path := notesPath() + "/" + dm.CurrentFolderPath() + "/people"

	files := readFilesInDirecory(path, "default", tm)

	filenames := []string{}
	for _, file := range files {
		if filterValue == "" {
			filenames = append(filenames, file.Name)
		} else if strings.Contains(strings.ToLower(file.Name), strings.ToLower(filterValue)) {
			filenames = append(filenames, file.Name)
		}
	}

	fm.PeopleSuggestions = filenames
	return filenames
}

func (fm *FileManager) TaskFilenames(dm *DirectoryManager, tm *TaskManager, filterValue string) []string {
	path := notesPath() + "/" + dm.CurrentFolderPath() + "/tasks"

	files := readFilesInDirecory(path, "default", tm)

	filenames := []string{}
	for _, file := range files {
		if filterValue == "" {
			filenames = append(filenames, file.Name)
		} else if strings.Contains(strings.ToLower(file.Name), strings.ToLower(filterValue)) {
			filenames = append(filenames, file.Name)
		}
	}

	fm.TaskSuggestions = filenames
	return filenames
}

func (fm *FileManager) GetActiveSuggestionsList(suggestionsListCursor int) []string {
	if suggestionsListCursor == 0 {
		return fm.PeopleSuggestions
	} else if suggestionsListCursor == 1 {
		return fm.TaskSuggestions
	}

	return []string{}
}

func (fm *FileManager) GetActiveSuggestion(suggestionsListCursor int, suggestionCursor int) string {
	activeSuggestionsList := fm.GetActiveSuggestionsList(suggestionsListCursor)
	if suggestionCursor < 0 || suggestionCursor >= len(activeSuggestionsList) {
		return ""
	}

	return strings.Split(activeSuggestionsList[suggestionCursor], fm.FileExtension)[0]
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

		if !strings.HasSuffix(file.Name(), tm.FileExtension) {
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

		// Extract title from YAML frontmatter if it exists
		contentStr := string(content)
		title := extractTitleFromYAML(contentStr)

		// Create display name with title if available
		displayName := ""
		if title != "" {
			displayName = title
		}

		// Remove YAML frontmatter
		contentStr = removeYAMLFrontmatter(contentStr)

		fileInfo, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}

		newFileInfo := FileInfo{
			Name:        file.Name(),
			DisplayName: displayName,
			Content:     contentStr,
			UpdatedAt:   fileInfo.ModTime(),
			FullPath:    fullPath,
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

// extractTitleFromYAML extracts the title field from YAML frontmatter if it exists
func extractTitleFromYAML(content string) string {
	if !strings.HasPrefix(strings.TrimSpace(content), "---") {
		return ""
	}

	lines := strings.Split(content, "\n")
	if len(lines) < 3 {
		return ""
	}

	// Look for title field in YAML frontmatter
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "---" {
			// End of YAML frontmatter
			break
		}

		if strings.HasPrefix(line, "title:") {
			// Extract title value
			titleParts := strings.SplitN(line, ":", 2)
			if len(titleParts) == 2 {
				title := strings.TrimSpace(titleParts[1])
				title = strings.ReplaceAll(title, "[[", "")
				title = strings.ReplaceAll(title, "]]", "")
				return title
			}
		}
	}

	return ""
}

// removeYAMLFrontmatter removes YAML frontmatter from content if it exists
func removeYAMLFrontmatter(content string) string {
	// Check if content starts with "---" which indicates YAML frontmatter
	if strings.HasPrefix(strings.TrimSpace(content), "---") {
		// Split the content by lines
		lines := strings.Split(content, "\n")
		if len(lines) < 2 {
			return content
		}

		// Skip the first line (which is "---")
		foundClosingMarker := false
		endIndex := 0

		// Look for the closing "---" marker
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" {
				foundClosingMarker = true
				endIndex = i
				break
			}
		}

		// If we found a closing marker, return everything after it
		if foundClosingMarker && endIndex < len(lines)-1 {
			return strings.Join(lines[endIndex+1:], "\n")
		}
	}
	return content
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

		iCompleted := tm.TaskCollection.IsCompleted(iFilename)
		jCompleted := tm.TaskCollection.IsCompleted(jFilename)

		if iCompleted && jCompleted {
			iUpdatedAt := tm.TaskCollection.LastUpdatedAt(iFilename)
			jUpdatedAt := tm.TaskCollection.LastUpdatedAt(jFilename)

			return iUpdatedAt > jUpdatedAt
		}

		iInactive := tm.TaskCollection.IsInactive(iFilename)
		jInactive := tm.TaskCollection.IsInactive(jFilename)

		if iCompleted {
			return false
		}

		if jCompleted {
			return true
		}

		if iInactive {
			return false
		}

		if jInactive {
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
