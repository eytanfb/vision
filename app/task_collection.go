package app

type TaskCollection struct {
	TasksByFile map[string][]Task
}

func (tc *TaskCollection) Add(filename string, tasks []Task) {
	// check if key for filename exists in TasksByFile
	// if it dooes, set its value to tasks
	// if it doesn't, create a new key-value pair
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
