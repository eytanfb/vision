# Comprehensive Refactoring Plan - Phase 2

**Status**: Planning Phase
**Target**: Architectural improvements for Vision task management app
**Timeline**: 4-6 weeks
**Risk Level**: Medium-High (significant architectural changes)

---

## Executive Summary

This plan outlines the remaining architectural improvements identified in CODE_QUALITY_IMPROVEMENTS.md. Phase 1 (completed) addressed critical safety and performance issues. Phase 2 focuses on improving the application's architecture to better align with Bubble Tea best practices.

**Phase 1 Completed:**
- ✅ Division by zero guards
- ✅ View rendering side effects fixed
- ✅ Performance optimizations (sorting, string building)
- ✅ Error handling improvements
- ✅ Magic numbers extracted to constants

**Phase 2 Objectives:**
1. Consolidate 35 key command files into logical groupings
2. Implement custom message types for Bubble Tea
3. Refactor Update() function to use message-based routing
4. Convert blocking operations to non-blocking (tea.ExecProcess)

---

## Current State Analysis

### Architecture Overview

**Current File Structure:**
```
app/
├── model.go                    # God object with all state
├── update.go                   # 90-line if-else chain
├── view.go, view_builder.go   # View rendering
├── *_key_command.go (35 files) # 916 lines total, avg 26 lines each
├── *_manager.go (4 files)      # DirectoryManager, TaskManager, FileManager, ViewManager
└── task*.go, file*.go          # Domain models
```

**Current Update Flow:**
```
User presses key
    ↓
Update(tea.KeyMsg) - 90-line if-else
    ↓
KeyCommandFactory.CreateKeyCommand(key)
    ↓
Specific KeyCommand.Execute(model)
    ↓
Direct model mutation
    ↓
View() renders mutated state
```

**Pain Points:**
1. **35 separate key command files** - Hard to see all shortcuts at once
2. **Massive Update() if-else chain** - Difficult to maintain, test, or extend
3. **Blocking operations** - UI freezes when opening external editor
4. **Direct state mutation** - Makes it hard to track state changes
5. **No message passing** - Doesn't leverage Bubble Tea's architecture

---

## Phase 2.1: Consolidate Key Command Files (Week 1-2)

### Objective
Reduce from 35 key command files to 5-6 logical groupings while maintaining all functionality.

### Current State
```
app/
├── a_key_command.go           (26 lines)
├── c_key_command.go           (23 lines)
├── d_key_command.go           (21 lines)
├── e_key_command.go           (32 lines)
├── enter_key_command.go       (67 lines)
├── esc_key_command.go         (28 lines)
├── f_key_command.go           (15 lines)
├── g_key_command.go           (18 lines)
├── h_key_command.go           (21 lines)
├── j_key_command.go           (17 lines)
├── k_key_command.go           (17 lines)
├── l_key_command.go           (25 lines)
├── m_key_command.go           (19 lines)
├── minus_key_command.go       (22 lines)
├── n_key_command.go           (19 lines)
├── o_key_command.go           (43 lines)
├── one_key_command.go         (21 lines)
├── p_key_command.go           (36 lines)
├── plus_key_command.go        (22 lines)
├── s_key_command.go           (35 lines)
├── shift_tab_key_command.go   (18 lines)
├── slash_key_command.go       (23 lines)
├── t_key_command.go           (26 lines)
├── tab_key_command.go         (18 lines)
├── three_key_command.go       (21 lines)
├── two_key_command.go         (21 lines)
├── uppercase_a_key_command.go (31 lines)
├── uppercase_c_key_command.go (36 lines)
├── uppercase_d_key_command.go (38 lines)
├── uppercase_l_key_command.go (36 lines)
├── uppercase_q_key_command.go (33 lines)
├── uppercase_s_key_command.go (33 lines)
├── uppercase_w_key_command.go (48 lines)
├── w_key_command.go           (44 lines)
└── nil_key_command.go         (12 lines)
```

### Proposed Structure

**Group by Functionality:**

```
app/commands/
├── navigation.go      (j, k, h, l, g, tab, shift+tab)  ~150 lines
├── file_operations.go (e, o, n, f)                      ~110 lines
├── task_operations.go (d, s, p, D, S, a, A)            ~200 lines
├── view_control.go    (c, w, W, 1, 2, 3, +, -, C, Q, L) ~250 lines
├── input_handling.go  (enter, esc, /, t, m)             ~150 lines
└── command.go         (interfaces and factory)          ~50 lines
```

