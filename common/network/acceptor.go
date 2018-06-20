package network

import (
	"errors"
	"fmt"
	"net"
)

type Acceptor struct {
	IP           string
	Port         int
	Listener     net.Listener
	MsgParser    *MsgParser
	MsgProcessor MsgProcessor
	// Output interface
	NewAgent func(sess *Session) Agent
}

func NewAcceptor(addr string, port int, p *MsgParser, proc MsgProcessor) (*Acceptor, error) {
	if p == nil || proc == nil {
		return nil, errors.New("nil pointer")
	}

	return &Acceptor{
		IP:           addr,
		Port:         port,
		MsgParser:    p,
		MsgProcessor: proc,
	}, nil
}

func (self *Acceptor) Start() error {
	// Socket Listen
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(self.IP), self.Port, ""})
	if err != nil {
		fmt.Println("Listen failed", err.Error())
		return nil
	}

	self.Listener = listen

	go self.Accept()

	return nil
}

func (self *Acceptor) Accept() error {
	for {
		conn, err := self.Listener.Accept()
		if err != nil {
			fmt.Printf("Accepter Accept() error: %v \n ", err)
			return err
		}

		go self.OnAccept(conn)
	}
}

func (self *Acceptor) OnAccept(conn net.Conn) error {
	defer conn.Close()

	sess, err := NewSession(conn, self.MsgParser, self.MsgProcessor)
	if err != nil {
		return err
	}
	if self.NewAgent == nil {
		return errors.New("Acceptor NewSession function is nil")
	}

	// new agent, who impl agent interface
	agent := self.NewAgent(sess)

	// recv messages, dispatch handler
	agent.Run()

	return nil
}
