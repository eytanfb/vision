package app

import (
	"strings"
	"testing"
	"time"
)

func TestNewRetrospectiveView(t *testing.T) {
	task := Task{
		Text:          "Test task",
		FileName:      "test-file.md",
		Company:       "TestCo",
		ScheduledDate: "2024-01-22",
		StartDate:     "2024-01-24",
		CompletedDate: "2024-01-27",
	}

	createdDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	retro := NewTaskRetrospective(task, createdDate)

	rv := NewRetrospectiveView(retro, 80, 40)

	if rv.Retro != retro {
		t.Error("expected Retro to be set")
	}
	if rv.Width != 80 {
		t.Errorf("expected Width 80, got %d", rv.Width)
	}
	if rv.Height != 40 {
		t.Errorf("expected Height 40, got %d", rv.Height)
	}
}

func TestRetrospectiveView_View_NilRetro(t *testing.T) {
	rv := NewRetrospectiveView(nil, 80, 40)

	view := rv.View()

	if view != "No retrospective data available" {
		t.Errorf("expected nil message, got %s", view)
	}
}

func TestRetrospectiveView_View_ContainsAllSections(t *testing.T) {
	task := Task{
		Text:          "Fix authentication bug",
		FileName:      "auth-system.md",
		Company:       "Clerky",
		Priority:      "high",
		ScheduledDate: "2024-01-22",
		StartDate:     "2024-01-24",
		CompletedDate: "2024-01-27",
	}

	createdDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	retro := NewTaskRetrospective(task, createdDate)

	rv := NewRetrospectiveView(retro, 80, 40)
	view := rv.View()

	if !strings.Contains(view, "Task Retrospective") {
		t.Error("expected view to contain 'Task Retrospective' title")
	}
	if !strings.Contains(view, "Timeline") {
		t.Error("expected view to contain 'Timeline' section")
	}
	if !strings.Contains(view, "Metrics") {
		t.Error("expected view to contain 'Metrics' section")
	}
	if !strings.Contains(view, "[r] Close") {
		t.Error("expected view to contain help bar")
	}
}

func TestRetrospectiveView_renderHeader_ContainsTaskInfo(t *testing.T) {
	task := Task{
		Text:          "Fix authentication bug",
		FileName:      "auth-system.md",
		Company:       "Clerky",
		Priority:      "high",
		ScheduledDate: "2024-01-22",
		StartDate:     "2024-01-24",
		CompletedDate: "2024-01-27",
	}

	createdDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	retro := NewTaskRetrospective(task, createdDate)

	rv := NewRetrospectiveView(retro, 80, 40)
	header := rv.renderHeader()

	if !strings.Contains(header, "Fix authentication bug") {
		t.Error("expected header to contain task name")
	}
	if !strings.Contains(header, "auth-system.md") {
		t.Error("expected header to contain file name")
	}
	if !strings.Contains(header, "Clerky") {
		t.Error("expected header to contain company name")
	}
}

func TestRetrospectiveView_renderTimeline_ContainsAllEvents(t *testing.T) {
	task := Task{
		Text:          "Test task",
		ScheduledDate: "2024-01-22",
		StartDate:     "2024-01-24",
		CompletedDate: "2024-01-27",
	}

	createdDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	retro := NewTaskRetrospective(task, createdDate)

	rv := NewRetrospectiveView(retro, 80, 40)
	timeline := rv.renderTimeline()

	if !strings.Contains(timeline, "Created") {
		t.Error("expected timeline to contain 'Created' event")
	}
	if !strings.Contains(timeline, "Scheduled") {
		t.Error("expected timeline to contain 'Scheduled' event")
	}
	if !strings.Contains(timeline, "Started") {
		t.Error("expected timeline to contain 'Started' event")
	}
	if !strings.Contains(timeline, "Completed") {
		t.Error("expected timeline to contain 'Completed' event")
	}
	if !strings.Contains(timeline, "Day 0") {
		t.Error("expected timeline to contain 'Day 0' for created")
	}
	if !strings.Contains(timeline, "Day 12") {
		t.Error("expected timeline to contain 'Day 12' for completed")
	}
}

func TestRetrospectiveView_renderMetrics_ContainsAllMetrics(t *testing.T) {
	task := Task{
		Text:          "Test task",
		ScheduledDate: "2024-01-22",
		StartDate:     "2024-01-24",
		CompletedDate: "2024-01-27",
	}

	createdDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	retro := NewTaskRetrospective(task, createdDate)

	rv := NewRetrospectiveView(retro, 80, 40)
	metrics := rv.renderMetrics()

	if !strings.Contains(metrics, "Total Time") {
		t.Error("expected metrics to contain 'Total Time'")
	}
	if !strings.Contains(metrics, "Wait Time") {
		t.Error("expected metrics to contain 'Wait Time'")
	}
	if !strings.Contains(metrics, "Lead Time") {
		t.Error("expected metrics to contain 'Lead Time'")
	}
	if !strings.Contains(metrics, "Active Time") {
		t.Error("expected metrics to contain 'Active Time'")
	}
	if !strings.Contains(metrics, "Efficiency") {
		t.Error("expected metrics to contain 'Efficiency'")
	}
}

func TestRetrospectiveView_renderHelpBar(t *testing.T) {
	rv := NewRetrospectiveView(nil, 80, 40)
	rv.Retro = &TaskRetrospective{}

	helpBar := rv.renderHelpBar()

	if !strings.Contains(helpBar, "[r] Close") {
		t.Error("expected help bar to contain '[r] Close'")
	}
	if !strings.Contains(helpBar, "[Esc] Back") {
		t.Error("expected help bar to contain '[Esc] Back'")
	}
	if !strings.Contains(helpBar, "[j/k] Scroll") {
		t.Error("expected help bar to contain '[j/k] Scroll'")
	}
}

func TestRetrospectiveView_MinimalTask(t *testing.T) {
	task := Task{
		Text:          "Quick fix",
		CompletedDate: "2024-01-15",
	}

	createdDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	retro := NewTaskRetrospective(task, createdDate)

	rv := NewRetrospectiveView(retro, 80, 40)
	view := rv.View()

	if !strings.Contains(view, "Quick fix") {
		t.Error("expected view to contain task name")
	}
	if !strings.Contains(view, "Day 0") {
		t.Error("expected same-day completion to show Day 0")
	}
}
