package app

import (
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/log"
)

// TaskOperations handles all task-related commands
type TaskOperations struct{}

// HandleKey routes task operation keys to their appropriate handlers
func (to TaskOperations) HandleKey(key string, m *Model) error {
	switch key {
	case "d":
		return to.CompleteTask(m)
	case "s":
		return to.ScheduleOrStartTask(m)
	case "p":
		return to.TogglePriority(m)
	case "D":
		return to.StartTaskOrCopyStandup(m)
	case "S":
		return to.ToggleScheduledUnscheduled(m)
	case "a":
		return to.AddTask(m)
	case "A":
		return to.AddSubTask(m)
	}
	return nil
}

// CompleteTask marks the current task as completed
func (to TaskOperations) CompleteTask(m *Model) error {
	if m.IsCategoryView() && m.ViewManager.HideSidebar {
		if err := m.TaskManager.UpdateTaskToCompleted(&m.FileManager, m.TaskManager.SelectedTask); err != nil {
			m.Errors = append(m.Errors, err.Error())
		}
		m.FileManager.FetchTasks(&m.DirectoryManager, &m.TaskManager)
	}
	return nil
}

// ScheduleOrStartTask schedules or starts a task depending on its current state
func (to TaskOperations) ScheduleOrStartTask(m *Model) error {
	if m.ViewManager.HideSidebar {
		log.Info("SKeyCommand: Show sidebar")
		if m.TaskManager.SelectedTask.Scheduled {
			log.Info("SKeyCommand: Update task to started")
			if err := m.TaskManager.UpdateTaskToStarted(&m.FileManager, m.TaskManager.SelectedTask); err != nil {
				m.Errors = append(m.Errors, err.Error())
			}
			m.ViewManager.KanbanListCursor = 1
			m.ViewManager.IsKanbanTaskUpdated = true
		} else {
			log.Info("SKeyCommand: Update task to scheduled")
			if err := m.TaskManager.UpdateTaskToScheduled(&m.FileManager, m.TaskManager.SelectedTask); err != nil {
				m.Errors = append(m.Errors, err.Error())
			}
			m.ViewManager.KanbanListCursor = 1
			m.ViewManager.IsKanbanTaskUpdated = true
		}
		m.FileManager.FetchFiles(&m.DirectoryManager, &m.TaskManager)
	}
	return nil
}

// TogglePriority toggles the priority marker on a task
func (to TaskOperations) TogglePriority(m *Model) error {
	if m.IsCategoryView() && m.ViewManager.HideSidebar {
		selectedTask := m.TaskManager.SelectedTask

		if !strings.Contains(selectedTask.Text, "🔺") {
			log.Info("Adding priority marker to task")
			if err := m.TaskManager.UpdateTaskToPriority(&m.FileManager, selectedTask); err != nil {
				m.Errors = append(m.Errors, err.Error())
			}
			m.FileManager.FetchTasks(&m.DirectoryManager, &m.TaskManager)
		} else {
			log.Info("Removing priority marker from task")
			if err := m.TaskManager.UpdateTaskToUnpriority(&m.FileManager, selectedTask); err != nil {
				m.Errors = append(m.Errors, err.Error())
			}
			m.FileManager.FetchTasks(&m.DirectoryManager, &m.TaskManager)
		}
	}
	return nil
}

// StartTaskOrCopyStandup starts a task in kanban view or copies standup in other views
func (to TaskOperations) StartTaskOrCopyStandup(m *Model) error {
	if m.IsCategoryView() && m.ViewManager.HideSidebar {
		if err := m.TaskManager.UpdateTaskToStarted(&m.FileManager, m.TaskManager.SelectedTask); err != nil {
			m.Errors = append(m.Errors, err.Error())
		}
		m.FileManager.FetchTasks(&m.DirectoryManager, &m.TaskManager)
		return nil
	}

	if !m.IsCategoryView() || m.ViewManager.IsWeeklyView {
		return nil
	}

	// Copy standup to clipboard
	slackMessage := m.TaskManager.SummaryForSlack(m.DirectoryManager.CurrentCompanyName())
	err := clipboard.WriteAll(slackMessage)
	if err != nil {
		log.Error("Failed to copy to clipboard", err)
	}

	return nil
}

