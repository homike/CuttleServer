package network

import (
	"bufio"
	"fmt"
	"net"
)

type session struct {
	Conn       net.Conn
	WriteChann chan []byte

	MsgParser    *MsgParser
	MsgProcessor MsgProcessor
	// User Data
}

func NewSession(conn net.Conn, parser *MsgParser, proc MsgProcessor) (*session, error) {
	sess := &session{
		Conn:         conn,
		WriteChann:   make(chan []byte, 128),
		MsgParser:    parser,
		MsgProcessor: proc,
	}

	go sess.sendLoop()

	return sess, nil
}

func (self *session) Run() {
	bufReader := bufio.NewReader(self.Conn)
	for {
		// handler messages
		msgID, msgBody, err := self.MsgParser.UnPack(bufReader)
		if err != nil {
			fmt.Println("message read error")
			return
		}
		self.MsgProcessor.Route(msgID, interface{}(self))
	}
}

func (self *session) sendLoop() {
	for b := range self.WriteChann {
		if b == nil {
			break
		}

		_, err := self.Conn.Write(b)
		if err != nil {
			break
		}
	}
}
