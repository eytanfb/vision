package app

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
	"vision/config"

	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	testWidth  = 160
	testHeight = 50
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]|\x1b\][^\x07]*\x07`)

func stripANSI(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
}

func writeNotesFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func newRenderTestModel(t *testing.T) *Model {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("PATH", t.TempDir())

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	nextWeek := time.Now().AddDate(0, 0, 7).Format("2006-01-02")

	notes := filepath.Join(home, "Notes")
	writeNotesFile(t, filepath.Join(notes, "lifeplus", "tasks", "project-alpha.md"), fmt.Sprintf(
		"# Project Alpha\n\n"+
			"- [x] Ship the release ⏰ 2024-01-15 🛫 2024-01-16 ✅ 2024-01-17\n"+
			"- [ ] Write the docs ⏰ %s 🛫 %s\n"+
			"- [ ] Fix the bug 🔺 ⏰ %s\n"+
			"- [ ] Plan next quarter ⏰ %s\n"+
			"- [x] Close the sprint 🛫 %s ✅ %s\n"+
			"- [ ] Unscheduled idea\n",
		yesterday, today, today, nextWeek, yesterday, today))
	writeNotesFile(t, filepath.Join(notes, "lifeplus", "tasks", "project-beta.md"),
		"# Project Beta\n\n- [x] Kickoff meeting ✅ 2024-02-01\n- [ ] Draft proposal\n")
	writeNotesFile(t, filepath.Join(notes, "lifeplus", "meetings", "2024-01-15.md"),
		"# Weekly sync\n\n- Discussed roadmap\n")
	writeNotesFile(t, filepath.Join(notes, "clerky", "tasks", "clerky-work.md"),
		fmt.Sprintf("- [ ] Clerky task ⏰ %s\n", today))

	cfg := &config.Config{
		Companies: []config.Company{
			{DisplayName: "Lifeplus", FolderPathName: "lifeplus", FullPath: filepath.Join(notes, "lifeplus"), SubFolders: []string{"tasks", "meetings"}, Color: "#00FF00"},
			{DisplayName: "Clerky", FolderPathName: "clerky", FullPath: filepath.Join(notes, "clerky"), SubFolders: []string{"tasks", "meetings"}, Color: "#0000FF"},
		},
		DefaultCompany:         "lifeplus",
		PreferredFileExtension: ".md",
	}
	cfg.Categories = config.LoadCategories(cfg)

	model, ok := InitialModel(cfg, []string{}).(*Model)
	if !ok {
		t.Fatalf("InitialModel did not return *Model")
	}
	model.FilterInput.Cursor.SetMode(cursor.CursorStatic)
	model.NewTaskInput.Cursor.SetMode(cursor.CursorStatic)
	return model
}

type renderHarness struct {
	t     *testing.T
	m     *Model
	trail []string
}

func newRenderHarness(t *testing.T) *renderHarness {
	t.Helper()
	h := &renderHarness{t: t, m: newRenderTestModel(t)}
	h.send(tea.WindowSizeMsg{Width: testWidth, Height: testHeight}, "resize")
	return h
}

func keyMsg(key string) tea.KeyMsg {
	switch key {
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	case "tab":
		return tea.KeyMsg{Type: tea.KeyTab}
	case "shift+tab":
		return tea.KeyMsg{Type: tea.KeyShiftTab}
	case "ctrl+c":
		return tea.KeyMsg{Type: tea.KeyCtrlC}
	}
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
}

func (h *renderHarness) press(keys ...string) tea.Cmd {
	h.t.Helper()
	var last tea.Cmd
	for _, k := range keys {
		last = h.send(keyMsg(k), k)
	}
	return last
}

func (h *renderHarness) send(msg tea.Msg, label string) tea.Cmd {
	h.t.Helper()
	h.trail = append(h.trail, label)

	var cmd tea.Cmd
	h.guard("Update", func() {
		_, cmd = h.m.Update(msg)
	})
	h.drain(cmd, 0)
	h.render()
	return cmd
}

func (h *renderHarness) guard(stage string, fn func()) {
	h.t.Helper()
	defer func() {
		if r := recover(); r != nil {
			h.t.Fatalf("%s panicked after keys %v: %v", stage, h.trail, r)
		}
	}()
	fn()
}

func runCmdWithTimeout(cmd tea.Cmd) (tea.Msg, bool) {
	result := make(chan tea.Msg, 1)
	go func() { result <- cmd() }()
	select {
	case msg := <-result:
		return msg, true
	case <-time.After(250 * time.Millisecond):
		return nil, false
	}
}

func (h *renderHarness) drain(cmd tea.Cmd, depth int) {
	h.t.Helper()
	if cmd == nil || depth > 5 {
		return
	}

	var msg tea.Msg
	var ok bool
	h.guard("Cmd", func() { msg, ok = runCmdWithTimeout(cmd) })
	if !ok || msg == nil {
		return
	}

	switch msg := msg.(type) {
	case tea.QuitMsg:
		return
	case tea.BatchMsg:
		for _, c := range msg {
			h.drain(c, depth+1)
		}
		return
	}

	var next tea.Cmd
	h.guard("Update", func() { _, next = h.m.Update(msg) })
	h.drain(next, depth+1)
}

func (h *renderHarness) render() string {
	h.t.Helper()
	var view string
	h.guard("View", func() { view = h.m.View() })

	if strings.TrimSpace(stripANSI(view)) == "" {
		h.t.Fatalf("View is empty after keys %v", h.trail)
	}
	if strings.Contains(strings.ToLower(view), "panic") {
		h.t.Fatalf("View contains 'panic' after keys %v:\n%s", h.trail, view)
	}
	if h.m.ViewManager.Width > 0 {
		h.assertMaxLineWidth(view, h.m.ViewManager.Width)
	}
	if h.m.ViewManager.Height > 0 {
		h.assertMaxLineCount(view, h.m.ViewManager.Height)
	}
	return view
}

func maxLineWidth(view string) int {
	max := 0
	for _, line := range strings.Split(view, "\n") {
		if w := lipgloss.Width(line); w > max {
			max = w
		}
	}
	return max
}

func (h *renderHarness) assertMaxLineWidth(view string, max int) {
	h.t.Helper()
	for i, line := range strings.Split(view, "\n") {
		if w := lipgloss.Width(line); w > max {
			h.t.Errorf("line %d is %d cells wide, exceeding %d, after keys %v", i, w, max, h.trail)
		}
	}
}

func (h *renderHarness) assertMaxLineCount(view string, max int) {
	h.t.Helper()
	if n := len(strings.Split(strings.TrimSuffix(view, "\n"), "\n")); n > max {
		h.t.Errorf("view is %d lines tall, exceeding %d, after keys %v", n, max, h.trail)
	}
}

func longTasksFile(taskCount int) string {
	var b strings.Builder
	b.WriteString("# Long Project\n\n")
	for i := 1; i <= taskCount; i++ {
		fmt.Fprintf(&b, "- [ ] Task number %03d with a deliberately long description that keeps going well past the width of the details pane so that it has to wrap onto further lines\n", i)
	}
	return b.String()
}

func openLongFile(t *testing.T, taskCount int) *renderHarness {
	t.Helper()
	h := newRenderHarness(t)
	writeNotesFile(t, filepath.Join(os.Getenv("HOME"), "Notes", "lifeplus", "tasks", "long-project.md"), longTasksFile(taskCount))

	h.press("L", "enter")
	for i := 0; i < 20 && h.m.FileManager.SelectedFile.Name != "long-project.md"; i++ {
		h.press("j")
	}
	if h.m.FileManager.SelectedFile.Name != "long-project.md" {
		t.Fatalf("could not select long-project.md, selected %q", h.m.FileManager.SelectedFile.Name)
	}
	return h
}

func TestModelLoadsTaskFilesFromNotesDirectory(t *testing.T) {
	h := newRenderHarness(t)

	if h.m.DirectoryManager.SelectedCompany.DisplayName != "Lifeplus" {
		t.Fatalf("expected default company Lifeplus, got %q", h.m.DirectoryManager.SelectedCompany.DisplayName)
	}
	if got := len(h.m.TaskManager.TaskCollection.GetAll()); got != 8 {
		t.Errorf("expected 8 tasks parsed from notes, got %d", got)
	}
	if !h.m.ViewManager.Ready {
		t.Error("expected model to be ready after WindowSizeMsg")
	}
}

func TestModelRendersCategoriesView(t *testing.T) {
	h := newRenderHarness(t)

	view := stripANSI(h.render())
	if !h.m.IsCategoryView() {
		t.Fatalf("expected categories view, got %q", h.m.ViewManager.CurrentView)
	}
	for _, want := range []string{"tasks", "meetings"} {
		if !strings.Contains(strings.ToLower(view), want) {
			t.Errorf("expected categories view to contain %q:\n%s", want, view)
		}
	}

	h.press("j", "k", "j", "k")
	if h.m.DirectoryManager.CategoriesCursor != 0 {
		t.Errorf("expected categories cursor back at 0, got %d", h.m.DirectoryManager.CategoriesCursor)
	}
}

func TestModelRendersDetailsViewWithFilesAndTasks(t *testing.T) {
	h := newRenderHarness(t)

	h.press("enter")
	if !h.m.IsDetailsView() {
		t.Fatalf("expected details view after enter, got %q", h.m.ViewManager.CurrentView)
	}
	view := stripANSI(h.render())
	if !strings.Contains(view, "project-alpha") && !strings.Contains(view, "Project Alpha") {
		t.Errorf("expected details view to list project-alpha:\n%s", view)
	}

	h.press("j", "k", "tab")
	h.press("j", "j", "k")
	h.press("shift+tab", "esc")
	h.press("l", "h")
}

func TestModelRendersKanbanViewViaSidebarToggle(t *testing.T) {
	h := newRenderHarness(t)

	h.press("f")
	if !h.m.ViewManager.HideSidebar || !h.m.IsCategoryView() {
		t.Fatal("expected kanban (category view with hidden sidebar) after f")
	}
	view := stripANSI(h.render())
	for _, want := range []string{"Fix the bug", "Write the docs", "Close the sprint"} {
		if !strings.Contains(view, want) {
			t.Errorf("expected kanban to show %q:\n%s", want, view)
		}
	}

	h.press("j", "k", "l", "j", "l", "h", "h")
	h.press("f")
	if h.m.ViewManager.HideSidebar {
		t.Error("expected sidebar visible after toggling f again")
	}
}

func TestModelRendersCalendarView(t *testing.T) {
	h := newRenderHarness(t)

	h.press("enter", "c")
	if !h.m.ViewManager.IsCalendarView {
		t.Fatal("expected calendar view after c")
	}
	h.press("j", "k", "+", "-")
	h.press("c")
	if h.m.ViewManager.IsCalendarView {
		t.Error("expected calendar view closed after second c")
	}
}

func TestModelRendersDailyAndWeeklySummaries(t *testing.T) {
	h := newRenderHarness(t)

	h.press("+", "-")
	h.press("w")
	if !h.m.ViewManager.IsWeeklyView {
		t.Fatal("expected weekly view after w")
	}
	h.press("+", "-")
	h.press("w")
	if h.m.ViewManager.IsWeeklyView {
		t.Error("expected weekly view closed after second w")
	}
}

func TestModelRendersFilterInput(t *testing.T) {
	h := newRenderHarness(t)

	h.press("f", "/")
	if !h.m.IsFilterView() {
		t.Fatal("expected filter view after /")
	}
	h.press("d", "o", "c", "s", "j", "k")
	if !h.m.IsFilterView() {
		t.Fatal("expected filter view to capture typed keys instead of navigating")
	}

	h.press("esc")
	if h.m.IsFilterView() {
		t.Error("expected filter view closed after esc")
	}
}

func TestModelFiltersKanbanByTypedText(t *testing.T) {
	h := newRenderHarness(t)

	h.press("f", "/", "d", "o", "c", "s")
	if got := h.m.TaskManager.TaskCollection.FilterValue; got != "docs" {
		t.Fatalf("expected filter value %q, got %q", "docs", got)
	}
	view := stripANSI(h.render())
	if !strings.Contains(view, "Write the docs") || strings.Contains(view, "Close the sprint") {
		t.Errorf("expected kanban filtered to matching task only:\n%s", view)
	}
}

func TestModelRendersRetrospectiveForCompletedTask(t *testing.T) {
	h := newRenderHarness(t)

	h.press("f", "l", "l")
	if h.m.ViewManager.KanbanListCursor != 2 {
		t.Fatalf("expected kanban cursor on completed list, got %d", h.m.ViewManager.KanbanListCursor)
	}
	h.press("r")
	if !h.m.ViewManager.IsRetrospectiveView {
		t.Fatalf("expected retrospective view for a completed task, selected %+v", h.m.TaskManager.SelectedTask)
	}

	view := h.render()
	if !strings.Contains(stripANSI(view), "Task Retrospective") {
		t.Errorf("expected retrospective title:\n%s", view)
	}

	h.press("r")
	if h.m.ViewManager.IsRetrospectiveView {
		t.Error("expected retrospective closed after second r")
	}
}

func TestModelSwitchesCompanies(t *testing.T) {
	h := newRenderHarness(t)

	h.press("n")
	if got := h.m.DirectoryManager.SelectedCompany.DisplayName; got != "Clerky" {
		t.Fatalf("expected Clerky after n, got %q", got)
	}
	h.press("enter")
	view := stripANSI(h.render())
	if !strings.Contains(view, "clerky-work") {
		t.Errorf("expected Clerky files in details view:\n%s", view)
	}

	h.press("esc", "L")
	if got := h.m.DirectoryManager.SelectedCompany.DisplayName; got != "Lifeplus" {
		t.Errorf("expected Lifeplus after L, got %q", got)
	}
	h.press("C")
	if got := h.m.DirectoryManager.SelectedCompany.DisplayName; got != "Clerky" {
		t.Errorf("expected Clerky after C, got %q", got)
	}
	h.press("Q", "1", "3")
}

func TestModelQuitsOnQAndCtrlC(t *testing.T) {
	for _, key := range []string{"q", "ctrl+c"} {
		t.Run(key, func(t *testing.T) {
			h := newRenderHarness(t)
			cmd := h.press(key)
			if cmd == nil {
				t.Fatalf("expected quit command for %s, got nil", key)
			}
			if _, ok := cmd().(tea.QuitMsg); !ok {
				t.Errorf("expected tea.QuitMsg for %s", key)
			}
		})
	}
}

func TestModelSurvivesResizeAcrossViews(t *testing.T) {
	h := newRenderHarness(t)

	sizes := []tea.WindowSizeMsg{{Width: 200, Height: 70}, {Width: 120, Height: 40}, {Width: testWidth, Height: testHeight}}
	steps := [][]string{{}, {"f"}, {"f", "enter"}, {"c"}, {"c", "f"}, {"tab"}}
	for _, keys := range steps {
		h.press(keys...)
		for _, size := range sizes {
			h.send(size, fmt.Sprintf("resize %dx%d", size.Width, size.Height))
		}
	}
}

func TestModelKeepsLongFileDetailsWithinWindowHeight(t *testing.T) {
	h := openLongFile(t, 120)

	view := stripANSI(h.render())
	if lines := len(strings.Split(strings.TrimSuffix(view, "\n"), "\n")); lines > testHeight {
		t.Errorf("expected at most %d lines, got %d", testHeight, lines)
	}
	if !strings.Contains(view, "Task number 001") {
		t.Errorf("expected the top of the file in the details pane:\n%s", view)
	}

	h.press("f")
	view = stripANSI(h.render())
	if lines := len(strings.Split(strings.TrimSuffix(view, "\n"), "\n")); lines > testHeight {
		t.Errorf("expected at most %d lines with hidden sidebar, got %d", testHeight, lines)
	}
}

func TestModelScrollsLongFileDetailsWithJAndK(t *testing.T) {
	h := openLongFile(t, 120)
	h.press("f")

	top := stripANSI(h.render())
	if h.m.Viewport.YOffset != 0 {
		t.Fatalf("expected viewport at top, got offset %d", h.m.Viewport.YOffset)
	}

	h.press("j", "j", "j")
	scrolled := h.m.Viewport.YOffset
	if scrolled <= 0 {
		t.Fatalf("expected j to scroll the viewport, got offset %d", scrolled)
	}
	moved := stripANSI(h.render())
	if moved == top {
		t.Error("expected the rendered pane to change after scrolling")
	}

	h.press("k")
	if h.m.Viewport.YOffset >= scrolled {
		t.Errorf("expected k to scroll back up from %d, got %d", scrolled, h.m.Viewport.YOffset)
	}

	h.press("k", "k", "k")
	if h.m.Viewport.YOffset != 0 {
		t.Errorf("expected viewport back at top, got offset %d", h.m.Viewport.YOffset)
	}
	if stripANSI(h.render()) != top {
		t.Error("expected the pane to match the original top view")
	}
}

func TestModelKeepsTaskCursorVisibleInTasksFocus(t *testing.T) {
	h := openLongFile(t, 120)
	h.press("f", "t")

	for i := 0; i < 60; i++ {
		h.press("j")
		cursor := h.m.TaskManager.TasksCursor
		want := fmt.Sprintf("Task number %03d", cursor+1)
		view := stripANSI(h.render())
		if !strings.Contains(view, want) {
			t.Fatalf("expected %q to be visible with cursor %d (offset %d):\n%s", want, cursor, h.m.TasksViewport.YOffset, view)
		}
	}

	if h.m.TaskManager.TasksCursor != 59 {
		t.Fatalf("expected cursor at 59, got %d", h.m.TaskManager.TasksCursor)
	}
	if h.m.TasksViewport.YOffset == 0 {
		t.Error("expected the tasks viewport to follow the cursor")
	}

	for i := 0; i < 60; i++ {
		h.press("k")
		want := fmt.Sprintf("Task number %03d", h.m.TaskManager.TasksCursor+1)
		if view := stripANSI(h.render()); !strings.Contains(view, want) {
			t.Fatalf("expected %q to be visible moving up (offset %d):\n%s", want, h.m.TasksViewport.YOffset, view)
		}
	}
}

func TestModelKeepsLongFileWithinWindowAcrossResizes(t *testing.T) {
	h := openLongFile(t, 120)
	h.press("f", "j", "j", "j", "j")

	for _, size := range []tea.WindowSizeMsg{{Width: 200, Height: 70}, {Width: 100, Height: 30}, {Width: testWidth, Height: testHeight}} {
		h.send(size, fmt.Sprintf("resize %dx%d", size.Width, size.Height))
		if h.m.Viewport.Width != h.m.ViewManager.DetailsViewWidth || h.m.Viewport.Height != h.m.ViewManager.DetailsViewHeight {
			t.Errorf("expected viewport to match details size after resize")
		}
	}

	h.press("t", "j", "j", "j", "j", "j", "j", "j", "j", "j", "j", "j", "j")
	h.send(tea.WindowSizeMsg{Width: 100, Height: 30}, "resize 100x30 in tasks focus")
	h.send(tea.WindowSizeMsg{Width: testWidth, Height: testHeight}, "resize back in tasks focus")
}

func TestModelFitsDetailsViewWithHiddenSidebarInWindow(t *testing.T) {
	h := newRenderHarness(t)

	h.press("enter", "f")
	h.assertMaxLineWidth(h.render(), testWidth)

	h.send(tea.WindowSizeMsg{Width: 120, Height: testHeight}, "resize 120 while hidden")
	h.assertMaxLineWidth(h.render(), 120)

	h.send(tea.WindowSizeMsg{Width: testWidth, Height: testHeight}, "resize back while hidden")
	h.assertMaxLineWidth(h.render(), testWidth)

	h.press("f")
	h.assertMaxLineWidth(h.render(), testWidth)
	h.send(tea.WindowSizeMsg{Width: 120, Height: testHeight}, "resize 120 while shown")
	h.assertMaxLineWidth(h.render(), 120)
}

func TestModelCompanyShortcutsLoadFiles(t *testing.T) {
	h := newRenderHarness(t)

	h.press("C", "enter")
	view := stripANSI(h.render())
	if !strings.Contains(view, "clerky-work") || strings.Contains(view, "project-alpha") {
		t.Errorf("expected only Clerky files after C:\n%s", view)
	}

	h.press("esc", "L", "enter")
	view = stripANSI(h.render())
	if !strings.Contains(view, "project-alpha") || strings.Contains(view, "clerky-work") {
		t.Errorf("expected only Lifeplus files after L:\n%s", view)
	}
}

func loadedFileDirs(h *renderHarness) []string {
	h.m.FileManager.filesMu.RLock()
	defer h.m.FileManager.filesMu.RUnlock()
	dirs := make([]string, 0, len(h.m.FileManager.Files))
	for _, f := range h.m.FileManager.Files {
		dirs = append(dirs, filepath.Base(filepath.Dir(f.FullPath)))
	}
	return dirs
}

func assertLoadedFiles(h *renderHarness, category string, count int) {
	h.t.Helper()
	dirs := loadedFileDirs(h)
	if len(dirs) != count {
		h.t.Fatalf("expected %d %s files, got %d: %v", count, category, len(dirs), dirs)
	}
	for _, d := range dirs {
		if d != category {
			h.t.Fatalf("expected only %s files, got %v", category, dirs)
		}
	}
}

func TestModelLKeyLoadsFilesForSelectedCategory(t *testing.T) {
	h := newRenderHarness(t)

	h.press("j", "l")

	if h.m.DirectoryManager.SelectedCategory != "meetings" {
		t.Fatalf("expected meetings category, got %q", h.m.DirectoryManager.SelectedCategory)
	}
	assertLoadedFiles(h, "meetings", 1)
}

func TestModelEnterLoadsFilesForSelectedCategory(t *testing.T) {
	h := newRenderHarness(t)

	h.press("j", "enter")

	if h.m.DirectoryManager.SelectedCategory != "meetings" {
		t.Fatalf("expected meetings category, got %q", h.m.DirectoryManager.SelectedCategory)
	}
	assertLoadedFiles(h, "meetings", 1)
}

func TestModelReloadsFilesWhenSwitchingCategoriesBack(t *testing.T) {
	h := newRenderHarness(t)

	h.press("j", "l")
	assertLoadedFiles(h, "meetings", 1)

	h.press("h", "k", "l")
	if h.m.DirectoryManager.SelectedCategory != "tasks" {
		t.Fatalf("expected tasks category, got %q", h.m.DirectoryManager.SelectedCategory)
	}
	assertLoadedFiles(h, "tasks", 2)
}

func TestRenderMarkdownProducesStyledOutput(t *testing.T) {
	doc := "# Heading\n\nSome **bold** text and a list:\n\n- one\n- two\n\n```go\nfmt.Println(\"hi\")\n```\n"

	out := renderMarkdown(doc)

	plain := stripANSI(out)
	if strings.TrimSpace(plain) == "" {
		t.Fatal("expected glamour to render non-empty markdown")
	}
	for _, want := range []string{"Heading", "bold", "one", "two", "fmt.Println"} {
		if !strings.Contains(plain, want) {
			t.Errorf("expected rendered markdown to contain %q:\n%s", want, plain)
		}
	}
	if !strings.Contains(plain, "•") {
		t.Errorf("expected list items rendered as bullets:\n%s", plain)
	}
}

