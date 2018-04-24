package network

import "net"
import "bufio"

type session struct {
	Conn       net.Conn
	WriteChann chan []byte
}

func NewSession(conn net.Conn) (session, error) {
	sess := &session{
		Conn:       conn,
		WriteChann: make(chan []byte, 128),
	}

	go sendLoop()

	return sess, nil
}

func (self *session) Run() {
	bufReader := bufio.NewReader(self.Conn)
	for {
		// handler messages
		_ = bufReader
	}
}

func (self *session) sendLoop() {
	for b := range WriteChann {
		if b == nil {
			break
		}

		_, err := self.Conn.Write(b)
		if err != nil {
			break
		}
	}
}
