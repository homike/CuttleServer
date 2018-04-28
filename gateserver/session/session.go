package session

import (
	"bufio"
	"cuttleserver/common/network"
	"fmt"
)

type Session struct {
	SocketSession *network.SocketSession
	// User Data
	Token string
}

func (self *Session) Run() {
	bufReader := bufio.NewReader(self.SocketSession.Conn)
	for {
		// handler messages
		msgID, msgBody, err := self.SocketSession.MsgParser.UnPack(bufReader)
		if err != nil {
			fmt.Println("message read error")
			return
		}
		self.SocketSession.MsgProcessor.Route(msgID, msgBody, interface{}(self))
	}
}

func (self *Session) OnClose() {
}
