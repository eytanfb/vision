package app

import (
	"testing"
	"time"
)

func TestNewTaskRetrospective_FullLifecycle(t *testing.T) {
	task := Task{
		Text:          "Complete feature implementation",
		ScheduledDate: "2024-01-05",
		StartDate:     "2024-01-08",
		CompletedDate: "2024-01-12",
		Completed:     true,
	}
	createdDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	retro := NewTaskRetrospective(task, createdDate)

	if retro == nil {
		t.Fatal("Expected non-nil retrospective")
	}

	// Total: Jan 1 to Jan 12 = 11 days
	if retro.TotalDays != 11 {
		t.Errorf("Expected TotalDays=11, got %d", retro.TotalDays)
	}

	// Wait: Jan 1 to Jan 5 = 4 days
	if retro.WaitDays != 4 {
		t.Errorf("Expected WaitDays=4, got %d", retro.WaitDays)
	}

	// Lead: Jan 5 to Jan 8 = 3 days
	if retro.LeadDays != 3 {
		t.Errorf("Expected LeadDays=3, got %d", retro.LeadDays)
	}

	// Active: Jan 8 to Jan 12 = 4 days
	if retro.ActiveDays != 4 {
		t.Errorf("Expected ActiveDays=4, got %d", retro.ActiveDays)
	}

	// Efficiency: 4/11 = 0.3636...
	expectedEfficiency := float64(4) / float64(11)
	if !floatEquals(retro.EfficiencyRatio, expectedEfficiency, 0.001) {
		t.Errorf("Expected EfficiencyRatio=%f, got %f", expectedEfficiency, retro.EfficiencyRatio)
	}
}

func TestNewTaskRetrospective_MissingScheduledDate(t *testing.T) {
	task := Task{
		Text:          "Quick fix",
		StartDate:     "2024-01-03",
		CompletedDate: "2024-01-05",
		Completed:     true,
	}
	createdDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	retro := NewTaskRetrospective(task, createdDate)

	if retro == nil {
		t.Fatal("Expected non-nil retrospective")
	}

	if retro.ScheduledDate != nil {
		t.Error("Expected ScheduledDate to be nil")
	}

	// Total: Jan 1 to Jan 5 = 4 days
	if retro.TotalDays != 4 {
		t.Errorf("Expected TotalDays=4, got %d", retro.TotalDays)
	}

	// Wait: Jan 1 to Jan 3 (start date used when no scheduled) = 2 days
	if retro.WaitDays != 2 {
		t.Errorf("Expected WaitDays=2, got %d", retro.WaitDays)
	}

	// Lead: 0 (no scheduled date)
	if retro.LeadDays != 0 {
		t.Errorf("Expected LeadDays=0, got %d", retro.LeadDays)
	}

	// Active: Jan 3 to Jan 5 = 2 days
	if retro.ActiveDays != 2 {
		t.Errorf("Expected ActiveDays=2, got %d", retro.ActiveDays)
	}
}

func TestNewTaskRetrospective_MissingStartDate(t *testing.T) {
	task := Task{
		Text:          "Scheduled but completed directly",
		ScheduledDate: "2024-01-05",
		CompletedDate: "2024-01-10",
		Completed:     true,
	}
	createdDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	retro := NewTaskRetrospective(task, createdDate)

	if retro == nil {
		t.Fatal("Expected non-nil retrospective")
	}

	if retro.StartDate != nil {
		t.Error("Expected StartDate to be nil")
	}

	// Total: Jan 1 to Jan 10 = 9 days
	if retro.TotalDays != 9 {
		t.Errorf("Expected TotalDays=9, got %d", retro.TotalDays)
	}

	// Wait: Jan 1 to Jan 5 = 4 days
	if retro.WaitDays != 4 {
		t.Errorf("Expected WaitDays=4, got %d", retro.WaitDays)
	}

	// Lead: 0 (no start date)
	if retro.LeadDays != 0 {
		t.Errorf("Expected LeadDays=0, got %d", retro.LeadDays)
	}

	// Active: 0 (no start date)
	if retro.ActiveDays != 0 {
		t.Errorf("Expected ActiveDays=0, got %d", retro.ActiveDays)
	}

	// Efficiency: 0 (no active time)
	if retro.EfficiencyRatio != 0.0 {
		t.Errorf("Expected EfficiencyRatio=0.0, got %f", retro.EfficiencyRatio)
	}
}

