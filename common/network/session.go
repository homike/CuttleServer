package network

import (
	"net"
)

type Agent interface {
	Run()
	OnClose()
}

type Session struct {
	Conn      net.Conn
	WriteChan chan []byte

	MsgParser    *MsgParser
	MsgProcessor MsgProcessor
}

func NewSession(conn net.Conn, parser *MsgParser, proc MsgProcessor) (*Session, error) {
	sess := &Session{
		Conn:         conn,
		WriteChan:    make(chan []byte, 128),
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

func (self *Session) sendLoop() {
	for b := range self.WriteChan {
		if b == nil {
			break
		}

		_, err := self.Conn.Write(b)
		if err != nil {
			break
		}
	}
}

func (self *Session) Write(msgid uint16, message interface{}) error {
	byteMessage, err := self.MsgProcessor.Marshal(message)
	if err != nil {
		return err
	}
	data, err := self.MsgParser.Pack(msgid, byteMessage)
	if err != nil {
		return err
	}

	self.WriteChan <- data
	return nil
}
