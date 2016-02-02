package task

import (
	"errors"
	"log"
	urlparse "net/url"
	"time"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

const (
	StatusTodo = iota
	StatusDoing
	StatusDone
	StatusFail
)

var (
	ErrNotFound = errors.New("your record not found")
)

type Task struct {
	Id             int `xorm:"pk autoincr"`
	Site           string
	StartUrls      []string
	AllowedDomains []string
	Headers        map[string]string
	MaxDepth       int
	Status         int
	CreatedAt      time.Time `xorm:"created`
}

type SubTask struct {
	Id        int `xorm:"pk autoincr"`
	TaskId    int `xorm:"index"`
	Site      string
	Url       string `xorm:"index"`
	Headers   map[string]string
	Depth     int
	Status    int
	Retries   int
	Result    string
	CreatedAt time.Time `xorm:"created"`
}

type TaskManager struct {
	db *xorm.Engine
}

func NewTaskManager(dbDriver string, dbSource string) (*TaskManager, error) {
	db, err := xorm.NewEngine(dbDriver, dbSource)
	if err != nil {
		return nil, err
	}

	db.Sync2(new(Task))
	db.Sync2(new(SubTask))
	return &TaskManager{db: db}, nil
}

func (tm *TaskManager) AddSubTask(subTask *SubTask) (int, error) {
	_, err := tm.db.Insert(subTask)
	if err != nil {
		return 0, err
	}
	return subTask.Id, err
}

func (tm *TaskManager) AddTask(task *Task) (int, error) {
	if task.MaxDepth <= 0 {
		task.MaxDepth = 3
	}

	_, err := tm.db.Insert(task)
	if err != nil {
		return 0, err
	}

	for _, url := range task.StartUrls {
		subTask := &SubTask{
			TaskId:  task.Id,
			Site:    task.Site,
			Url:     url,
			Headers: task.Headers,
		}
		_, err := tm.AddSubTask(subTask)
		if err != nil {
			return 0, err
		}
	}

	return task.Id, nil
}

func (tm *TaskManager) CreateTask(url string) (*Task, error) {
	uri, err := urlparse.Parse(url)
	if err != nil {
		return nil, err
	}

	task := &Task{
		Site:      uri.Host,
		StartUrls: []string{url},
		Headers:   make(map[string]string),
	}
	_, err = tm.AddTask(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (tm *TaskManager) GetTask(taskId int) (*Task, error) {
	var task Task
	has, err := tm.db.Id(taskId).Get(&task)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrNotFound
	}

	return &task, nil
}

func (tm *TaskManager) GetTodo() (*SubTask, error) {
	var subTask SubTask
	has, err := tm.db.Where("status=?", StatusTodo).OrderBy("created_at").Get(&subTask)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrNotFound
	}

	subTask.Status = StatusDoing
	affected, err := tm.db.Id(subTask.Id).Update(&subTask)
	if err != nil {
		return nil, err
	} else if affected == 0 {
		return nil, ErrNotFound
	}

	return &subTask, err
}

func (tm *TaskManager) Done(subTask *SubTask) error {
	subTask.Status = StatusDone
	affected, err := tm.db.Id(subTask.Id).Update(&subTask)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}

	return nil
}

func (tm *TaskManager) ListTask() error {
	var tasks []Task
	err := tm.db.Limit(10).Find(&tasks)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		log.Printf("task: %v", task)
	}

	var subTasks []SubTask
	err = tm.db.OrderBy("created_at").Limit(100).Find(&subTasks)
	if err != nil {
		return err
	}

	for _, subTask := range subTasks {
		log.Printf("subTask: %v", subTask)
	}
	return nil
}
