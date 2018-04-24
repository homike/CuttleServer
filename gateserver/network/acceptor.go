package network

import (
	"fmt"
	"net"
)

type Acceptor struct {
	Listener net.Listener
}

func (self *Acceptor) Start(addr string) error {

	ln, err := net.Listen("TCP", addr)
	if err != nil {
		fmt.Printf("Acceptor Start(%v) error: %v \n", addr, err)
		return err
	}

	self.Listener = ln

	go self.Accept()
}

func (self *Acceptor) Accept() error {
	for {
		conn, err := self.Listener.Accept()
		if err != nil {
			fmt.Printf("Accepter Accept() error: %v \n ", addr, err)
			return err
		}

		go self.OnAccept(conn)
	}
}

func (self *Acceptor) OnAccept(conn net.Conn) error {
	defer conn.Close()

	sess, err := NewSession(conn)
	if err != nil {
		return err
	}

	// recv messages, dispatch handler
	sess.Run()

	return nil
}
