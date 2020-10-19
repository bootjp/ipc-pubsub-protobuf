package daemon

import (
	"bufio"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	client "github.com/bootjp/ipc-pubsub-protobuf/build"
)

const producerSock = "ipc-pubsub-producer.sock"
const consumerSock = "ipc-pubsub-consumer.sock"

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

func (d *Daemon) serviceProducer(fd net.Conn, ch chan client.Message) {
	bufReader := bufio.NewReader(fd)
	scanner := bufio.NewScanner(bufReader)
	for scanner.Scan() {
		b := scanner.Bytes()
		if len(b) == 0 {
			continue
		}

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

func (d *Daemon) serviceConsumer(fd net.Conn, ch chan client.Message) {

	for v := range ch {
		for _, c := range d.consumers {
			conn, err := net.DialUnix("unix", nil, c.con)
			if err != nil {
				log.Println(err)
			}

			_, err = conn.Write([]byte(v.ProtoReflect()))
			if err != nil {
				log.Println(err)
			}
			conn.Close()
		}
	}

}

func (d *Daemon) Start() error {

	pSockPath := filepath.Join(os.TempDir(), producerSock)
	cSockPath := filepath.Join(os.TempDir(), consumerSock)

	pcon, err := net.Listen("unix", pSockPath)
	if err != nil {
		return err
	}

	ccon, err := net.Listen("unix", cSockPath)
	if err != nil {
		return err
	}

	// todo 個別に
	defer func() {
		_ = os.Remove(filepath.Join(os.TempDir(), producerSock))
		_ = os.Remove(filepath.Join(os.TempDir(), consumerSock))
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	ch := make(chan client.Message)
	for {
		pfd, err := pcon.Accept()
		if err != nil {
			return nil
		}
		go d.serviceProducer(pfd, ch)

		cfd, err := ccon.Accept()
		if err != nil {
			return nil
		}
		go d.serviceProducer(cfd, ch)
	}
}
