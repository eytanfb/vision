# Code Quality Improvement Recommendations

This document outlines recommendations for improving code quality in the Vision codebase, with specific focus on Bubble Tea best practices, Go idioms, and maintainability.

## Critical Issues

### 1. Update Function - Massive If-Else Chain

**Current Problem (app/update.go:9-91):**
```go
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Deeply nested if-else chains handling all key events
    // Mixing concerns: input handling, state transitions, view logic
    // Magic numbers for viewport sizing (lines 79-84)
}
```

**Impact:**
- Hard to test individual key handlers
- Difficult to add new keyboard shortcuts
- Poor separation of concerns
- Doesn't leverage Bubble Tea's message-passing architecture

**Recommended Solution:**

Use custom message types to decouple actions from key handling:

```go
// Define custom message types
type (
    FileSelectedMsg struct{ file FileInfo }
    TaskCompletedMsg struct{ task Task }
    ViewChangedMsg struct{ view string }
    ErrorOccurredMsg struct{ err error }
)

// Update becomes a clean message router
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    case FileSelectedMsg:
        m.FileManager.SelectedFile = msg.file
        return m, m.loadFileContentCmd()
    case TaskCompletedMsg:
        return m, m.updateTaskCmd(msg.task)
    case tea.WindowSizeMsg:
        return m.handleWindowResize(msg)
    }
    return m, nil
}
```

**Benefits:**
- Testable message handlers
- Clear separation of input → action → effect
- Enables async operations with commands
- Follows Elm Architecture pattern properly

---

### 2. View Rendering Side Effects

**Current Problem (app/view_builder.go:100-144):**

The `renderKanbanList` function **mutates model state during rendering**:

```go
func renderKanbanList(m *Model, kanbanList []KanbanItem, ...) string {
    // SIDE EFFECT: Mutating model during view rendering!
    m.SelectTask(task)                           // Line 117
    m.FileManager.SelectFile(filename)           // Line 118
    m.ViewManager.IsKanbanTaskUpdated = false    // Line 119
    m.ViewManager.KanbanTaskCursor = totalIndex  // Line 120
}
```

**Impact:**
- Violates functional programming principles
- Makes debugging extremely difficult
- View rendering should be pure (output only)
- Can cause race conditions if rendering happens multiple times

**Recommended Solution:**

Separate state updates from rendering:

```go
// Update determines state, View just renders it
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Update state here
    if m.ViewManager.IsKanbanTaskUpdated {
        m.updateSelectedKanbanTask()
    }
    return m, nil
}

// Pure view function - no mutations!
func renderKanbanList(m *Model, kanbanList []KanbanItem, ...) string {
    // Only read state, never write
    return buildKanbanHTML(kanbanList, m.selectedTaskIndex)
}
```

---

### 3. Excessive File Granularity - Key Commands

**Current Problem:**
- 35 separate `*_key_command.go` files (916 total lines)
- Average 26 lines per file
- Most commands have identical structure
- Hard to see all keyboard shortcuts at once

**Files:**
```
app/a_key_command.go
app/c_key_command.go
app/d_key_command.go
... (32 more files)
```

**Recommended Solution:**

Consolidate into logical groupings:

```go
// app/navigation_commands.go
type NavigationCommands struct{}

func (nc NavigationCommands) HandleKey(key string, m *Model) error {
    switch key {
    case "j": return nc.moveDown(m)
    case "k": return nc.moveUp(m)
    case "h": return nc.moveLeft(m)
    case "l": return nc.moveRight(m)
    case "g": return nc.goToTop(m)
    }
    return nil
}

// app/file_commands.go
type FileCommands struct{}

func (fc FileCommands) HandleKey(key string, m *Model) error {
    switch key {
    case "e": return fc.editFile(m)
    case "o": return fc.openInObsidian(m)
    case "n": return fc.newFile(m)
    }
    return nil
}

// app/commands.go - Central registry
var commandHandlers = []KeyCommandHandler{
    NavigationCommands{},
    FileCommands{},
    TaskCommands{},
    ViewCommands{},
}
```

