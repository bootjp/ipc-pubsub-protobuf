package daemon

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

const producerSock = "ipc-pubsub-producer.sock"

type Consumer struct {
	conn *IPCConn
}

type Subscriber struct {
	Channel map[string][]*IPCConn
}

type Daemon struct {
	sync.Mutex
	consumers    *Subscriber
	produceSock  net.Conn
	consumerSock net.Conn
	sockPath     string
}

func NewDaemon() *Daemon {
	return &Daemon{
		consumers: &Subscriber{map[string][]*IPCConn{}},
	}
}

type IPCConn struct {
	mu   sync.Mutex
	conn net.Conn

	readTimeout time.Duration
	Br          *bufio.Reader
	Bw          *bufio.Writer
}

const maxPhase = 3

func (c *IPCConn) ReadCommand() (*Command, error) {
	p := NewParser()

	for i := 0; i < maxPhase; i++ {
		data, _, err := c.Br.ReadLine()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		p.Add(data)
		if p.Command.Name != PUBLISH {
			break
		}
	}
	if !p.IsValid() {
		var err error
		for _, e := range p.errors {
			err = fmt.Errorf("any func: %w", e)
		}
		return nil, err
	}

	return p.Command, nil
}

func NewConn(c net.Conn) *IPCConn {
	return &IPCConn{
		mu:   sync.Mutex{},
		conn: c,
		Bw:   bufio.NewWriter(c),
		Br:   bufio.NewReader(c),
	}
}

func (d *Daemon) serviceProducer(conn_ net.Conn, ch chan *Command) {

	c := NewConn(conn_)
	command, err := c.ReadCommand()
	if err != nil {
		log.Println(err)
	}
	if command == nil {
		log.Printf("command is nil")
		return
	}

	switch command.Name {
	case SUBSCRIBE:
		d.Lock()
		for _, channel := range command.Channel {
			_, ok := d.consumers.Channel[channel]
			if ok {
				d.consumers.Channel[channel] = append(d.consumers.Channel[channel], c)
			} else {
				d.consumers.Channel[channel] = []*IPCConn{c}
			}
		}
		d.Unlock()
	case UNSUBSCRIBE:
		d.Lock()
		for _, channel := range command.Channel {
			chs, ok := d.consumers.Channel[channel]
			if !ok {
				log.Printf("not found subscribe channel for unsubscribe %s\n", command.Channel)
				continue
			}
			for _, conn := range chs {
				if conn != c {
					d.consumers.Channel[channel] = append(d.consumers.Channel[channel], c)
				}
			}
		}
		d.Unlock()
	}
	ch <- command
}

func (d *Daemon) serviceConsumer(ch chan *Command) {
	for command := range ch {
		d.Lock()

		for _, channels := range command.Channel {
			fmt.Println("channel " + channels)
			conns, ok := d.consumers.Channel[channels]
			if !ok {
				continue
			}

			for _, conn := range conns {
				_, err := conn.conn.Write(command.GetByte())
				if err != nil {
					log.Println(err)
				}
			}
		}

		d.Unlock()
	}
}

func (d *Daemon) Start() error {

	d.sockPath = filepath.Join(os.TempDir(), producerSock)
	listener, err := net.Listen("unix", d.sockPath)

	if err != nil {
		_ = os.Remove(d.sockPath)
		return err
	}

	defer func() {
		fmt.Println("deleting sockets" + d.sockPath)
		err = os.Remove(d.sockPath)
		if err != nil {
			fmt.Println(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go d.signalHandler(sig)
	ch := make(chan *Command, 1)

	for {
		pfd, err := listener.Accept()

		if err != nil {
			_ = pfd.Close()
		} else {
			go d.serviceProducer(pfd, ch)
			go d.serviceConsumer(ch)
		}
	}
}

func (d *Daemon) signalHandler(c chan os.Signal) {
	sig := <-c
	log.Printf("Caught signal %s: shutting down.", sig)
	err := os.Remove(d.sockPath)
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(0)
}
