package rchan

import (
	"log"
	"net"
	"os"
	"time"
)

type Server struct {
	addr       string
	bufSize    int
	listener   net.Listener
	logger     *log.Logger
	sessionMgr *SessionMgr
	SendCh     chan []byte
	RecvCh     chan []byte
}

func NewServer(addr string, bufSize int) *Server {
	return &Server{
		addr:       addr,
		bufSize:    bufSize,
		logger:     log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile),
		sessionMgr: NewSessionMgr(),
		SendCh:     make(chan []byte, bufSize),
		RecvCh:     make(chan []byte, bufSize),
	}
}
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.listener = listener
	go s.listenLoop()
	go s.sendLoop()
	return nil
}

func (s *Server) listenLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.Printf("accept error: %s", err)
			continue
		}
		s.logger.Printf("%s connecting", conn.RemoteAddr())
		go s.handleClient(conn)
	}
}

func (s *Server) sendLoop() {
	for {
		data := <-s.SendCh
		session := s.sessionMgr.GetRandom()
		if session == nil {
			s.logger.Print("no available session")
			time.Sleep(5 * time.Second)
			continue
		}
		_, err := session.Send(data)
		if err != nil {
			s.logger.Printf("send to %s error: %s, will retry", session.RemoteAddr(), err)
			s.SendCh <- data
		}
	}
}

func (s *Server) handleClient(conn net.Conn) {
	session := NewSession(conn)
	defer session.Close()
	s.sessionMgr.Add(session)

	buf := make([]byte, 1024*1024)
	for {
		length, err := session.Recv(buf)
		if err != nil {
			s.logger.Printf("recv from %s error: %s", session.RemoteAddr(), err)
			break
		}

		s.RecvCh <- buf[:length]
	}
}