**Benefits:**
- Reduce from 35 files to ~5-6 logical files
- Easier to see all shortcuts in one place
- Better code navigation
- Shared logic between similar commands

---

### 4. String Concatenation in View Rendering

**Current Problem (throughout view_builder.go):**

```go
func BuildFilesView(m *Model, hiddenSidebar bool) (string, string) {
    list := ""
    // String concatenation in loop - O(n²) complexity
    for index, file := range m.FileManager.Files {
        list = joinVertical(list, style.Render(line))  // Copies entire string each iteration
    }
}
```

**Impact:**
- Quadratic time complexity for large lists
- Excessive memory allocations
- Slower rendering performance

**Recommended Solution:**

Use `strings.Builder` for efficient string building:

```go
func BuildFilesView(m *Model, hiddenSidebar bool) (string, string) {
    var listBuilder strings.Builder
    listBuilder.Grow(len(m.FileManager.Files) * 50) // Pre-allocate capacity

    for index, file := range m.FileManager.Files {
        style := m.getFileStyle(index)
        listBuilder.WriteString(style.Render(file.Name))
        listBuilder.WriteString("\n")
    }

    return listBuilder.String(), itemDetails
}
```

**Benefits:**
- O(n) complexity instead of O(n²)
- Single allocation when capacity is known
- 3-5x faster for large lists

---

### 5. Missing Error Handling

**Current Problem (app/file_manager.go:110-119):**

```go
func (fm FileManager) CreateStandup(company string) {
    err := copyFile(templatePath, filePath)
    if err != nil {
        log.Fatal(err)  // Crashes entire application!
    }
}
```

**Impact:**
- Application crashes instead of gracefully handling errors
- No way for UI to show error messages to user
- Lost work if error occurs

**Recommended Solution:**

Return errors and use Bubble Tea's message pattern:

```go
// Return error instead of crashing
func (fm FileManager) CreateStandup(company string) error {
    if err := copyFile(templatePath, filePath); err != nil {
        return fmt.Errorf("failed to create standup: %w", err)
    }
    return nil
}

// In Update, convert to message
func (m *Model) handleCreateStandup() tea.Cmd {
    return func() tea.Msg {
        if err := m.FileManager.CreateStandup(company); err != nil {
            return ErrorOccurredMsg{err: err}
        }
        return StandupCreatedMsg{}
    }
}

// Handle error message in Update
case ErrorOccurredMsg:
    m.Errors = append(m.Errors, msg.err.Error())
    return m, nil
```

---

### 6. God Object - Model Struct

**Current Problem (app/model.go:16-26):**

```go
type Model struct {
    DirectoryManager DirectoryManager
    TaskManager      TaskManager
    FileManager      FileManager
    ViewManager      ViewManager
    MindMapUpdater   mindmap.MindMapUpdaterInterface
    Viewport         viewport.Model
    NewTaskInput     textinput.Model
    FilterInput      textinput.Model
    Errors           []string
}
```

**Impact:**
- Single struct with too many responsibilities
- Hard to test individual components
- Tight coupling between all managers
- Doesn't leverage Bubble Tea's component pattern

**Recommended Solution:**

Use Bubble Tea's component pattern with embedded models:

```go
// Each manager becomes a tea.Model
type DirectoryManager struct {
    // ... existing fields
}

func (dm DirectoryManager) Init() tea.Cmd { return nil }
func (dm DirectoryManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) { ... }
func (dm DirectoryManager) View() string { ... }

// Main model composes components
type Model struct {
    directory DirectoryManager  // Each is independently updateable
    tasks     TaskManager
    files     FileManager
    viewport  viewport.Model

    // Current focus determines which component receives updates
    focusedComponent string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch m.focusedComponent {
    case "directory":
        var cmd tea.Cmd
        m.directory, cmd = m.directory.Update(msg)
        return m, cmd
    case "tasks":
        var cmd tea.Cmd
        m.tasks, cmd = m.tasks.Update(msg)
        return m, cmd
    }
    return m, nil
}
```