**Total: 6 files, ~910 lines** (vs 35 files, 916 lines)

### Implementation Steps

**Step 1: Create Command Interface (1 day)**
```go
// app/commands/command.go
package commands

import "vision/app"

// Command represents a keyboard command
type Command interface {
    Execute(model *app.Model) error
    Description() string
    Contexts() []string // Which views this command is available in
}

// CommandRegistry maps keys to commands
type CommandRegistry struct {
    commands map[string]Command
}

func NewRegistry() *CommandRegistry {
    return &CommandRegistry{
        commands: make(map[string]Command),
    }
}

func (r *CommandRegistry) Register(key string, cmd Command) {
    r.commands[key] = cmd
}

func (r *CommandRegistry) Get(key string) (Command, bool) {
    cmd, ok := r.commands[key]
    return cmd, ok
}
```

**Step 2: Implement Navigation Commands (2 days)**
```go
// app/commands/navigation.go
package commands

import "vision/app"

// NavigationCommands handles all cursor movement
type NavigationCommands struct{}

func (nc NavigationCommands) HandleKey(key string, m *app.Model) error {
    switch key {
    case "j":
        return nc.moveDown(m)
    case "k":
        return nc.moveUp(m)
    case "h":
        return nc.moveLeft(m)
    case "l":
        return nc.moveRight(m)
    case "g":
        return nc.goToTop(m)
    case "tab":
        return nc.nextPane(m)
    case "shift+tab":
        return nc.previousPane(m)
    }
    return nil
}

func (nc NavigationCommands) moveDown(m *app.Model) error {
    if m.IsCategoryView() {
        m.GoToNextCategory()
    } else if m.IsDetailsView() {
        if m.IsTaskDetailsFocus() {
            m.GoToNextTask()
        } else {
            m.GoToNextFile()
        }
    } else if m.IsKanbanView() {
        m.GoToNextKanbanTask()
    }
    return nil
}

// ... similar methods for other directions
```

**Step 3: Migrate Existing Commands (3 days)**
- Move each existing command to appropriate group
- Preserve all existing functionality
- Update imports

**Step 4: Update KeyCommandFactory (1 day)**
```go
// app/key_command_factory.go
package app

import "vision/app/commands"

type KeyCommandFactory struct {
    registry *commands.CommandRegistry
}

func NewKeyCommandFactory() *KeyCommandFactory {
    registry := commands.NewRegistry()

    // Register all commands
    nav := commands.NavigationCommands{}
    registry.Register("j", nav.MoveDownCommand())
    registry.Register("k", nav.MoveUpCommand())
    // ... etc

    return &KeyCommandFactory{registry: registry}
}

func (kcf KeyCommandFactory) CreateKeyCommand(key string) KeyCommand {
    if cmd, ok := kcf.registry.Get(key); ok {
        return cmd
    }
    return NilKeyCommand{}
}
```

**Step 5: Testing and Verification (1 day)**
- Build and run application
- Test all 35+ keyboard shortcuts
- Verify no regressions

### Migration Strategy

**Phase A: Create new structure (Days 1-2)**
- Create `app/commands/` directory
- Implement base interfaces and registry
- NO changes to existing files yet

**Phase B: Parallel implementation (Days 3-5)**
- Implement new grouped commands
- Keep old commands intact
- Switch can happen atomically

**Phase C: Cutover (Day 6)**
- Update KeyCommandFactory to use new commands
- Run full test suite
- Delete old command files only after verification

**Phase D: Cleanup (Day 7)**
- Remove old `*_key_command.go` files
- Update documentation
- Final testing

### Success Criteria
- [ ] All keyboard shortcuts work identically
- [ ] Code compiles without errors
- [ ] All tests pass
- [ ] File count reduced from 35 to 6
- [ ] No change in user-facing behavior

### Risks & Mitigations

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| Missing keyboard shortcut | High | Medium | Create comprehensive test matrix of all shortcuts |
| Regression in functionality | High | Medium | Keep old files until new system fully verified |
| Different behavior | Medium | Low | Careful method-by-method migration |

---

## Phase 2.2: Custom Message Types (Week 2-3)

### Objective
Implement Bubble Tea's message-passing pattern to decouple actions from state mutations.

