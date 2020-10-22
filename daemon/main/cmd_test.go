package main

import (
	"fmt"

	"log"
	"net"
	"os"
	"path/filepath"
	"testing"

	pb "github.com/bootjp/ipc-pubsub-protobuf/build"
	_ "github.com/golang/protobuf/jsonpb"
	"google.golang.org/protobuf/proto"
)

const producerSock = "ipc-pubsub-producer.sock"
const consumerSock = "ipc-pubsub-consumer.sock"

var breakLine = make([]byte, len("\n"))

func TestName(t *testing.T) {

	addr, err := net.ResolveUnixAddr("unix", filepath.Join(os.TempDir(), producerSock))
	if err != nil {
		t.Errorf("%v", err)
	}

	conn, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		t.Error(err)
	}

	cmd := fmt.Sprintf("SEND PROTOBUF\n")
	_, err = conn.Write([]byte(cmd))

	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}

	m := pb.MessageContainer{}
	//m = m.ProtoReflect()
	mp := proto.MarshalOptions{Deterministic: true}
	b, err := mp.Marshal(&m)
	fmt.Println(b)
	if err != nil {
		t.Error(err)
	}

	_, err = conn.Write(append(b, breakLine...))
	if err != nil {
		t.Fatal(err)
	}
	err = conn.CloseWrite()
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}

	err = conn.Close()
	if err != nil {
		log.Println(err)
	}
}
