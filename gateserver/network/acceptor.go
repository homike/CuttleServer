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
		fmt.Printf("Acceptor Start(%v) error (%v) \n", addr, err)
	}

}
