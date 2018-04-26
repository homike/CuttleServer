package service

import (
	"cuttleserver/common/network/"
	"cuttleserver/common/network/cproto"
)

type Gate struct {
	Acceptor  *network.Acceptor
	Parser    *network.Parser
	Processor network.Processor
}

func NewGate() (*Gate, error) {
	gate := New(Gate)

	parser := network.NewMsgParser(6, true)
	proc := New(cproto.CProto)
	acceptor, err := network.NewAcceptor("127.0.0.1", 6370, parser, proc)
	if err != nil {
		return nil, err
	}

	gate.Parser = parser
	gate.Processor = proc
	gate.Acceptor = acceptor

	return gate, nil
}

func (self *Gate) Run() error {
	self.Acceptor.Run()
}
