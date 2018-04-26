package cproto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

type CProto struct {
	mapHandlers map[uint16]MsgInfo
}

type MsgHandler func([]interface{})

type MsgInfo struct {
	MsgType    reflect.Type
	MsgHandler MsgHandler
}

func (self *CProto) SetHandler(msgID uint16, msgInfo MsgInfo) error {
	_, ok := self.mapHandlers[msgID]
	if ok {
		return errors.New("exist handler")
	}
	self.mapHandlers[msgID] = msgInfo

	return nil
}

func (self *CProto) Route(msgID int16, msgBody []byte, userdata interface{}) error {
	msgInfo, ok := self.mapHandlers[msgID]
	if !ok {
		return errors.New("not exist handler")
	}

	// UnMarshal
	msgEntry := reflect.New(msgInfo.MsgType.Elem())
	err := self.UnMarshal(msgBody, msgEntry)
	if err != nil {
		return err
	}

	// Route
	msgInfo.MsgHandler([]interface{}{userdata, msgEntry})

	return nil
}

// Gorutine safe
func (p *CProto) UnMarshal(msgBody []byte, msgStruct interface{}) error {

	readIndex := 0
	v := reflect.ValueOf(msgStruct).Elem(i)
	//vType := v.Type()
	for i := 0; i < v.NumField(); i++ {
		vf := v.Field(i)
		//tf := vType.Field(i)

		//fmt.Println(tf.Name, vf.Kind())
		switch vf.Kind() {
		case reflect.String:
			for i := readIndex; i < len(msgBody); i++ {
				if msgBody[i] == byte(0) {
					//fmt.Println(readIndex, "string :", string(msgBody[readIndex:i]))
					canSetValue := reflect.ValueOf(string(msgBody[readIndex:i]))
					vf.Set(canSetValue)
					readIndex = i + 1
					break
				}
			}

		case reflect.Int32:
			//fmt.Println(readIndex, "int :", msgBody[readIndex:readIndex+4])
			var intValue int32
			bytesBuffer := bytes.NewBuffer(msgBody[readIndex : readIndex+4])
			binary.Read(bytesBuffer, binary.LittleEndian, &intValue)
			canSetValue := reflect.ValueOf(int32(intValue))
			vf.Set(canSetValue)
			readIndex = readIndex + 4
		default:
		}
	}
	return nil
}

// Gorutine safe
func (self *CProto) Marshal(msgStruct interface{}) ([]byte, error) {

	bytesBuffer := bytes.NewBuffer([]byte{})

	v := reflect.ValueOf(msgStruct).Elem()
	for i := 0; i < v.NumField(); i++ {
		vf := v.Field(i)

		vfBytes := marshal(vf)
		binary.Write(bytesBuffer, binary.LittleEndian, vfBytes)
	}

	return bytesBuffer.Bytes(), nil
}

func marshal(v reflect.Value) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})

	switch v.Kind() {
	case reflect.String:
		err := binary.Write(bytesBuffer, binary.LittleEndian, []byte(v.String()))
		if err != nil {
			fmt.Println("czx@@@ write string error :", err)
		}
		binary.Write(bytesBuffer, binary.LittleEndian, byte(0))

	case reflect.Uint8:
		binary.Write(bytesBuffer, binary.LittleEndian, uint8(v.Uint()))

	case reflect.Int8:
		binary.Write(bytesBuffer, binary.LittleEndian, int8(v.Int()))

	case reflect.Int32:
		binary.Write(bytesBuffer, binary.LittleEndian, int32(v.Int()))

	case reflect.Uint32:
		binary.Write(bytesBuffer, binary.LittleEndian, uint32(v.Uint()))

	case reflect.Int64:
		binary.Write(bytesBuffer, binary.LittleEndian, v.Int())

	case reflect.Uint64:
		binary.Write(bytesBuffer, binary.LittleEndian, v.Uint())

	case reflect.Bool:
		b := 0
		if v.Bool() {
			b = 1
		}
		binary.Write(bytesBuffer, binary.LittleEndian, uint8(b))

	case reflect.Slice:
		binary.Write(bytesBuffer, binary.LittleEndian, int32(v.Len()))
		for j := 0; j < v.Len(); j++ {
			data := v.Slice(j, j+1).Index(0)
			sliceBytes := marshal(data)
			binary.Write(bytesBuffer, binary.LittleEndian, sliceBytes)
		}

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			vf := v.Field(i)
			vfBytes := marshal(vf)
			binary.Write(bytesBuffer, binary.LittleEndian, vfBytes)
		}

	case reflect.Ptr:
		ve := v.Elem()
		for i := 0; i < ve.NumField(); i++ {
			vf := ve.Field(i)
			vfBytes := marshal(vf)
			binary.Write(bytesBuffer, binary.LittleEndian, vfBytes)
		}

	default:
		binary.Write(bytesBuffer, binary.LittleEndian, v.Bytes())
	}

	return bytesBuffer.Bytes()
}
