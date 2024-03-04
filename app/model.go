package app

import (
	"strings"
	"vision/config"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	CompaniesView  = "companies"
	CategoriesView = "categories"
	DetailsView    = "details"
)

type Model struct {
	DirectoryManager DirectoryManager
	TaskManager      TaskManager
	FileManager      FileManager
	ViewManager      ViewManager
	Viewport         viewport.Model
	Errors           []string
}

func InitialModel(cfg *config.Config, args []string) tea.Model {
	companies := CompaniesFromConfig(cfg.Companies)

	var clerky Company
	for _, company := range companies {
		if strings.ToLower(company.DisplayName) == "clerky" {
			clerky = company
		}
	}

	m := Model{
		DirectoryManager: DirectoryManager{
			Companies:        companies,
			Categories:       cfg.Categories,
			SelectedCompany:  clerky,
			SelectedCategory: "",
			CompaniesCursor:  0,
			CategoriesCursor: 0,
		},
		TaskManager: TaskManager{
			TaskCollection: TaskCollection{
				TasksByFile: make(map[string][]Task),
			},
			TasksCursor: 0,
		},
		FileManager: FileManager{
			FilesCursor: 0,
			Files:       []FileInfo{},
			FileCache:   make(map[string][]FileInfo),
			TaskCache:   make(map[string][]Task),
		},
		ViewManager: ViewManager{
			CurrentView:      CategoriesView,
			Width:            0,
			Height:           0,
			Ready:            false,
			TaskDetailsFocus: false,
			ItemDetailsFocus: false,
		},
		Viewport: viewport.Model{},
	}

	SetArgs(&m, args)
	m.FetchFiles()

	return &m
}

func SetArgs(m *Model, args []string) {
	if len(args) > 0 {
		requestedCompany := args[0]
		found := m.DirectoryManager.SelectCompany(requestedCompany)

		if !found {
			return
		}

		m.GoToNextView()
	}

	if len(args) > 1 {
		requestedCategory := strings.ToLower(args[1])
		found := m.DirectoryManager.SelectCategory(requestedCategory)

		if !found {
			return
		}

		m.GoToNextView()
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) IsCompanyView() bool {
	return m.ViewManager.IsCompanyView()
}

func (m Model) IsCategoryView() bool {
	return m.ViewManager.IsCategoryView()
}

func (m Model) IsDetailsView() bool {
	return m.ViewManager.IsDetailsView()
}

func (m Model) IsTaskDetailsFocus() bool {
	return m.ViewManager.IsTaskDetailsFocus()
}

func (m Model) IsItemDetailsFocus() bool {
	return m.ViewManager.IsItemDetailsFocus()
}

func (m *Model) GoToNextCompany() {
	if m.DirectoryManager.CompaniesCursor == len(m.DirectoryManager.Companies)-1 {
		m.DirectoryManager.CompaniesCursor = 0
	} else {
		m.DirectoryManager.CompaniesCursor++
	}
	m.DirectoryManager.SelectedCompany = m.DirectoryManager.Companies[m.DirectoryManager.CompaniesCursor]
	m.FileManager.ResetCache()
	m.TaskManager.TaskCollection.Flush()
	m.FetchFiles()
}

func (m *Model) GoToNextCategory() {
	goToNext(&m.DirectoryManager.CategoriesCursor, len(m.DirectoryManager.Categories))
}

func (m *Model) GoToNextTask() {
	goToNext(&m.TaskManager.TasksCursor, m.TaskManager.TaskCollection.Size(m.FileManager.currentFileName()))
}

func (m *Model) GoToNextFile() {
	goToNext(&m.FileManager.FilesCursor, len(m.FileManager.Files))
}

func (m *Model) GoToPreviousCompany() {
	goToPrevious(&m.DirectoryManager.CompaniesCursor)
}

func (m *Model) GoToPreviousCategory() {
	goToPrevious(&m.DirectoryManager.CategoriesCursor)
}

func (m *Model) GoToPreviousTask() {
	goToPrevious(&m.TaskManager.TasksCursor)
}

func (m *Model) GoToPreviousFile() {
	goToPrevious(&m.FileManager.FilesCursor)
}

func (m *Model) GotoCategoryView() {
	m.ViewManager.CurrentView = CategoriesView
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
	m.ViewManager.GoToNextView(&m.FileManager, &m.DirectoryManager, &m.TaskManager)
}

func (m *Model) GoToNextViewWithCategory(category string) {
	m.DirectoryManager.SelectCategory(category)
	m.GoToNextView()
}

func (m *Model) GoToPreviousView() {
	m.ViewManager.GoToPreviousView()
}

func (m *Model) Select() {
	m.ViewManager.Select(&m.FileManager, &m.DirectoryManager, &m.TaskManager)
}

func (m *Model) LoseDetailsFocus() {
	m.ViewManager.ItemDetailsFocus = false
	m.ViewManager.TaskDetailsFocus = false
}

func (m *Model) GainDetailsFocus() {
	m.ViewManager.ItemDetailsFocus = true
}

func (m *Model) ShowTasks() {
	//fileTasks := utils.ExtractTasksFromText(m.FileManager.CurrentFileContent())
	//taskCollection := CreateTaskCollectionFromFileTasks(fileTasks)
	//m.TaskManager.TaskCollection = taskCollection
	m.ViewManager.TaskDetailsFocus = true
	m.ViewManager.TaskDetailsFocus = true
}

func (m Model) GetCurrentCursor() int {
	if m.IsCompanyView() {
		return m.DirectoryManager.CompaniesCursor
	} else if m.IsCategoryView() {
		return m.DirectoryManager.CategoriesCursor
	} else if m.IsDetailsView() {
		return m.FileManager.FilesCursor
	}
	return 0
}

func (m Model) HasFiles() bool {
	return len(m.FileManager.Files) > 0
}

func (m Model) GetCurrentCompanyName() string {
	return m.DirectoryManager.CurrentCompanyName()
}

func (m Model) GetCurrentFilePath() string {
	return m.FileManager.GetCurrentFilePath(m.DirectoryManager.CurrentCompanyName(), m.DirectoryManager.SelectedCategory)
}

func (m Model) CompanyNames() []string {
	return m.DirectoryManager.CompanyNames()
}

func (m Model) CategoryNames() []string {
	return m.DirectoryManager.Categories
}

func (m Model) FetchFiles() []FileInfo {
	return m.FileManager.FetchFiles(&m.DirectoryManager, &m.TaskManager)
}