### Current Problem
```go
// Current: Direct mutation
func (cmd DKeyCommand) Execute(m *Model) error {
    if m.IsCategoryView() && m.ViewManager.HideSidebar {
        if err := m.TaskManager.UpdateTaskToCompleted(&m.FileManager, m.TaskManager.SelectedTask); err != nil {
            m.Errors = append(m.Errors, err.Error())
        }
        m.FileManager.FetchTasks(&m.DirectoryManager, &m.TaskManager)  // Side effect!
    }
    return nil
}
```

### Proposed Solution: Message Types

```go
// app/messages.go
package app

import tea "github.com/charmbracelet/bubbletea"

// View Navigation Messages
type (
    ViewChangedMsg struct {
        From string
        To   string
    }

    CompanySelectedMsg struct {
        Company Company
    }

    CategorySelectedMsg struct {
        Category string
    }
)

// File Operations Messages
type (
    FileSelectedMsg struct {
        File FileInfo
    }

    FileLoadedMsg struct {
        File    FileInfo
        Content string
        Err     error
    }

    FileCreatedMsg struct {
        Filename string
        Err      error
    }
)

// Task Operations Messages
type (
    TaskSelectedMsg struct {
        Task Task
    }

    TaskUpdatedMsg struct {
        Task   Task
        Status string
        Err    error
    }

    TasksRefreshedMsg struct {
        Tasks []Task
    }
)

// External Operations Messages
type (
    EditorOpenedMsg struct{}

    EditorClosedMsg struct {
        Err error
    }

    StandupGeneratedMsg struct {
        Content string
        Err     error
    }
)

// Error Messages
type (
    ErrorOccurredMsg struct {
        Err     error
        Context string
    }
)
```

### Implementation Steps

**Step 1: Define Message Types (1 day)**
- Create `app/messages.go`
- Document each message type
- Define message contracts

**Step 2: Create Command Generators (2 days)**
```go
// app/commands.go
package app

import tea "github.com/charmbracelet/bubbletea"

// Commands return messages after async operations
func (m *Model) loadFileCmd(filename string) tea.Cmd {
    return func() tea.Msg {
        content, err := m.FileManager.LoadFile(filename)
        return FileLoadedMsg{
            File:    FileInfo{Name: filename},
            Content: content,
            Err:     err,
        }
    }
}

func (m *Model) updateTaskCmd(task Task, status string) tea.Cmd {
    return func() tea.Msg {
        err := m.TaskManager.UpdateTaskToCompleted(&m.FileManager, task)
        return TaskUpdatedMsg{
            Task:   task,
            Status: status,
            Err:    err,
        }
    }
}

func (m *Model) refreshTasksCmd() tea.Cmd {
    return func() tea.Msg {
        tasks := m.FileManager.FetchTasks(&m.DirectoryManager, &m.TaskManager)
        return TasksRefreshedMsg{Tasks: tasks}
    }
}
```

**Step 3: Refactor Update() to Handle Messages (3 days)**
```go
// app/update.go
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    // Keyboard input routes to commands
    case tea.KeyMsg:
        return m.handleKeyPress(msg)

    // Window events
    case tea.WindowSizeMsg:
        return m.handleWindowResize(msg)

    // File operations
    case FileSelectedMsg:
        m.FileManager.SelectedFile = msg.File
        return m, m.loadFileCmd(msg.File.Name)

    case FileLoadedMsg:
        if msg.Err != nil {
            m.Errors = append(m.Errors, msg.Err.Error())
            return m, nil
        }
        m.FileManager.SelectedFile.Content = msg.Content
        return m, nil

    case FileCreatedMsg:
        if msg.Err != nil {
            m.Errors = append(m.Errors, msg.Err.Error())
            return m, nil
        }
        return m, m.refreshTasksCmd()

    // Task operations
    case TaskSelectedMsg:
        m.TaskManager.SelectedTask = msg.Task
        return m, nil

    case TaskUpdatedMsg:
        if msg.Err != nil {
            m.Errors = append(m.Errors, msg.Err.Error())
            return m, nil
        }
        return m, m.refreshTasksCmd()

    case TasksRefreshedMsg:
        // Tasks already updated by command
        return m, nil

    // View changes
    case ViewChangedMsg:
        m.ViewManager.CurrentView = msg.To
        return m, nil

    // External editor
    case EditorClosedMsg:
        if msg.Err != nil {
            m.Errors = append(m.Errors, msg.Err.Error())
            return m, nil
        }
        // Reload current file after editing
        return m, m.loadFileCmd(m.FileManager.SelectedFile.Name)

    // Errors
    case ErrorOccurredMsg:
        m.Errors = append(m.Errors, msg.Context+": "+msg.Err.Error())
        return m, nil
    }

    return m, nil
}

func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    key := msg.String()

    // Special keys
    if key == "ctrl+c" || key == "q" {
        return m, tea.Quit
    }

    // Delegate to command factory
    cmd := KeyCommandFactory{}.CreateKeyCommand(key)
    if err := cmd.Execute(m); err != nil {
        m.Errors = append(m.Errors, err.Error())
    }

    return m, nil
}
```

