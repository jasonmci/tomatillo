package task

type TaskManager struct {
	tasks  []Task
	nextID int
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:  []Task{},
		nextID: 1,
	}
}

func (tm *TaskManager) AddTask(title string) {
	task := Task{
		ID:    tm.nextID,
		Title: title,
		Done:  false,
	}

	tm.tasks = append(tm.tasks, task)
	tm.nextID++
}

func (tm *TaskManager) GetTasks() []Task {
	return tm.tasks
}
