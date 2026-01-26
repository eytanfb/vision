package app

import (
	"strings"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

// TaskOperations handles all task-related commands
type TaskOperations struct{}

// HandleKey routes task operation keys to their appropriate handlers
func (to TaskOperations) HandleKey(key string, m *Model) tea.Cmd {
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
func (to TaskOperations) CompleteTask(m *Model) tea.Cmd {
	if m.IsCategoryView() && m.ViewManager.HideSidebar {
		selectedTask := m.TaskManager.SelectedTask
		return m.updateTaskCmd(selectedTask, "completed")
	}
	return nil
}

// ScheduleOrStartTask schedules or starts a task depending on its current state
func (to TaskOperations) ScheduleOrStartTask(m *Model) tea.Cmd {
	if m.ViewManager.HideSidebar {
		selectedTask := m.TaskManager.SelectedTask
		if selectedTask.Scheduled {
			log.Info("SKeyCommand: Update task to started")
			return m.updateTaskCmd(selectedTask, "started")
		} else {
			log.Info("SKeyCommand: Update task to scheduled")
			return m.updateTaskCmd(selectedTask, "scheduled")
		}
	}
	return nil
}

// TogglePriority toggles the priority marker on a task
func (to TaskOperations) TogglePriority(m *Model) tea.Cmd {
	if m.IsCategoryView() && m.ViewManager.HideSidebar {
		selectedTask := m.TaskManager.SelectedTask

		if !strings.Contains(selectedTask.Text, "🔺") {
			log.Info("Adding priority marker to task")
			return m.updateTaskCmd(selectedTask, "priority")
		} else {
			log.Info("Removing priority marker from task")
			return m.updateTaskCmd(selectedTask, "unpriority")
		}
	}
	return nil
}

// StartTaskOrCopyStandup starts a task in kanban view or copies standup in other views
func (to TaskOperations) StartTaskOrCopyStandup(m *Model) tea.Cmd {
	if m.IsCategoryView() && m.ViewManager.HideSidebar {
		selectedTask := m.TaskManager.SelectedTask
		return m.updateTaskCmd(selectedTask, "started")
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
func (to TaskOperations) ToggleScheduledUnscheduled(m *Model) tea.Cmd {
	if m.IsCategoryView() && m.ViewManager.HideSidebar {
		selectedTask := m.TaskManager.SelectedTask
		if selectedTask.Started {
			log.Info("Updating started task to scheduled")
			return m.updateTaskCmd(selectedTask, "scheduled")
		} else {
			log.Info("Updating task to unscheduled")
			return m.updateTaskCmd(selectedTask, "unscheduled")
		}
	}
	return nil
}

// AddTask opens the add task dialog
func (to TaskOperations) AddTask(m *Model) tea.Cmd {
	m.ViewManager.IsAddTaskView = true
	m.NewTaskInput.Reset()
	m.NewTaskInput.Prompt = ""
	m.NewTaskInput.Placeholder = "Add a task..."
	m.NewTaskInput.Focus()
	return nil
}

// AddSubTask opens the add subtask dialog
func (to TaskOperations) AddSubTask(m *Model) tea.Cmd {
	if m.FileManager.SelectedFile.Name == "" {
		m.AddError("No file selected")
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

func (cmd DKeyCommand) Execute(m *Model) tea.Cmd {
	return TaskOperations{}.CompleteTask(m)
}

func (cmd DKeyCommand) Description() string {
	return "Mark task as completed"
}

func (cmd DKeyCommand) Contexts() []string {
	return []string{ContextKanban}
}

type SKeyCommand struct{}

func (cmd SKeyCommand) Execute(m *Model) tea.Cmd {
	return TaskOperations{}.ScheduleOrStartTask(m)
}

func (cmd SKeyCommand) Description() string {
	return "Schedule or start task"
}

func (cmd SKeyCommand) Contexts() []string {
	return []string{ContextKanban}
}

type PKeyCommand struct{}

func (cmd PKeyCommand) Execute(m *Model) tea.Cmd {
	return TaskOperations{}.TogglePriority(m)
}

func (cmd PKeyCommand) Description() string {
	return "Toggle task priority"
}

func (cmd PKeyCommand) Contexts() []string {
	return []string{ContextKanban}
}

type UppercaseDKeyCommand struct{}

func (cmd UppercaseDKeyCommand) Execute(m *Model) tea.Cmd {
	return TaskOperations{}.StartTaskOrCopyStandup(m)
}

func (cmd UppercaseDKeyCommand) Description() string {
	return "Start task or copy standup to clipboard"
}

func (cmd UppercaseDKeyCommand) Contexts() []string {
	return []string{}
}

type UppercaseSKeyCommand struct{}

func (cmd UppercaseSKeyCommand) Execute(m *Model) tea.Cmd {
	return TaskOperations{}.ToggleScheduledUnscheduled(m)
}

func (cmd UppercaseSKeyCommand) Description() string {
	return "Toggle between scheduled and unscheduled"
}

func (cmd UppercaseSKeyCommand) Contexts() []string {
	return []string{ContextKanban}
}

type AKeyCommand struct{}

func (cmd AKeyCommand) Execute(m *Model) tea.Cmd {
	return TaskOperations{}.AddTask(m)
}

func (cmd AKeyCommand) Description() string {
	return "Add new task"
}

func (cmd AKeyCommand) Contexts() []string {
	return []string{}
}

type UppercaseAKeyCommand struct{}

func (cmd UppercaseAKeyCommand) Execute(m *Model) tea.Cmd {
	return TaskOperations{}.AddSubTask(m)
}

func (cmd UppercaseAKeyCommand) Description() string {
	return "Add subtask to current file"
}

func (cmd UppercaseAKeyCommand) Contexts() []string {
	return []string{}
}
