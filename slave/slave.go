package main

import (
	"encoding/json"
	"log"
	"os"

	"kspider/extract"
	"kspider/rchan"
	"kspider/task"
)

type Slave struct {
	masterAddr  string
	rchanClient *rchan.Client
	logger      *log.Logger
}

func NewSlave(masterAddr string) *Slave {
	return &Slave{
		masterAddr:  masterAddr,
		rchanClient: rchan.NewClient(masterAddr, 1024),
		logger:      log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile),
	}
}

func (s *Slave) Start() {
	s.rchanClient.Connect()
	go s.recvLoop()
}

func (s *Slave) Send(message []byte) {
	s.rchanClient.SendCh <- message
}

func (s *Slave) recvLoop() {
	for {
		message := <-s.rchanClient.RecvCh
		s.logger.Printf("recv: %d -> %s", len(message), message)

		var subTask task.SubTask
		err := json.Unmarshal(message, &subTask)
		if err != nil {
			s.logger.Printf("json decode error: %s", err)
			continue
		}

		subTask.Url
		s.Send([]byte(s.rchanClient.LocalAddr() + ":recvd"))
	}
}
