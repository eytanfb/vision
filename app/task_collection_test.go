package app

import (
	"reflect"
	"testing"
)

func TestFilteredForDay(t *testing.T) {
	// Define test cases
	cases := []struct {
		name     string
		tc       TaskCollection
		date     string
		expected map[string][]Task
	}{
		{
			name: "Filter tasks for specific day",
			tc: TaskCollection{
				TasksByFile: map[string][]Task{
					"file1.md": {
						{Text: "Task 1", ScheduledDate: "2022-01-05", StartDate: "2022-01-06", CompletedDate: "2022-01-07"},
						{Text: "Task 2", ScheduledDate: "2022-01-06", StartDate: "2022-01-09"},
						{Text: "Task 3", ScheduledDate: "2022-01-10"},
					},
					"file2.md": {
						{Text: "Task 4", ScheduledDate: "2022-01-05", StartDate: "2022-01-06"},
						{Text: "Task 5", CompletedDate: "2022-01-04"},
						{Text: "Task 6", StartDate: "2022-01-04"},
					},
				},
			},
			date: "2022-01-06",
			expected: map[string][]Task{
				"file1.md": {
					{Text: "Task 1", ScheduledDate: "2022-01-05", StartDate: "2022-01-06", CompletedDate: "2022-01-07"},
					{Text: "Task 2", ScheduledDate: "2022-01-06", StartDate: "2022-01-09"},
				},
				"file2.md": {
					{Text: "Task 4", ScheduledDate: "2022-01-05", StartDate: "2022-01-06"},
					{Text: "Task 6", StartDate: "2022-01-04"},
				},
			},
		},
		// Add more cases as necessary...
		{
			name: "Filter by historical event dates",
			tc: TaskCollection{
				TasksByFile: map[string][]Task{
					"January6Insurrection.md": {
						{Text: "Reflect on political dynamics", ScheduledDate: "2021-01-06", StartDate: "2021-01-07", CompletedDate: "2021-01-08"},
						{Text: "Prepare report on event", ScheduledDate: "2021-01-09", StartDate: "2021-01-10"},
					},
					"COVID19PandemicStart.md": {
						{Text: "Gather initial data", ScheduledDate: "2020-03-11", StartDate: "2020-03-12", CompletedDate: "2020-03-20"},
						{Text: "Draft initial response plans", ScheduledDate: "2020-03-15", StartDate: "2020-03-16"},
					},
				},
			},
			date: "2021-01-07",
			expected: map[string][]Task{
				"January6Insurrection.md": {
					{Text: "Reflect on political dynamics", ScheduledDate: "2021-01-06", StartDate: "2021-01-07", CompletedDate: "2021-01-08"},
				},
				"COVID19PandemicStart.md": {
					{Text: "Draft initial response plans", ScheduledDate: "2020-03-15", StartDate: "2020-03-16"},
				},
			},
		},
	}

	// Iterate through test cases
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.tc.FilteredForDay(c.date)
			if !reflect.DeepEqual(actual, c.expected) {
				t.Errorf("Test %s failed: expected %v, got %v", c.name, c.expected, actual)
			}
		})
	}
}
