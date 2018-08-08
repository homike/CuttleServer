package msghandler

import (
	"cuttleserver/common/network"
	"cuttleserver/gameserver/service"
	"fmt"
	"reflect"
)

const (
	Protocol_Ping                    = 1
	Protocol_Pong                    = 2
	Protocol_GetSystemTime_Req       = 3
	Protocol_GetSystemTime_Resp      = 4
	Protocol_LoginServerResult_Ntf   = 1001
	Protocol_LoginServerPlatform_Req = 1007
)

type Ping struct {
	Value int32
}

type Pong struct {
	Value int32
}

func Init() error {
	service.MessageHandlers.RegisterHandler(Protocol_Ping, network.MsgInfo{reflect.TypeOf(&Ping{}), TestReqProcess})
	return nil
}

func TestReqProcess(args []interface{}) (uint16, interface{}) {
	req := args[0].(*Ping)
	sess := args[1].(*service.Session)

	resp := &Pong{
		Value: req.Value + 1,
	}

	_ = sess
	//sess.Send(Protocol_Pong, resp)
	fmt.Println("req", req.Value)

	return uint16(Protocol_Pong), resp
}
