package utils

import (
	"strings"
)

type FileTask struct {
	IsDone     bool
	Text       string
	LineNumber int
}

func ExtractTasksFromText(text string) []FileTask {
	lines := strings.Split(text, "\n")
	tasks := []FileTask{}

	for index, line := range lines {
		if strings.HasPrefix(line, "- [ ]") || strings.HasPrefix(line, "- [x]") {
			text := strings.TrimPrefix(line, "- [ ]")
			text = strings.TrimPrefix(text, "- [x]")

			task := FileTask{
				IsDone:     strings.HasPrefix(line, "- [x]"),
				Text:       text,
				LineNumber: index + 1,
			}
			tasks = append(tasks, task)
		}
	}

	return tasks
}