func TestMarkdownStyleIsStandardName(t *testing.T) {
	if got := markdownStyle(); got != "dark" && got != "light" {
		t.Errorf("expected dark or light, got %q", got)
	}
}

func TestRenderMarkdownWithWidthWrapsAtNarrowerWidth(t *testing.T) {
	doc := strings.Repeat("word ", 60)

	narrow := stripANSI(renderMarkdownWithWidth(doc, 40))
	wide := stripANSI(renderMarkdownWithWidth(doc, 120))

	if strings.TrimSpace(narrow) == "" || strings.TrimSpace(wide) == "" {
		t.Fatal("expected non-empty output at both widths")
	}
	if got := maxLineWidth(narrow); got > 40 {
		t.Errorf("expected lines within 40 cells, got %d:\n%s", got, narrow)
	}
	if strings.Count(narrow, "\n") <= strings.Count(wide, "\n") {
		t.Errorf("expected narrower width to produce more lines")
	}
}

func TestStyleHelpersRenderInputText(t *testing.T) {
	const text = "vision"
	styles := map[string]lipgloss.Style{
		"defaultTextStyle":                    defaultTextStyle,
		"selectedCompanyStyle":                selectedCompanyStyle,
		"completedTextStyle":                  completedTextStyle,
		"highlightedTextStyle":                highlightedTextStyle,
		"navbarTextStyle":                     navbarTextStyle("#FF0000"),
		"navbarContainerStyle":                navbarContainerStyle(40),
		"navbarStyle":                         navbarStyle(40),
		"sidebarStyle":                        sidebarStyle(30, 10),
		"summaryContainerStyle":               summaryContainerStyle(40, 10),
		"summaryContainerStyleNoBorder":       summaryContainerStyleNoBorder(40, 10),
		"contentContainerStyle":               contentContainerStyle(40, 10),
		"contentContainerStyleNoBorder":       contentContainerStyleNoBorder(40, 10),
		"taskSummaryContainerStyle":           taskSummaryContainerStyle(40, 10),
		"taskTitleContainer":                  taskTitleContainer(40),
		"companiesContainerStyle":             companiesContainerStyle(40),
		"listContainerStyle focused":          listContainerStyle(30, 10, true),
		"listContainerStyle unfocused":        listContainerStyle(30, 10, false),
		"filesItemDetailsContainerStyle":      filesItemDetailsContainerStyle(40, 10),
		"progressTextStyle":                   progressTextStyle(defaultTextStyle),
		"taskStyle":                           taskStyle(40),
		"kanbanTaskStyle":                     kanbanTaskStyle(40),
		"highlightedKanbanTaskStyle":          highlightedKanbanTaskStyle(40),
		"kanbanBoardTitleStyle":               kanbanBoardTitleStyle(completedColor),
		"datesContainerStyle":                 datesContainerStyle(40),
		"dateStyle":                           dateStyle(),
		"tasksStringStyle":                    tasksStringStyle(40),
		"summaryTitleStyle":                   summaryTitleStyle(40),
		"inactiveTitleStyle":                  inactiveTitleStyle(),
		"completedTitleStyle":                 completedTitleStyle(),
		"renderedActiveListStyle":             renderedActiveListStyle(),
		"renderedInactiveListStyle":           renderedInactiveListStyle(),
		"renderedCompletedListStyle":          renderedCompletedListStyle(),
		"taskItemDetailsContainerStyle":       taskItemDetailsContainerStyle(40, 10),
		"boardContainerStyle selected":        boardContainerStyle(40, 10, true),
		"boardContainerStyle unselected":      boardContainerStyle(40, 10, false),
		"filterInputStyle":                    filterInputStyle("#FFFFFF"),
		"taskOptionSuggestionsContainerStyle": taskOptionSuggestionsContainerStyle(),
	}

	for name, style := range styles {
		t.Run(name, func(t *testing.T) {
			out := style.Render(text)
			if !strings.Contains(stripANSI(out), text) {
				t.Errorf("expected %s to render %q, got %q", name, text, out)
			}
		})
	}

	joined := stripANSI(joinVertical("top", "bottom") + joinHorizontal("left", "right"))
	for _, want := range []string{"top", "bottom", "left", "right"} {
		if !strings.Contains(joined, want) {
			t.Errorf("expected joined output to contain %q, got %q", want, joined)
		}
	}
}

