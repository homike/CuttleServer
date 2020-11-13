package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

var Services map[string]*Service
var Handlers map[string]*Handler

func init() {
	Services = make(map[string]*Service)
	Handlers = make(map[string]*Handler)
}

func Register(comp interface{}, opts []Option) error {
	s := NewService(comp, opts)

	if _, ok := Services[s.Name]; ok {
		return fmt.Errorf("handler: service already defined: %s", s.Name)
	}

	if err := s.ExtractHandler(); err != nil {
		return err
	}

	// register all localHandlers
	Services[s.Name] = s
	for name, handler := range s.Handlers {
		n := fmt.Sprintf("%s.%s", s.Name, name)
		log.Println("Register local handler", n)
		Handlers[n] = handler
	}
	return nil
}

func ProcessMessage(msgkey string, msgdata []byte) {
	handler, ok := Handlers[msgkey]
	if !ok {
		fmt.Println("no handler")
		return
	}

	var payload = msgdata
	var data interface{}
	if handler.IsRawArg {
		data = payload
	} else {
		data = reflect.New(handler.Type.Elem()).Interface()
		err := json.Unmarshal(payload, data)
		if err != nil {
			log.Println(fmt.Sprintf("Deserialize to %T failed: %+v (%v)", data, err, payload))
			return
		}
	}

	session := &Session{}
	args := []reflect.Value{handler.Receiver, reflect.ValueOf(session), reflect.ValueOf(data)}

	result := handler.Method.Func.Call(args)
	if len(result) > 0 {
		if err := result[0].Interface(); err != nil {
			log.Println(fmt.Sprintf("Service error: %+v", err))
		}
	}
}
