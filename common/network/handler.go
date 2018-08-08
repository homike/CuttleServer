package network

import (
	"cuttleserver/common/network/cproto"
	"errors"
	"reflect"
)

type MsgHandler func([]interface{}) (uint16, interface{})
type MsgInfo struct {
	MsgType    reflect.Type
	MsgHandler MsgHandler
}

type MessageHandlers struct {
	Handlers map[uint16]MsgInfo
	Proc     MsgProcessor
}

func NewMessageHandlers(proc MsgProcessor) *MessageHandlers {
	return &MessageHandlers{
		Handlers: make(map[uint16]MsgInfo),
		Proc:     proc,
	}
}

func (self *MessageHandlers) RegisterHandler(msgID uint16, msgInfo MsgInfo) error {
	_, ok := self.Handlers[msgID]
	if ok {
		return errors.New("exist handler")
	}
	self.Handlers[msgID] = msgInfo

	return nil
}

func (self *MessageHandlers) DispatchMessage(msgID uint16, msgBody []byte, userData interface{}) int8 {
	msgInfo, ok := self.Handlers[msgID]
	if !ok {
		return cproto.CPROTORET_NO_HANDLER
	}

	// struct
	msgEntry := reflect.New(msgInfo.MsgType.Elem()).Interface()

	// unmarshal
	err := self.Proc.UnMarshal(msgBody, msgEntry)
	if err != nil {
		return cproto.CPROTORET_MSG_FORMAT_ERROR
	}

	// dispatch message
	msgInfo.MsgHandler([]interface{}{msgEntry, userData})

	return cproto.CPROTORET_OK
}
