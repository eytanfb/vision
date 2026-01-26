package app

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func hashString(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func (m *Model) View() string {
	return ViewHandler(m)
}

func ViewHandler(m *Model) string {
	// On window resize or Everything flag, force full re-render by clearing viewport cache
	if m.ViewManager.Dirty.Everything {
		m.ViewCache = ViewCache{} // Clear all caches
		// Note: We don't call ClearDirty() here because dirty flags are not yet
		// used for selective rendering optimization. They are set by Update()
		// but currently only Everything is checked above.
	}

	// Check for retrospective view early - it takes over the entire screen
	if m.IsRetrospectiveView() {
		return renderRetrospectiveView(m)
	}

	content := "Something is wrong"

	if m.IsCategoryView() {
		content = renderList(m, m.CategoryNames())
	} else if m.IsDetailsView() {
		if m.IsTaskDetailsFocus() {
			Debug("Task details focus")
			content = renderTasks(m)
		} else {
			content = renderFiles(m)
		}
	}

	result := joinVertical(renderNavbar(m), renderFilterInput(m), content, renderErrors(m))

	return result
}

func renderRetrospectiveView(m *Model) string {
	selectedTask := m.TaskManager.SelectedTask

	// Get the earliest date as the creation date fallback
	createdDate := getTaskCreationDate(selectedTask)

	retro := NewTaskRetrospective(selectedTask, createdDate)
	if retro == nil {
		return "Unable to generate retrospective - task may not be completed"
	}

	retroView := NewRetrospectiveView(retro, m.ViewManager.Width, m.ViewManager.Height)

	// Center the retrospective view content
	content := retroView.View()
	centeredContent := lipgloss.Place(
		m.ViewManager.Width,
		m.ViewManager.Height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)

	return centeredContent
}

func getTaskCreationDate(task Task) time.Time {
	// Try to find the earliest date from the task's dates
	var earliestDate time.Time

	if task.ScheduledDate != "" {
		if d, err := time.Parse("2006-01-02", task.ScheduledDate); err == nil {
			if earliestDate.IsZero() || d.Before(earliestDate) {
				earliestDate = d
			}
		}
	}

	if task.StartDate != "" {
		if d, err := time.Parse("2006-01-02", task.StartDate); err == nil {
			if earliestDate.IsZero() || d.Before(earliestDate) {
				earliestDate = d
			}
		}
	}

	if task.CompletedDate != "" {
		if d, err := time.Parse("2006-01-02", task.CompletedDate); err == nil {
			if earliestDate.IsZero() || d.Before(earliestDate) {
				earliestDate = d
			}
		}
	}

	// If no dates found, use today as fallback
	if earliestDate.IsZero() {
		earliestDate = time.Now()
	}

	return earliestDate
}

func renderErrors(m *Model) string {
	var errors strings.Builder
	for _, err := range m.Errors {
		errors.WriteString(err + "\n")
	}

	return errors.String()
}

func renderCompanies(m *Model, companies []string) string {
	result := ""

	for index, company := range companies {
		textStyle := companyTextStyle
		if index == m.DirectoryManager.CompaniesCursor {
			textStyle = selectedCompanyStyle
			result = joinHorizontal(result, textStyle.Render("["+company+"]"))
		}
	}

	return companiesContainerStyle(m.ViewManager.Width).Render(result)
}

func renderList(m *Model, items []string) string {
	list := ""

	cursor := m.GetCurrentCursor()

	for index, item := range items {
		list = joinVertical(list, createListItem(item, index, cursor))
	}

	sidebar := sidebarStyle(m.ViewManager.SidebarWidth, m.ViewManager.SidebarHeight).Render(joinVertical(list))

	summaryView := buildSummaryView(m, m.ViewManager.HideSidebar)

	// Create a temporary viewport for rendering - don't mutate model state
	tempViewport := viewport.Model{
		Width:  m.Viewport.Width,
		Height: m.Viewport.Height,
	}

	// Only update viewport content if it changed
	// NOTE: This mutates ViewCache in View(), which technically violates
	// the pure View() principle. However, this is safe because:
	// 1. Bubble Tea's Update/View cycle is single-threaded
	// 2. This is purely a rendering optimization (doesn't affect business logic)
	// 3. The cache state doesn't influence application behavior
	hash := hashString(summaryView)
	if hash != m.ViewCache.LastSummaryViewHash {
		tempViewport.SetContent(summaryView)
		m.ViewCache.LastSummaryViewHash = hash
		m.ViewCache.LastViewportContent = summaryView
	} else {
		// Reuse cached content
		tempViewport.SetContent(m.ViewCache.LastViewportContent)
	}

	summary := summaryContainerStyle(m.ViewManager.DetailsViewWidth, m.ViewManager.SummaryViewHeight).Render(tempViewport.View())

	if m.ViewManager.HideSidebar {
		sidebar = ""
		summary = summaryContainerStyleNoBorder(m.ViewManager.DetailsViewWidth, m.ViewManager.DetailsViewHeight).Render(summaryView)
	}

	return joinHorizontal(sidebar, summary)
}

func buildSummaryView(m *Model, hiddenSidebar bool) string {
	summaryView := ""
	period := "daily"

	if m.IsAddTaskView() || m.IsAddSubTaskView() {
		summaryView = m.NewTaskInput.View() + "\n"

		if hasUnclosedDoubleSquareBrackets(m.NewTaskInput.Value()) {
			// IsSuggestionsActive is managed in Update(), not View()
			optionsView := linkingSuggestionsDisplay(m)
			summaryView = joinVertical(summaryView, optionsView)
		}
	} else {
		if m.ViewManager.IsWeeklyView {
			period = "weekly"
		}

		if hiddenSidebar {
			summaryView = kanbanSummaryView(m, period)
		} else {
			summaryView = taskSummaryToView(m, period)
		}
	}

	return summaryView
}

func kanbanSummaryView(m *Model, period string) string {
	tasksByFile, summaryDate := setDailySummaryValues(m)

	if period == "weekly" {
		tasksByFile, summaryDate = setWeeklySummaryValues(m)
	}

	keys := sortTaskKeys(tasksByFile)
	viewSort(keys, m)

	view := BuildKanbanSummaryView(m, keys, tasksByFile, m.ViewManager.DetailsViewWidth, summaryDate)

	return view
}

func taskSummaryToView(m *Model, period string) string {
	if m.ViewManager.IsCalendarView {
		calendarView := NewCalendarView(
			m.TaskManager.TaskCollection.allTasks(),
			m.ViewManager.DetailsViewWidth,
			m.ViewManager.DetailsViewHeight,
		)
		return calendarView.View()
	}

	tasksByFile, summaryDate := setDailySummaryValues(m)

	if period == "weekly" {
		tasksByFile, summaryDate = setWeeklySummaryValues(m)
	}

	keys := sortTaskKeys(tasksByFile)
	viewSort(keys, m)

	height := m.ViewManager.SummaryViewHeight
	containerTitleHeight := 2
	viewHeight := height - containerTitleHeight

	view := BuildSummaryView(m, keys, tasksByFile, m.ViewManager.DetailsViewWidth, summaryDate)

	containerTitle := taskSummaryContainerStyle(m.ViewManager.DetailsViewWidth, containerTitleHeight).Height(2).PaddingBottom(0).Render(summaryTitle(m, period))
	renderedView := taskSummaryContainerStyle(m.ViewManager.DetailsViewWidth, viewHeight).PaddingTop(0).Render(view)

	return joinVertical(containerTitle, renderedView)
}

func setDailySummaryValues(m *Model) (map[string][]Task, string) {
	tasksByFile := m.TaskManager.Summary(m.GetCurrentCompanyName())
	summaryDate := m.TaskManager.DailySummaryDate

	return tasksByFile, summaryDate
}

func setWeeklySummaryValues(m *Model) (map[string][]Task, string) {
	startDate := m.TaskManager.WeeklySummaryStartDate
	endDate := m.TaskManager.WeeklySummaryEndDate
	tasksByFile := m.TaskManager.WeeklySummary(m.GetCurrentCompanyName(), startDate, endDate)
	summaryDate := m.TaskManager.WeeklySummaryEndDate

	return tasksByFile, summaryDate
}

func summaryTitle(m *Model, period string) string {
	dailySummaryDateText := m.TaskManager.DailySummaryDate

	today := time.Now()
	tomorrow := today.Add(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)
	dailySummaryDate, _ := time.Parse("2006-01-02", dailySummaryDateText)

	if dailySummaryDate.Day() == today.Day() && dailySummaryDate.Month() == today.Month() && dailySummaryDate.Year() == today.Year() {
		dailySummaryDateText = "Today"
	} else if dailySummaryDate.Day() == tomorrow.Day() && dailySummaryDate.Month() == tomorrow.Month() && dailySummaryDate.Year() == tomorrow.Year() {
		dailySummaryDateText = "Tomorrow"
	} else if dailySummaryDate.Day() == yesterday.Day() && dailySummaryDate.Month() == yesterday.Month() && dailySummaryDate.Year() == yesterday.Year() {
		dailySummaryDateText = "Yesterday"
	} else if dailySummaryDate.After(today) {
		diff := dailySummaryDate.Sub(today)
		diffInDaysString := int(diff.Hours() / 24)

		dailySummaryDateText = fmt.Sprintf("%s (+%d)", dailySummaryDate.Format("Monday, January 2"), diffInDaysString)
	} else {
		diff := today.Sub(dailySummaryDate)
		diffInDaysString := int(diff.Hours() / 24)

		dailySummaryDateText = fmt.Sprintf("%s (-%d)", dailySummaryDate.Format("Monday, January 2"), diffInDaysString)
	}

	title := "Daily Tasks for " + dailySummaryDateText
	if period == "weekly" {
		title = "Weekly Tasks for " + m.TaskManager.WeeklySummaryStartDate + " - " + m.TaskManager.WeeklySummaryEndDate
	}

	return title
}

func sortTaskKeys(tasksByFile map[string][]Task) []string {
	keys := make([]string, 0, len(tasksByFile))
	for k := range tasksByFile {
		keys = append(keys, k)
	}

	return keys
}

func renderFiles(m *Model) string {
	if !m.HasFiles() {
		return "No files found"
	}

	listContainerStyle := listContainerStyle(m.ViewManager.SidebarWidth, m.ViewManager.SidebarHeight, m.IsItemDetailsFocus())
	listContainer := ""

	list, itemDetails := BuildFilesView(m, m.ViewManager.HideSidebar)

	itemDetailsContainer := filesItemDetailsContainerStyle(m.ViewManager.DetailsViewWidth, m.ViewManager.DetailsViewHeight).Render(itemDetails)

	if list != "" {
		listContainer = listContainerStyle.Render(list)
	}

	container := joinHorizontal(listContainer, itemDetailsContainer)

	return joinVertical(container)
}

func renderNavbar(m *Model) string {
	companyColor := m.DirectoryManager.SelectedCompany.Color
	textStyle := navbarTextStyle(companyColor)
	container := navbarContainerStyle(m.ViewManager.NavbarWidth)
	style := navbarStyle(m.ViewManager.NavbarWidth)
	navbar := textStyle.Render("Vision")

	if m.IsCategoryView() {
		navbar = textStyle.Render(m.GetCurrentCompanyName())
	} else if m.IsDetailsView() {
		navbar = textStyle.Render(m.GetCurrentCompanyName() + " > " + m.DirectoryManager.SelectedCategory + " > " + m.FileManager.SelectedFile.Name)
	}

	navbarView := joinVertical(style.Render(navbar))

	if m.ViewManager.ShowCompanies {
		navbarView = joinHorizontal(navbar, renderCompanies(m, m.CategoryNames()))
	}

	view := container.Render(navbarView)

	return view
}

func renderFilterInput(m *Model) string {
	// Focus is managed in Update(), not View()
	// The input will be focused when FilterModeEnteredMsg is handled
	return filterInputStyle(m.DirectoryManager.SelectedCompany.Color).Render(m.FilterInput.View())
}

func renderTasks(m *Model) string {
	if m.ViewManager.IsAddTaskView || m.ViewManager.IsAddSubTaskView {
		addTaskView := m.NewTaskInput.View()

		if hasUnclosedDoubleSquareBrackets(m.NewTaskInput.Value()) {
			// IsSuggestionsActive is set in Update(), not View()
			optionsView := linkingSuggestionsDisplay(m)
			addTaskView = joinVertical(addTaskView, optionsView)
		}

		// Create temporary viewport for rendering
		tempViewport := viewport.Model{
			Width:  m.Viewport.Width,
			Height: m.Viewport.Height,
		}
		tempViewport.SetContent("Add new subtask\n" + addTaskView + "\n")

		return contentContainerStyle(m.ViewManager.DetailsViewWidth, m.ViewManager.DetailsViewHeight).Render(tempViewport.View())
	}

	date := m.TaskManager.DailySummaryDate

	tasks := m.TaskManager.TaskCollection.GetTasks(m.FileManager.currentFileName())

	tasksView := BuildTasksForFileView(m, tasks, date, m.TaskManager.TasksCursor)

	renderedTasksView := taskItemDetailsContainerStyle(m.ViewManager.DetailsViewWidth, m.ViewManager.DetailsViewHeight).Render(tasksView)

	return joinVertical(renderedTasksView)
}

func renderMarkdown(content string) string {
	out, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
	)

	markdown, err := out.Render(content)
	if err != nil {
		return ""
	}

	return markdown
}

