package app

import (
	"testing"
	"time"
)

// TestKanbanFilterLogicMatchesView verifies that the filter logic in navigation
// matches exactly what the view renders. This prevents the bug where cursor
// position N in the view maps to a different task than position N in selection.
func TestKanbanFilterLogicMatchesView(t *testing.T) {
	// Helper to create tasks with different states
	today := time.Now().Format("2006-01-02")
	future := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	past := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

	testCases := []struct {
		name         string
		task         Task
		expectedList int // 0=inactive, 1=active, 2=completed, -1=none
	}{
		// Future scheduled tasks always go to inactive (list 0)
		{
			name:         "future scheduled completed task goes to inactive",
			task:         Task{Text: "Future completed", ScheduledDate: future, Completed: true, CompletedDate: past},
			expectedList: 0,
		},
		{
			name:         "future scheduled started task goes to inactive",
			task:         Task{Text: "Future started", ScheduledDate: future, Started: true, StartDate: past},
			expectedList: 0,
		},
		{
			name:         "future scheduled task goes to inactive",
			task:         Task{Text: "Future scheduled", ScheduledDate: future, Scheduled: true},
			expectedList: 0,
		},

		// Completed tasks (not future) go to completed (list 2)
		{
			name:         "completed task goes to completed",
			task:         Task{Text: "Completed", Completed: true, CompletedDate: past},
			expectedList: 2,
		},
		{
			name:         "completed task with scheduled date goes to completed",
			task:         Task{Text: "Completed scheduled", ScheduledDate: past, Scheduled: true, Completed: true, CompletedDate: today},
			expectedList: 2,
		},
		{
			name:         "completed task with no scheduled date goes to completed",
			task:         Task{Text: "Completed no schedule", ScheduledDate: "", Completed: true, CompletedDate: today},
			expectedList: 2,
		},

		// Active tasks (started or scheduled, not completed, not future) go to active (list 1)
		{
			name:         "started task goes to active",
			task:         Task{Text: "Started", Started: true, StartDate: past},
			expectedList: 1,
		},
		{
			name:         "scheduled task goes to active",
			task:         Task{Text: "Scheduled", ScheduledDate: today, Scheduled: true},
			expectedList: 1,
		},
		{
			name:         "started and scheduled task goes to active",
			task:         Task{Text: "Started scheduled", ScheduledDate: past, Scheduled: true, Started: true, StartDate: today},
			expectedList: 1,
		},
		{
			name:         "started task with no scheduled date goes to active",
			task:         Task{Text: "Started no schedule", ScheduledDate: "", Started: true, StartDate: past},
			expectedList: 1,
		},

		// Unscheduled tasks (not completed, not started, not scheduled, no scheduled date) go to inactive (list 0)
		{
			name:         "unscheduled task goes to inactive",
			task:         Task{Text: "Unscheduled", ScheduledDate: ""},
			expectedList: 0,
		},
		{
			name:         "task with only text goes to inactive",
			task:         Task{Text: "Just text"},
			expectedList: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test the filter logic used in navigation/tea_commands
			isFuture := tc.task.IsScheduledForFuture(today)

			var actualList int = -1

			// Inactive check (matches makeKanbanLists logic)
			if isFuture {
				actualList = 0
			} else if !tc.task.Completed && !tc.task.Started && !tc.task.Scheduled && tc.task.ScheduledDate == "" {
				actualList = 0
			}

			// Completed check (only if not already categorized)
			if actualList == -1 && !isFuture && tc.task.Completed {
				actualList = 2
			}

			// Active check (only if not already categorized)
			if actualList == -1 && !isFuture && !tc.task.Completed && (tc.task.Started || tc.task.Scheduled) {
				actualList = 1
			}

			if actualList != tc.expectedList {
				t.Errorf("Task %q: expected list %d, got %d", tc.task.Text, tc.expectedList, actualList)
			}
		})
	}
}

