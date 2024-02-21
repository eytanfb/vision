package app

import (
	"strings"
	"vision/utils"
)

type TKeyCommand struct{}

func (j TKeyCommand) Execute(m *Model) error {
	if m.CurrentView == "categories" {
		m.SelectedCategory = "tasks"
		m.CurrentView = "details"
		m.FilesCursor = 0
		m.Files = m.FetchFiles()
	} else if m.CurrentView == "details" && m.ItemDetailsFocus {
		fileTasks := utils.ExtractTasksFromText(m.Files[m.FilesCursor].Content)
		tasks := []Task{}
		for _, task := range fileTasks {
			tasks = append(tasks, Task{
				IsDone:        task.IsDone,
				Text:          task.Text,
				StartDate:     ExtractStartDateFromText(task.Text),
				ScheduledDate: ExtractScheduledDateFromText(task.Text),
				CompletedDate: ExtractCompletedDateFromText(task.Text),
				LineNumber:    task.LineNumber,
			})
		}
		m.Tasks = tasks
		m.TaskDetailsFocus = true
	}

	return nil
}

func (j TKeyCommand) HelpText() string {
	return "TKeyCommand help text"
}

func (j TKeyCommand) AllowedStates() []string {
	return []string{}
}

func ExtractStartDateFromText(text string) string {
	startIcon := "🛫 "
	return ExtractDateFromText(text, startIcon)
}

func ExtractScheduledDateFromText(text string) string {
	scheduledIcon := "⏳"
	return ExtractDateFromText(text, scheduledIcon)
}

func ExtractCompletedDateFromText(text string) string {
	completedIcon := "✅ "
	return ExtractDateFromText(text, completedIcon)
}

func ExtractDateFromText(text string, icon string) string {
	index := strings.Index(text, icon)
	if index == -1 {
		return ""
	}
	// read date from the next 10 characters
	date := text[index : index+14]
	return date
}
