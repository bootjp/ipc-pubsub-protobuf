package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	_ "github.com/golang/protobuf/jsonpb"
)

const producerSock = "ipc-pubsub-producer.sock"

func TestName(t *testing.T) {

	addr, err := net.ResolveUnixAddr("unix", filepath.Join(os.TempDir(), producerSock))
	if err != nil {
		t.Errorf("%v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
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

		fmt.Println("waiting data")

		buf := make([]byte, 100)
		_, err = conn.Read(buf)

		if err != nil {
			t.Error(err)
		}
		fmt.Println("receive")
		fmt.Printf("%v\n", buf)
		fmt.Printf("%s\n", buf)

		buf = make([]byte, 100)
		_, err = conn.Read(buf)

		if err != nil {
			t.Error(err)
		}
		fmt.Println("receive data")
		fmt.Printf("%v\n", buf)
		fmt.Printf("%s\n", buf)

		err = conn.Close()
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()
	go func() {
		time.Sleep(1 * time.Second)
		conn, err := net.DialUnix("unix", nil, addr)
		if err != nil {
			t.Error(err)
		}
		cmd := fmt.Sprintf("+PUBLISH PROTOBUF\r\n$10\r\nPROTO_DATA\r\n")
		fmt.Println("write")
		fmt.Printf("%v\n", []byte(cmd))
		fmt.Printf("%s\n", []byte(cmd))

		_, err = conn.Write([]byte(cmd))

		if err != nil {
			log.Printf("error: %v\n", err)
			return
		}

		err = conn.CloseWrite()
		if err != nil {
			log.Println(err)
		}
		wg.Done()
		fmt.Println("sender done")
	}()
	wg.Wait()
}