func TestDerivedStylesLeaveBaseStylesUnchanged(t *testing.T) {
	baseFg := taskFileTitleStyle.GetForeground()
	_ = inactiveTitleStyle()
	_ = completedTitleStyle()
	if taskFileTitleStyle.GetForeground() != baseFg {
		t.Errorf("expected taskFileTitleStyle foreground to stay %v, got %v", baseFg, taskFileTitleStyle.GetForeground())
	}

	baseWidth := defaultTextStyle.GetWidth()
	_ = progressTextStyle(defaultTextStyle)
	if defaultTextStyle.GetWidth() != baseWidth {
		t.Errorf("expected defaultTextStyle width to stay %d, got %d", baseWidth, defaultTextStyle.GetWidth())
	}

	unfocused := listContainerStyle(30, 10, false)
	focused := listContainerStyle(30, 10, true)
	if unfocused.GetBorderTopForeground() == focused.GetBorderTopForeground() {
		t.Error("expected unfocused list container to have a distinct border color from focused")
	}
}

func writeManyMeetingFiles(t *testing.T, count int) {
	t.Helper()
	dir := filepath.Join(os.Getenv("HOME"), "Notes", "lifeplus", "meetings")
	for i := 1; i <= count; i++ {
		writeNotesFile(t, filepath.Join(dir, fmt.Sprintf("meeting-%03d.md", i)), fmt.Sprintf("# Meeting %d\n\n- Notes\n", i))
	}
}

