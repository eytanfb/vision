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
		case "priority_toggled":
			err = m.TaskManager.UpdateTaskToPriority(&m.FileManager, task)
		}

		return TaskUpdatedMsg{
			Task:   task,
			Action: action,
			Err:    err,
		}
	}
}

// refreshTasksCmd refreshes the task list
func (m *Model) refreshTasksCmd() tea.Cmd {
	return func() tea.Msg {
		m.FileManager.FetchTasks(&m.DirectoryManager, &m.TaskManager)
		return TasksRefreshedMsg{}
	}
}

// createTaskCmd creates a new task
func (m *Model) createTaskCmd(company string, taskName string) tea.Cmd {
	return func() tea.Msg {
		err := m.FileManager.CreateTask(company, taskName)
		return TaskCreatedMsg{
			TaskName: taskName,
			Err:      err,
		}
	}
}

// createSubTaskCmd creates a new subtask
func (m *Model) createSubTaskCmd(parentFile FileInfo, subtaskName string) tea.Cmd {
	return func() tea.Msg {
		err := m.FileManager.CreateSubTask(m.GetCurrentCompanyName(), parentFile, subtaskName)
		return SubTaskCreatedMsg{
			ParentTask: Task{}, // Empty task for now
			SubTask:    subtaskName,
			Err:        err,
		}
	}
}

// File Operation Commands

// loadFileCmd loads a file's content
func (m *Model) loadFileCmd(filename string) tea.Cmd {
	return func() tea.Msg {
		// In the actual implementation, this would load file content
		// For now, return a message indicating the file was selected
		return FileLoadedMsg{
			File:    FileInfo{Name: filename},
			Content: "",
			Err:     nil,
		}
	}
}

// createStandupCmd creates a new standup file
func (m *Model) createStandupCmd(company string) tea.Cmd {
	return func() tea.Msg {
		err := m.FileManager.CreateStandup(company)
		return FileCreatedMsg{
			Filename: "standup",
			Err:      err,
		}
	}
}

// refreshFilesCmd refreshes the file list
func (m *Model) refreshFilesCmd() tea.Cmd {
	return func() tea.Msg {
		// Files are refreshed by FetchTasks or LoadFile operations
		return FilesRefreshedMsg{
			Files: m.FileManager.Files,
		}
	}
}

// Clipboard Commands

// copyToClipboardCmd copies content to clipboard
func (m *Model) copyToClipboardCmd(content string) tea.Cmd {
	return func() tea.Msg {
		// The actual clipboard operation will be done in the message handler
		return ClipboardCopiedMsg{
			Content: content,
			Err:     nil,
		}
	}
}

// View Navigation Commands

// changeViewCmd changes the current view
func (m *Model) changeViewCmd(from string, to string) tea.Cmd {
	return func() tea.Msg {
		return ViewChangedMsg{
			From: from,
			To:   to,
		}
	}
}

// selectCompanyCmd selects a company
func (m *Model) selectCompanyCmd(company Company) tea.Cmd {
	return func() tea.Msg {
		return CompanySelectedMsg{
			Company: company,
		}
	}
}

// selectCategoryCmd selects a category
func (m *Model) selectCategoryCmd(category string) tea.Cmd {
	return func() tea.Msg {
		return CategorySelectedMsg{
			Category: category,
		}
	}
}

// toggleSidebarCmd toggles the sidebar
func (m *Model) toggleSidebarCmd() tea.Cmd {
	return func() tea.Msg {
		return SidebarToggledMsg{
			Hidden: !m.ViewManager.HideSidebar,
		}
	}
}

// Input Mode Commands

// enterFilterModeCmd enters filter mode
func (m *Model) enterFilterModeCmd() tea.Cmd {
	return sendMsg(FilterModeEnteredMsg{})
}

// exitFilterModeCmd exits filter mode
func (m *Model) exitFilterModeCmd() tea.Cmd {
	return sendMsg(FilterModeExitedMsg{})
}

// enterAddTaskModeCmd enters add task mode
func (m *Model) enterAddTaskModeCmd() tea.Cmd {
	return sendMsg(AddTaskModeEnteredMsg{})
}

// exitAddTaskModeCmd exits add task mode
func (m *Model) exitAddTaskModeCmd() tea.Cmd {
	return sendMsg(AddTaskModeExitedMsg{})
}

// enterAddSubTaskModeCmd enters add subtask mode
func (m *Model) enterAddSubTaskModeCmd() tea.Cmd {
	return sendMsg(AddSubTaskModeEnteredMsg{})
}

// exitAddSubTaskModeCmd exits add subtask mode
func (m *Model) exitAddSubTaskModeCmd() tea.Cmd {
	return sendMsg(AddSubTaskModeExitedMsg{})
}

// Error Commands

// errorCmd sends an error message
func (m *Model) errorCmd(err error, context string) tea.Cmd {
	return func() tea.Msg {
		return ErrorOccurredMsg{
			Err:     err,
			Context: context,
		}
	}
}

// Batch command helper - executes multiple commands
func batch(cmds ...tea.Cmd) tea.Cmd {
	return tea.Batch(cmds...)
}
