package network

import "net"

type SocketOption struct {
	readBufferSize  int
	writeBufferSize int
	noDelay         bool
}

func (self *SocketOption) SetSocketOption(readSize, writeSize int, noDelay bool) {
	self.readBufferSize = readSize
	self.writeBufferSize = writeBufferSize
	self.noDelay = noDelay
}

func (self *SocketOption) ApplySocketOption(conn *net.Conn) {
	if tcp, ok := conn.(*net.TCPConn); ok {
		if self.readBufferSize > 0 {
			tcp.SetReadBuffer(self.readBufferSize)
		}
		if self.writeBufferSize > 0 {
			tcp.SetWriteBuffer(self.writeBufferSize)
		}
		tcp.SetNoDelay(self.noDelay)
	}
}

func (self *SocketOption) Init() {
	self.readBufferSize = -1
	self.writeBufferSize = -1
	self.noDelay = false
}
