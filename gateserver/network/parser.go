package network

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
)

type Parser struct {
}

func (m *MsgParser) Read(bufReader *bufio.Reader) (uint16, []byte, error) {
	var headerSize uint32
	err := binary.Read(bufReader, binary.LittleEndian, &headerSize)
	if err != nil {
		log.Println("read headsize error")
		return 0, nil, err
	}

	var msgID uint16
	err = binary.Read(bufReader, binary.LittleEndian, &msgID)
	if err != nil {
		log.Println("read msgid error")
		return 0, nil, err
	}

	bodySize := headerSize - uint32(m.MsgLen)
	bodyData := make([]byte, bodySize)
	err = binary.Read(bufReader, binary.LittleEndian, &bodyData)
	if err != nil {
		log.Println("read body error")
		return 0, nil, err
	}

	return msgID, bodyData, nil
}

func (m *MsgParser) Write(msgID uint16, msgStruct interface{}) []byte {
	message := m.MsgProcessor.Marshal(msgStruct)

	w := bytes.NewBuffer([]byte{})
	binary.Write(w, binary.LittleEndian, uint32(len(message)+m.MsgLen))
	binary.Write(w, binary.LittleEndian, msgID)
	binary.Write(w, binary.LittleEndian, message)

	return w.Bytes()
	//client.Write(w.Bytes())
}