**Step 4: Update Commands to Return Messages (2 days)**
```go
// Old approach: Direct mutation
func (cmd DKeyCommand) Execute(m *Model) error {
    m.TaskManager.UpdateTaskToCompleted(...)
    m.FileManager.FetchTasks(...)  // Side effect
    return nil
}

// New approach: Return command that generates message
func (cmd DKeyCommand) Execute(m *Model) tea.Cmd {
    if !m.IsCategoryView() || !m.ViewManager.HideSidebar {
        return nil
    }

    task := m.TaskManager.SelectedTask
    return m.updateTaskCmd(task, "completed")
}
```

**Step 5: Testing (2 days)**
- Unit tests for message handlers
- Integration tests for command flow
- Verify state changes happen correctly

### Benefits
1. **Testability**: Each message handler can be tested independently
2. **Clarity**: State changes are explicit and traceable
3. **Async-ready**: Easy to add loading states, progress indicators
4. **Debugging**: Can log all messages for debugging
5. **Time-travel**: Could implement undo/redo by replaying messages

### Success Criteria
- [ ] All state changes go through messages
- [ ] No direct mutations in key commands
- [ ] Commands return tea.Cmd instead of mutating directly
- [ ] All tests pass
- [ ] Application behavior unchanged

---

## Phase 2.3: Non-Blocking External Commands (Week 3-4)

### Objective
Replace blocking `exec.Command().Run()` calls with `tea.ExecProcess()` for non-blocking execution.

### Current Problems

**Problem 1: Vim Editor Blocks UI**
```go
// app/e_key_command.go
func (j EKeyCommand) Execute(m *Model) error {
    filePath := m.FileManager.SelectedFile.FullPath
    cmd := exec.Command("vim", "-u", "~/.dotfiles/.vimrc", filePath)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Run()  // UI freezes until vim closes!
    return nil
}
```

**Problem 2: Obsidian Open Blocks UI**
```go
// app/o_key_command.go
func (j OKeyCommand) Execute(m *Model) error {
    cmd := exec.Command("open", "-a", "Obsidian", obsidianPath)
    cmd.Run()  // Blocks!
    return nil
}
```

### Proposed Solution

**Use tea.ExecProcess for External Commands:**

```go
// app/commands/file_operations.go
func (fc FileCommands) openInVim(m *app.Model) tea.Cmd {
    filePath := m.FileManager.SelectedFile.FullPath

    c := exec.Command("vim", "-u", "~/.dotfiles/.vimrc", filePath)

    return tea.ExecProcess(c, func(err error) tea.Msg {
        if err != nil {
            return EditorClosedMsg{Err: err}
        }
        return EditorClosedMsg{}
    })
}

func (fc FileCommands) openInObsidian(m *app.Model) tea.Cmd {
    filePath := m.FileManager.SelectedFile.FullPath
    obsidianURL := constructObsidianURL(filePath, notesPath())

    c := exec.Command("open", "-a", "Obsidian", obsidianURL)

    return tea.ExecProcess(c, func(err error) tea.Msg {
        if err != nil {
            return ErrorOccurredMsg{
                Err:     err,
                Context: "opening in Obsidian",
            }
        }
        return nil  // Obsidian opened successfully
    })
}
```

### Implementation Steps

**Step 1: Identify All External Commands (1 day)**
Current external commands:
- `e` - Open in Vim
- `o` - Open in Obsidian
- `h-m-m` - Mind map integration (mindmap/updater.go)

**Step 2: Implement Non-Blocking Vim (2 days)**
```go
// Before
case "e":
    cmd := exec.Command("vim", filePath)
    cmd.Run()  // Blocks

// After
case "e":
    return m, m.openInVimCmd()

// In model.go
func (m *Model) openInVimCmd() tea.Cmd {
    return tea.ExecProcess(
        exec.Command("vim", "-u", "~/.dotfiles/.vimrc", m.FileManager.SelectedFile.FullPath),
        func(err error) tea.Msg {
            return EditorClosedMsg{Err: err}
        },
    )
}
```

