package network

import (
	"net"
)

type Session interface {
	Run()
	OnClose()
}

type SocketSession struct {
	Conn       net.Conn
	WriteChann chan []byte

	MsgParser    *MsgParser
	MsgProcessor MsgProcessor
	// User Data
}

func NewSocketSession(conn net.Conn, parser *MsgParser, proc MsgProcessor) (*SocketSession, error) {
	sess := &SocketSession{
		Conn:         conn,
		WriteChann:   make(chan []byte, 128),
		MsgParser:    parser,
		MsgProcessor: proc,
	}

	go sess.sendLoop()

	return sess, nil
}

//func (self *SocketSession) RecvLoop() {
//	bufReader := bufio.NewReader(self.Conn)
//	for {
//		// handler messages
//		msgID, msgBody, err := self.MsgParser.UnPack(bufReader)
//		if err != nil {
//			fmt.Println("message read error")
//			return
//		}
//		self.MsgProcessor.Route(msgID, msgBody, interface{}(self))
//	}
//}

func (self *SocketSession) sendLoop() {
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
