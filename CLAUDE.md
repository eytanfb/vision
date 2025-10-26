# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Vision is a terminal-based task management hub for managing work across multiple clients/companies. It provides a keyboard-driven interface for viewing, organizing, and tracking tasks stored in markdown files.

**Core Features:**
- Multi-company task management with configurable clients
- Multiple view modes: Categories, Kanban, Calendar, Task Details
- Task parsing from markdown files with date tracking (scheduled, started, completed)
- Daily and weekly standup generation for Slack
- Integration with h-m-m mind mapping tool (optional)

## Technology Stack

- **Language:** Go 1.21.4
- **TUI Framework:** Bubble Tea (charmbracelet/bubbletea)
- **UI Components:** Bubble Tea Bubbles, Lipgloss, Glamour (markdown rendering)
- **Architecture:** Model-View pattern with Command pattern for keyboard bindings

## Common Commands

### Running the Application
```bash
# Run with default company
go run main.go

# Run with specific company
go run main.go clerky

# Run with company and category
go run main.go clerky tasks
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./app
go test ./utils
go test ./config

# Run tests with verbose output
go test -v ./...

# Run a specific test
go test ./app -run TestTaskExtraction
```

### Building
```bash
# Build binary
go build -o vision

# Run the binary
./vision
```

## Architecture

### Core Components

The application follows a modular architecture with clear separation of concerns:

**Model (app/model.go)**
- Central application state containing all managers
- Implements Bubble Tea's Model interface (Init, Update, View)
- Coordinates between DirectoryManager, TaskManager, FileManager, ViewManager

**DirectoryManager (app/directory_manager.go)**
- Manages company/client selection and folder navigation
- Handles category selection (tasks, standups, meetings, etc.)
- Maintains cursor state for company/category navigation

**TaskManager (app/task_manager.go)**
- Parses markdown files to extract tasks with date annotations
- Tracks task states: scheduled (⏳), started (🛫), completed (✅), priority (🔺)
- Generates daily and weekly summaries for Slack standups
- Task format: `- [ ] Task description ⏰ 2024-01-15 🛫 2024-01-16 ✅ 2024-01-17`

**FileManager (app/file_manager.go)**
- Handles file I/O operations for markdown task files
- Implements caching for files and tasks (FileCache, TaskCache)
- Provides autocomplete suggestions for people and task names
- Manages file sorting by completion progress and modification time

**ViewManager (app/view_manager.go)**
- Controls view state transitions (Companies → Categories → Details)
- Manages viewport dimensions and layout calculations
- Handles Kanban view with list/task cursor positions
- Coordinates focus states for different UI elements

### Command Pattern for Key Bindings

Keyboard input is handled through a command pattern with organized command groups:

**Core Files:**
- **app/command_interface.go**: Defines Command interface with Execute(), Description(), Contexts()
- **app/command_registry.go**: CommandRegistry for mapping keys to commands
- **app/key_command_factory.go**: Creates registry and registers all commands

**Command Groups:**
- **app/navigation.go**: Movement commands (j, k, h, l, g, tab, shift+tab)
- **app/file_operations.go**: File handling (e, o, n, f - edit, open, next company, sidebar)
- **app/task_operations.go**: Task management (d, s, p, D, S, a, A - complete, schedule, priority, etc.)
- **app/view_control.go**: View switching (c, w, W, 1-3, +/-, C, Q, L - calendar, weekly, companies)
- **app/input_handling.go**: Input modes (enter, esc, /, t, m - selection, filtering, task/meeting views)

Each command group contains:
1. A struct (e.g., `NavigationCommands{}`)
2. Handler methods for each operation
3. Command wrapper types implementing the Command interface
4. Execute(), Description(), and Contexts() methods for each wrapper

### Message-Passing Architecture (Phase 2.2)

The application follows Bubble Tea's message-passing pattern (Elm Architecture):

**Core Files:**
- **app/messages.go**: Custom message type definitions for all state changes
- **app/tea_commands.go**: Command generator functions that return tea.Cmd
- **app/update.go**: Message handlers for processing state changes

**Message Categories:**
- **View Navigation**: ViewChangedMsg, CompanySelectedMsg, CategorySelectedMsg, SidebarToggledMsg
- **File Operations**: FileSelectedMsg, FileLoadedMsg, FileCreatedMsg, FilesRefreshedMsg
- **Task Operations**: TaskSelectedMsg, TaskUpdatedMsg, TasksRefreshedMsg, TaskCreatedMsg
- **External Operations**: EditorClosedMsg, StandupGeneratedMsg, ClipboardCopiedMsg
- **Input Modes**: FilterModeEnteredMsg, AddTaskModeEnteredMsg, AddSubTaskModeEnteredMsg
- **Errors**: ErrorOccurredMsg

**Command Pattern:**
Commands return `tea.Cmd` instead of mutating state directly:
```go
// Command returns tea.Cmd
func (cmd DKeyCommand) Execute(m *Model) tea.Cmd {
    return m.updateTaskCmd(task, "completed")
}

// Command generator creates message
func (m *Model) updateTaskCmd(task Task, action string) tea.Cmd {
    return func() tea.Msg {
        err := m.TaskManager.UpdateTaskToCompleted(&m.FileManager, task)
        return TaskUpdatedMsg{Task: task, Action: action, Err: err}
    }
}

// Update() handles message
case TaskUpdatedMsg:
    if msg.Err != nil {
        m.Errors = append(m.Errors, msg.Err.Error())
    }
    return m, m.refreshTasksCmd()
```

