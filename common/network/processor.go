package network

type MsgProcessor interface {
	// Must be gorutine safe
	Marshal(msg interface{}) ([]byte, error)
	// Must be gorutine safe
	// TODO: paramter "msg" can not input?
	UnMarshal(data []byte, msg interface{}) error
	// Must be gorutine safe
	Route(msgID uint16, msgBody []byte, userData interface{}) int8
}
