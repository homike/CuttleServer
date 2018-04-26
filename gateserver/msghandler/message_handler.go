package msghandler

import (
	"cuttleserver/gateserver/common/cproto"
	"fmt"
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

func NewMsgProcessor() (*cproto.CProto, error) {
	proto := New(cproto.CProto)

	proto.SetHandler(1, cproto.MsgInfo{MsgType: LoginServerPlatformReq, MsgHandler: ProcessLoginServerPlatformReq})

	return proto, nil
}

// 1
type LoginServerPlatformReq struct {
	Takon     string
	Version   int32
	ChannelID string
}

func ProcessLoginServerPlatformReq(args []interface{}) {

	msg := args[0].(LoginServerPlatformReq)
	//sess := args[1].(

	fmt.Println("args ", msg)
}