**Step 3: Implement Non-Blocking Obsidian (1 day)**
```go
func (m *Model) openInObsidianCmd() tea.Cmd {
    url := constructObsidianURL(m.FileManager.SelectedFile.FullPath, notesPath())

    return tea.ExecProcess(
        exec.Command("open", "-a", "Obsidian", url),
        func(err error) tea.Msg {
            if err != nil {
                return ErrorOccurredMsg{Err: err, Context: "Obsidian"}
            }
            return nil
        },
    )
}
```

**Step 4: Handle EditorClosedMsg (1 day)**
```go
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case EditorClosedMsg:
        if msg.Err != nil {
            m.Errors = append(m.Errors, "Editor error: "+msg.Err.Error())
            return m, nil
        }

        // Reload file after editing
        return m, m.reloadCurrentFileCmd()
    }
    return m, nil
}

func (m *Model) reloadCurrentFileCmd() tea.Cmd {
    return func() tea.Msg {
        content, err := os.ReadFile(m.FileManager.SelectedFile.FullPath)
        if err != nil {
            return ErrorOccurredMsg{Err: err, Context: "reloading file"}
        }

        return FileLoadedMsg{
            File:    m.FileManager.SelectedFile,
            Content: string(content),
        }
    }
}
```

**Step 5: Add Loading Indicators (2 days)**
```go
type ExternalCommandRunningMsg struct {
    Command string
}

func (m *Model) View() string {
    if m.externalCommandRunning {
        return renderNavbar(m) + "\n\nOpening external editor..."
    }
    // ... normal view
}
```

**Step 6: Testing (1 day)**
- Test vim opens and closes correctly
- Test file reload after editing
- Test error handling
- Test Obsidian integration

### Benefits
1. **Responsive UI**: Application doesn't freeze
2. **Better UX**: Can show loading state
3. **Error Handling**: Proper error messages if editor fails
4. **File Reload**: Automatically reload file after editing

### Success Criteria
- [ ] Vim opens without blocking UI
- [ ] Obsidian opens without blocking UI
- [ ] Files reload after editing
- [ ] Error messages display correctly
- [ ] No UI freeze during external commands

---

## Phase 2.4: Integration and Polish (Week 4)

### Objective
Integrate all Phase 2 improvements and ensure everything works together.

### Tasks

**Integration Testing (2 days)**
- Test all keyboard shortcuts with new command system
- Test message flow for all operations
- Test external commands with message system
- Verify no regressions

**Documentation (1 day)**
- Update CLAUDE.md with new architecture
- Document message types
- Update command system documentation
- Add architecture diagrams

**Performance Validation (1 day)**
- Benchmark key operations
- Verify no performance regressions
- Measure improvement in external command handling

**Code Cleanup (1 day)**
- Remove old command files
- Remove unused code
- Update imports
- Format and lint

---

## Testing Strategy

### Unit Tests

**Test Coverage Requirements:**
- Command handlers: 80%+ coverage
- Message handlers: 90%+ coverage
- State transitions: 100% coverage

**Test Structure:**
```go
func TestNavigationCommands(t *testing.T) {
    tests := []struct {
        name      string
        key       string
        initState func(*Model)
        validate  func(*testing.T, *Model)
    }{
        {
            name: "j moves down in category view",
            key:  "j",
            initState: func(m *Model) {
                m.ViewManager.CurrentView = CategoriesView
                m.DirectoryManager.CategoriesCursor = 0
            },
            validate: func(t *testing.T, m *Model) {
                if m.DirectoryManager.CategoriesCursor != 1 {
                    t.Errorf("Expected cursor at 1, got %d", m.DirectoryManager.CategoriesCursor)
                }
            },
        },
        // ... more tests
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := createTestModel()
            tt.initState(m)

            cmd := NavigationCommands{}.HandleKey(tt.key, m)
            if cmd != nil {
                // Execute command if returned
            }

            tt.validate(t, m)
        })
    }
}
```

### Integration Tests

**Test Scenarios:**
1. Complete task workflow (schedule → start → complete)
2. File creation and editing workflow
3. View navigation workflow
4. External editor workflow
5. Error handling workflow

