package daemon

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const producerSock = "ipc-pubsub-producer.sock"

type Consumer struct {
	con *net.UnixAddr
}

type Daemon struct {
	sync.Mutex
	consumers    []*conn
	produceSock  net.Conn
	consumerSock net.Conn
}

func NewDaemon() *Daemon {
	return &Daemon{}
}

type conn struct {
	mu   sync.Mutex
	conn net.Conn

	readTimeout time.Duration
	br          *bufio.Reader
	bw          *bufio.Writer
}

var ErrInvalidProtocol = errors.New("invalid commands %v")

func (c *conn) ReadCommand() (*Command, []byte, error) {
	p := NewParser()
	for i := 0; i < 3; i++ {
		data, _, err := c.br.ReadLine()
		if err != nil {
			return nil, nil, ErrInvalidProtocol
		}
		p.Add(data)
		if p.command.Name != PUBLISH {
			break
		}
	}
	if !p.IsValid() {
		var err error
		// todo fix

		for _, e := range p.errors {
			err = errors.Unwrap(e)
		}
		return nil, nil, err
	}

	return p.command, p.rowData, nil
}

func NewConn(c net.Conn) *conn {
	return &conn{
		mu:   sync.Mutex{},
		conn: c,
		bw:   bufio.NewWriter(c),
		br:   bufio.NewReader(c),
	}
}

func (d *Daemon) serviceProducer(conn_ net.Conn, ch chan []byte) {

	c := NewConn(conn_)
	command, raw, err := c.ReadCommand()
	if err != nil {
		log.Println(err)
	}

	if command.Name == SUBSCRIBE {
		d.Lock()
		d.consumers = append(d.consumers, c)
		d.Unlock()
	}

	if command != nil {
		fmt.Println(command)
		ch <- raw
	}
}

func (d *Daemon) serviceConsumer(fd net.Conn, ch chan []byte) {

	for command := range ch {
		d.Lock()
		for _, consumer := range d.consumers {
			fmt.Printf("sender: %v\n", command)
			//bytes.Buffer.
			//err := binary.Write(consumer.conn, binary.LittleEndian, command)
			_, err := consumer.conn.Write(append(command, []byte("\r\n")...))
			if err != nil {
				log.Println(err)
			}
		}
		d.Unlock()
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

	ch := make(chan []byte, 10)

	for {
		pfd, err := pcon.Accept()

		if err != nil {
			return nil
		}
		go d.serviceProducer(pfd, ch)
		go d.serviceConsumer(nil, ch)
	}

}
