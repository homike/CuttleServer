package service

import (
	"cuttleserver/common/network/cproto"
	"errors"
	"reflect"
)

type MsgHandler func([]interface{})
type MsgInfo struct {
	MsgType    reflect.Type
	MsgHandler MsgHandler
}

type MessageHandlers struct {
	Handlers map[uint16]MsgInfo
}

func NewMessageHandlers() *MessageHandlers {
	return &MessageHandlers{
		Handlers: make(map[uint16]MsgInfo),
	}
}

func RegisterHandler(msgID uint16, msgInfo MsgInfo) error {
	_, ok := messageHandlers.Handlers[msgID]
	if ok {
		return errors.New("exist handler")
	}
	messageHandlers.Handlers[msgID] = msgInfo

	return nil
}

func DispatchMessage(msgID uint16, msgBody []byte, userData interface{}) int8 {
	msgInfo, ok := messageHandlers.Handlers[msgID]
	if !ok {
		return cproto.CPROTORET_NO_HANDLER
	}

	// struct
	msgEntry := reflect.New(msgInfo.MsgType.Elem()).Interface()

	// unmarshal
	err := proc.UnMarshal(msgBody, msgEntry)
	if err != nil {
		return cproto.CPROTORET_MSG_FORMAT_ERROR
	}

	// dispatch message
	msgInfo.MsgHandler([]interface{}{msgEntry, userData})

	return cproto.CPROTORET_OK
}
