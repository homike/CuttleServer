package agent

import (
	"cuttleserver/common/network"
	"cuttleserver/common/network/cproto"
	"fmt"
)

type ConnectStatus int32

//const (
//	Connect_init           ConnectStatus = 0x1 << 0
//	Connect_verify_ing     ConnectStatus = 0x2 << 0
//	Connect_verify_succ    ConnectStatus = 0x3 << 0
//	Connect_entergame_ing  ConnectStatus = 0x4 << 0
//	Connect_entergame_succ ConnectStatus = 0x5 << 0
//)

type Agent struct {
	network.Session
	// User Data
	Token string
}

func (self *Agent) OnRecv(msgid uint16, msgdata []byte) {
	ret := self.MsgProcessor.Route(msgid, msgdata, interface{}(self))
	if ret == cproto.CPROTORET_NO_HANDLER {
		// czxdo: forward to gameserver
	} else if ret != cproto.CPROTORET_OK {
		fmt.Println("message route error")
	}

	// inrc qps count
	InrcRequestCount()
}

func (self *Agent) Send(msgid uint16, message interface{}) error {
	self.WriteMessage(msgid, message)
	return nil
}
