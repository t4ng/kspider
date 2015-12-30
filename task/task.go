package task

type Task struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Site string `json:"site"`
}

type SubTask struct {
	TaskId int `json:"task_id"`
}

type TaskManager struct {
	tasks map[int]*Task
}