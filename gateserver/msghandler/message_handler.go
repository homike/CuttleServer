package msghandler

import (
	"cuttleserver/common/network/cproto"
	"cuttleserver/gateserver/session"
	"fmt"
	"reflect"
)

//func Dispatch(msgID uint16, msgBody []byte, tc *network.TCPClient) {
//
//	sess, err := sessions.SessionMgr.FindSession(tc.AccountID)
//	if err != nil {
//		sess = sessions.NewSession(tc)
//	}
//
//	processFunc, ok := MapFunc[msgID]
//	if ok {
//		processFunc(sess, msgBody)
//	}
//}

const (
	Protocol_Test_Req                = 1
	Protocol_Test_Resp               = 2
	Protocol_GetSystemTime_Req       = 3
	Protocol_GetSystemTime_Resp      = 4
	Protocol_LoginServerResult_Ntf   = 1001
	Protocol_LoginServerPlatform_Req = 1007
)

type TestReq struct {
	Value int32
}

func NewMsgProcessor() (*cproto.CProto, error) {
	proto := cproto.NewCProto()

	proto.SetHandler(1, cproto.MsgInfo{MsgType: reflect.TypeOf(&TestReq{}), MsgHandler: TestReqProcess})

	return proto, nil
}

func TestReqProcess(args []interface{}) {
	msg := args[0].(*TestReq)
	sess := args[1].(*session.Session)

	fmt.Println("args ", msg.Value, "Token", sess.Token)
}
