package agent

import (
	"bufio"
	"cuttleserver/common/network"
	"cuttleserver/common/network/cproto"
	"fmt"
)

type ConnectStatus int32

const (
	Connect_init           ConnectStatus = 0x1 << 0
	Connect_verify_ing     ConnectStatus = 0x2 << 0
	Connect_verify_succ    ConnectStatus = 0x3 << 0
	Connect_entergame_ing  ConnectStatus = 0x4 << 0
	Connect_entergame_succ ConnectStatus = 0x5 << 0
)

type Agent struct {
	Session    *network.Session
	ConnStatus ConnectStatus
	// User Data
	Token string
}

func (self *Agent) Run() {
	self.ConnStatus = Connect_init

	bufReader := bufio.NewReader(self.Session.Conn)
	for {
		// handler messages
		msgID, msgBody, err := self.Session.MsgParser.UnPack(bufReader)
		if err != nil {
			fmt.Println("message read error")
			return
		}
		ret := self.Session.MsgProcessor.Route(msgID, msgBody, interface{}(self))
		InrcRequestCount()
		if ret == cproto.CPROTORET_NO_HANDLER {
			// czxdo: forward to gameserver
		} else if ret != cproto.CPROTORET_OK {
			fmt.Println("message route error")
		}
	}
}

func (self *Agent) Send(msgid uint16, message interface{}) error {
	self.Session.Write(msgid, message)
	return nil
}

func (self *Agent) OnClose() {
}