**Benefits:**
- Better separation of concerns
- Each component testable independently
- Components can manage their own state
- Follows Bubble Tea's design philosophy

---

### 7. Magic Numbers and Constants

**Current Problem (app/update.go:79-85, app/view_manager.go:38-40):**

```go
if m.ViewManager.DetailsViewHeight > 65 {
    m.ViewManager.KanbanViewLineDownFactor = 5
} else if m.ViewManager.DetailsViewHeight > 45 {
    m.ViewManager.KanbanViewLineDownFactor = 10
} else {
    m.ViewManager.KanbanViewLineDownFactor = 15
}

// In view_manager.go
const (
    heightOffset       = 12  // What does this represent?
    detailsWidthOffset = 9   // Why 9?
    navbarWidthOffset  = 5   // Why 5?
)
```

**Recommended Solution:**

Use well-named constants with comments:

```go
const (
    // Terminal size thresholds for optimal scrolling
    largeTerminalHeight  = 65  // Full HD monitors
    mediumTerminalHeight = 45  // Laptop screens

    // Scroll speed factors (higher = slower scrolling)
    largeTerminalScrollFactor  = 5   // Fast scroll for large screens
    mediumTerminalScrollFactor = 10  // Medium scroll
    smallTerminalScrollFactor  = 15  // Slow scroll for small screens

    // Layout spacing (in characters/lines)
    navbarHeight        = 3
    filterInputHeight   = 2
    statusBarHeight     = 1
    verticalPadding     = 2
    horizontalPadding   = 4

    totalHeightOffset = navbarHeight + filterInputHeight +
                       statusBarHeight + verticalPadding  // = 8
)

func (m *Model) updateScrollFactor(height int) {
    switch {
    case height > largeTerminalHeight:
        m.ViewManager.KanbanViewLineDownFactor = largeTerminalScrollFactor
    case height > mediumTerminalHeight:
        m.ViewManager.KanbanViewLineDownFactor = mediumTerminalScrollFactor
    default:
        m.ViewManager.KanbanViewLineDownFactor = smallTerminalScrollFactor
    }
}
```

---

### 8. Potential Division by Zero

**Current Problem (app/view_builder.go:405-410, app/view.go:324-328):**

```go
func buildProgressText(m *Model, category string) string {
    completedTasksCount, totalTasksCount := m.TaskManager.TaskCollection.Progress(category)
    percentage := float64(completedTasksCount) / float64(totalTasksCount)  // Division by zero!
    roundedUpPercentage := int(percentage*10) * 10
    return progressBar(completedTasksCount, totalTasksCount) + " " + fmt.Sprintf("%d%%", roundedUpPercentage)
}
```

**Recommended Solution:**

Guard against division by zero:

```go
func buildProgressText(m *Model, category string) string {
    completed, total := m.TaskManager.TaskCollection.Progress(category)

    if total == 0 {
        return "[          ] 0%"  // Empty progress bar
    }

    percentage := float64(completed) / float64(total)
    roundedPercentage := int(percentage*10) * 10

    return fmt.Sprintf("%s %d%%", progressBar(completed, total), roundedPercentage)
}

// Also add validation in Progress method
func (tc TaskCollection) Progress(filename string) (completed, total int) {
    tasks, exists := tc.TasksByFile[filename]
    if !exists {
        return 0, 0  // Explicit zero values
    }

    for _, task := range tasks {
        total++
        if task.Completed {
            completed++
        }
    }
    return completed, total
}
```

---

### 9. Not Using Bubble Tea Commands for Async Operations

**Current Problem (app/e_key_command.go:12-23):**

