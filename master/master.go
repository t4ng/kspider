package main

import (
	"os"
	"log"
	
	"kspider/rchan"
)

type Master struct {
	addr string
	rchanServer *rchan.Server
	logger *log.Logger
}

func NewMaster(addr string) *Master {
	return &Master {
		addr: addr,
		rchanServer: rchan.NewServer(addr, 1024),
		logger: log.New(os.Stderr, "", log.LstdFlags | log.Lshortfile),
	}
}

func (m *Master) Start() {
	go m.recvLoop()
	m.rchanServer.Start()
}

func (m *Master) onRecv(message []byte) {
	m.logger.Print(string(message))
}

func (m *Master) recvLoop() {
	for {
		message := <- m.rchanServer.RecvCh
		m.onRecv(message)
	}
}

func (m *Master) Send(message []byte) {
	m.rchanServer.SendCh <- message
}
