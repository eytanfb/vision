package app

import (
	"vision/utils"

	"github.com/charmbracelet/log"
)

type TaskManager struct {
	TaskCollection TaskCollection
	TasksCursor    int
}

type TaskCollectionSummary struct {
	StartedTasks     []Task
	CompletedTasks   []Task
	ScheduledTasks   []Task
	UnscheduledTasks []Task
}

func (tm *TaskManager) ExtractTasks(name string, content string) []Task {
	var tasks []Task

	fileTasks := utils.ExtractTasksFromText(content)

	for _, fileTask := range fileTasks {
		task := createTaskFromFileTask(name, fileTask)
		tasks = append(tasks, task)
	}

	return tasks
}

func (tm *TaskManager) Summary(companyName string) TaskCollectionSummary {
	log.Info("Summary for " + companyName)
	startedTasks := tm.TaskCollection.GetStartedTasks()
	completedTasks := tm.TaskCollection.GetCompletedTasks()
	scheduledTasks := tm.TaskCollection.GetScheduledTasks()

	return TaskCollectionSummary{
		StartedTasks:   startedTasks,
		CompletedTasks: completedTasks,
		ScheduledTasks: scheduledTasks,
	}
}

func createTaskFromFileTask(name string, task utils.FileTask) Task {
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
		FileName:      name,
	}
}
