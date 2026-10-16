package app

import tea "github.com/charmbracelet/bubbletea"

// messages.go defines custom Bubble Tea message types for state changes.
// These messages decouple actions from state mutations, following the Elm Architecture.

// View Navigation Messages
// These messages are sent when the user navigates between different views

type (
	// ViewChangedMsg indicates a view transition
	ViewChangedMsg struct {
		From string
		To   string
	}

	// CompanySelectedMsg indicates a company was selected
	CompanySelectedMsg struct {
		Company Company
	}

	// SidebarToggledMsg indicates the sidebar was toggled
	SidebarToggledMsg struct {
		Hidden bool
	}
)

// File Operations Messages
// These messages are sent for file-related operations

type (
	// FileSelectedMsg indicates a file was selected
	FileSelectedMsg struct {
		File FileInfo
	}

	// FileLoadedMsg indicates a file was loaded from disk
	FileLoadedMsg struct {
		File    FileInfo
		Content string
		Err     error
	}

	// FileCreatedMsg indicates a new file was created
	FileCreatedMsg struct {
		Filename string
		Err      error
	}

	// FilesRefreshedMsg indicates the file list was refreshed
	FilesRefreshedMsg struct {
		Files []FileInfo
	}

	// FilesLoadedMsg indicates files were loaded asynchronously
	FilesLoadedMsg struct {
		Files []FileInfo
		Err   error
	}
)

// Task Operations Messages
// These messages are sent for task-related operations

type (
	// TaskSelectedMsg indicates a task was selected
	TaskSelectedMsg struct {
		Task Task
	}

	// TaskUpdatedMsg indicates a task was updated
	TaskUpdatedMsg struct {
		Task   Task
		Action string // "completed", "scheduled", "started", "priority_toggled"
		Err    error
	}

	// TasksRefreshedMsg indicates tasks were reloaded
	TasksRefreshedMsg struct{}

	// TaskCreatedMsg indicates a new task was created
	TaskCreatedMsg struct {
		TaskName string
		Err      error
	}

	// SubTaskCreatedMsg indicates a new subtask was created
	SubTaskCreatedMsg struct {
		ParentTask Task
		SubTask    string
		Err        error
	}

	// TasksLoadedMsg indicates tasks were loaded asynchronously
	TasksLoadedMsg struct {
		Tasks []Task
		Err   error
	}
)

// External Operations Messages
// These messages are sent for operations involving external programs

type (
	// EditorOpenedMsg indicates an external editor was opened
	EditorOpenedMsg struct {
		Editor string
		File   string
	}

	// EditorClosedMsg indicates the external editor closed
	EditorClosedMsg struct {
		Err error
	}

	// StandupGeneratedMsg indicates a standup was generated
	StandupGeneratedMsg struct {
		Content string
		Err     error
	}

	// ClipboardCopiedMsg indicates content was copied to clipboard
	ClipboardCopiedMsg struct {
		Content string
		Err     error
	}
)

// Input Mode Messages
// These messages are sent when input modes change

type (
	// FilterModeEnteredMsg indicates filter mode was activated
	FilterModeEnteredMsg struct{}

	// FilterModeExitedMsg indicates filter mode was deactivated
	FilterModeExitedMsg struct{}

	// AddTaskModeEnteredMsg indicates add task mode was activated
	AddTaskModeEnteredMsg struct{}

	// AddTaskModeExitedMsg indicates add task mode was deactivated
	AddTaskModeExitedMsg struct{}

	// AddSubTaskModeEnteredMsg indicates add subtask mode was activated
	AddSubTaskModeEnteredMsg struct{}

	// AddSubTaskModeExitedMsg indicates add subtask mode was deactivated
	AddSubTaskModeExitedMsg struct{}
)

// Cursor Movement Messages
// These messages are sent when cursor position changes

type (
	// CursorMovedMsg indicates the cursor moved
	CursorMovedMsg struct {
		Direction string // "up", "down", "left", "right"
		NewIndex  int
	}

	// SuggestionSelectedMsg indicates a suggestion was selected
	SuggestionSelectedMsg struct {
		Suggestion string
		Index      int
	}

	// FileNavigatedMsg indicates cursor moved to a different file
	FileNavigatedMsg struct {
		File      FileInfo
		Direction string // "up", "down"
		NewCursor int
	}

	// KanbanNavigatedMsg indicates movement within kanban view
	KanbanNavigatedMsg struct {
		ListCursor int
		TaskCursor int
		Direction  string // "up", "down", "left", "right"
	}
)

// Loading Messages
// These messages indicate loading states

type (
	// LoadingStartedMsg indicates a loading operation started
	LoadingStartedMsg struct {
		Operation string // "tasks", "files", "standup", etc.
	}

	// LoadingCompletedMsg indicates a loading operation completed
	LoadingCompletedMsg struct {
		Operation string
	}
)

// Error Messages
// These messages are sent when errors occur

type (
	// ErrorOccurredMsg indicates an error occurred
	ErrorOccurredMsg struct {
		Err     error
		Context string
	}
)

// NoOp returns a command that does nothing
func NoOp() tea.Cmd {
	return nil
}