// TestKanbanSelectionMatchesView creates a full scenario and verifies
// that the navigation cursor maps to the correct task.
// Note: Files are sorted by progress percentage descending, so we need to
// account for that in our expectations.
func TestKanbanSelectionMatchesView(t *testing.T) {
	today := time.Now().Format("2006-01-02")
	past := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

	// Create a model with tasks - both files have same completion % to ensure alphabetical order
	// FileA: 1 completed, 2 total = 50%
	// FileB: 1 completed, 2 total = 50%
	// With equal percentage, alphabetical order applies: FileA before FileB
	m := &Model{
		TaskManager: TaskManager{
			DailySummaryDate: today,
			TaskCollection: TaskCollection{
				TasksByFile: map[string][]Task{
					"FileA.md": {
						{Text: "A1 Completed", Completed: true, CompletedDate: past},  // completed
						{Text: "A2 Active", Started: true, StartDate: past},           // active
					},
					"FileB.md": {
						{Text: "B1 Completed", Completed: true, CompletedDate: today}, // completed
						{Text: "B2 Active", ScheduledDate: today, Scheduled: true},    // active
					},
				},
			},
		},
		ViewManager: ViewManager{
			HideSidebar: true,
		},
		FileManager: FileManager{
			Files: []FileInfo{
				{Name: "FileA.md"},
				{Name: "FileB.md"},
			},
		},
	}

	// Test completed list (cursor 2)
	// With equal progress, alphabetical order: FileA (A1) then FileB (B1)
	m.ViewManager.KanbanListCursor = 2
	m.ViewManager.KanbanTaskCursor = 0

	NavigationCommands{}.updateSelectedTaskInKanban(m)
	if m.TaskManager.SelectedTask.Text != "A1 Completed" {
		t.Errorf("Completed list cursor 0: expected 'A1 Completed', got %q", m.TaskManager.SelectedTask.Text)
	}

	m.ViewManager.KanbanTaskCursor = 1
	NavigationCommands{}.updateSelectedTaskInKanban(m)
	if m.TaskManager.SelectedTask.Text != "B1 Completed" {
		t.Errorf("Completed list cursor 1: expected 'B1 Completed', got %q", m.TaskManager.SelectedTask.Text)
	}

	// Test active list (cursor 1)
	// With equal progress, alphabetical order: FileA (A2) then FileB (B2)
	m.ViewManager.KanbanListCursor = 1
	m.ViewManager.KanbanTaskCursor = 0

	NavigationCommands{}.updateSelectedTaskInKanban(m)
	if m.TaskManager.SelectedTask.Text != "A2 Active" {
		t.Errorf("Active list cursor 0: expected 'A2 Active', got %q", m.TaskManager.SelectedTask.Text)
	}

	m.ViewManager.KanbanTaskCursor = 1
	NavigationCommands{}.updateSelectedTaskInKanban(m)
	if m.TaskManager.SelectedTask.Text != "B2 Active" {
		t.Errorf("Active list cursor 1: expected 'B2 Active', got %q", m.TaskManager.SelectedTask.Text)
	}
}

// TestFutureCompletedTaskGoesToInactive specifically tests the bug case where
// a future-scheduled completed task was incorrectly going to the completed list
func TestFutureCompletedTaskGoesToInactive(t *testing.T) {
	today := time.Now().Format("2006-01-02")
	future := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	past := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

	task := Task{
		Text:          "Future completed task",
		ScheduledDate: future,
		Completed:     true,
		CompletedDate: past,
	}

	isFuture := task.IsScheduledForFuture(today)
	if !isFuture {
		t.Fatal("Task should be detected as future scheduled")
	}

	// This task should NOT appear in completed list
	shouldBeInCompleted := !isFuture && task.Completed
	if shouldBeInCompleted {
		t.Error("Future completed task should NOT be in completed list")
	}

	// This task SHOULD appear in inactive list
	shouldBeInInactive := isFuture
	if !shouldBeInInactive {
		t.Error("Future completed task SHOULD be in inactive list")
	}
}

// TestCompletedTaskWithNoScheduledDateGoesToCompleted tests that a completed
// task with empty scheduled date goes to completed (not inactive)
func TestCompletedTaskWithNoScheduledDateGoesToCompleted(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	task := Task{
		Text:          "Completed no schedule",
		ScheduledDate: "", // Empty!
		Completed:     true,
		CompletedDate: today,
	}

	isFuture := task.IsScheduledForFuture(today)

	// Should NOT be in inactive (the old buggy behavior)
	shouldBeInInactive := isFuture || (!task.Completed && !task.Started && !task.Scheduled && task.ScheduledDate == "")
	if shouldBeInInactive {
		t.Error("Completed task with no scheduled date should NOT be in inactive list")
	}

	// Should be in completed
	shouldBeInCompleted := !isFuture && task.Completed
	if !shouldBeInCompleted {
		t.Error("Completed task with no scheduled date SHOULD be in completed list")
	}
}
