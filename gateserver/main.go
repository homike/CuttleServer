package main

import (
	"cuttleserver/gateserver/service"
)

func main() {
	gateServer, err := service.NewGate()

	acceptor.Start()
}