```go
func (j EKeyCommand) Execute(m *Model) error {
    if m.IsDetailsView() || m.IsKanbanView() {
        filePath := m.FileManager.SelectedFile.FullPath
        cmd := exec.Command("vim", "-u", "~/.dotfiles/.vimrc", filePath)
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Run()  // Blocks entire UI!
    }
    return nil
}
```

**Impact:**
- UI freezes while external editor is open
- No way to show loading state
- Can't cancel operation
- Not following Bubble Tea patterns

**Recommended Solution:**

Use `tea.ExecProcess` for external commands:

```go
func (m Model) openInEditor() tea.Cmd {
    return tea.ExecProcess(
        exec.Command("vim", "-u", "~/.dotfiles/.vimrc", m.FileManager.SelectedFile.FullPath),
        func(err error) tea.Msg {
            if err != nil {
                return ErrorOccurredMsg{err: err}
            }
            return EditorClosedMsg{}
        },
    )
}

// In Update
case tea.KeyMsg:
    if msg.String() == "e" {
        return m, m.openInEditor()  // Returns command, doesn't block
    }

case EditorClosedMsg:
    // Reload file after editing
    return m, m.reloadFileCmd()
```

**Benefits:**
- Non-blocking operation
- Proper Bubble Tea lifecycle management
- Can show loading states
- Better error handling

---

### 10. Inefficient File Sorting

**Current Problem (app/view.go:312-347):**

```go
func viewSort(filenames []string, m *Model) {
    sort.Slice(filenames, func(i, j int) bool {
        // Calls Progress() on every comparison
        // If sorting 100 items, this runs ~664 times (n log n)
        iCompletedTasks, iTotalTasks := m.TaskManager.TaskCollection.Progress(filenames[i])
        jCompletedTasks, jTotalTasks := m.TaskManager.TaskCollection.Progress(filenames[j])

        iPercentage := float64(iCompletedTasks) / float64(iTotalTasks)
        jPercentage := float64(jCompletedTasks) / float64(jTotalTasks)
        // ... more calculations
    })
}
```

**Impact:**
- O(n log n) calls to Progress() - expensive
- Recalculates same values multiple times
- Potential division by zero in comparator

**Recommended Solution:**

Pre-calculate sort keys (Schwartzian transform):

```go
type fileSortKey struct {
    filename   string
    percentage int
    isInactive bool
    isComplete bool
}

func viewSort(filenames []string, m *Model) {
    // Pre-calculate all sort keys once - O(n)
    sortKeys := make([]fileSortKey, len(filenames))
    for i, filename := range filenames {
        completed, total := m.TaskManager.TaskCollection.Progress(filename)

        var percentage int
        if total > 0 {
            percentage = int(float64(completed)/float64(total)*100)
        }

        sortKeys[i] = fileSortKey{
            filename:   filename,
            percentage: percentage,
            isInactive: m.TaskManager.TaskCollection.IsInactive(filename),
            isComplete: total > 0 && completed == total,
        }
    }

    // Sort using pre-calculated keys - O(n log n)
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

    // Copy back to original slice
    for i := range sortKeys {
        filenames[i] = sortKeys[i].filename
    }
}
```

**Performance:**
- Reduces from ~664 Progress() calls to 100
- 6x faster for 100 items
- No risk of inconsistent comparisons

---

## Medium Priority Issues

### 11. Improve Cache Management

**Current Problem:**
- `FileCache` and `TaskCache` in FileManager
- No cache invalidation strategy
- No size limits
- Manual cache management

**Recommended Solution:**
```go
type CacheEntry struct {
    data      interface{}
    timestamp time.Time
    size      int
}

type LRUCache struct {
    entries map[string]*list.Element
    list    *list.List
    maxSize int
}

func (c *LRUCache) Get(key string) (interface{}, bool) { ... }
func (c *LRUCache) Set(key string, value interface{}) { ... }
func (c *LRUCache) Invalidate(pattern string) { ... }
```

