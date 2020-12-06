package client_test

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	client "github.com/bootjp/ipc-pubsub-protobuf/build/proto"

	"github.com/golang/protobuf/proto"
)

const producerSock = "ipc-pubsub-producer.sock"

func TestName(t *testing.T) {
	container := client.MessageContainer{
		Dummy: "aaa",
	}

	data, err := proto.Marshal(&container)
	if err != nil {
		log.Fatal(err)
	}

	c := &client.MessageContainer{}
	err = proto.Unmarshal(data, c)
	if err != nil {
		log.Fatal(err)
	}

	addr, err := net.ResolveUnixAddr("unix", filepath.Join(os.TempDir(), producerSock))
	if err != nil {
		t.Errorf("%v", err)
	}

	testCmdData := "+PUBLISH %s\r\n$%d\r\n"
	testCmdSubscribe := "+SUBSCRIBE %s\r\n"

	dataB := append([]byte(fmt.Sprintf(testCmdData, reflect.TypeOf(container).Name(), len(data))), append(data, []byte("\r\n")...)...)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		conn, err := net.DialUnix("unix", nil, addr)
		if err != nil {
			t.Error(err)
		}

		{
			cmd := fmt.Sprintf(testCmdSubscribe, reflect.TypeOf(container).Name())

			_, err = conn.Write([]byte(cmd))
			if err != nil {
				t.Error(err)
			}

			buf := make([]byte, len(cmd))
			_, err = conn.Read(buf)

			if err != nil {
				t.Error(err)
			}

			if cmd != string(buf) {
				t.Errorf("commands response miss match expected %v, got %v", testCmdSubscribe, buf)
			}
			t.Logf("%s", buf)
		}
		{
			buf := make([]byte, len(dataB))
			fmt.Println("waiting data")
			_, err = conn.Read(buf)
			fmt.Println("done")

			if err != nil {
				t.Error(err)
			}

			err = conn.Close()
			if err != nil {
				log.Println(err)
			}

			if !reflect.DeepEqual(dataB, buf) {
				t.Errorf("commands response miss match expected %v, got %v", testCmdSubscribe, buf)
			}
			t.Logf("%s", buf)

			ct := &client.MessageContainer{}
			err = proto.Unmarshal(buf, ct)

			if err != nil {
				log.Println(err)
			}

			if !reflect.DeepEqual(ct, container) {
				t.Error("miss match")
			}

			wg.Done()
		}
	}()
	go func() {
		time.Sleep(1 * time.Second)

		conn, err := net.DialUnix("unix", nil, addr)
		if err != nil {
			t.Error(err)
		}

		fmt.Println("writing data")

		_, err = conn.Write(dataB)
		if err != nil {
			t.Error(err)
			return
		}

		err = conn.CloseWrite()
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()
	wg.Wait()
}
