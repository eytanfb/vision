package app

import (
	"os/exec"
	"strings"
	"time"
	"vision/config"
	"vision/mindmap"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type Model struct {
	DirectoryManager DirectoryManager
	TaskManager      TaskManager
	FileManager      FileManager
	ViewManager      ViewManager
	MindMapUpdater   mindmap.MindMapUpdaterInterface
	Viewport         viewport.Model
	NewTaskInput     textinput.Model
	FilterInput      textinput.Model
	Errors           []string
}

func InitialModel(cfg *config.Config, args []string) tea.Model {
	companies := CompaniesFromConfig(cfg.Companies)

	var defaultCompany Company
	for _, company := range companies {
		if strings.ToLower(company.DisplayName) == cfg.DefaultCompany {
			defaultCompany = company
		}
	}

	textInput := textinput.New()
	textInput.Placeholder = "Add a task..."
	filterInput := textinput.New()
	filterInput.Placeholder = "Filter... (/ to start)"

	monday := time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1).Format("2006-01-02")
	friday := time.Now().AddDate(0, 0, 5-int(time.Now().Weekday())).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")

	// Check if h-m-m exists in PATH
	var mindMapUpdater mindmap.MindMapUpdaterInterface
	_, err := exec.LookPath("h-m-m")
	if err != nil {
		log.Info("h-m-m not found in PATH, using NullMindMapUpdater")
		mindMapUpdater = mindmap.NewNullUpdater()
	} else {
		log.Info("h-m-m found in PATH, using MindMapUpdater")
		mindMapUpdater = mindmap.NewUpdater(notesPath() + "/personal/daily_mind_maps")
	}

	m := Model{
		DirectoryManager: DirectoryManager{
			Companies:        companies,
			Categories:       cfg.Categories,
			SelectedCompany:  defaultCompany,
			SelectedCategory: "tasks",
			CompaniesCursor:  0,
			CategoriesCursor: 0,
		},
		TaskManager: TaskManager{
			TaskCollection: TaskCollection{
				TasksByFile: make(map[string][]Task),
			},
			TasksCursor:            -1,
			WeeklySummaryStartDate: monday,
			WeeklySummaryEndDate:   friday,
			DailySummaryDate:       today,
			FileExtension:          cfg.PreferredFileExtension,
		},
		FileManager: FileManager{
			FilesCursor:       0,
			Files:             []FileInfo{},
			FileCache:         make(map[string][]FileInfo),
			TaskCache:         make(map[string]map[string][]Task),
			PeopleSuggestions: []string{},
			TaskSuggestions:   []string{},
			FileExtension:     cfg.PreferredFileExtension,
		},
		MindMapUpdater: mindMapUpdater,
		ViewManager: ViewManager{
			CurrentView:              CategoriesView,
			Width:                    0,
			Height:                   0,
			Ready:                    false,
			TaskDetailsFocus:         false,
			ItemDetailsFocus:         false,
			SidebarWidth:             40,
			SidebarHeight:            40,
			HideSidebar:              false,
			NavbarWidth:              40,
			NavbarHeight:             12,
			DetailsViewWidth:         40,
			IsAddTaskView:            false,
			IsAddSubTaskView:         false,
			IsWeeklyView:             false,
			ShowCompanies:            false,
			KanbanListCursor:         0,
			KanbanTaskCursor:         0,
			KanbanTasksCount:         0,
			KanbanViewLineDownFactor: 3,
			SuggestionsListsCursor:   -1,
			SuggestionCursor:         -1,
			IsSuggestionsActive:      false,
		},
		Viewport:     viewport.Model{},
		NewTaskInput: textInput,
		FilterInput:  filterInput,
	}

	// Initialize today's mind-map if using real updater
	if _, ok := mindMapUpdater.(*mindmap.NullMindMapUpdater); !ok {
		mindMapUpdater.InitializeDailyMindMap(time.Now())
	}

	// Provide updater reference to FileManager
	m.FileManager.Updater = mindMapUpdater

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

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) IsCompanyView() bool {
	return m.ViewManager.IsCompanyView()
}

func (m *Model) IsCategoryView() bool {
	return m.ViewManager.IsCategoryView()
}

func (m *Model) IsDetailsView() bool {
	return m.ViewManager.IsDetailsView()
}

func (m *Model) IsAddTaskView() bool {
	return m.ViewManager.IsAddTaskView
}

func (m *Model) IsAddSubTaskView() bool {
	return m.ViewManager.IsAddSubTaskView
}

func (m *Model) IsFilterView() bool {
	return m.ViewManager.IsFilterView
}

func (m *Model) IsTaskDetailsFocus() bool {
	return m.ViewManager.IsTaskDetailsFocus()
}

func (m *Model) IsItemDetailsFocus() bool {
	return m.ViewManager.IsItemDetailsFocus()
}

func (m *Model) IsKanbanView() bool {
	return m.IsCategoryView() && m.ViewManager.HideSidebar
}

func (m *Model) IsSuggestionsActive() bool {
	return m.ViewManager.IsSuggestionsActive
}

func (m *Model) GoToCompany(companyName string) {
	m.DirectoryManager.SelectCompany(companyName)
	m.FileManager.ResetCache()
	m.TaskManager.TaskCollection.Flush()
	m.FetchFiles()
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
	if !m.ViewManager.HideSidebar {
		goToNext(&m.FileManager.FilesCursor, len(m.FileManager.Files))
	}
}

func (m *Model) GoToNextKanbanTask() {
	goToNext(&m.ViewManager.KanbanTaskCursor, m.ViewManager.KanbanTasksCount)
}

func (m *Model) GoToPreviousKanbanTask() {
	goToPrevious(&m.ViewManager.KanbanTaskCursor)
}

func (m *Model) GoToNextKanbanList() {
	goToNext(&m.ViewManager.KanbanListCursor, 3)
}

func (m *Model) GoToPreviousKanbanList() {
	goToPrevious(&m.ViewManager.KanbanListCursor)
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

func (m *Model) SelectTask(task Task) {
	m.TaskManager.SelectTask(task)
}

func (m *Model) LoseDetailsFocus() {
	if !m.ViewManager.HideSidebar {
		m.ViewManager.ItemDetailsFocus = false
		m.ViewManager.TaskDetailsFocus = false
	}
}

func (m *Model) GainDetailsFocus() {
	m.ViewManager.ItemDetailsFocus = true
}

func (m *Model) ShowTasks() {
	m.ViewManager.TaskDetailsFocus = true
}

func (m *Model) GetCurrentCursor() int {
	if m.IsCompanyView() {
		return m.DirectoryManager.CompaniesCursor
	} else if m.IsCategoryView() {
		return m.DirectoryManager.CategoriesCursor
	} else if m.IsDetailsView() {
		return m.FileManager.FilesCursor
	}
	return 0
}

func (m *Model) HasFiles() bool {
	return len(m.FileManager.Files) > 0
}

func (m *Model) GetCurrentCompanyName() string {
	return m.DirectoryManager.CurrentCompanyName()
}

func (m *Model) GetCurrentFilePath() string {
	return m.FileManager.GetCurrentFilePath(m.DirectoryManager.CurrentCompanyName(), m.DirectoryManager.SelectedCategory)
}

func (m *Model) CompanyNames() []string {
	return m.DirectoryManager.CompanyNames()
}

func (m *Model) CategoryNames() []string {
	return m.DirectoryManager.Categories
}

func (m *Model) FetchFiles() []FileInfo {
	return m.FileManager.FetchFiles(&m.DirectoryManager, &m.TaskManager)
}
