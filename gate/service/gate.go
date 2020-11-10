package service

import (
	"cuttleserver/common/network"
	"cuttleserver/gateserver/msghandler"
	"cuttleserver/gateserver/session"
	"fmt"
)

type Gate struct {
	Addr       string
	Port       int
	IsLittle   bool
	MsgHeadLen int

	Acceptor *network.Acceptor
}

func (self *Gate) Run() error {

	parser := network.NewMsgParser(self.MsgHeadLen, self.IsLittle)
	cprotoProc, err := msghandler.NewMsgProcessor()
	if err != nil {
		return err
	}
	proc := interface{}(cprotoProc).(network.MsgProcessor)

	// Create Acceptor
	acceptor, err := network.NewAcceptor(self.Addr, self.Port, parser, proc)
	if err != nil {
		return err
	}

	acceptor.NewSession = func(socketSession *network.SocketSession) network.Session {
		sess := &session.Session{SocketSession: socketSession}
		return sess
	}

	fmt.Println("[GateServer] Start, Listen", self.Addr, self.Port)
	acceptor.Start()

	return nil
}
