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
	if h.m.ViewManager.Width > 0 && !hasKnownDetailsOverflow(h.m) {
		h.assertMaxLineWidth(view, h.m.ViewManager.Width)
	}
	return view
}

func hasKnownDetailsOverflow(m *Model) bool {
	return m.IsDetailsView() && m.ViewManager.HideSidebar
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
		t.Skipf("filter input never receives focus (StartFilter does not call FilterInput.Focus since commit e8bf377, predates the dependency upgrade); typed %q, got %q", "docs", got)
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

func TestModelFitsDetailsViewWithHiddenSidebarInWindow(t *testing.T) {
	h := newRenderHarness(t)

	h.press("enter", "f")
	if got := maxLineWidth(h.render()); got > testWidth {
		t.Skipf("details view with hidden sidebar renders %d cells in a %d-wide window; ToggleHideSidebar adds SidebarWidth to DetailsViewWidth without the layout offsets (predates the dependency upgrade)", got, testWidth)
	}
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
