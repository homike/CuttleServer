package network

import "net"

type session struct {
	Conn net.Conn

	WriteChann chan []byte
}

func NewSession(conn net.Conn) (session, error) {
	sess := &session{
		Conn: conn,
	}
	return sess, nil
}

func (self *session) Run() {

	go sendLoop()
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