### Manual Testing Checklist

**Phase 2.1 - Commands:**
- [ ] All 35+ keyboard shortcuts work
- [ ] Help screen shows all commands
- [ ] Commands work in correct contexts
- [ ] No duplicate key bindings

**Phase 2.2 - Messages:**
- [ ] All state changes go through messages
- [ ] Error messages display correctly
- [ ] Async operations work
- [ ] No race conditions

**Phase 2.3 - External Commands:**
- [ ] Vim opens and closes correctly
- [ ] Obsidian opens correctly
- [ ] Files reload after editing
- [ ] UI stays responsive

**Phase 2.4 - Integration:**
- [ ] All features work together
- [ ] No performance regressions
- [ ] Documentation is accurate
- [ ] Code is clean and well-organized

---

## Risk Management

### High-Risk Items

| Item | Risk | Mitigation |
|------|------|------------|
| Message System | Breaking existing workflows | Extensive testing, parallel implementation |
| External Commands | Platform-specific behavior | Test on macOS/Linux, graceful fallbacks |
| Command Consolidation | Missing functionality | Comprehensive test matrix |

### Rollback Strategy

**For Each Phase:**
1. Work in feature branch
2. Keep old code until new code verified
3. Can revert commit if issues found
4. Staged rollout to production

**Rollback Triggers:**
- Critical bug found
- >5% performance regression
- User-facing behavior change
- Test coverage drops below 70%

---

## Success Metrics

### Code Quality Metrics
- [ ] Test coverage ≥ 80%
- [ ] Zero `go vet` warnings
- [ ] Zero `golint` warnings
- [ ] Cyclomatic complexity < 15 per function

### Performance Metrics
- [ ] No measurable performance regression
- [ ] External commands non-blocking
- [ ] UI responsive at all times

### Architecture Metrics
- [ ] File count reduced from 35 to 6 (key commands)
- [ ] Update() function < 100 lines
- [ ] All state changes through messages
- [ ] Zero blocking external calls

### User Experience Metrics
- [ ] All keyboard shortcuts work
- [ ] No change in visible behavior
- [ ] Error messages are helpful
- [ ] UI never freezes

---

## Timeline

### Week 1: Command Consolidation
- Days 1-2: Create new structure
- Days 3-5: Migrate commands
- Days 6-7: Testing and verification

### Week 2: Message Types (Part 1)
- Day 1: Define message types
- Days 2-3: Create command generators
- Days 4-5: Refactor Update() function

### Week 3: Message Types (Part 2) + External Commands
- Days 1-2: Update commands to use messages
- Days 3-4: Implement non-blocking external commands
- Day 5: Testing

### Week 4: Integration and Polish
- Days 1-2: Integration testing
- Day 3: Documentation
- Day 4: Performance validation
- Day 5: Code cleanup and final review

**Total: 4 weeks (20 working days)**

---

## Dependencies

### Technical Dependencies
- Bubble Tea >= 0.25.0
- Go >= 1.21.4
- All dependencies in go.mod

### External Dependencies
- Vim installed (for editor integration)
- Obsidian installed (for Obsidian integration)
- h-m-m installed (optional, for mind map integration)

### Team Dependencies
- Code review availability
- Testing resources
- Deployment access

---

## Post-Refactoring Improvements

### Possible Future Enhancements (Not in Scope)

**Component Pattern:**
- Convert managers to tea.Model components
- Each manager handles its own Update/View
- Better separation of concerns

**State Management:**
- Centralized state store
- State history for undo/redo
- State persistence

**Performance:**
- Lazy loading of files
- Virtual scrolling for large lists
- Caching improvements

**Features:**
- Keyboard shortcut help screen
- Custom keybinding configuration
- Plugin system for extensions

---

## Conclusion

This refactoring plan provides a structured approach to improving the Vision codebase architecture. By following the phased approach and testing strategy, we can safely migrate to a more maintainable, testable, and Bubble Tea-idiomatic architecture.

**Expected Outcomes:**
1. Cleaner, more maintainable codebase
2. Better alignment with Bubble Tea best practices
3. Improved testability
4. Enhanced user experience (non-blocking operations)
5. Foundation for future improvements

**Next Steps:**
1. Review and approve this plan
2. Create feature branch for Phase 2
3. Begin Phase 2.1: Command Consolidation
4. Weekly progress reviews
5. Continuous testing and validation
