package network

import (
	"errors"
	"fmt"
	"io"
	"net"
)

type Acceptor struct {
	IP        string
	Port      int
	Listener  net.Listener
	MsgParser *MsgParser
	SocketOption

	// constructor
	NewAgent func(sess *Session) SessionInterface
}

func NewAcceptor(addr string, port int, parser *MsgParser) (*Acceptor, error) {
	if parser == nil {
		return nil, errors.New("nil pointer")
	}

	return &Acceptor{
		IP:        addr,
		Port:      port,
		MsgParser: parser,
	}, nil
}

func (self *Acceptor) Start() error {
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

	sess, err := NewSession(conn, self.MsgParser)
	if err != nil {
		return err
	}

	if self.NewAgent == nil {
		return errors.New("Acceptor NewSession function is nil")
	}
	agent := self.NewAgent(sess)

	self.Run(agent)

	return nil
}

func isEOFOrNetReadError(err error) bool {
	if err == io.EOF {
		return true
	}
	ne, ok := err.(*net.OpError)
	return ok && ne.Op == "read"
}

func (self *Acceptor) Run(sess SessionInterface) {

	sess.OnEstablish()

	for sess.ConnectAvailable() {
		msgid, msgdata, err := sess.ReadMessage()
		if err != nil {
			if !isEOFOrNetReadError(err) {
				fmt.Sprintf("session closed, err: %s \n", err)
			}
			sess.Close()
			break
		}
		sess.OnRecv(msgid, msgdata)
	}
}
