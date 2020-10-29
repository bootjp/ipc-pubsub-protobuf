package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"

	_ "github.com/golang/protobuf/jsonpb"
)

const producerSock = "ipc-pubsub-producer.sock"
const consumerSock = "ipc-pubsub-consumer.sock"

//var breakLine = make([]byte, len("\n"))

func TestName(t *testing.T) {

	addr, err := net.ResolveUnixAddr("unix", filepath.Join(os.TempDir(), producerSock))
	if err != nil {
		t.Errorf("%v", err)
	}

	wg := sync.WaitGroup{}

	{
		wg.Add(1)
		conn, err := net.DialUnix("unix", nil, addr)
		if err != nil {
			t.Error(err)
		}

		cmd := fmt.Sprintf("+SUBSCRIBE PROTOBUF\r\n")
		_, err = conn.Write([]byte(cmd))

		if err != nil {
			log.Printf("error: %v\n", err)
			return
		}
		var response = make([]byte, 10)
		_, err = conn.Read(response)
		if err != nil {
			t.Error(err)
		}

		fmt.Println("res " + string(response))

		err = conn.Close()
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}
	{
		wg.Add(1)
		conn, err := net.DialUnix("unix", nil, addr)
		if err != nil {
			t.Error(err)
		}
		cmd := fmt.Sprintf("+PUBLISH PROTOBUF\r\n$10\r\nPROTO_DATA\r\n")
		_, err = conn.Write([]byte(cmd))

		if err != nil {
			log.Printf("error: %v\n", err)
			return
		}

		err = conn.Close()
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}
	wg.Wait()
}
