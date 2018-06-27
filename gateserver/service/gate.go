package service

import (
	"cuttleserver/common/network"
	"cuttleserver/gateserver/agent"
	"cuttleserver/gateserver/msghandler"
	"fmt"
)

type Gate struct {
	Addr       string
	Port       int
	IsLittle   bool
	MsgHeadLen int

	Acceptor *network.Acceptor
}

func (self *Gate) Run(close chan bool) error {
	agent.InitAgentManager()

	// parser
	parser := network.NewMsgParser(self.MsgHeadLen, self.IsLittle)

	// processor
	cprotoProc, err := msghandler.NewMsgProcessor()
	if err != nil {
		return err
	}
	proc := network.MsgProcessor(cprotoProc)

	// acceptor
	acceptor, err := network.NewAcceptor(self.Addr, self.Port, parser, proc)
	if err != nil {
		return err
	}
	acceptor.NewAgent = func(sess *network.Session) network.SessionInterface {
		agent := &agent.Agent{Session: *sess}
		return agent
	}
	fmt.Println("[GateServer] Start, Listen", self.Addr, self.Port)
	acceptor.Start()

	<-close
	fmt.Println("[GateServer] Close")

	return nil
}
