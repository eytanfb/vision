package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type Model struct {
	companies        []Company
	categories       []string
	currentView      string
	selectedCompany  Company
	selectedCategory string
	cursor           int
	files            []FileInfo
	cache            map[string][]FileInfo
	Width            int
	Height           int
	viewport         viewport.Model
	itemDetailsFocus bool
}

type ListItem struct {
	Title    string
	FullPath string
}

func InitialModel(cfg *Config, args []string) tea.Model {
	items := []list.Item{} // Placeholder for list items initialization
	listModel := list.New(items, list.NewDefaultDelegate(), 0, 0)
	listModel.Title = "Select an Item"

	m := Model{
		companies:        cfg.Companies,
		categories:       cfg.Categories,
		currentView:      "companies",
		cursor:           0,
		selectedCompany:  Company{},
		selectedCategory: "",
		files:            []FileInfo{},
		cache:            make(map[string][]FileInfo),
		viewport:         viewport.Model{},
	}

	if len(args) > 0 {
		requestedCompany := args[0]
		for _, company := range m.companies {
			if strings.ToLower(company.DisplayName) == strings.ToLower(requestedCompany) {
				m.selectedCompany = company
				m.currentView = "categories"
				break
			}
		}
	}

	if len(args) > 1 {
		requestedCategory := strings.ToLower(args[1])
		for _, category := range m.categories {
			if strings.ToLower(category) == requestedCategory {
				m.selectedCategory = category
				m.currentView = "details"
				m.files = m.FetchFiles()
				break
			}
		}
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

type FileInfo struct {
	Name    string
	Content string
}

func (m Model) FetchFiles() []FileInfo {
	var files []FileInfo
	path := "/Users/eytananjel/Notes/" + m.selectedCompany.FolderPathName + "/" + strings.ToLower(m.selectedCategory)

	cacheKey := m.selectedCompany.FolderPathName + ":" + m.selectedCategory
	cached, ok := m.cache[cacheKey]

	if !ok {
		files = readFilesInDirecory(path)
		m.cache[cacheKey] = files
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