// ToggleScheduledUnscheduled toggles between scheduled and unscheduled states
func (to TaskOperations) ToggleScheduledUnscheduled(m *Model) error {
	if m.IsCategoryView() && m.ViewManager.HideSidebar {
		if m.TaskManager.SelectedTask.Started {
			if err := m.TaskManager.UpdateTaskToScheduled(&m.FileManager, m.TaskManager.SelectedTask); err != nil {
				m.Errors = append(m.Errors, err.Error())
			}
			m.ViewManager.KanbanListCursor = 1
			m.ViewManager.IsKanbanTaskUpdated = true
		} else {
			if err := m.TaskManager.UpdateTaskToUnscheduled(&m.FileManager, m.TaskManager.SelectedTask); err != nil {
				m.Errors = append(m.Errors, err.Error())
			}
			m.ViewManager.KanbanListCursor = 0
			m.ViewManager.IsKanbanTaskUpdated = true
			log.Info("Updating task to unscheduled ", m.ViewManager.IsKanbanTaskUpdated)
		}
		m.FileManager.FetchTasks(&m.DirectoryManager, &m.TaskManager)
	}
	return nil
}

// AddTask opens the add task dialog
func (to TaskOperations) AddTask(m *Model) error {
	m.ViewManager.IsAddTaskView = true
	m.NewTaskInput.Reset()
	m.NewTaskInput.Prompt = ""
	m.NewTaskInput.Placeholder = "Add a task..."
	m.NewTaskInput.Focus()
	return nil
}

// AddSubTask opens the add subtask dialog
func (to TaskOperations) AddSubTask(m *Model) error {
	if m.FileManager.SelectedFile.Name == "" {
		m.Errors = append(m.Errors, "No file selected")
		return nil
	}
	m.ViewManager.IsAddSubTaskView = true
	m.NewTaskInput.Reset()

	prompt := m.FileManager.SelectedFile.FileNameWithoutExtension(m.FileManager.FileExtension) + "\n"
	m.NewTaskInput.Prompt = prompt
	m.NewTaskInput.Placeholder = "> Add a subtask..."
	m.NewTaskInput.Focus()
	return nil
}

// Command implementations for registry

type DKeyCommand struct{}

func (cmd DKeyCommand) Execute(m *Model) error {
	return TaskOperations{}.CompleteTask(m)
}

func (cmd DKeyCommand) Description() string {
	return "Mark task as completed"
}

func (cmd DKeyCommand) Contexts() []string {
	return []string{"kanban"}
}

type SKeyCommand struct{}

func (cmd SKeyCommand) Execute(m *Model) error {
	return TaskOperations{}.ScheduleOrStartTask(m)
}

func (cmd SKeyCommand) Description() string {
	return "Schedule or start task"
}

func (cmd SKeyCommand) Contexts() []string {
	return []string{"kanban"}
}

type PKeyCommand struct{}

func (cmd PKeyCommand) Execute(m *Model) error {
	return TaskOperations{}.TogglePriority(m)
}

func (cmd PKeyCommand) Description() string {
	return "Toggle task priority"
}

func (cmd PKeyCommand) Contexts() []string {
	return []string{"kanban"}
}

type UppercaseDKeyCommand struct{}

func (cmd UppercaseDKeyCommand) Execute(m *Model) error {
	return TaskOperations{}.StartTaskOrCopyStandup(m)
}

func (cmd UppercaseDKeyCommand) Description() string {
	return "Start task or copy standup to clipboard"
}

func (cmd UppercaseDKeyCommand) Contexts() []string {
	return []string{}
}

type UppercaseSKeyCommand struct{}

func (cmd UppercaseSKeyCommand) Execute(m *Model) error {
	return TaskOperations{}.ToggleScheduledUnscheduled(m)
}

func (cmd UppercaseSKeyCommand) Description() string {
	return "Toggle between scheduled and unscheduled"
}

func (cmd UppercaseSKeyCommand) Contexts() []string {
	return []string{"kanban"}
}

type AKeyCommand struct{}

func (cmd AKeyCommand) Execute(m *Model) error {
	return TaskOperations{}.AddTask(m)
}

func (cmd AKeyCommand) Description() string {
	return "Add new task"
}

func (cmd AKeyCommand) Contexts() []string {
	return []string{}
}

type UppercaseAKeyCommand struct{}

func (cmd UppercaseAKeyCommand) Execute(m *Model) error {
	return TaskOperations{}.AddSubTask(m)
}

func (cmd UppercaseAKeyCommand) Description() string {
	return "Add subtask to current file"
}

func (cmd UppercaseAKeyCommand) Contexts() []string {
	return []string{}
}