type fileSortKey struct {
	filename   string
	percentage int
	isInactive bool
	isComplete bool
}

func viewSort(filenames []string, m *Model) {
	// Create a map of filename to FileInfo for fast lookup
	m.FileManager.filesMu.RLock()
	fileMap := make(map[string]*FileInfo)
	for i := range m.FileManager.Files {
		fileMap[m.FileManager.Files[i].Name] = &m.FileManager.Files[i]
	}
	m.FileManager.filesMu.RUnlock()

	// Pre-calculate all sort keys once using cached progress - O(n) instead of O(n log n)
	sortKeys := make([]fileSortKey, len(filenames))
	for i, filename := range filenames {
		var progress ProgressInfo

		// Try to use cached progress if available
		if file, ok := fileMap[filename]; ok {
			progress = m.FileManager.getOrCalculateProgress(file, &m.TaskManager)
		} else {
			// Fallback to direct calculation if file not found
			completed, total := m.TaskManager.TaskCollection.Progress(filename)
			progress.Percentage = 0
			if total > 0 {
				progress.Percentage = int(float64(completed) / float64(total) * 100)
			}
			progress.IsInactive = m.TaskManager.TaskCollection.IsInactive(filename)
			progress.IsComplete = total > 0 && completed == total
		}

		sortKeys[i] = fileSortKey{
			filename:   filename,
			percentage: progress.Percentage,
			isInactive: progress.IsInactive,
			isComplete: progress.IsComplete,
		}
	}

	// Sort using pre-calculated keys - O(n log n) comparisons only
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

func createListItem(item string, index int, cursor int) string {
	line := ""
	style := defaultTextStyle

	if index == cursor {
		style = highlightedTextStyle
	}

	line += item + "\n"

	return style.Render(line)
}

func hasUnclosedDoubleSquareBrackets(input string) bool {
	openBrackets := 0

	// Iterate over the string
	for i := 0; i < len(input)-1; i++ {
		if input[i:i+2] == "[[" {
			openBrackets++ // Increment for each opening [[
			i++            // Skip the next character to avoid re-checking
		} else if input[i:i+2] == "]]" {
			if openBrackets > 0 {
				openBrackets-- // Decrement for each closing ]] if there's an opening
			}
			i++ // Skip the next character
		}
	}

	return openBrackets > 0
}

func peopleFilterValue(input string) string {
	openBrackets := 0
	lastOpenIndex := -1

	// Iterate over the string to find the latest unmatched [[
	for i := 0; i < len(input)-1; i++ {
		if input[i:i+2] == "[[" {
			openBrackets++
			lastOpenIndex = i + 2 // Position right after this [[
			i++                   // Skip the next character to avoid re-checking
		} else if input[i:i+2] == "]]" {
			if openBrackets > 0 {
				openBrackets-- // Close the most recent unmatched [[
			}
			i++ // Skip the next character
		}
	}

	// If no unmatched [[ is found, return an empty string
	if openBrackets == 0 || lastOpenIndex == -1 {
		return ""
	}

	// Extract characters after the last unmatched [[
	var result []rune
	for _, r := range input[lastOpenIndex:] {
		if unicode.IsLetter(r) || unicode.IsSpace(r) { // Only include alphabetic characters
			result = append(result, r)
		}
	}

	return string(result)
}

func linkingSuggestionsDisplay(m *Model) string {
	filterValue := peopleFilterValue(m.NewTaskInput.Value())
	// SuggestionsFilterValue is set elsewhere if needed for state persistence

	peopleOptions := m.FileManager.PeopleFilenames(&m.DirectoryManager, &m.TaskManager, filterValue)
	taskOptions := m.FileManager.TaskFilenames(&m.DirectoryManager, &m.TaskManager, filterValue)

	var peopleBuilder strings.Builder
	for index, option := range peopleOptions {
		person := strings.Split(option, m.FileManager.FileExtension)[0]

		style := suggestionTextStyle

		if m.ViewManager.SuggestionsListsCursor == 0 && index == m.ViewManager.SuggestionCursor {
			style = selectedSuggestionTextStyle
		}

		peopleBuilder.WriteString(style.Render(person))
		if index < len(peopleOptions)-1 {
			peopleBuilder.WriteRune('\n')
		}
	}

	peopleOptionsView := ""
	if peopleBuilder.Len() > 0 {
		peopleOptionViewTitle := suggestionTitleStyle.Render("People")
		peopleOptionsView = peopleOptionViewTitle + "\n" + peopleBuilder.String()
	}

	var taskBuilder strings.Builder
	for index, option := range taskOptions {
		task := strings.Split(option, m.FileManager.FileExtension)[0]

		style := suggestionTextStyle

		if m.ViewManager.SuggestionsListsCursor == 1 && index == m.ViewManager.SuggestionCursor {
			style = selectedSuggestionTextStyle
		}

		taskBuilder.WriteString(style.Render(task))
		if index < len(taskOptions)-1 {
			taskBuilder.WriteRune('\n')
		}
	}

	taskOptionsView := ""
	if taskBuilder.Len() > 0 {
		taskOptionViewTitle := suggestionTitleStyle.Render("Tasks")
		taskOptionsView = taskOptionViewTitle + "\n" + taskBuilder.String()
	}

	return joinHorizontal(peopleOptionsView, taskOptionSuggestionsContainerStyle().Render(taskOptionsView))
}