func openManyMeetings(t *testing.T, count int) *renderHarness {
	t.Helper()
	h := newRenderHarness(t)
	writeManyMeetingFiles(t, count)
	h.press("L", "j", "enter")
	return h
}

func selectedDisplayName(h *renderHarness) string {
	return h.m.FileManager.SelectedFile.FileNameWithoutExtension(h.m.FileManager.FileExtension)
}

func sidebarShows(h *renderHarness, view, name string) bool {
	for _, line := range strings.Split(view, "\n") {
		runes := []rune(line)
		if n := h.m.ViewManager.SidebarWidth; len(runes) > n {
			runes = runes[:n]
		}
		if strings.Contains(string(runes), name) {
			return true
		}
	}
	return false
}

func TestModelKeepsManyFilesSidebarWithinWindow(t *testing.T) {
	h := openManyMeetings(t, 80)

	if got := len(h.m.FileManager.Files); got != 81 {
		t.Fatalf("expected 81 meeting files, got %d", got)
	}
	h.render()
}

func TestModelKeepsFileCursorVisibleInSidebar(t *testing.T) {
	h := openManyMeetings(t, 80)

	for i := 0; i < 60; i++ {
		h.press("j")
		name := selectedDisplayName(h)
		if view := stripANSI(h.render()); !sidebarShows(h, view, name) {
			t.Fatalf("selected file %q not visible after %d presses of j:\n%s", name, i+1, view)
		}
	}
	for i := 0; i < 60; i++ {
		h.press("k")
		name := selectedDisplayName(h)
		if view := stripANSI(h.render()); !sidebarShows(h, view, name) {
			t.Fatalf("selected file %q not visible after %d presses of k:\n%s", name, i+1, view)
		}
	}
}

