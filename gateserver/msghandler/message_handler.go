package msghandler

import (
	"cuttleserver/common/network/cproto"
	"cuttleserver/gateserver/agent"
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

func NewMsgProcessor() (*cproto.CProto, error) {
	proto := cproto.NewCProto()

	proto.SetHandler(Protocol_Ping, cproto.MsgInfo{MsgType: reflect.TypeOf(&Ping{}), MsgHandler: TestReqProcess})

	return proto, nil
}

func TestReqProcess(args []interface{}) {
	req := args[0].(*Ping)
	agent := args[1].(*agent.Agent)

	resp := &Pong{
		Value: req.Value + 1,
	}
	agent.Send(Protocol_Pong, resp)
	//fmt.Println("args ", msg.Value, "Token", agent.Token)
}
