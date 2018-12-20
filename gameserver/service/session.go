package service

import (
	"cuttleserver/common/network"
	"reflect"
)

type Session struct {
	network.Session
	AccountID uint32
	Token     string
}

func NewSession(accountID uint32) (*Session, error) {
	return &Session{
		AccountID: accountID,
	}, nil
}

func (self *Session) Handler(msgID uint16, msgBody []byte) (uint16, []byte) {
	msgInfo, ok := MessageHandlers.Handlers[msgID]
	if !ok {
		return 0, nil
	}

	// struct
	msgEntry := reflect.New(msgInfo.MsgType.Elem()).Interface()

	// unmarshal
	err := MessageHandlers.Proc.UnMarshal(msgBody, msgEntry)
	if err != nil {
		return 0, nil
	}

	// handler message
	respID, resp := msgInfo.MsgHandler([]interface{}{msgEntry, self})

	byteMessage, err := MessageHandlers.Proc.Marshal(resp)
	if err != nil {
		return 0, nil
	}

	return respID, byteMessage
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
