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

func (tc *TaskCollection) GetStartedTasks() []Task {
	log.Info("Getting started tasks")
	tasks := tc.allTasks()
	var startedTasks []Task
	for _, task := range tasks {
		if task.Started && !task.Completed {
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
