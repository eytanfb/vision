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

		// Get the task count for the target list to enforce bounds
		taskCount := m.getKanbanListTaskCount(newListCursor)

		switch direction {
		case "down":
			newTaskCursor++
			// Upper bounds check - don't go past last task
			if taskCount > 0 && newTaskCursor >= taskCount {
				newTaskCursor = taskCount - 1
			}
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

// getKanbanListTaskCount returns the number of tasks in a specific kanban list
func (m *Model) getKanbanListTaskCount(listCursor int) int {
	allTasksMap := m.TaskManager.TaskCollection.TasksByFile
	count := 0

	for _, tasks := range allTasksMap {
		for _, task := range tasks {
			if listCursor == 0 {
				// Inactive: future scheduled or unscheduled
				if task.IsScheduledForFuture(m.TaskManager.DailySummaryDate) || task.ScheduledDate == "" {
					count++
				}
			} else if listCursor == 1 {
				// Active: started or scheduled (not completed, not future)
				if !task.Completed && !task.IsScheduledForFuture(m.TaskManager.DailySummaryDate) && (task.Started || task.Scheduled) {
					count++
				}
			} else if listCursor == 2 {
				// Completed
				if task.Completed {
					count++
				}
			}
		}
	}
	return count
}
