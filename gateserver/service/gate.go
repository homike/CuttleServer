package service

import (
	"cuttleserver/common/network"
	"cuttleserver/common/network/cproto"
	"fmt"
)

type Gate struct {
	Addr       string
	Port       int
	IsLittle   bool
	MsgHeadLen int

	Acceptor *network.Acceptor
}

var messageHandlers *MessageHandlers
var proc network.MsgProcessor
var gateServer *Gate
var gameServer1 *GameServer
var sessionManager *SessionManager

func Init() {
	messageHandlers = NewMessageHandlers()

	proc = network.MsgProcessor(cproto.NewCProto())

	gateServer = &Gate{
		Addr:       "0.0.0.0",
		Port:       9110,
		IsLittle:   true,
		MsgHeadLen: 4,
	}

	gameServer1 = NewGameServer()

	sessionManager = NewSessionManager()
}

func Run(close chan bool) error {
	StartGateWatch()
	// parser
	parser := network.NewMsgParserWithOption(gateServer.MsgHeadLen, gateServer.IsLittle)
	// acceptor
	acceptor, err := network.NewAcceptor(gateServer.Addr, gateServer.Port, parser)
	if err != nil {
		return err
	}
	acceptor.NewAgent = func(socketSess *network.Session) network.SessionInterface {
		sess := &Session{Session: *socketSess}
		return sess
	}
	fmt.Println("[GateServer] Start, Listen", gateServer.Addr, gateServer.Port)
	acceptor.Start()

	<-close
	fmt.Println("[GateServer] Close")

	return nil
}
