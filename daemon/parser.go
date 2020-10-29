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
	command       *Command
	channelLength int
	dataLength    int
}

type CommandName int

const (
	SUBSCRIBE CommandName = iota + 10
	UNSUBSCRIBE
	PUBLISH
	UnknownCommand = 100
)

type Command struct {
	Name    CommandName
	Channel []string
	Data    []byte
}

var ErrCommandLengthTooMin = errors.New("command name is too min")
var ErrInvalidLength = errors.New("invalid specify length")

func NewParser() *Parser {
	return &Parser{
		command: &Command{},
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
		p.command.Name = p.parseCommandName(commands[0][1:])
		p.command.Channel = commands[1:]

	case first == "$" && p.dataStage == 1: // bulk string
		n, err := strconv.Atoi(s)
		if n < 0 && err != nil {
			p.errors = append(p.errors, ErrInvalidLength)
			return
		}
		p.command.Data = make([]byte, n)
	case p.dataStage == 3:
		p.command.Data = b
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

	if p.command.Name == 0 {
		return false
	}

	if len(p.command.Channel) == 0 {
		return false
	}

	if p.command.Name == PUBLISH && p.command.Data == nil {
		return false
	}

	return true
}

func (p *Parser) GetError() []error {
	return p.errors
}

func (p *Parser) GetCommand() *Command {
	return p.command
}
