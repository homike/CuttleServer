package main

import (
	"cuttleserver/gateserver/msghandler"
	"cuttleserver/gateserver/service"
)

func main() {
	service.Init()

	msghandler.Init()

	close := make(chan bool)
	service.Run(close)
}
