package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	timelineConnectorColor = lipgloss.Color("#6C6C6C")
	waitPeriodColor        = lipgloss.Color("#A0A0A0")
	activePeriodColor      = lipgloss.Color("#0AAFC7")
	retroCompletedColor    = lipgloss.Color("#4CD137")
)

var (
	retroTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(summaryTitleColor).
			Align(lipgloss.Center)

	retroTaskNameStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(white)

	retroMetadataStyle = lipgloss.NewStyle().
				Foreground(inactiveFileColor)

	retroPriorityStyle = lipgloss.NewStyle().
				Foreground(priorityTextColor)

	retroSectionTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(white)

	retroTimelineEventStyle = lipgloss.NewStyle().
				Foreground(white)

	retroTimelineConnectorStyle = lipgloss.NewStyle().
					Foreground(timelineConnectorColor)

	retroWaitPeriodStyle = lipgloss.NewStyle().
				Foreground(waitPeriodColor)

	retroActivePeriodStyle = lipgloss.NewStyle().
				Foreground(activePeriodColor)

	retroCompletedStyle = lipgloss.NewStyle().
				Foreground(retroCompletedColor)

	retroMetricLabelStyle = lipgloss.NewStyle().
				Foreground(white).
				Width(18)

	retroMetricValueStyle = lipgloss.NewStyle().
				Foreground(white).
				Width(10).
				Align(lipgloss.Right)

	retroMetricBarStyle = lipgloss.NewStyle().
				Width(22)

	retroMetricPercentStyle = lipgloss.NewStyle().
				Foreground(inactiveFileColor).
				Width(5).
				Align(lipgloss.Right)

	retroHelpBarStyle = lipgloss.NewStyle().
				Foreground(inactiveFileColor)

	retroDividerStyle = lipgloss.NewStyle().
				Foreground(timelineConnectorColor)
)

type RetrospectiveView struct {
	Retro  *TaskRetrospective
	Width  int
	Height int
}

func NewRetrospectiveView(retro *TaskRetrospective, width, height int) *RetrospectiveView {
	return &RetrospectiveView{
		Retro:  retro,
		Width:  width,
		Height: height,
	}
}

