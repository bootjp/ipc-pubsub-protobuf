package main

import (
	"log"

	"github.com/bootjp/ipc-pubsub-protobuf/daemon"
)

func main() {
	err := daemon.NewDaemon().Start()
	if err != nil {
		log.Fatal(err)
	}
}
