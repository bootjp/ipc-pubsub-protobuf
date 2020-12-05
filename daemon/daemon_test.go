package daemon

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestData(t *testing.T) {
	addr, err := net.ResolveUnixAddr("unix", filepath.Join(os.TempDir(), producerSock))
	if err != nil {
		t.Errorf("%v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	testCmdData := "+PUBLISH PROTOBUF\r\n$10\r\nPROTO_DATA\r\n"
	testCmdSubscribe := "+SUBSCRIBE PROTOBUF\r\n"

	go func() {
		conn, err := net.DialUnix("unix", nil, addr)
		if err != nil {
			t.Error(err)
		}

		{
			_, err = conn.Write([]byte(testCmdSubscribe))
			if err != nil {
				t.Error(err)
			}

			buf := make([]byte, len(testCmdSubscribe))
			_, err = conn.Read(buf)

			if err != nil {
				t.Error(err)
			}

			if testCmdSubscribe != string(buf) {
				t.Errorf("commands response miss match expected %v, got %v", testCmdSubscribe, buf)
			}
			t.Logf("%s", buf)
		}
		{
			buf := make([]byte, len(testCmdData))
			_, err = conn.Read(buf)

			if err != nil {
				t.Error(err)
			}

			err = conn.Close()
			if err != nil {
				log.Println(err)
			}

			if testCmdData != string(buf) {
				t.Errorf("commands response miss match expected %v, got %v", testCmdSubscribe, buf)
			}
			t.Logf("%s", buf)
			wg.Done()
		}
	}()
	go func() {
		time.Sleep(1 * time.Second)
		conn, err := net.DialUnix("unix", nil, addr)
		if err != nil {
			t.Error(err)
		}

		_, err = conn.Write([]byte(testCmdData))
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
