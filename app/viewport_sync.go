package app

import "github.com/charmbracelet/bubbles/viewport"

type detailsSyncState struct {
	fileName string
	category string
	hash     uint64
	width    int
}

type filesSyncState struct {
	category string
	company  string
}

type tasksSyncState struct {
	fileName string
	cursor   int
	width    int
	height   int
}

func (m *Model) syncViewports() {
	if !m.ViewManager.Ready || !m.IsDetailsView() {
		return
	}

	m.resizeViewports()
	m.syncFilesViewport()

	if m.IsTaskDetailsFocus() {
		m.syncTasksViewport()
		return
	}

	m.syncDetailsViewport()
}

func (m *Model) resizeViewports() {
	m.Viewport.Width = m.ViewManager.DetailsViewWidth
	m.Viewport.Height = m.ViewManager.DetailsViewHeight
	m.TasksViewport.Width = m.ViewManager.DetailsViewWidth
	m.TasksViewport.Height = m.ViewManager.DetailsViewHeight
	m.FilesViewport.Width, m.FilesViewport.Height = m.ViewManager.sidebarContentSize()
}

func (m *Model) syncFilesViewport() {
	if m.ViewManager.HideSidebar {
		return
	}

	content, start, end := buildSidebarFilesContent(m, m.FilesViewport.Width)
	m.FilesViewport.SetContent(content)

	next := filesSyncState{category: m.DirectoryManager.SelectedCategory, company: m.DirectoryManager.SelectedCompany.DisplayName}
	if next != m.filesSync {
		m.filesSync = next
		m.FilesViewport.GotoTop()
	}

	if start < 0 {
		return
	}

	scrollIntoView(&m.FilesViewport, start, end)
}

func (m *Model) syncDetailsViewport() {
	m.FileManager.filesMu.RLock()
	files := m.FileManager.Files
	cursor := m.FileManager.FilesCursor
	var file FileInfo
	found := cursor >= 0 && cursor < len(files)
	if found {
		file = files[cursor]
	}
	m.FileManager.filesMu.RUnlock()

	next := detailsSyncState{width: m.Viewport.Width, category: m.DirectoryManager.SelectedCategory}
	content := ""
	if found {
		next.fileName = file.Name
		next.hash = hashString(file.Content)
		content = file.Content
	}

	if next == m.detailsSync {
		return
	}

	fileChanged := next.fileName != m.detailsSync.fileName || next.category != m.detailsSync.category
	m.detailsSync = next

	m.Viewport.SetContent(renderMarkdownWithWidth(content, max(m.Viewport.Width, 1)))
	if fileChanged {
		m.Viewport.GotoTop()
	}
}

func (m *Model) syncTasksViewport() {
	if m.ViewManager.IsAddTaskView || m.ViewManager.IsAddSubTaskView {
		return
	}

	fileName := m.FileManager.currentFileName()
	tasks := m.TaskManager.TaskCollection.GetTasks(fileName)
	cursor := m.TaskManager.TasksCursor

	content, start, end := buildTaskRowsForFileView(m, tasks, m.TaskManager.DailySummaryDate, cursor)
	m.TasksViewport.SetContent(content)

	next := tasksSyncState{fileName: fileName, cursor: cursor, width: m.TasksViewport.Width, height: m.TasksViewport.Height}
	if next == m.tasksSync {
		return
	}
	m.tasksSync = next

	if cursor < 0 || start < 0 {
		m.TasksViewport.GotoTop()
		return
	}

	scrollIntoView(&m.TasksViewport, start, end)
}

func scrollIntoView(vp *viewport.Model, start, end int) {
	height := vp.Height
	switch {
	case start < vp.YOffset:
		vp.SetYOffset(start)
	case end > vp.YOffset+height:
		vp.SetYOffset(min(end-height, start))
	}
}
