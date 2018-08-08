package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// --------------
// | len | data |
// --------------
type MsgParser struct {
	MsgLen         int
	IsLittleEndian bool
}

func NewMsgParser() *MsgParser {
	return &MsgParser{
		MsgLen:         4,
		IsLittleEndian: false,
	}
}

func NewMsgParserWithOption(len int, isLittle bool) *MsgParser {
	return &MsgParser{
		MsgLen:         len,
		IsLittleEndian: isLittle,
	}
}

func (self *MsgParser) UnPack(reader io.Reader) (uint16, []byte, error) {
	// length
	var size int
	sizeBuf := make([]byte, self.MsgLen)
	fmt.Println("len ", self.MsgLen)
	_, err := io.ReadFull(reader, sizeBuf)
	if err != nil {
		fmt.Println("read size error")
		return 0, nil, err
	}
	if self.MsgLen == 4 {
		if self.IsLittleEndian {
			size = int(binary.LittleEndian.Uint32(sizeBuf))
		} else {
			size = int(binary.BigEndian.Uint32(sizeBuf))
		}
	} else if self.MsgLen == 2 {
		if self.IsLittleEndian {
			size = int(binary.LittleEndian.Uint16(sizeBuf))
		} else {
			size = int(binary.BigEndian.Uint16(sizeBuf))
		}
	}

	fmt.Println("size ", size)
	// data
	dataBuf := make([]byte, size-self.MsgLen)
	_, err = io.ReadFull(reader, dataBuf)
	if err != nil {
		fmt.Println("read data error")
		return 0, nil, err
	}

	// message id
	var msgID uint16
	if self.IsLittleEndian {
		msgID = binary.LittleEndian.Uint16(dataBuf)
	} else {
		msgID = binary.BigEndian.Uint16(dataBuf)
	}

	// message data
	msgData := dataBuf[2:]

	fmt.Println("msgid ", msgID)
	return msgID, msgData, nil
}

func (self *MsgParser) Pack(msgID uint16, message []byte) ([]byte, error) {
	if len(message) <= 0 {
		return nil, errors.New("empty message")
	}

	pkt := make([]byte, self.MsgLen+2+len(message))
	// length
	if self.MsgLen == 4 {
		binary.LittleEndian.PutUint32(pkt, uint32(self.MsgLen+2+len(message)))
	} else if self.MsgLen == 2 {
		binary.LittleEndian.PutUint16(pkt, uint16(self.MsgLen+2+len(message)))
	}
	// message id
	binary.LittleEndian.PutUint16(pkt[self.MsgLen:], uint16(msgID))
	// message data
	copy(pkt[self.MsgLen+2:], message)

	return pkt, nil
}
