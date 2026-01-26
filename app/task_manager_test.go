package app

import (
	"testing"
)

func TestTaskManager_ExtractTasksCreatesTastsFromInput(t *testing.T) {
	file := "test_input.txt"
	content := "- [ ] Task 1\n- [ ] Task 2\n- [ ] Task 3\n"
	company := "test_project"

	tm := TaskManager{}
	tasks := tm.ExtractTasks(company, file, content)

	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}
}

func TestTaskManager_ExtractTasksReturnsEmptyArrayForEmptyInput(t *testing.T) {
	file := "test_input.txt"
	content := ""
	company := "test_project"

	tm := TaskManager{}
	tasks := tm.ExtractTasks(company, file, content)

	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}
}