**Benefits:**
- Async-ready architecture for non-blocking operations
- Clear separation between actions and state mutations
- Better testability (message handlers tested independently)
- Foundation for undo/redo and time-travel debugging

### Non-Blocking External Commands (Phase 2.3)

All external command executions use `tea.ExecProcess` for non-blocking execution:

**External Commands:**
- **vim editor** (e key): Opens file, returns EditorClosedMsg, auto-reloads tasks
- **Obsidian app** (o key): Opens file in Obsidian, returns ErrorOccurredMsg on failure
- **gh dash** (g key): Opens GitHub dashboard, returns ErrorOccurredMsg on failure

**Pattern:**
```go
func (fo FileOperations) OpenInVim(m *Model) tea.Cmd {
    c := exec.Command("vim", "-u", "~/.dotfiles/.vimrc", filePath)

    return tea.ExecProcess(c, func(err error) tea.Msg {
        if err != nil {
            return EditorClosedMsg{Err: err}
        }
        return EditorClosedMsg{}  // Triggers task reload
    })
}

// In Update()
case EditorClosedMsg:
    if msg.Err != nil {
        m.Errors = append(m.Errors, "Editor error: "+msg.Err.Error())
        return m, nil
    }
    m.FileManager.FetchTasks(&m.DirectoryManager, &m.TaskManager)
    return m, nil
```

**Benefits:**
- UI remains responsive while external programs are open
- Automatic file/task reload after editing
- Proper error handling for external command failures
- Foundation for loading indicators

### Task Parsing

Tasks are extracted from markdown files using special date annotations:
- `⏰ YYYY-MM-DD` - Scheduled date
- `🛫 YYYY-MM-DD` - Start date
- `✅ YYYY-MM-DD` - Completion date
- `🔺` - Priority marker
- Format: Standard markdown checkbox `- [ ]` for incomplete, `- [x]` for complete

### Configuration

**config/config.json** defines:
- Companies: displayName, folderPathName, fullPath, subFolders, color
- Categories: tasks, standups, meetings, projects, people, teams, estimates, other, onboarding
- DefaultCompany: which company to load on startup
- PreferredFileExtension: .md or .txt

Files are stored at `~/Notes/{company}/{category}/{filename}.md`

### Views and Navigation

**View Hierarchy:**
1. **CompaniesView** - Select client/company (if ShowCompanies is true)
2. **CategoriesView** - Select category (tasks, meetings, etc.)
3. **DetailsView** - View files and tasks with details pane

**Special Views:**
- **Kanban Mode** - Triggered by hiding sidebar (HideSidebar=true), groups tasks by status
- **Calendar View** - Shows tasks organized by date
- **Weekly View** - Week-at-a-glance for planning

**Key Bindings** (see key_command_factory.go for full list):
- `j/k` - Navigate up/down
- `h/l` - Navigate left/right or change views
- `enter` - Select item
- `e` - Edit current file
- `o` - Open file in editor
- `tab/shift+tab` - Switch between sidebar and details
- `w/W` - Daily/weekly standup views
- `c` - Calendar view
- `/` - Filter tasks

## Development Patterns

### Adding a New Key Command

Commands are organized into logical groups. To add a new command:

1. **Identify the appropriate group** (navigation, file_operations, task_operations, view_control, or input_handling)
2. **Add a handler method** to the group's struct (e.g., `func (nc NavigationCommands) NewMethod(m *Model) error`)
3. **Create a command wrapper** implementing the Command interface:
   ```go
   type NewKeyCommand struct{}

   func (cmd NewKeyCommand) Execute(m *Model) error {
       return NavigationCommands{}.NewMethod(m)
   }

   func (cmd NewKeyCommand) Description() string {
       return "Description of what this command does"
   }

   func (cmd NewKeyCommand) Contexts() []string {
       return []string{} // Empty for all contexts, or specify like []string{"details_view"}
   }
   ```
4. **Register in KeyCommandFactory** (app/key_command_factory.go):
   ```go
   registry.Register("x", NewKeyCommand{})
   ```
5. **Update documentation** if it's a user-facing feature

For a completely new command category, create a new group file following the pattern of existing group files.

### Adding a New View

1. Add view constant in app/view_manager.go
2. Implement view rendering logic in app/view_builder.go
3. Add navigation logic in ViewManager.GoToNextView/GoToPreviousView
4. Create key command to trigger the view if needed

### Working with Tasks

- Task extraction happens in utils/file_utils.go (ExtractTasksFromText)
- Date parsing is handled by utils/date_parser.go
- Task state is determined by presence of date annotations
- Always preserve task order and line numbers for file updates

### Testing Strategy

- Unit tests exist for core logic: task parsing, date utilities, config loading
- Test files follow Go convention: `*_test.go`
- Mock interfaces where needed (e.g., MindMapUpdaterInterface has NullMindMapUpdater)
- Focus tests on business logic, not Bubble Tea UI rendering

## Notes Path Configuration

The application expects task files at `~/Notes/{company}/{category}/` by default. The notes path is determined by:
- FileManager operations read from DirectoryManager's SelectedCompany.FullPath
- Config.json defines fullPath with `~` expansion
- MindMap integration uses `~/Notes/personal/daily_mind_maps` if h-m-m is installed

## External Dependencies

**Optional:**
- `h-m-m` - Mind mapping tool for daily brain dumps. If not in PATH, NullMindMapUpdater is used.

**Required:**
- All dependencies in go.mod are required for core functionality
- Glamour provides markdown rendering in the terminal
- Teacup provides file picker component
