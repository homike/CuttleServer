package session

import (
	"bufio"
	"cuttleserver/common/network"
	"cuttleserver/common/network/cproto"
	"fmt"
)

type SessionConnectStatus int32

const (
	Connect_init           SessionConnectStatus = 0x1 << 0
	Connect_verify_ing     SessionConnectStatus = 0x2 << 0
	Connect_verify_succ    SessionConnectStatus = 0x3 << 0
	Connect_entergame_ing  SessionConnectStatus = 0x4 << 0
	Connect_entergame_succ SessionConnectStatus = 0x5 << 0
)

type Session struct {
	SocketSession *network.SocketSession
	ConnStatus    SessionConnectStatus
	// User Data
	Token string
}

func (self *Session) Run() {
	self.ConnStatus = Connect_init

	bufReader := bufio.NewReader(self.SocketSession.Conn)
	for {
		// handler messages
		msgID, msgBody, err := self.SocketSession.MsgParser.UnPack(bufReader)
		if err != nil {
			fmt.Println("message read error")
			return
		}
		ret := self.SocketSession.MsgProcessor.Route(msgID, msgBody, interface{}(self))
		if ret == cproto.CPROTORET_NO_HANDLER {
			// forward to gameserver
		} else if ret != cproto.CPROTORET_OK {
			fmt.Println("message route error")
		}
	}
}

func (self *Session) OnClose() {
}
