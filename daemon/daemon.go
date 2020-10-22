package daemon

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"

	client "github.com/bootjp/ipc-pubsub-protobuf/build"
)

const producerSock = "ipc-pubsub-producer.sock"

type Consumer struct {
	con *net.UnixAddr
}

type Daemon struct {
	sync.Mutex
	consumers    []Consumer
	produceSock  net.Conn
	consumerSock net.Conn
}

func NewDaemon() *Daemon {
	return &Daemon{}
}

func (d *Daemon) serviceProducer(fd net.Conn, ch chan client.MessageContainer) {
	bufReader := bufio.NewReader(fd)
	scanner := bufio.NewScanner(bufReader)

	for scanner.Scan() {
		b := scanner.Bytes()
		if len(b) == 0 {
			continue
		}
		fmt.Println(string(b))
	}
}

func (d *Daemon) registerConsumer(path string) {
	addr, err := net.ResolveUnixAddr("unix", path)
	if err != nil {
	}
	d.Lock()
	defer d.Unlock()

	d.consumers = append(d.consumers, Consumer{addr})
}

func (d *Daemon) serviceConsumer(fd net.Conn, ch chan client.MessageContainer) {

	for _ = range ch {
		for _, c := range d.consumers {
			conn, err := net.DialUnix("unix", nil, c.con)
			if err != nil {
				log.Println(err)
				continue
			}

			err = conn.Close()
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (d *Daemon) Start() error {

	pSockPath := filepath.Join(os.TempDir(), producerSock)
	//cSockPath := filepath.Join(os.TempDir(), consumerSock)

	pcon, err := net.Listen("unix", pSockPath)
	if err != nil {
		return err
	}
	defer func() {
		err = os.Remove(pSockPath)
		if err != nil {
			fmt.Println(err)
		}

	}()

	//ccon, err := net.Listen("unix", cSockPath)
	//if err != nil {
	//	return err
	//}
	//defer func() {
	//	err = os.Remove(cSockPath)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//
	//}()

	//sigc := make(chan os.Signal, 1)
	//signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	// todo handling signal

	wg := &sync.WaitGroup{}
	ch := make(chan client.MessageContainer)
	pfd, err := pcon.Accept()

	if err != nil {
		return nil
	}
	go d.serviceProducer(pfd, ch)
	//
	//cfd, err := ccon.Accept()
	//if err != nil {
	//	return nil
	//}
	//go d.serviceProducer(cfd, ch)

	wg.Add(1)
	wg.Wait()
	return nil
}
