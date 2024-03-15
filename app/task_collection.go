package app

import (
	"github.com/charmbracelet/log"
)

type TaskCollection struct {
	TasksByFile map[string][]Task
}

func (tc *TaskCollection) Add(filename string, tasks []Task) {
	if _, ok := tc.TasksByFile[filename]; ok {
		tc.TasksByFile[filename] = tasks
	} else {
		tc.TasksByFile[filename] = tasks
	}
}

func (tc *TaskCollection) Size(filename string) int {
	return len(tc.TasksByFile[filename])
}

func (tc *TaskCollection) GetTasks(filename string) []Task {
	return tc.TasksByFile[filename]
}

func (tc *TaskCollection) FilteredByDates(startDate, endDate string) map[string][]Task {
	filteredTasks := make(map[string][]Task)
	for filename, tasks := range tc.TasksByFile {
		var filtered []Task
		for _, task := range tasks {
			if task.ScheduledDate >= startDate && task.ScheduledDate <= endDate {
				filtered = append(filtered, task)
			} else if task.StartDate >= startDate && task.StartDate <= endDate {
				filtered = append(filtered, task)
			} else if task.CompletedDate >= startDate && task.CompletedDate <= endDate {
				filtered = append(filtered, task)
			}
		}
		if len(filtered) > 0 {
			filteredTasks[filename] = filtered
		}
	}
	return filteredTasks
}

func (tc *TaskCollection) Progress(filename string) (int, int) {
	tasks := tc.TasksByFile[filename]
	var completed int
	for _, task := range tasks {
		if task.Completed {
			completed++
		}
	}
	return completed, len(tasks)
}

func (tc *TaskCollection) IsInactive(filename string) bool {
	tasks := tc.TasksByFile[filename]
	for _, task := range tasks {
		if !task.IsInactive() {
			return false
		}
	}
	return true
}

func (tc *TaskCollection) GetStartedTasksByDate(startDate, endDate string) []Task {
	log.Info("Getting started tasks")
	tasks := tc.allTasks()
	var startedTasks []Task
	for _, task := range tasks {
		// if task started date is between start and end date
		if task.Started && !task.Completed && task.StartDate >= startDate && task.StartDate <= endDate {
			startedTasks = append(startedTasks, task)
		}
	}
	return startedTasks
}

func (tc *TaskCollection) GetCompletedTasksByDate(startDate, endDate string) []Task {
	log.Info("Getting completed tasks")
	tasks := tc.allTasks()
	var completedTasks []Task
	for _, task := range tasks {
		if task.Completed && task.CompletedDate >= startDate && task.CompletedDate <= endDate {
			completedTasks = append(completedTasks, task)
		}
	}
	return completedTasks
}

func (tc *TaskCollection) GetScheduledTasksByDate(startDate, endDate string) []Task {
	log.Info("Getting scheduled tasks")
	tasks := tc.allTasks()
	var scheduledTasks []Task
	for _, task := range tasks {
		if task.IsScheduled() && task.ScheduledDate >= startDate && task.ScheduledDate <= endDate {
			scheduledTasks = append(scheduledTasks, task)
		}
	}
	return scheduledTasks
}

func (tc *TaskCollection) GetUnscheduledTasksByDate(startDate, endDate string) []Task {
	log.Info("Getting unscheduled tasks")
	tasks := tc.allTasks()
	var unscheduledTasks []Task
	for _, task := range tasks {
		if !task.Scheduled && !task.Completed && !task.Started {
			unscheduledTasks = append(unscheduledTasks, task)
		}
	}
	return unscheduledTasks
}

func (tc *TaskCollection) GetStartedTasksByDay(date string) []Task {
	log.Info("Getting started tasks")
	tasks := tc.allTasks()
	var startedTasks []Task
	for _, task := range tasks {
		if task.IsStarted() && task.StartDate <= date {
			startedTasks = append(startedTasks, task)
		}
	}
	return startedTasks
}

func (tc *TaskCollection) GetCompletedTasks() []Task {
	log.Info("Getting completed tasks")
	tasks := tc.allTasks()
	var completedTasks []Task
	for _, task := range tasks {
		if task.Completed {
			completedTasks = append(completedTasks, task)
		}
	}
	return completedTasks
}

func (tc *TaskCollection) GetScheduledTasks() []Task {
	log.Info("Getting scheduled tasks")
	tasks := tc.allTasks()
	var scheduledTasks []Task

	for _, task := range tasks {
		if task.Scheduled && !task.Completed && !task.Started {
			scheduledTasks = append(scheduledTasks, task)
		}
	}

	return scheduledTasks
}

func (tc *TaskCollection) allTasks() []Task {
	log.Info("Getting all tasks")
	var allTasks []Task
	for _, tasks := range tc.TasksByFile {
		allTasks = append(allTasks, tasks...)
	}
	log.Info("All tasks count: ", len(allTasks))
	return allTasks
}

func (tc *TaskCollection) GetAll() []Task {
	return tc.allTasks()
}

func (tc *TaskCollection) Flush() {
	tc.TasksByFile = make(map[string][]Task)
}
