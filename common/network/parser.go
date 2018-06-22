package network

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
)

// --------------
// | len | data |
// --------------
type MsgParser struct {
	MsgLen         int
	IsLittleEndian bool
}

func NewMsgParser(len int, isLittle bool) *MsgParser {
	return &MsgParser{
		MsgLen:         len,
		IsLittleEndian: isLittle,
	}
}

func (self *MsgParser) UnPack(bufReader *bufio.Reader) (uint16, []byte, error) {
	// Header
	var headerSize uint32
	err := errors.New("read error")
	if self.IsLittleEndian {
		err = binary.Read(bufReader, binary.LittleEndian, &headerSize)
	} else {
		err = binary.Read(bufReader, binary.BigEndian, &headerSize)
	}
	if err != nil {
		return 0, nil, err
	}

	// MsgID
	var msgID uint16
	if self.IsLittleEndian {
		err = binary.Read(bufReader, binary.LittleEndian, &msgID)
	} else {
		err = binary.Read(bufReader, binary.BigEndian, &msgID)
	}
	if err != nil {
		return 0, nil, err
	}

	bodySize := headerSize - uint32(self.MsgLen)

	// BodyData
	bodyData := make([]byte, bodySize)
	if self.IsLittleEndian {
		err = binary.Read(bufReader, binary.LittleEndian, &bodyData)
	} else {
		err = binary.Read(bufReader, binary.BigEndian, &bodyData)
	}
	if err != nil {
		return 0, nil, err
	}

	return msgID, bodyData, nil
}

func (self *MsgParser) Pack(msgID uint16, message []byte) ([]byte, error) {

	if len(message) <= 0 {
		return nil, errors.New("empty message")
	}

	w := bytes.NewBuffer([]byte{})
	if self.IsLittleEndian {
		binary.Write(w, binary.LittleEndian, uint32(len(message)+self.MsgLen))
		binary.Write(w, binary.LittleEndian, msgID)
		binary.Write(w, binary.LittleEndian, message)
	} else {
		binary.Write(w, binary.BigEndian, uint32(len(message)+self.MsgLen))
		binary.Write(w, binary.BigEndian, msgID)
		binary.Write(w, binary.BigEndian, message)
	}
	return w.Bytes(), nil
}
