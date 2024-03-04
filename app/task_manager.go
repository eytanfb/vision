package app

import (
	"fmt"
	"vision/utils"
)

type TaskManager struct {
	TaskCollection TaskCollection
	TasksCursor    int
}

func (tm *TaskManager) ExtractTasks(content string) []Task {
	var tasks []Task

	fileTasks := utils.ExtractTasksFromText(content)

	for _, fileTask := range fileTasks {
		task := createTaskFromFileTask(fileTask)
		tasks = append(tasks, task)
	}

	return tasks
}

func (tm *TaskManager) Summary(companyName string) string {
	startedTasks := tm.TaskCollection.GetStartedTasks()
	completedTasks := tm.TaskCollection.GetCompletedTasks()
	scheduledTasks := tm.TaskCollection.GetScheduledTasks()

	return fmt.Sprintf("You have %d tasks started, %d tasks completed, and %d tasks scheduled.", len(startedTasks), len(completedTasks), len(scheduledTasks))
}

func createTaskFromFileTask(task utils.FileTask) Task {
	completedDate := extractCompletedDateFromText(task.Text)
	startDate := extractStartDateFromText(task.Text)
	scheduledDate := extractScheduledDateFromText(task.Text)
	completed := completedDate != ""
	started := startDate != ""
	scheduled := scheduledDate != ""

	return Task{
		IsDone:        task.IsDone,
		Text:          task.Text,
		StartDate:     startDate,
		ScheduledDate: scheduledDate,
		CompletedDate: completedDate,
		LineNumber:    task.LineNumber,
		Completed:     completed,
		Started:       started,
		Scheduled:     scheduled,
	}
}
