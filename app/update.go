package app

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := m.update(msg)
	m.syncViewports()
	return model, cmd
}

func (m *Model) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.Errors = []string{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		Debug("Key pressed: ", key)
		if key == "ctrl+c" {
			return m, tea.Quit
		} else if key == "q" {
			if m.IsAddTaskView() {
				m.NewTaskInput, cmd = m.NewTaskInput.Update(msg)
				return m, cmd
			} else if m.IsAddSubTaskView() {
				m.NewTaskInput, cmd = m.NewTaskInput.Update(msg)
				return m, cmd
			} else if m.IsFilterView() {
				m.FilterInput, cmd = m.FilterInput.Update(msg)
				return m, cmd
			}

			return m, tea.Quit
		} else if m.IsAddTaskView() || m.IsFilterView() || m.IsAddSubTaskView() {
			factory := NewKeyCommandFactory()
			if key == "esc" {
				cmdResult := factory.CreateKeyCommand("esc").Execute(m)
				cmds = append(cmds, cmdResult)
			} else if key == "enter" {
				cmdResult := factory.CreateKeyCommand("enter").Execute(m)
				cmds = append(cmds, cmdResult)
			}

			if m.IsFilterView() {
				m.FilterInput, cmd = m.FilterInput.Update(msg)
				m.TaskManager.TaskCollection.FilterValue = m.FilterInput.Value()
				cmds = append(cmds, cmd)
			} else {
				if key == "[" {
					cmdResult := factory.CreateKeyCommand("[").Execute(m)
					cmds = append(cmds, cmdResult)
				} else if key == "]" {
					cmdResult := factory.CreateKeyCommand("]").Execute(m)
					cmds = append(cmds, cmdResult)
				}

				m.NewTaskInput, cmd = m.NewTaskInput.Update(msg)
				cmds = append(cmds, cmd)
			}

			if key == "tab" {
				cmdResult := factory.CreateKeyCommand("tab").Execute(m)
				cmds = append(cmds, cmdResult)
			} else if key == "shift+tab" {
				cmdResult := factory.CreateKeyCommand("shift+tab").Execute(m)
				cmds = append(cmds, cmdResult)
			}
		} else {
			keyCommandFactory := NewKeyCommandFactory()
			keyCommand := keyCommandFactory.CreateKeyCommand(key)
			cmdResult := keyCommand.Execute(m)
			cmds = append(cmds, cmdResult)
		}
	case tea.WindowSizeMsg:
		m.ViewManager.SetWidth(msg.Width)
		m.ViewManager.SetHeight(msg.Height)

		if !m.ViewManager.Ready {
			m.Viewport = viewport.New(m.ViewManager.DetailsViewWidth, m.ViewManager.DetailsViewHeight)
			m.TasksViewport = viewport.New(m.ViewManager.DetailsViewWidth, m.ViewManager.DetailsViewHeight)
			m.FilesViewport = viewport.New(m.ViewManager.sidebarContentSize())
			m.ViewManager.Ready = true
		} else {
			m.Viewport.Width = m.ViewManager.DetailsViewWidth
			m.Viewport.Height = m.ViewManager.DetailsViewHeight

			// Adjust scroll speed based on terminal height
			switch {
			case m.ViewManager.DetailsViewHeight > largeTerminalHeight:
				m.ViewManager.KanbanViewLineDownFactor = largeTerminalScrollFactor
			case m.ViewManager.DetailsViewHeight > mediumTerminalHeight:
				m.ViewManager.KanbanViewLineDownFactor = mediumTerminalScrollFactor
			default:
				m.ViewManager.KanbanViewLineDownFactor = smallTerminalScrollFactor
			}
		}

		// Mark everything dirty on window resize - forces complete re-render
		m.ViewManager.MarkDirty("all")
		// Invalidate style cache since dimensions changed
		globalStyleCache.InvalidateAll()

	// Handle custom messages from Phase 2.2 and 2.3
	case EditorClosedMsg:
		if msg.Err != nil {
			m.AddError("Editor error: " + msg.Err.Error())
			return m, nil
		}
		// Reload file and tasks after editing - async
		log.Info("Editor closed, reloading tasks")
		return m, m.loadTasksCmd()

	case ErrorOccurredMsg:
		m.AddError(msg.Context + ": " + msg.Err.Error())
		return m, nil

	// Handle async file/task loading messages
	case FilesLoadedMsg:
		m.IsLoadingFiles = false
		if msg.Err != nil {
			m.AddError("Error loading files: " + msg.Err.Error())
			return m, nil
		}
		m.FileManager.SetFiles(msg.Files)
		// Mark files dirty when files are loaded
		m.ViewManager.MarkDirty("files")
		return m, nil

	case TasksLoadedMsg:
		m.IsLoadingTasks = false
		if msg.Err != nil {
			m.AddError("Error loading tasks: " + msg.Err.Error())
			return m, nil
		}
		// Tasks are already loaded into TaskManager by FetchTasks
		// Mark tasks and kanban dirty when tasks are loaded
		m.ViewManager.MarkDirty("tasks")
		// Sync SelectedTask after tasks are reloaded (important after task operations)
		if m.ViewManager.HideSidebar {
			NavigationCommands{}.updateSelectedTaskInKanban(m)
		}
		return m, nil

	case LoadingStartedMsg:
		if msg.Operation == "tasks" {
			m.IsLoadingTasks = true
		} else if msg.Operation == "files" {
			m.IsLoadingFiles = true
		}
		return m, nil

	case LoadingCompletedMsg:
		if msg.Operation == "tasks" {
			m.IsLoadingTasks = false
		} else if msg.Operation == "files" {
			m.IsLoadingFiles = false
		}
		return m, nil

	case TaskUpdatedMsg:
		if msg.Err != nil {
			m.AddError("Error updating task: " + msg.Err.Error())
			return m, nil
		}

		// Update kanban list cursor based on action
		if msg.Action == "completed" {
			m.ViewManager.KanbanListCursor = 2 // Move to completed list
			m.ViewManager.KanbanTaskCursor = 0
		} else if msg.Action == "started" || msg.Action == "scheduled" {
			m.ViewManager.KanbanListCursor = 1 // Move to active list
			m.ViewManager.KanbanTaskCursor = 0
		} else if msg.Action == "unscheduled" {
			m.ViewManager.KanbanListCursor = 0 // Move to todo list
			m.ViewManager.KanbanTaskCursor = 0
		}

		m.ViewManager.IsKanbanTaskUpdated = true
		// Mark tasks dirty after update
		m.ViewManager.MarkDirty("tasks")
		return m, m.loadTasksCmd()

	case TaskCreatedMsg:
		if msg.Err != nil {
			m.AddError("Error creating task: " + msg.Err.Error())
		}
		m.ViewManager.MarkDirty("tasks")
		return m, m.loadFilesCmd()

	case SubTaskCreatedMsg:
		if msg.Err != nil {
			m.AddError("Error creating subtask: " + msg.Err.Error())
		}
		m.ViewManager.MarkDirty("tasks")
		return m, m.loadTasksCmd()

	case ViewChangedMsg:
		m.ViewManager.MarkDirty("all")
		return m, nil

	case SidebarToggledMsg:
		m.ViewManager.MarkDirty("all")
		return m, nil

	case FileSelectedMsg:
		m.ViewManager.MarkDirty("tasks")
		return m, nil

	case CompanySelectedMsg:
		m.ViewManager.MarkDirty("all")
		return m, m.loadFilesCmd()

	case FileNavigatedMsg:
		m.FileManager.FilesCursor = msg.NewCursor
		m.FileManager.SelectedFile = msg.File
		m.Viewport.GotoTop()
		m.ViewManager.MarkDirty("files")
		return m, nil

	case KanbanNavigatedMsg:
		m.ViewManager.KanbanListCursor = msg.ListCursor
		m.ViewManager.KanbanTaskCursor = msg.TaskCursor
		NavigationCommands{}.updateSelectedTaskInKanban(m)
		m.ViewManager.MarkDirty("kanban")
		return m, nil
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
