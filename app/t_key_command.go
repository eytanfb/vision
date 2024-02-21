package app

import "vision/utils"

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
				IsDone: task.IsDone,
				Text:   task.Text,
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