func TestNewTaskRetrospective_SameDayCompletion(t *testing.T) {
	task := Task{
		Text:          "Quick task",
		StartDate:     "2024-01-01",
		CompletedDate: "2024-01-01",
		Completed:     true,
	}
	createdDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	retro := NewTaskRetrospective(task, createdDate)

	if retro == nil {
		t.Fatal("Expected non-nil retrospective")
	}

	// Total: 0 days (same day)
	if retro.TotalDays != 0 {
		t.Errorf("Expected TotalDays=0, got %d", retro.TotalDays)
	}

	// Wait: 0 days
	if retro.WaitDays != 0 {
		t.Errorf("Expected WaitDays=0, got %d", retro.WaitDays)
	}

	// Active: 0 days (same day)
	if retro.ActiveDays != 0 {
		t.Errorf("Expected ActiveDays=0, got %d", retro.ActiveDays)
	}

	// Efficiency: 1.0 (100% for same-day completion)
	if retro.EfficiencyRatio != 1.0 {
		t.Errorf("Expected EfficiencyRatio=1.0, got %f", retro.EfficiencyRatio)
	}
}

func TestNewTaskRetrospective_NoCompletedDate(t *testing.T) {
	task := Task{
		Text:      "Incomplete task",
		StartDate: "2024-01-01",
		Completed: false,
	}
	createdDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	retro := NewTaskRetrospective(task, createdDate)

	if retro != nil {
		t.Error("Expected nil retrospective for task without completed date")
	}
}

func TestNewTaskRetrospective_InvalidCompletedDate(t *testing.T) {
	task := Task{
		Text:          "Bad date task",
		CompletedDate: "invalid-date",
		Completed:     true,
	}
	createdDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	retro := NewTaskRetrospective(task, createdDate)

	if retro != nil {
		t.Error("Expected nil retrospective for task with invalid completed date")
	}
}

func TestTaskRetrospective_GetTimeline_FullLifecycle(t *testing.T) {
	task := Task{
		Text:          "Timeline test",
		ScheduledDate: "2024-01-05",
		StartDate:     "2024-01-08",
		CompletedDate: "2024-01-12",
		Completed:     true,
	}
	createdDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	retro := NewTaskRetrospective(task, createdDate)
	timeline := retro.GetTimeline()

	if len(timeline) != 4 {
		t.Fatalf("Expected 4 timeline events, got %d", len(timeline))
	}

	// Verify order: created, scheduled, started, completed
	expectedTypes := []string{"created", "scheduled", "started", "completed"}
	for i, event := range timeline {
		if event.EventType != expectedTypes[i] {
			t.Errorf("Event %d: expected type %s, got %s", i, expectedTypes[i], event.EventType)
		}
	}

	// Verify day numbers
	if timeline[0].DayNumber != 0 {
		t.Errorf("Created event should be day 0, got %d", timeline[0].DayNumber)
	}
	if timeline[1].DayNumber != 4 {
		t.Errorf("Scheduled event should be day 4, got %d", timeline[1].DayNumber)
	}
	if timeline[2].DayNumber != 7 {
		t.Errorf("Started event should be day 7, got %d", timeline[2].DayNumber)
	}
	if timeline[3].DayNumber != 11 {
		t.Errorf("Completed event should be day 11, got %d", timeline[3].DayNumber)
	}
}

func TestTaskRetrospective_GetTimeline_MinimalLifecycle(t *testing.T) {
	task := Task{
		Text:          "Minimal task",
		CompletedDate: "2024-01-05",
		Completed:     true,
	}
	createdDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	retro := NewTaskRetrospective(task, createdDate)
	timeline := retro.GetTimeline()

	if len(timeline) != 2 {
		t.Fatalf("Expected 2 timeline events, got %d", len(timeline))
	}

	if timeline[0].EventType != "created" {
		t.Errorf("First event should be 'created', got %s", timeline[0].EventType)
	}
	if timeline[1].EventType != "completed" {
		t.Errorf("Second event should be 'completed', got %s", timeline[1].EventType)
	}
}

func TestTaskRetrospective_GetTimeline_OutOfOrderDates(t *testing.T) {
	// Simulate a task where scheduled date comes after start date (edge case)
	task := Task{
		Text:          "Out of order dates",
		ScheduledDate: "2024-01-10",
		StartDate:     "2024-01-05",
		CompletedDate: "2024-01-15",
		Completed:     true,
	}
	createdDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	retro := NewTaskRetrospective(task, createdDate)
	timeline := retro.GetTimeline()

	// Timeline should be sorted by date regardless of logical order
	for i := 1; i < len(timeline); i++ {
		if timeline[i].Date.Before(timeline[i-1].Date) {
			t.Errorf("Timeline not sorted: event %d (%s) before event %d (%s)",
				i, timeline[i].Date, i-1, timeline[i-1].Date)
		}
	}
}

