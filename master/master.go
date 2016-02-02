package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"kspider/rchan"
	"kspider/task"
)

type Master struct {
	addr        string
	logger      *log.Logger
	rchanServer *rchan.Server
	taskMgr     *task.TaskManager
}

func NewMaster(addr string) (*Master, error) {
	taskMgr, err := task.NewTaskManager("sqlite3", "./ks.db")
	if err != nil {
		return nil, err
	}

	return &Master{
		addr:        addr,
		logger:      log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile),
		rchanServer: rchan.NewServer(addr, 1024),
		taskMgr:     taskMgr,
	}, nil
}

func (m *Master) Start() {
	go m.recvLoop()
	go m.sendLoop()
	m.rchanServer.Start()
}

func (m *Master) onRecv(message []byte) {
	m.logger.Print(string(message))
}

func (m *Master) recvLoop() {
	for {
		message := <-m.rchanServer.RecvCh
		m.onRecv(message)
	}
}

func (m *Master) Send(message []byte) {
	m.rchanServer.SendCh <- message
}

func (m *Master) sendLoop() {
	for {
		subTask, err := m.taskMgr.GetTodo()
		if err != nil {
			if err != task.ErrNotFound {
				m.logger.Printf("get todo error: %s", err)
			}
			time.Sleep(2 * time.Second)
			continue
		}

		message, err := json.Marshal(subTask)
		if err != nil {
			m.logger.Printf("json encode error: %s", err)
			continue
		}

		m.Send(message)
	}
}
