package network

import (
	"io"
	"net"
)

type SessionInterface interface {
	OnEstablish()
	OnTerminate()
	OnError(err error)
	OnRecv(msgid uint16, msgdata []byte)

	ConnectAvailable() bool
	ReadMessage() (uint16, []byte, error)
	Close()
}

type Session struct {
	Conn      net.Conn
	WriteChan chan []byte

	MsgParser *MsgParser
}

func NewSession(conn net.Conn, parser *MsgParser) (*Session, error) {
	sess := &Session{
		Conn:      conn,
		WriteChan: make(chan []byte, 128),
		MsgParser: parser,
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

func (self *Session) WriteMessage(msgid uint16, message []byte) error {
	data, err := self.MsgParser.Pack(msgid, message)
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

func (self *Session) ConnectAvailable() bool {
	return self.Conn != nil
}

func (self *Session) ReadMessage() (uint16, []byte, error) {
	return self.MsgParser.UnPack((self.Conn).(io.Reader))
}

func (self *Session) Close() {
	self.Conn.Close()
}

func (self *Session) OnEstablish() {
}

func (self *Session) OnTerminate() {
}

func (self *Session) OnError(err error) {
}

func (self *Session) OnRecv(msgid uint16, msgdata []byte) {
}
