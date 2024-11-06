package app

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/log"
)

func (m *Model) View() string {
	return ViewHandler(m)
}

func ViewHandler(m *Model) string {
	content := "Something is wrong"

	if m.IsCategoryView() {
		content = renderList(m, m.CategoryNames())
	} else if m.IsDetailsView() {
		if m.IsTaskDetailsFocus() {
			log.Info("Task details focus")
			content = renderTasks(m)
		} else {
			content = renderFiles(m)
		}
	}

	return joinVertical(renderNavbar(m), renderFilterInput(m), content, renderErrors(m))
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

	m.Viewport.SetContent(summaryView)

	summary := summaryContainerStyle(m.ViewManager.DetailsViewWidth, m.ViewManager.SummaryViewHeight).Render(m.Viewport.View())

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
	if m.ViewManager.IsFilterView {
		m.FilterInput.Focus()
	}

	return filterInputStyle(m.DirectoryManager.SelectedCompany.Color).Render(m.FilterInput.View())
}

func renderTasks(m *Model) string {
	if m.ViewManager.IsAddTaskView || m.ViewManager.IsAddSubTaskView {
		addTaskView := m.NewTaskInput.View()

		if hasUnclosedDoubleSquareBrackets(m.NewTaskInput.Value()) {
			optionsView := linkingSuggestionsDisplay(m)
			addTaskView = joinVertical(addTaskView, optionsView)
		}

		m.Viewport.SetContent("Add new subtask\n" + addTaskView + "\n")

		return contentContainerStyle(m.ViewManager.DetailsViewWidth, m.ViewManager.DetailsViewHeight).Render(m.Viewport.View())
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

func viewSort(filenames []string, m *Model) {
	sort.Slice(filenames, func(i, j int) bool {
		iInactive := m.TaskManager.TaskCollection.IsInactive(filenames[i])
		jInactive := m.TaskManager.TaskCollection.IsInactive(filenames[j])

		if iInactive {
			return false
		}

		if jInactive {
			return true
		}
		iCompletedTasks, iTotalTasks := m.TaskManager.TaskCollection.Progress(filenames[i])
		jCompletedTasks, jTotalTasks := m.TaskManager.TaskCollection.Progress(filenames[j])

		iPercentage := float64(iCompletedTasks) / float64(iTotalTasks)
		jPercentage := float64(jCompletedTasks) / float64(jTotalTasks)

		iRoundedUp := int(iPercentage*10) * 10
		jRoundedUp := int(jPercentage*10) * 10

		if iRoundedUp == 100 {
			return false
		}

		if jRoundedUp == 100 {
			return true
		}

		if iRoundedUp == jRoundedUp {
			return filenames[i] < filenames[j]
		}

		return iRoundedUp > jRoundedUp
	})
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
	peopleOptions := m.FileManager.PeopleFilenames(&m.DirectoryManager, &m.TaskManager, filterValue)
	taskOptions := m.FileManager.TaskFilenames(&m.DirectoryManager, &m.TaskManager, filterValue)

	peopleOptionsView := ""
	for _, option := range peopleOptions {
		person := strings.Split(option, ".md")[0]
		peopleOptionsView = joinVertical(peopleOptionsView, suggestionTextStyle.Render(person))
	}

	peopleOptionViewTitle := suggestionTitleStyle.Render("People")

	if peopleOptionsView != "" {
		peopleOptionsView = joinVertical(peopleOptionViewTitle, peopleOptionsView)
	}

	taskOptionsView := ""
	for _, option := range taskOptions {
		task := strings.Split(option, ".md")[0]
		taskOptionsView = joinVertical(taskOptionsView, suggestionTextStyle.Render(task))
	}

	taskOptionViewTitle := suggestionTitleStyle.Render("Tasks")

	if taskOptionsView != "" {
		taskOptionsView = joinVertical(taskOptionViewTitle, taskOptionsView)
	}

	return joinHorizontal(peopleOptionsView, taskOptionsView)
}
