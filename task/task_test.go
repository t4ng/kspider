package task

import (
	"log"
	"testing"
)

func TestTaskManager(t *testing.T) {
	taskMgr, err := NewTaskManager("sqlite3", "./test.db")
	if err != nil {
		t.Fatalf("create task manager error: %s", err)
	}

	task := &Task{
		Site:      "baidu.com",
		StartUrls: []string{"http://www.baidu.com", "http://www.baidu.com/1"},
	}
	taskId, err := taskMgr.AddTask(task)
	if err != nil {
		t.Fatalf("add task error: %s", err)
	}
	log.Printf("create task: %d", taskId)

	taskMgr.ListTask()

	subTask, err := taskMgr.GetTodo()
	if err != nil {
		t.Fatalf("get todo error: %s", err)
	}

	log.Printf("todo: %v", subTask)
	taskMgr.ListTask()
}