func TestModelListsAllTaskFilesInSidebar(t *testing.T) {
	h := newRenderHarness(t)
	dir := filepath.Join(os.Getenv("HOME"), "Notes", "lifeplus", "tasks")
	for i := 1; i <= 25; i++ {
		writeNotesFile(t, filepath.Join(dir, fmt.Sprintf("extra-%02d.md", i)), fmt.Sprintf("# Extra %d\n\n- [ ] Item\n", i))
	}
	h.press("L", "enter")

	total := len(h.m.FileManager.Files)
	if total != 27 {
		t.Fatalf("expected 27 task files, got %d", total)
	}
	for i := 0; i < total-1; i++ {
		h.press("j")
	}

	name := selectedDisplayName(h)
	if view := stripANSI(h.render()); !sidebarShows(h, view, name) {
		t.Fatalf("last task file %q not visible:\n%s", name, view)
	}
}

func TestModelSelectsFirstFileWhenCategoryLoads(t *testing.T) {
	h := newRenderHarness(t)
	h.press("j", "l")

	if got := h.m.FileManager.SelectedFile.Name; got != "2024-01-15.md" {
		t.Fatalf("expected 2024-01-15.md selected, got %q", got)
	}
	if view := stripANSI(h.render()); !strings.Contains(view, "meetings > 2024-01-15") {
		t.Fatalf("navbar does not show the selected file:\n%s", view)
	}
}