func TestTaskRetrospective_GetMetrics(t *testing.T) {
	task := Task{
		Text:          "Metrics test",
		ScheduledDate: "2024-01-05",
		StartDate:     "2024-01-10",
		CompletedDate: "2024-01-20",
		Completed:     true,
	}
	createdDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	retro := NewTaskRetrospective(task, createdDate)
	metrics := retro.GetMetrics()

	// Should have: Total, Wait, Lead, Active, Efficiency = 5 metrics
	if len(metrics) != 5 {
		t.Fatalf("Expected 5 metrics, got %d", len(metrics))
	}

	// Verify Total metric
	if metrics[0].Label != "Total Time" {
		t.Errorf("First metric should be 'Total Time', got %s", metrics[0].Label)
	}
	if metrics[0].Value != 19 { // Jan 1 to Jan 20
		t.Errorf("Total Time should be 19, got %d", metrics[0].Value)
	}
	if metrics[0].BarLength != 20 {
		t.Errorf("Total Time bar length should be 20, got %d", metrics[0].BarLength)
	}

	// Verify Efficiency metric is last
	effMetric := metrics[len(metrics)-1]
	if effMetric.Label != "Efficiency" {
		t.Errorf("Last metric should be 'Efficiency', got %s", effMetric.Label)
	}
}

func TestTaskRetrospective_GetMetrics_SameDayCompletion(t *testing.T) {
	task := Task{
		Text:          "Quick task",
		StartDate:     "2024-01-01",
		CompletedDate: "2024-01-01",
		Completed:     true,
	}
	createdDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	retro := NewTaskRetrospective(task, createdDate)
	metrics := retro.GetMetrics()

	// Find efficiency metric
	var effMetric *RetrospectiveMetric
	for i := range metrics {
		if metrics[i].Label == "Efficiency" {
			effMetric = &metrics[i]
			break
		}
	}

	if effMetric == nil {
		t.Fatal("Efficiency metric not found")
	}

	if effMetric.Value != 100 {
		t.Errorf("Same-day completion efficiency should be 100%%, got %d%%", effMetric.Value)
	}
	if effMetric.BarLength != 20 {
		t.Errorf("Same-day completion efficiency bar should be 20, got %d", effMetric.BarLength)
	}
}

func TestCalculateEfficiencyRatio(t *testing.T) {
	tests := []struct {
		name     string
		active   int
		total    int
		expected float64
	}{
		{"same day completion", 0, 0, 1.0},
		{"no active time", 0, 10, 0.0},
		{"full efficiency", 10, 10, 1.0},
		{"half efficiency", 5, 10, 0.5},
		{"quarter efficiency", 25, 100, 0.25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateEfficiencyRatio(tt.active, tt.total)
			if !floatEquals(result, tt.expected, 0.001) {
				t.Errorf("calculateEfficiencyRatio(%d, %d) = %f, want %f",
					tt.active, tt.total, result, tt.expected)
			}
		})
	}
}

func TestDaysBetween(t *testing.T) {
	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected int
	}{
		{
			"same day",
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 1, 1, 23, 59, 59, 0, time.UTC),
			0,
		},
		{
			"consecutive days",
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			1,
		},
		{
			"one week",
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			7,
		},
		{
			"end before start returns 0",
			time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			0,
		},
		{
			"across month boundary",
			time.Date(2024, 1, 30, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC),
			3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := daysBetween(tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("daysBetween() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestPercentageToBarLength(t *testing.T) {
	tests := []struct {
		pct      float64
		expected int
	}{
		{0, 0},
		{5, 1},
		{10, 2},
		{25, 5},
		{50, 10},
		{75, 15},
		{100, 20},
		{-10, 0},   // Negative clamps to 0
		{150, 20},  // Over 100 clamps to 20
		{2.5, 1},   // Rounds
		{7.5, 2},   // Rounds
		{12.5, 3},  // Rounds
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := percentageToBarLength(tt.pct)
			if result != tt.expected {
				t.Errorf("percentageToBarLength(%f) = %d, want %d", tt.pct, result, tt.expected)
			}
		})
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		count    int
		singular string
		plural   string
		expected string
	}{
		{0, "day", "days", "days"},
		{1, "day", "days", "day"},
		{2, "day", "days", "days"},
		{100, "day", "days", "days"},
	}

	for _, tt := range tests {
		result := pluralize(tt.count, tt.singular, tt.plural)
		if result != tt.expected {
			t.Errorf("pluralize(%d, %s, %s) = %s, want %s",
				tt.count, tt.singular, tt.plural, result, tt.expected)
		}
	}
}

// floatEquals compares two floats within a tolerance
func floatEquals(a, b, tolerance float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < tolerance
}
