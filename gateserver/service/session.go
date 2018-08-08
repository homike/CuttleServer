package service

import (
	"cuttleserver/common/network"
)

type Session struct {
	network.Session
	// User Data
	SessionID uint32
	AccountID uint32
	Token     string
}

func (self *Session) OnEstablish() {
	sessionManager.AssociateSession(self)
}

func (self *Session) OnRecv(msgid uint16, msgdata []byte) {
	task := &Tasks{
		AccountID: self.SessionID,
		MessageID: uint32(msgid),
		Content:   msgdata,
	}
	gameServer1.tasks <- task

	//ret := DispatchMessage(msgid, msgdata, interface{}(self))
	// _ = ret

	/*
		if ret == cproto.CPROTORET_NO_HANDLER {
			// czxdo: forward to gameserver
		} else if ret != cproto.CPROTORET_OK {
			fmt.Println("message route error")
		}
	*/
	// inrc qps count
	InrcRequestCount()
}

func (self *Session) Send(msgid uint16, message interface{}) error {
	byteMessage, err := MessageHandlers.Proc.Marshal(message)
	if err != nil {
		return nil
	}

	self.WriteMessage(msgid, byteMessage)
	return nil
}

func (self *Session) SetSessionID(sessID uint32) error {
	self.SessionID = sessID
	return nil
}
