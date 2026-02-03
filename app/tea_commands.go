package app

import (
	"sort"

	tea "github.com/charmbracelet/bubbletea"
)

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
	// Use same sorted keys as view to ensure consistent counting
	keys := m.getSortedKanbanKeys(allTasksMap)
	count := 0

	for _, key := range keys {
		tasks := allTasksMap[key]
		for _, task := range tasks {
			// Match EXACTLY the same logic as makeKanbanLists in view_builder.go
			// The view checks IsScheduledForFuture FIRST and puts those in inactive
			// regardless of completion status
			isFuture := task.IsScheduledForFuture(m.TaskManager.DailySummaryDate)

			if listCursor == 0 {
				// Inactive: future scheduled, OR (not completed AND not started AND not scheduled)
				if isFuture {
					count++
				} else if !task.Completed && !task.Started && !task.Scheduled {
					count++
				}
			} else if listCursor == 1 {
				// Active: not future, not completed, and (started OR scheduled)
				if !isFuture && !task.Completed && (task.Started || task.Scheduled) {
					count++
				}
			} else if listCursor == 2 {
				// Completed: not future, and completed
				if !isFuture && task.Completed {
					count++
				}
			}
		}
	}
	return count
}

// getSortedKanbanKeys returns task keys sorted the same way as the kanban view
// This ensures cursor position maps to the same task in both view and selection logic
func (m *Model) getSortedKanbanKeys(tasksByFile map[string][]Task) []string {
	keys := make([]string, 0, len(tasksByFile))
	for k := range tasksByFile {
		keys = append(keys, k)
	}

	// Create a map of filename to FileInfo for fast lookup
	m.FileManager.filesMu.RLock()
	fileMap := make(map[string]*FileInfo)
	for i := range m.FileManager.Files {
		fileMap[m.FileManager.Files[i].Name] = &m.FileManager.Files[i]
	}
	m.FileManager.filesMu.RUnlock()

	// Sort using the same logic as viewSort in view.go
	type sortKey struct {
		filename   string
		percentage int
		isInactive bool
		isComplete bool
	}

	sortKeys := make([]sortKey, len(keys))
	for i, filename := range keys {
		var percentage int
		var isInactive, isComplete bool

		if file, ok := fileMap[filename]; ok {
			progress := m.FileManager.getOrCalculateProgress(file, &m.TaskManager)
			percentage = progress.Percentage
			isInactive = progress.IsInactive
			isComplete = progress.IsComplete
		} else {
			completed, total := m.TaskManager.TaskCollection.Progress(filename)
			if total > 0 {
				percentage = int(float64(completed) / float64(total) * 100)
				isComplete = completed == total
			}
			isInactive = m.TaskManager.TaskCollection.IsInactive(filename)
		}

		sortKeys[i] = sortKey{
			filename:   filename,
			percentage: percentage,
			isInactive: isInactive,
			isComplete: isComplete,
		}
	}

	sort.Slice(sortKeys, func(i, j int) bool {
		// Inactive files go last
		if sortKeys[i].isInactive != sortKeys[j].isInactive {
			return !sortKeys[i].isInactive
		}
		// Complete files go last
		if sortKeys[i].isComplete != sortKeys[j].isComplete {
			return !sortKeys[i].isComplete
		}
		// Sort by percentage descending
		if sortKeys[i].percentage != sortKeys[j].percentage {
			return sortKeys[i].percentage > sortKeys[j].percentage
		}
		// Alphabetical tiebreaker
		return sortKeys[i].filename < sortKeys[j].filename
	})

	// Extract sorted filenames
	result := make([]string, len(sortKeys))
	for i := range sortKeys {
		result[i] = sortKeys[i].filename
	}
	return result
}
