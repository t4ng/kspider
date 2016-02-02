package rchan

import (
	"math/rand"
	"net"
	"sync"
)

type Session struct {
	mu     sync.Mutex
	conn   net.Conn
	closed bool
}

type SessionMgr struct {
	mu         sync.RWMutex
	sessionMap map[string]*Session
	addresses  []string
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		conn: conn,
	}
}

func (s *Session) Send(buf []byte) (length int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	length, err = s.conn.Write(buf)
	return
}

func (s *Session) Recv(buf []byte) (length int, err error) {
	length, err = s.conn.Read(buf)
	return
}

func (s *Session) Close() {
	s.conn.Close()
	s.closed = true
}

func (s *Session) LocalAddr() string {
	return s.conn.LocalAddr().String()
}

func (s *Session) RemoteAddr() string {
	return s.conn.RemoteAddr().String()
}

func NewSessionMgr() *SessionMgr {
	return &SessionMgr{
		sessionMap: make(map[string]*Session),
	}
}

func (m *SessionMgr) Get(key string) *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	session, ok := m.sessionMap[key]
	if ok && session != nil && session.closed {
		delete(m.sessionMap, key)
		return nil
	}
	return session
}

func (m *SessionMgr) Add(session *Session) {
	key := session.RemoteAddr()
	oldSession := m.Get(key)
	if oldSession != nil {
		oldSession.Close()
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessionMap[key] = session
	m.addresses = append(m.addresses, key)
}

func (m *SessionMgr) GetRandom() *Session {
	sessionCount := len(m.addresses)
	if sessionCount == 0 {
		return nil
	} else if sessionCount == 1 {
		return m.Get(m.addresses[0])
	} else {
		return m.Get(m.addresses[rand.Intn(sessionCount)])
	}
}