### 12. Extract Viewport Management

**Current Problem:**
- Viewport logic scattered across multiple files
- Viewport created in Update function
- Inconsistent viewport sizing

**Recommended Solution:**
```go
type ViewportManager struct {
    viewport viewport.Model
    content  string
    scroll   int
}

func (vm *ViewportManager) Update(msg tea.Msg) tea.Cmd { ... }
func (vm *ViewportManager) SetContent(content string) { ... }
func (vm *ViewportManager) View() string { ... }
```

### 13. Consolidate Style Definitions

**Current Problem:**
- Styles defined inline throughout view code
- Duplicate style definitions
- Hard to maintain consistent theme

**Recommended Solution:**
```go
// app/theme.go
type Theme struct {
    Primary        lipgloss.Color
    Secondary      lipgloss.Color
    Background     lipgloss.Color
    Highlight      lipgloss.Color

    TitleStyle     lipgloss.Style
    ListItemStyle  lipgloss.Style
    SelectedStyle  lipgloss.Style
}

func (t Theme) ApplyToModel(m *Model) { ... }
```

### 14. Add Keyboard Shortcut Help

**Current Problem:**
- No help screen showing available shortcuts
- Users must read code to discover features

**Recommended Solution:**
```go
type KeyBinding struct {
    Key         string
    Description string
    Context     string
}

var keyBindings = []KeyBinding{
    {"j/k", "Navigate up/down", "all"},
    {"e", "Edit file", "file view"},
    // ...
}

func (m Model) helpView() string {
    // Render help screen with all keybindings
}
```

### 15. Use Interfaces for Testing

**Current Problem:**
- Tight coupling to file system
- Hard to test without real files
- No dependency injection

**Recommended Solution:**
```go
type FileSystemInterface interface {
    ReadFile(path string) ([]byte, error)
    WriteFile(path string, data []byte) error
    ListFiles(path string) ([]string, error)
}

type FileManager struct {
    fs FileSystemInterface  // Can be mocked
}

// For tests
type MockFileSystem struct {
    files map[string][]byte
}
```

---

## Refactoring Priorities

### Phase 1 - Critical Safety (Week 1)
1. Fix view rendering side effects
2. Add division-by-zero guards
3. Improve error handling (remove log.Fatal calls)

### Phase 2 - Performance (Week 2)
4. Replace string concatenation with strings.Builder
5. Optimize file sorting with pre-calculated keys
6. Add cache invalidation strategy

### Phase 3 - Architecture (Week 3-4)
7. Consolidate key command files (35 → 6 files)
8. Implement custom message types
9. Convert managers to Bubble Tea components
10. Use tea.ExecProcess for external commands

### Phase 4 - Polish (Week 5-6)
11. Extract constants for magic numbers
12. Add keyboard shortcut help screen
13. Consolidate style definitions into theme
14. Add comprehensive unit tests

---

## Testing Strategy

### What to Test

**Unit Tests:**
- Task parsing logic (utils/file_utils.go)
- Date calculations (utils/date_parser.go)
- Progress calculations (avoid division by zero)
- Sort functions with pre-calculated keys

**Integration Tests:**
- Key command handlers with mock model
- File operations with mock filesystem
- View rendering with known state

**Example Test:**
```go
func TestProgressCalculation(t *testing.T) {
    tests := []struct {
        name      string
        completed int
        total     int
        want      int
    }{
        {"empty", 0, 0, 0},
        {"half", 5, 10, 50},
        {"complete", 10, 10, 100},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := calculatePercentage(tt.completed, tt.total)
            if got != tt.want {
                t.Errorf("got %d, want %d", got, tt.want)
            }
        })
    }
}
```

---

## Additional Resources

- [Bubble Tea Best Practices](https://github.com/charmbracelet/bubbletea/tree/master/examples)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Bubble Tea Component Pattern](https://github.com/charmbracelet/bubbletea/discussions/293)
