package app

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"vision/config"
	"vision/utils"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

const (
	CompaniesView  = "companies"
	CategoriesView = "categories"
	DetailsView    = "details"
)

type Model struct {
	Companies        []Company
	Categories       []string
	currentView      string
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
		currentView:      "companies",
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
				m.GoToNextView()
				break
			}
		}
	}

	if len(args) > 1 {
		requestedCategory := strings.ToLower(args[1])
		for _, category := range m.Categories {
			if strings.ToLower(category) == requestedCategory {
				m.assignCategory(category)
				m.GoToNextView()
				break
			}
		}
	}

	return &m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) IsCompanyView() bool {
	return m.currentView == CompaniesView
}

func (m Model) IsCategoryView() bool {
	return m.currentView == CategoriesView
}

func (m Model) IsDetailsView() bool {
	return m.currentView == DetailsView
}

func (m Model) IsTaskDetailsFocus() bool {
	return m.IsDetailsView() && m.TaskDetailsFocus
}

func (m Model) IsItemDetailsFocus() bool {
	return m.IsDetailsView() && m.ItemDetailsFocus
}

func (m *Model) GoToNextCompany() {
	goToNext(&m.CompaniesCursor, len(m.Companies))
}

func (m *Model) GoToNextCategory() {
	goToNext(&m.CategoriesCursor, len(m.Categories))
}

func (m *Model) GoToNextTask() {
	goToNext(&m.TasksCursor, len(m.Tasks))
}

func (m *Model) GoToNextFile() {
	goToNext(&m.FilesCursor, len(m.Files))
}

func (m *Model) GoToPreviousCompany() {
	goToPrevious(&m.CompaniesCursor)
}

func (m *Model) GoToPreviousCategory() {
	goToPrevious(&m.CategoriesCursor)
}

func (m *Model) GoToPreviousTask() {
	goToPrevious(&m.TasksCursor)
}

func (m *Model) GoToPreviousFile() {
	goToPrevious(&m.FilesCursor)
}

func goToNext(cursor *int, length int) {
	*cursor++
	if *cursor >= length {
		*cursor = length - 1
	}
}

func goToPrevious(cursor *int) {
	*cursor--
	if *cursor < 0 {
		*cursor = 0
	}
}

func (m *Model) GoToNextView() {
	if m.IsCompanyView() {
		m.currentView = CategoriesView
	} else if m.IsCategoryView() {
		m.currentView = DetailsView
		m.FilesCursor = 0
		m.Files = m.FetchFiles()
	}
}

func (m *Model) GoToNextViewWithCategory(category string) {
	m.assignCategory(category)
	m.GoToNextView()
}

func (m *Model) GoToPreviousView() {
	if m.IsCategoryView() {
		m.currentView = CompaniesView
	} else if m.IsDetailsView() {
		m.currentView = CategoriesView
	}
}

func (m *Model) selectCompany() {
	m.SelectedCompany = m.Companies[m.CompaniesCursor]
}

func (m *Model) selectCategory() {
	m.SelectedCategory = m.Categories[m.CategoriesCursor]
}

func (m *Model) Select() {
	if m.IsCompanyView() {
		m.selectCompany()
	} else if m.IsCategoryView() {
		m.selectCategory()
	}
	m.GoToNextView()
}

func (m *Model) MoveDown() {
	if m.IsCompanyView() {
		m.GoToNextCompany()
	} else if m.IsCategoryView() {
		m.GoToNextCategory()
	} else if m.IsDetailsView() {
		if m.IsItemDetailsFocus() {
			if m.IsTaskDetailsFocus() {
				m.GoToNextTask()
			} else {
				m.Viewport.LineDown(10)
			}
		} else {
			m.GoToNextFile()
		}
	}
}

func (m *Model) MoveUp() {
	if m.IsDetailsView() {
		if m.IsItemDetailsFocus() {
			if m.IsTaskDetailsFocus() {
				m.GoToPreviousTask()
			} else {
				m.Viewport.LineUp(10)
			}
		} else {
			m.GoToPreviousFile()
		}
	} else if m.IsCategoryView() {
		m.GoToPreviousCategory()
	} else if m.IsCompanyView() {
		m.GoToPreviousCompany()
	}
}

func (m *Model) assignCategory(category string) {
	m.SelectedCategory = category
}

func (m *Model) LoseDetailsFocus() {
	m.ItemDetailsFocus = false
	m.TaskDetailsFocus = false
}

func (m *Model) GainDetailsFocus() {
	m.ItemDetailsFocus = true
}

func (m *Model) ShowTasks() {
	fileTasks := utils.ExtractTasksFromText(m.Files[m.FilesCursor].Content)
	tasks := []Task{}
	for _, task := range fileTasks {
		tasks = append(tasks, Task{
			IsDone:        task.IsDone,
			Text:          task.Text,
			StartDate:     ExtractStartDateFromText(task.Text),
			ScheduledDate: ExtractScheduledDateFromText(task.Text),
			CompletedDate: ExtractCompletedDateFromText(task.Text),
			LineNumber:    task.LineNumber,
		})
	}
	m.Tasks = tasks
	m.TaskDetailsFocus = true
	m.TaskDetailsFocus = true
}

func (m Model) GetCurrentCursor() int {
	if m.IsCompanyView() {
		return m.CompaniesCursor
	} else if m.IsCategoryView() {
		return m.CategoriesCursor
	} else if m.IsDetailsView() {
		return m.FilesCursor
	}
	return 0
}

func (m Model) HasFiles() bool {
	return len(m.Files) > 0
}

func (m Model) GetCurrentCompanyName() string {
	return m.SelectedCompany.DisplayName
}

func (m Model) CompanyNames() []string {
	var names []string
	for _, company := range m.Companies {
		names = append(names, company.DisplayName)
	}
	return names
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

func RemoveDatesFromText(text string) string {
	datesRegex := regexp.MustCompile(`[✅, ⏳, 🛫]\s+\d{4}-\d{2}-\d{2}`)

	text = datesRegex.ReplaceAllString(text, "")

	return strings.Trim(text, " ")
}

func ExtractStartDateFromText(text string) string {
	startIcon := "🛫 "
	return ExtractDateFromText(text, startIcon)
}

func ExtractScheduledDateFromText(text string) string {
	scheduledIcon := "⏳"
	return ExtractDateFromText(text, scheduledIcon)
}

func ExtractCompletedDateFromText(text string) string {
	completedIcon := "✅ "
	return ExtractDateFromText(text, completedIcon)
}

func ExtractDateFromText(text string, icon string) string {
	index := strings.Index(text, icon)
	if index == -1 {
		return ""
	}
	// read date from the next 10 characters
	date := text[index : index+14]
	return date
}
