package rchan

import (
	"log"
	"net"
	"os"
	"time"
)

type Client struct {
	serverAddr string
	session    *Session
	logger     *log.Logger
	SendCh     chan []byte
	RecvCh     chan []byte
}

func NewClient(serverAddr string, bufSize int) *Client {
	return &Client{
		serverAddr: serverAddr,
		logger:     log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile),
		SendCh:     make(chan []byte, bufSize),
		RecvCh:     make(chan []byte, bufSize),
	}
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.serverAddr)
	if err != nil {
		c.logger.Printf("connect to %s error: %s", c.serverAddr, err)
		return err
	}

	c.session = NewSession(conn)
	go c.recvLoop()
	go c.sendLoop()
	return nil
}

func (c *Client) reconnect() {
	for {
		conn, err := net.Dial("tcp", c.serverAddr)
		if err == nil {
			c.session = NewSession(conn)
			break
		}

		c.logger.Printf("reconnect to %s error: %s", c.serverAddr, err)
		time.Sleep(5 * time.Second)
	}
}

func (c *Client) LocalAddr() string {
	if c.session == nil {
		return ""
	}
	return c.session.LocalAddr()
}

func (c *Client) recvLoop() {
	defer c.session.Close()
	buf := make([]byte, 1024*1024)
	for {
		length, err := c.session.Recv(buf)
		if err != nil {
			c.logger.Printf("recv from %s error: %s", c.serverAddr, err)
			c.reconnect()
			continue
		}

		c.RecvCh <- buf[:length]
	}
}

func (c *Client) sendLoop() {
	for {
		data := <-c.SendCh
		_, err := c.session.Send(data)
		if err != nil {
			c.logger.Printf("send to %s error: %s", c.serverAddr, err)
			c.SendCh <- data
			time.Sleep(1 * time.Second)
			continue
		}
	}
}
