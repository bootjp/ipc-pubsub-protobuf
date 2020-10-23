package daemon

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestCommandName(t *testing.T) {

	p := Parser{}

	testData := []struct {
		send string
		want CommandName
	}{
		{
			"publish",
			PUBLISH,
		},
		{
			"Subscribe",
			SUBSCRIBE,
		},
		{
			"unsubscribe",
			UNSUBSCRIBE,
		},
	}

	for _, td := range testData {
		if got := p.parseCommandName(td.send); got != td.want {
			t.Errorf("command miss match expected publish got %v", got)
		}
	}
}

func TestChannelName(t *testing.T) {

	p := Parser{}

	testData := []struct {
		send   []byte
		length int
		want   string
		hasErr bool
	}{
		{
			[]byte("Message"),
			len("Message"),
			"Message",
			false,
		},
		{
			[]byte("PLAYER JOIN"),
			len("PLAYER JOIN"),
			"PLAYER JOIN",
			false,
		},
		{
			[]byte("Message"),
			len("Message"),
			"Message",
			false,
		},
		{
			[]byte("Message"),
			len("Messag"),
			"Messag",
			false,
		},
		{
			[]byte("Messag"),
			len("Messagaaa"),
			"",
			true,
		},
	}

	for _, td := range testData {
		if got, err := p.parseChannelName(td.send, td.length); got != td.want || (err != nil) != td.hasErr {
			t.Errorf("channel name miss match expected %v got %v hasError %v", td.want, got, err != nil)
		}
	}
}

func TestParse(t *testing.T) {

	testData := []struct {
		send   []byte
		want   *Command
		valid  bool
		hasErr bool
	}{
		{
			[]byte("PUBLISH 3 10\r\nMES\r\nPROTO_DATA\r\n"),
			&Command{
				PUBLISH,
				"MES",
				[]byte("PROTO_DATA"),
			},
			true,
			false,
		},
		{
			[]byte("PUBLISH 3 11\r\nddddd\r\nPROTO_DATA\r\n"),
			&Command{
				PUBLISH,
				"ddddd",
				[]byte("PROTO_DATA"),
			},
			false,
			true,
		},
		{
			[]byte("PUBLISH 4 11\r\nddddd\r\nPROTO_DATA\r\n"),
			&Command{
				PUBLISH,
				"dddd",
				[]byte("PROTO_DATA"),
			},
			false,
			true,
		},

		{
			[]byte("SUBSCRIBE 5\r\nddddd\r\n"),
			&Command{
				SUBSCRIBE,
				"ddddd",
				nil,
			},
			true,
			false,
		},
		{
			[]byte("UNSUBSCRIBE 4\r\nxxxx\r\n"),
			&Command{
				UNSUBSCRIBE,
				"xxxx",
				nil,
			},
			true,
			false,
		},
	}
	for _, td := range testData {
		p := NewParser()
		bb := bytes.NewBuffer(td.send)

		scanner := bufio.NewScanner(bb)
		scanner.Split(ScanCRLF)

		for scanner.Scan() {
			b := scanner.Bytes()
			p.Add(b)
		}

		if p.IsValid() != td.valid {
			t.Errorf("invalid valid condition expect %v got %v", td.valid, p.IsValid())
			for _, e := range p.GetError() {
				t.Error(e)
			}
		}

		cond := reflect.DeepEqual(*p.GetCommand(), *td.want)
		if !cond && td.valid {
			t.Errorf("type miss match expect %v got %v", p.GetCommand(), td.want)
		}

	}
}

func ScanCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{'\r', '\n'}); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
