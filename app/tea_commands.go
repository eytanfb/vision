package app

import tea "github.com/charmbracelet/bubbletea"

// tea_commands.go contains Bubble Tea command generators.
// These functions return tea.Cmd that perform operations and send messages back.

// Task Operation Commands

// updateTaskCmd updates a task and returns a message
func (m *Model) updateTaskCmd(task Task, action string) tea.Cmd {
	return func() tea.Msg {
		var err error

		switch action {
		case "completed":
			err = m.TaskManager.UpdateTaskToCompleted(&m.FileManager, task)
		case "scheduled":
			err = m.TaskManager.UpdateTaskToScheduled(&m.FileManager, task)
		case "started":
			err = m.TaskManager.UpdateTaskToStarted(&m.FileManager, task)
		case "priority":
			err = m.TaskManager.UpdateTaskToPriority(&m.FileManager, task)
		case "unpriority":
			err = m.TaskManager.UpdateTaskToUnpriority(&m.FileManager, task)
		case "unscheduled":
			err = m.TaskManager.UpdateTaskToUnscheduled(&m.FileManager, task)
		}

		return TaskUpdatedMsg{
			Task:   task,
			Action: action,
			Err:    err,
		}
	}
}

// loadTasksCmd loads tasks asynchronously
func (m *Model) loadTasksCmd() tea.Cmd {
	return func() tea.Msg {
		tasks := m.FileManager.FetchTasks(&m.DirectoryManager, &m.TaskManager)
		return TasksLoadedMsg{
			Tasks: tasks,
			Err:   nil,
		}
	}
}

// File Operation Commands

// loadFilesCmd loads files asynchronously
func (m *Model) loadFilesCmd() tea.Cmd {
	return func() tea.Msg {
		files := m.FileManager.FetchFiles(&m.DirectoryManager, &m.TaskManager)
		return FilesLoadedMsg{
			Files: files,
			Err:   nil,
		}
	}
}

// Navigation Commands

// navigateFileCmd moves cursor and sends file navigation message
func (m *Model) navigateFileCmd(direction string) tea.Cmd {
	return func() tea.Msg {
		m.FileManager.filesMu.RLock()
		filesLen := len(m.FileManager.Files)
		m.FileManager.filesMu.RUnlock()

		if filesLen == 0 {
			return nil
		}

		var newCursor int
		if direction == "down" {
			newCursor = m.FileManager.FilesCursor + 1
			if newCursor >= filesLen {
				newCursor = filesLen - 1
			}
		} else {
			newCursor = m.FileManager.FilesCursor - 1
			if newCursor < 0 {
				newCursor = 0
			}
		}

		m.FileManager.filesMu.RLock()
		file := m.FileManager.Files[newCursor]
		m.FileManager.filesMu.RUnlock()

		return FileNavigatedMsg{
			File:      file,
			Direction: direction,
			NewCursor: newCursor,
		}
	}
}

// navigateKanbanCmd moves within kanban and sends navigation message
func (m *Model) navigateKanbanCmd(direction string) tea.Cmd {
	return func() tea.Msg {
		newListCursor := m.ViewManager.KanbanListCursor
		newTaskCursor := m.ViewManager.KanbanTaskCursor

		switch direction {
		case "down":
			newTaskCursor++
		case "up":
			newTaskCursor--
			if newTaskCursor < 0 {
				newTaskCursor = 0
			}
		case "left":
			newListCursor--
			if newListCursor < 0 {
				newListCursor = 0
			}
			newTaskCursor = 0
		case "right":
			newListCursor++
			if newListCursor > 2 {
				newListCursor = 2
			}
			newTaskCursor = 0
		}

		return KanbanNavigatedMsg{
			ListCursor: newListCursor,
			TaskCursor: newTaskCursor,
			Direction:  direction,
		}
	}
}
