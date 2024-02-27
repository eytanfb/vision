package app

type TaskCollection struct {
	Tasks []Task
}

func CreateTaskCollection(tasks []Task) TaskCollection {
	return TaskCollection{Tasks: tasks}
}
