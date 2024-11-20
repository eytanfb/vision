package app

import (
	"strings"
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
	SelectedTask           Task
	FileExtension          string
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

func (tm *TaskManager) ExtractTasks(company string, name string, content string) []Task {
	var tasks []Task

	fileTasks := utils.ExtractTasksFromText(content)

	for _, fileTask := range fileTasks {
		task := createTaskFromFileTask(company, name, fileTask)
		tasks = append(tasks, task)
	}

	return tasks
}

func (tm *TaskManager) Summary(companyName string) map[string][]Task {
	currentDayDate := tm.DailySummaryDate

	log.Info("Summary for " + companyName + " for " + currentDayDate)

	return tm.TaskCollection.FilteredForDay(currentDayDate)
}

func (tm *TaskManager) SummaryForSlack(companyName string) string {
	slackMessage := strings.Builder{}
	previousDayString := previousDayString(tm.DailySummaryDate)
	previousDaySummary := tm.Summary(previousDayString)
	summary := tm.Summary(tm.DailySummaryDate)

	previousDayKeys := sortTaskKeys(previousDaySummary)

	slackMessage.WriteString("*Daily Update*" + "\n")
	slackMessage.WriteString("Previously" + "\n")
	for _, key := range previousDayKeys {
		category := key
		tasks := previousDaySummary[key]

		taskTitle := category[0 : len(category)-len(tm.FileExtension)]
		slackMessage.WriteString("• " + taskTitle + "\n")

		for _, task := range tasks {
			if task.Completed {
				slackMessage.WriteString("  • Finished " + task.textWithoutDates() + "\n")
			} else if task.Started && !task.IsScheduledForDay(previousDayString) {
				slackMessage.WriteString("  • Kept working on " + task.textWithoutDates() + "\n")
			} else if task.Scheduled {
				slackMessage.WriteString("  • Starting to work on " + task.textWithoutDates() + "\n")
			}
		}
	}

	slackMessage.WriteString("Today" + "\n")
	keys := sortTaskKeys(summary)

	for _, key := range keys {
		category := key
		tasks := summary[key]

		taskTitle := category[0 : len(category)-len(tm.FileExtension)]
		slackMessage.WriteString("• " + taskTitle + "\n")

		for _, task := range tasks {
			if task.Completed {
				slackMessage.WriteString("  • Finished " + task.textWithoutDates() + "\n")
			} else if task.Started && !task.IsScheduledForDay(tm.DailySummaryDate) {
				slackMessage.WriteString("  • Kept working on " + task.textWithoutDates() + "\n")
			} else if task.Scheduled {
				slackMessage.WriteString("  • Starting to work on " + task.textWithoutDates() + "\n")
			}
		}
	}

	return slackMessage.String()
}

func (tm *TaskManager) WeeklySummaryForSlack(companyName string) string {
	slackMessage := strings.Builder{}
	summary := tm.WeeklySummary(companyName, tm.WeeklySummaryStartDate, tm.WeeklySummaryEndDate)

	previousDayKeys := sortTaskKeys(summary)

	slackMessage.WriteString("*Daily Update*" + "\n")
	slackMessage.WriteString("Previously" + "\n")
	for _, key := range previousDayKeys {
		category := key
		tasks := summary[key]

		taskTitle := category[0 : len(category)-len(tm.FileExtension)]
		slackMessage.WriteString("• " + taskTitle + "\n")

		for _, task := range tasks {
			if task.Completed {
				slackMessage.WriteString("  • Finished " + task.textWithoutDates() + "\n")
			} else if task.Started {
				slackMessage.WriteString("  • Kept working on " + task.textWithoutDates() + "\n")
			} else if task.Scheduled {
				slackMessage.WriteString("  • Starting to work on " + task.textWithoutDates() + "\n")
			}
		}
	}

	slackMessage.WriteString("Today" + "\n")
	keys := sortTaskKeys(summary)

	for _, key := range keys {
		category := key
		tasks := summary[key]

		taskTitle := category[0 : len(category)-len(tm.FileExtension)]
		slackMessage.WriteString("• " + taskTitle + "\n")

		for _, task := range tasks {
			if task.Completed {
				slackMessage.WriteString("  • Finished " + task.textWithoutDates() + "\n")
			} else if task.Started {
				slackMessage.WriteString("  • Kept working on " + task.textWithoutDates() + "\n")
			} else if task.Scheduled {
				slackMessage.WriteString("  • Starting to work on " + task.textWithoutDates() + "\n")
			}
		}
	}

	return slackMessage.String()
}

func (tm *TaskManager) WeeklySummary(companyName string, startDate string, endDate string) map[string][]Task {
	log.Info("Weekly Summary for " + companyName)

	tm.WeeklySummaryStartDate = startDate
	tm.WeeklySummaryEndDate = endDate

	return tm.TaskCollection.FilteredByDates(startDate, endDate)
}

func (tm *TaskManager) ChangeDailySummaryDateToNextDay() {
	log.Info("Changing daily summary day to next day")

	currentDate, _ := time.Parse("2006-01-02", tm.DailySummaryDate)

	nextDay := currentDate.AddDate(0, 0, 1)

	tm.DailySummaryDate = nextDay.Format("2006-01-02")
}

func (tm *TaskManager) ChangeDailySummaryDateToPreviousDay() {
	log.Info("Changing daily summary day to previous day")

	currentDate, _ := time.Parse("2006-01-02", tm.DailySummaryDate)

	previousDay := currentDate.AddDate(0, 0, -1)

	tm.DailySummaryDate = previousDay.Format("2006-01-02")
}

func (tm *TaskManager) ChangeDailySummaryDateToToday() {
	log.Info("Changing daily summary day to today")

	tm.DailySummaryDate = time.Now().Format("2006-01-02")
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

func (tm *TaskManager) FridayOfWeekFromDay(day string) string {
	parsedDay, _ := time.Parse("2006-01-02", day)

	offset := (5 + 7 - int(parsedDay.Weekday())) % 7

	friday := parsedDay.Add(time.Hour * 24 * time.Duration(offset))

	return friday.Format("2006-01-02")
}

func (tm *TaskManager) SelectTask(task Task) {
	tm.SelectedTask = task
}

func (tm *TaskManager) UpdateTaskToUnscheduled(fm FileManager, task Task) {
	fm.UpdateTask(task, "unscheduled")
}

func (tm *TaskManager) UpdateTaskToScheduled(fm FileManager, task Task) {
	fm.UpdateTask(task, "scheduled")
}

func (tm *TaskManager) UpdateTaskToStarted(fm FileManager, task Task) {
	fm.UpdateTask(task, "started")
}

func (tm *TaskManager) UpdateTaskToCompleted(fm FileManager, task Task) {
	fm.UpdateTask(task, "completed")
}

func createTaskFromFileTask(company string, name string, task utils.FileTask) Task {
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
		Company:       company,
	}
}

func previousDayString(date string) string {
	currentDate, _ := time.Parse("2006-01-02", date)

	previousDay := currentDate.AddDate(0, 0, -1)

	return previousDay.Format("2006-01-02")
}

func nextDayString(date string) string {
	currentDate, _ := time.Parse("2006-01-02", date)

	nextDay := currentDate.AddDate(0, 0, 1)

	return nextDay.Format("2006-01-02")
}
