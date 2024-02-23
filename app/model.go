package app

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"vision/config"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type Company struct {
	DisplayName    string   `json:"displayName"`
	FolderPathName string   `json:"folderPathName"`
	FullPath       string   `json:"fullPath"`
	SubFolders     []string `json:"subFolders"`
}

type Model struct {
	Companies        []Company
	Categories       []string
	CurrentView      string
	SelectedCompany  Company
	SelectedCategory string
	CompaniesCursor  int
	CategoriesCursor int
	FilesCursor      int
	Files            []FileInfo
	Cache            map[string][]FileInfo
	Width            int
	Height           int
	Viewport         viewport.Model
	ItemDetailsFocus bool
	Ready            bool
	Tasks            []Task
	TaskDetailsFocus bool
	TasksCursor      int
}

type Task struct {
	IsDone        bool
	Text          string
	StartDate     string
	CompletedDate string
	ScheduledDate string
	LineNumber    int
}

type ListItem struct {
	Title    string
	FullPath string
}

type FileInfo struct {
	Name    string
	Content string
}

func InitialModel(cfg *config.Config, args []string) tea.Model {
	items := []list.Item{} // Placeholder for list items initialization
	listModel := list.New(items, list.NewDefaultDelegate(), 0, 0)
	listModel.Title = "Select an Item"

	companies := convertCompanies(cfg.Companies)

	m := Model{
		Companies:        companies,
		Categories:       cfg.Categories,
		CurrentView:      "companies",
		CompaniesCursor:  0,
		CategoriesCursor: 0,
		FilesCursor:      0,
		SelectedCompany:  Company{},
		SelectedCategory: "",
		Files:            []FileInfo{},
		Cache:            make(map[string][]FileInfo),
		Viewport:         viewport.Model{},
		Ready:            false,
	}

	if len(args) > 0 {
		requestedCompany := args[0]
		for _, company := range m.Companies {
			if strings.ToLower(company.DisplayName) == strings.ToLower(requestedCompany) {
				m.SelectedCompany = company
				m.CurrentView = "categories"
				break
			}
		}
	}

	if len(args) > 1 {
		requestedCategory := strings.ToLower(args[1])
		for _, category := range m.Categories {
			if strings.ToLower(category) == requestedCategory {
				m.SelectedCategory = category
				m.CurrentView = "details"
				m.Files = m.FetchFiles()
				break
			}
		}
	}

	return &m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) FetchFiles() []FileInfo {
	var files []FileInfo
	path := "/Users/eytananjel/Notes/" + m.SelectedCompany.FolderPathName + "/" + strings.ToLower(m.SelectedCategory)

	cacheKey := m.SelectedCompany.FolderPathName + ":" + m.SelectedCategory
	cached, ok := m.Cache[cacheKey]

	if !ok {
		files = readFilesInDirecory(path)
		m.Cache[cacheKey] = files
	} else {
		log.Info("Read from cache")
		files = cached
	}

	return files
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

func convertCompanies(companies []config.Company) []Company {
	var items []Company
	for _, company := range companies {
		items = append(items, Company{
			DisplayName:    company.DisplayName,
			FolderPathName: company.FolderPathName,
			FullPath:       company.FullPath,
			SubFolders:     company.SubFolders,
		})
	}
	return items
}

func (t Task) String() string {
	var stringBuilder strings.Builder

	if t.IsDone {
		stringBuilder.WriteString("- [x] ")
	} else {
		stringBuilder.WriteString("- [ ] ")
	}

	stringBuilder.WriteString(t.Text)
	result := stringBuilder.String()
	resultWithoutDates := RemoveDatesFromText(result)
	stringBuilder.Reset()
	stringBuilder.WriteString(resultWithoutDates)

	if t.StartDate != "" || t.CompletedDate != "" || t.ScheduledDate != "" {
		stringBuilder.WriteString("\n")
		if t.ScheduledDate != "" {
			stringBuilder.WriteString("Scheduled: " + strings.Trim(t.ScheduledDate, " ") + "\n")
		}
		if t.StartDate != "" {
			stringBuilder.WriteString("Start: " + strings.Trim(t.StartDate, " ") + "\n")
		}
		if t.CompletedDate != "" {
			stringBuilder.WriteString("Completed: " + strings.Trim(t.CompletedDate, " ") + "\n")
		}
	}

	return stringBuilder.String()
}

func RemoveDatesFromText(text string) string {
	datesRegex := regexp.MustCompile(`[✅, ⏳, 🛫]\s+\d{4}-\d{2}-\d{2}`)

	text = datesRegex.ReplaceAllString(text, "")

	return strings.Trim(text, " ")
}
