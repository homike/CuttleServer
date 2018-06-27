package network

import (
	"bytes"
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

func NewMsgParser(len int, isLittle bool) *MsgParser {
	return &MsgParser{
		MsgLen:         len,
		IsLittleEndian: isLittle,
	}
}

func (self *MsgParser) UnPack(reader io.Reader) (uint16, []byte, error) {
	// buffer size
	sizeBuf := make([]byte, 4)
	_, err := io.ReadFull(reader, sizeBuf)
	if err != nil {
		fmt.Println("read size error")
		return 0, nil, err
	}
	var size uint32
	if self.IsLittleEndian {
		size = binary.LittleEndian.Uint32(sizeBuf)
	} else {
		size = binary.BigEndian.Uint32(sizeBuf)
	}

	// data
	dataBuf := make([]byte, size-4)
	_, err = io.ReadFull(reader, dataBuf)
	if err != nil {
		fmt.Println("read data error")
		return 0, nil, err
	}
	//fmt.Println("//message ", dataBuf)

	// MsgID
	var msgID uint16
	if self.IsLittleEndian {
		msgID = binary.LittleEndian.Uint16(dataBuf)
	} else {
		msgID = binary.BigEndian.Uint16(dataBuf)
	}

	// BodyData
	msgData := dataBuf[2:]

	return msgID, msgData, nil
}

func (self *MsgParser) UnPack2(bufReader io.Reader) (uint16, []byte, error) {
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

/*
func RecvLTVPacket(reader io.Reader, maxPacketSize int) (msg interface{}, err error) {
	// Size为uint16，占2字节
	var sizeBuffer = make([]byte, 2)

	// 持续读取Size直到读到为止
	_, err = io.ReadFull(reader, sizeBuffer)

	// 发生错误时返回
	if err != nil {
		return
	}

	// 用小端格式读取Size
	size := binary.LittleEndian.Uint16(sizeBuffer)

	if maxPacketSize > 0 && size >= uint16(maxPacketSize) {
		return nil, ErrMaxPacket
	}

	// 分配包体大小
	body := make([]byte, size)

	// 读取包体数据
	_, err = io.ReadFull(reader, body)

	// 发生错误时返回
	if err != nil {
		return
	}

	msgid := binary.LittleEndian.Uint16(body)

	msgData := body[2:]

	// 将字节数组和消息ID用户解出消息
	msg, _, err = codec.DecodeMessage(int(msgid), msgData)
	if err != nil {
		// TODO 接收错误时，返回消息
		return nil, err
	}

	return
}
*/

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
