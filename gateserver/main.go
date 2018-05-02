package main

import (
	"cuttleserver/gateserver/service"
)

func main() {
	gateServer := &service.Gate{
		Addr:       "127.0.0.1",
		Port:       9110,
		IsLittle:   true,
		MsgHeadLen: 6,
	}

	gateServer.Run()

	for {
	}
}
