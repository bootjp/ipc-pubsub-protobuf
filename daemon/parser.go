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
	SUBSCRIBE CommandName = iota
	UNSUBSCRIBE
	PUBLISH
	UnknownCommand = 100
)

type Command struct {
	Name    CommandName
	Channel string
	Data    []byte
}

var ErrCommandLengthTooBig = errors.New("command name is too big")
var ErrInvalidCommandNameArgs = errors.New("command name args is missing")
var ErrInvalidLength = errors.New("invalid specify length")

func NewParser() *Parser {
	return &Parser{
		command: &Command{},
	}
}

func (p *Parser) Add(b []byte) {
	p.dataStage++

	var err error
	switch p.dataStage {
	case 1:
		commands := strings.Fields(string(b))

		// rejected invalid big command.
		if len(commands[0]) > 30 {
			p.errors = append(p.errors, ErrCommandLengthTooBig)
		}

		// required command CHANNEL length and DATA length.
		if len(commands) < 2 {
			p.errors = append(p.errors, ErrInvalidCommandNameArgs)
		}

		p.channelLength, err = strconv.Atoi(commands[1])
		if err != nil {
			p.errors = append(p.errors, err)
		}

		if len(commands) == 3 {
			p.dataLength, err = strconv.Atoi(commands[2])
			if err != nil {
				p.errors = append(p.errors, err)
			}
		}
		p.command.Name = p.parseCommandName(commands[0])

	case 2:
		p.command.Channel, err = p.parseChannelName(b, p.channelLength)
		if err != nil {
			p.errors = append(p.errors, err)
		}

	case 3:

		if len(b) != p.dataLength {
			p.errors = append(p.errors, ErrInvalidLength)
		}
		p.command.Data = b[:p.dataLength]
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

var ErrInvalidChannel = errors.New("invalid channel name length")

func (p *Parser) parseChannelName(b []byte, length int) (string, error) {
	if len(b) != length {
		return "", ErrInvalidChannel
	}

	return string(b[:length]), nil
}

func (p *Parser) IsValid() bool {
	return len(p.errors) == 0
}

func (p *Parser) GetError() []error {
	return p.errors
}

func (p *Parser) GetCommand() *Command {
	return p.command
}
