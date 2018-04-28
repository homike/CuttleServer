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

func NewMsgProcessor() (*cproto.CProto, error) {
	proto := cproto.NewCProto()

	proto.SetHandler(1, cproto.MsgInfo{MsgType: reflect.TypeOf(LoginServerPlatformReq{}), MsgHandler: ProcessLoginServerPlatformReq})

	return proto, nil
}

type LoginServerPlatformReq struct {
	Takon     string
	Version   int32
	ChannelID string
}

func ProcessLoginServerPlatformReq(args []interface{}) {
	msg := args[0].(LoginServerPlatformReq)
	sess := args[1].(session.Session)

	fmt.Println("args ", msg, sess.Token)
}
