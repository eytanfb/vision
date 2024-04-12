package app

import (
	"testing"
	"time"
)

func TestTask_IsInactive(t *testing.T) {
	t.Run("when task is inactive", func(t *testing.T) {
		task := Task{
			Scheduled: false,
			Started:   false,
			Completed: false,
		}

		if !task.IsInactive() {
			t.Error("Expected task to be inactive")
		}
	})

	t.Run("when task is active", func(t *testing.T) {
		task := Task{
			Scheduled: true,
			Started:   false,
			Completed: false,
		}

		if task.IsInactive() {
			t.Error("Expected task to be active")
		}
	})

	t.Run("when task is completed", func(t *testing.T) {
		task := Task{
			Scheduled: false,
			Started:   false,
			Completed: true,
		}

		if !task.IsInactive() {
			t.Error("Expected task to be inactive")
		}
	})

	t.Run("when task is scheduled for future", func(t *testing.T) {
		task := Task{
			Scheduled:     true,
			Started:       false,
			Completed:     false,
			ScheduledDate: time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
		}

		if !task.IsInactive() {
			t.Error("Expected task to be inactive")
		}
	})
}

func TestStatusAtDate(t *testing.T) {
	// Define a series of test cases
	cases := []struct {
		name     string
		task     Task
		date     string
		expected status
	}{
		{
			name: "Task scheduled for the future",
			task: Task{
				ScheduledDate: "2022-01-10",
			},
			date:     "2022-01-05",
			expected: unscheduled,
		},
		{
			name: "Task started but not completed",
			task: Task{
				StartDate:     "2022-01-05",
				ScheduledDate: "2022-01-05",
			},
			date:     "2022-01-06",
			expected: started,
		},
		{
			name: "Task completed",
			task: Task{
				CompletedDate: "2022-01-10",
				StartDate:     "2022-01-05",
				ScheduledDate: "2022-01-05",
			},
			date:     "2022-01-15",
			expected: completed,
		},
		{
			name: "Task scheduled and started",
			task: Task{
				ScheduledDate: "2022-01-05",
				StartDate:     "2022-01-10",
			},
			date:     "2022-01-07",
			expected: scheduled,
		},
		{
			name: "Task scheduled, started, and completed",
			task: Task{
				ScheduledDate: "2022-01-05",
				StartDate:     "2022-01-10",
				CompletedDate: "2022-01-15",
			},
			date:     "2022-01-12",
			expected: started,
		},
		{
			name: "Task is scheduled and overdue",
			task: Task{
				ScheduledDate: "2022-01-05",
			},
			date:     "2022-01-20",
			expected: overdue,
		},
	}

	// Iterate through test cases, running the StatusAtDate method and comparing the result to the expected status.
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.task.StatusAtDate(tc.date)
			if actual != tc.expected {
				t.Errorf("For %s, expected status %v but got %v", tc.name, tc.expected, actual)
			}
		})
	}
}
