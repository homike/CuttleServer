package handler

import (
	"encoding/json"
	"fmt"
	"testing"
)

// A Service
type AReq struct {
	Uid  int64  `json:"Uid"`
	Data string `json:"data"`
}

type AService struct {
}

func (self *AService) Request(s *Session, msg *AReq) error {
	fmt.Println("AService.Request, msg:  ", msg)
	return nil
}

// B Service
type BReq struct {
	Uid  int64  `json:"Uid"`
	Data string `json:"data"`
}

type BService struct {
}

func (self *BService) Request(s *Session, msg *BReq) error {
	fmt.Println("BService.Request, msg:  ", msg)
	return nil
}

func TestService(t *testing.T) {
	a := &AService{}
	err := Register(a, nil)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	b := &BService{}
	err = Register(b, nil)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	// A
	{
		req := AReq{
			Uid:  111,
			Data: "aaa",
		}
		reqBytes, _ := json.Marshal(req)
		ProcessMessage("AService.Request", reqBytes)
	}
	// B
	{
		req := BReq{
			Uid:  222,
			Data: "bbb",
		}
		reqBytes, _ := json.Marshal(req)
		ProcessMessage("BService.Request", reqBytes)
	}
}