func (rv *RetrospectiveView) View() string {
	if rv.Retro == nil {
		return "No retrospective data available"
	}

	var sections []string

	sections = append(sections, rv.renderHeader())
	sections = append(sections, "")
	sections = append(sections, rv.renderTimeline())
	sections = append(sections, "")
	sections = append(sections, rv.renderMetrics())
	sections = append(sections, "")
	sections = append(sections, rv.renderHelpBar())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (rv *RetrospectiveView) renderHeader() string {
	var lines []string

	titleWidth := rv.Width - 4
	title := retroTitleStyle.Width(titleWidth).Render("Task Retrospective")

	divider := retroDividerStyle.Render(strings.Repeat("=", min(56, rv.Width-4)))

	lines = append(lines, title)
	lines = append(lines, divider)
	lines = append(lines, "")

	taskName := retroTaskNameStyle.Render("  " + rv.Retro.Task.Summary())
	lines = append(lines, taskName)

	var metaParts []string
	if rv.Retro.Task.FileName != "" {
		metaParts = append(metaParts, fmt.Sprintf("Project: %s", rv.Retro.Task.FileName))
	}
	if rv.Retro.Task.Company != "" {
		metaParts = append(metaParts, fmt.Sprintf("Company: %s", rv.Retro.Task.Company))
	}
	if rv.Retro.Task.Priority != "" {
		metaParts = append(metaParts, retroPriorityStyle.Render(PriorityIcon+"Priority"))
	}

	if len(metaParts) > 0 {
		metadata := retroMetadataStyle.Render("  " + strings.Join(metaParts, " | "))
		lines = append(lines, metadata)
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (rv *RetrospectiveView) renderTimeline() string {
	var lines []string

	sectionTitle := retroSectionTitleStyle.Render("Timeline")
	sectionDivider := retroDividerStyle.Render(strings.Repeat("-", min(56, rv.Width-4)))

	lines = append(lines, sectionTitle)
	lines = append(lines, sectionDivider)
	lines = append(lines, "")

	events := rv.Retro.GetTimeline()

	for i, event := range events {
		eventLine := rv.renderTimelineEvent(event)
		lines = append(lines, eventLine)

		if i < len(events)-1 {
			nextEvent := events[i+1]
			periodLines := rv.renderTimelinePeriod(event, nextEvent)
			lines = append(lines, periodLines...)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (rv *RetrospectiveView) renderTimelineEvent(event TimelineEvent) string {
	dateStr := event.Date.Format("Jan 02, 2006")
	dayStr := fmt.Sprintf("Day %d", event.DayNumber)

	var iconStyle lipgloss.Style
	switch event.EventType {
	case "created":
		iconStyle = retroTimelineEventStyle
	case "scheduled":
		iconStyle = scheduledTextStyle
	case "started":
		iconStyle = startedTextStyle
	case "completed":
		iconStyle = retroCompletedStyle
	default:
		iconStyle = retroTimelineEventStyle
	}

	icon := iconStyle.Render(event.Icon)
	label := iconStyle.Render(fmt.Sprintf("%-12s", event.Label))
	date := retroTimelineEventStyle.Render(fmt.Sprintf("%-16s", dateStr))
	day := retroMetadataStyle.Render(dayStr)

	return fmt.Sprintf("  %s %s %s %s", icon, label, date, day)
}

func (rv *RetrospectiveView) renderTimelinePeriod(from, to TimelineEvent) []string {
	var lines []string

	connector := retroTimelineConnectorStyle.Render("     |")
	lines = append(lines, connector)

	days := to.DayNumber - from.DayNumber
	if days == 0 {
		return lines
	}

	var periodType string
	var barStyle lipgloss.Style

	switch to.EventType {
	case "scheduled":
		periodType = "waiting"
		barStyle = retroWaitPeriodStyle
	case "started":
		if from.EventType == "scheduled" {
			periodType = "lead time"
		} else {
			periodType = "waiting"
		}
		barStyle = retroWaitPeriodStyle
	case "completed":
		periodType = "active"
		barStyle = retroActivePeriodStyle
	default:
		periodType = "time"
		barStyle = retroWaitPeriodStyle
	}

	barLength := min(days, 10)
	var barChar string
	if periodType == "active" {
		barChar = "\u2588" // Full block character for active
	} else {
		barChar = "\u2591" // Light shade character for waiting
	}
	bar := strings.Repeat(barChar, barLength)

	periodDesc := fmt.Sprintf("%d %s %s", days, pluralize(days, "day", "days"), periodType)
	periodLine := fmt.Sprintf("     | %s %s", barStyle.Render(bar), retroMetadataStyle.Render(periodDesc))
	lines = append(lines, periodLine)

	lines = append(lines, connector)

	return lines
}

func (rv *RetrospectiveView) renderMetrics() string {
	var lines []string

	sectionTitle := retroSectionTitleStyle.Render("Metrics")
	sectionDivider := retroDividerStyle.Render(strings.Repeat("-", min(56, rv.Width-4)))

	lines = append(lines, sectionTitle)
	lines = append(lines, sectionDivider)
	lines = append(lines, "")

	metrics := rv.Retro.GetMetrics()

	for _, metric := range metrics {
		metricLine := rv.renderMetricLine(metric)
		lines = append(lines, metricLine)
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (rv *RetrospectiveView) renderMetricLine(metric RetrospectiveMetric) string {
	label := retroMetricLabelStyle.Render("  " + metric.Label)

	var valueStr string
	if metric.Label == "Efficiency" {
		valueStr = ""
	} else {
		valueStr = fmt.Sprintf("%d %s", metric.Value, metric.Unit)
	}
	value := retroMetricValueStyle.Render(valueStr)

	bar := rv.renderMetricBar(metric)

	var percentStr string
	if metric.Label == "Total Time" {
		percentStr = ""
	} else {
		percentStr = fmt.Sprintf("%.0f%%", metric.Percentage)
	}
	percent := retroMetricPercentStyle.Render(percentStr)

	return fmt.Sprintf("%s%s   %s %s", label, value, bar, percent)
}

func (rv *RetrospectiveView) renderMetricBar(metric RetrospectiveMetric) string {
	const barWidth = 20

	var barStyle lipgloss.Style
	var fillChar string

	switch metric.Label {
	case "Wait Time":
		barStyle = retroWaitPeriodStyle
		fillChar = "\u2591" // Light shade
	case "Lead Time":
		barStyle = retroWaitPeriodStyle
		fillChar = "\u2591" // Light shade
	case "Active Time":
		barStyle = retroActivePeriodStyle
		fillChar = "\u2588" // Full block
	case "Efficiency":
		barStyle = retroActivePeriodStyle
		fillChar = "\u2588" // Full block
	default:
		barStyle = retroMetadataStyle
		fillChar = "\u2588" // Full block
	}

	filledLength := metric.BarLength
	if filledLength > barWidth {
		filledLength = barWidth
	}
	emptyLength := barWidth - filledLength

	filled := barStyle.Render(strings.Repeat(fillChar, filledLength))
	empty := retroMetadataStyle.Render(strings.Repeat("\u2591", emptyLength))

	return filled + empty
}

func (rv *RetrospectiveView) renderHelpBar() string {
	divider := retroDividerStyle.Render(strings.Repeat("-", min(56, rv.Width-4)))
	helpText := retroHelpBarStyle.Render("  [r] Close   [Esc] Back   [j/k] Scroll")

	return lipgloss.JoinVertical(lipgloss.Left, divider, helpText)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
