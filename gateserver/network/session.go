package network

import "net"

type session struct {
	Conn net.Conn
}

func NewSession(conn net.Conn) (session, error) {
	sess := &session{
		Conn: conn,
	}
	return sess, nil
}
