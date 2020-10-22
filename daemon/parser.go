package daemon

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
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

const commandMinLength = 10

var ErrCommandLengthTooMin = errors.New("command name is too min")
var ErrCommandLengthTooBig = errors.New("command name is too big")
var ErrInvalidCommandNameArgs = errors.New("command name args is missing")

func (p *Parser) parse(b ...[]byte) (*Command, error) {
	if len(b) > commandMinLength {
		return nil, ErrCommandLengthTooMin
	}

	c := &Command{}

	var channelLength, dataLength int
	var err error
	for i, d := range b {

		switch i + 1 {
		case 1:
			commands := strings.Fields(string(d))

			// rejected invalid big command.
			if i != 3 && len(commands[0]) > 30 {
				return nil, ErrCommandLengthTooBig
			}

			// required command CHANNEL length and DATA length.
			if len(commands) >= 2 {
				fmt.Println(commands)
				return nil, ErrInvalidCommandNameArgs
			}

			channelLength, err = strconv.Atoi(commands[1])
			if err != nil {
				return nil, err
			}

			if len(commands) == 3 {
				dataLength, err = strconv.Atoi(commands[2])
				if err != nil {
					return nil, err
				}
			}
			c.Name = p.parseCommandName(commands[0])
		case 2:
			c.Channel, err = p.parseChannelName(d, channelLength)
			if err != nil {
				return nil, err
			}
		case 3:
			if len(c.Data) >= dataLength {
				return nil, err
			}
			c.Data = d[:dataLength]
		}
	}

	return c, nil
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

var ErrInvalidChannel = errors.New("")

func (p *Parser) parseChannelName(b []byte, length int) (string, error) {
	if len(b) < length {
		return "", ErrInvalidChannel
	}

	return string(b[:length]), nil
}
