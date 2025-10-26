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

Keyboard input is handled through a command pattern (app/key_command_factory.go):
- Each key mapping has a dedicated command file (e.g., j_key_command.go, enter_key_command.go)
- KeyCommandFactory routes key presses to appropriate command handlers
- Commands implement the KeyCommand interface with Execute(model) method
- This makes it easy to add new keyboard shortcuts without modifying core logic

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

1. Create new file in app/ following pattern: `{key}_key_command.go`
2. Implement KeyCommand interface with Execute method
3. Add mapping in KeyCommandFactory.CreateKeyCommand()
4. Update documentation if it's a user-facing feature

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
