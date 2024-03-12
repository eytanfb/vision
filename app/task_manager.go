package app

import (
	"time"
	"vision/utils"

	"github.com/charmbracelet/log"
)

type TaskManager struct {
	TaskCollection         TaskCollection
	TasksCursor            int
	WeeklySummaryStartDate string
	WeeklySummaryEndDate   string
	DailySummaryDate       string
}

type TaskCollectionSummary struct {
	StartedTasks     []Task
	CompletedTasks   []Task
	ScheduledTasks   []Task
	UnscheduledTasks []Task
}

type TaskCollectionWeeklySummary struct {
	StartDate        string
	EndDate          string
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

func (tm *TaskManager) Summary(companyName string) map[string][]Task {
	log.Info("Summary for " + companyName)

	return tm.TaskCollection.TasksByFile
}

func (tm *TaskManager) WeeklySummary(companyName string, startDate string, endDate string) map[string][]Task {
	log.Info("Weekly Summary for " + companyName)

	tm.WeeklySummaryStartDate = startDate
	tm.WeeklySummaryEndDate = endDate

	return tm.TaskCollection.FilteredByDates(startDate, endDate)
}

func (tm *TaskManager) ChangeWeeklySummaryToPreviousWeek() {
	log.Info("Changing weekly summary to previous week")

	currentStartDate, _ := time.Parse("2006-01-02", tm.WeeklySummaryStartDate)
	currentEndDate, _ := time.Parse("2006-01-02", tm.WeeklySummaryEndDate)

	previousStartDate := currentStartDate.AddDate(0, 0, -7)
	previousEndDate := currentEndDate.AddDate(0, 0, -7)

	tm.WeeklySummaryStartDate = previousStartDate.Format("2006-01-02")
	tm.WeeklySummaryEndDate = previousEndDate.Format("2006-01-02")
}

func (tm *TaskManager) ChangeWeeklySummaryToNextWeek() {
	log.Info("Changing weekly summary to next week")

	currentStartDate, _ := time.Parse("2006-01-02", tm.WeeklySummaryStartDate)
	currentEndDate, _ := time.Parse("2006-01-02", tm.WeeklySummaryEndDate)

	nextStartDate := currentStartDate.AddDate(0, 0, 7)
	nextEndDate := currentEndDate.AddDate(0, 0, 7)

	tm.WeeklySummaryStartDate = nextStartDate.Format("2006-01-02")
	tm.WeeklySummaryEndDate = nextEndDate.Format("2006-01-02")
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
