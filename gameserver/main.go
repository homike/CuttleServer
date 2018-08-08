package main

import (
	"cuttleserver/gameserver/msghandler"
	"cuttleserver/gameserver/service"
)

func main() {
	service.Init()

	msghandler.Init()

	service.StartGRPC()
}
