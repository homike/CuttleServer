package service

import (
	"cuttleserver/common/network"
	"cuttleserver/common/network/cproto"
)

var MessageHandlers *network.MessageHandlers
var SessionMgr *SessionManager

func Init() {
	MessageHandlers = network.NewMessageHandlers(network.MsgProcessor(cproto.NewCProto()))

	SessionMgr = NewSessionManager()
}
