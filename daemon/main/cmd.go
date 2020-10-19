package main

import (
	"github.com/bootjp/ipc-pubsub-protobuf/daemon"
	"log"
)

func main() {
	err := daemon.NewDaemon().Start()
	if err != nil {
		log.Fatal()
	}
}