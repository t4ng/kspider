package main

import (
	"os"
	"log"
	
	"kspider/rchan"
)

type Slave struct {
	masterAddr string
	rchanClient *rchan.Client
	logger *log.Logger
}

func NewSlave(masterAddr string) *Slave {
	return &Slave {
		masterAddr: masterAddr,
		rchanClient: rchan.NewClient(masterAddr, 1024),
		logger: log.New(os.Stderr, "", log.LstdFlags | log.Lshortfile),
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
		message := <- s.rchanClient.RecvCh
		s.logger.Printf("recv: %d -> %s", len(message), message)
		s.Send([]byte(s.rchanClient.LocalAddr() + ":recvd"))
	}
}