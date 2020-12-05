package daemon

import (
	"errors"
	"strconv"
	"strings"
)

type Parser struct {
	dataStage     int
	hasError      bool
	errors        []error
	Command       *Command
	channelLength int
	dataLength    int
	rowData       []byte
}

type CommandName int

const (
	SUBSCRIBE CommandName = iota + 10
	UNSUBSCRIBE
	PUBLISH
	UnknownCommand = 100
)

var CommandMaps = map[CommandName]string{
	SUBSCRIBE:   "SUBSCRIBE",
	UNSUBSCRIBE: "UNSUBSCRIBE",
	PUBLISH:     "PUBLISH",
}

type Command struct {
	Name    CommandName
	Channel []string
	Data    []byte
}

func (c *Command) GetByte() []byte {
	name := CommandMaps[c.Name]
	s := []byte(name + " " + strings.Join(c.Channel, " "))
	return append(s, c.Data...)
}

var ErrCommandLengthTooMin = errors.New("command name is too min")
var ErrInvalidLength = errors.New("invalid specify length")

func NewParser() *Parser {
	return &Parser{
		Command: &Command{},
	}
}

func (p *Parser) Add(b []byte) {
	p.dataStage++

	if len(b) == 0 {
		p.errors = append(p.errors, ErrCommandLengthTooMin)
		return
	}

	s := string(b)
	first := string(s[0])

	switch {
	case first == "+" && p.dataStage == 1:
		commands := strings.Fields(s)
		if len(commands) < 2 {
			p.errors = append(p.errors, ErrCommandLengthTooMin)
			return
		}
		p.Command.Name = p.parseCommandName(commands[0][1:])
		p.Command.Channel = commands[1:]

	case first == "$" && p.dataStage == 1: // bulk string
		n, err := strconv.Atoi(s)
		if n < 0 && err != nil {
			p.errors = append(p.errors, ErrInvalidLength)
			return
		}
		p.Command.Data = make([]byte, n)
	case p.dataStage == 3:
		p.Command.Data = b
	}
}

func (p *Parser) parseCommandName(b string) CommandName {
	switch strings.ToUpper(b) {
	case "SUBSCRIBE":
		return SUBSCRIBE
	case "UNSUBSCRIBE":
		return UNSUBSCRIBE
	case "PUBLISH":
		return PUBLISH
	}

	return UnknownCommand
}

func (p *Parser) IsValid() bool {
	if len(p.errors) != 0 {
		return false
	}

	if p.Command.Name == 0 {
		return false
	}

	if len(p.Command.Channel) == 0 {
		return false
	}

	if p.Command.Name == PUBLISH && p.Command.Data == nil {
		return false
	}

	return true
}

func (p *Parser) GetError() []error {
	return p.errors
}

func (p *Parser) GetCommand() *Command {
	return p.Command
}

func (p *Parser) GetRowData() []byte {
	return p.Command.GetByte()
}
