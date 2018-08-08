package service

import (
	"reflect"
)

type Session struct {
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
