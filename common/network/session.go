package network

import (
	"net"
)

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

//func (self *Session) recvLoop() {
//	for self.Conn != nil {
//		msg, err := self.ReadMessage(self)
//
//		if err != nil {
//			if !util.IsEOFOrNetReadError(err) {
//				log.Errorf("session closed, sesid: %d, err: %s", self.ID(), err)
//			}
//
//			self.Close()
//
//			break
//		}
//	}
//
//	self.cleanup()
//}

func (self *Session) WriteMessage(msgid uint16, message interface{}) error {
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

func (self *Session) WriteBytes(msgid uint16, byteMessage []byte) error {
	data, err := self.MsgParser.Pack(msgid, byteMessage)
	if err != nil {
		return err
	}

	self.WriteChan <- data
	return nil
}
