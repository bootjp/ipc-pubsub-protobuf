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
			[]byte("PLAYER_JOIN"),
			len("PLAYER_JOIN"),
			"PLAYER_JOIN",
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
			"",
			true,
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
				[]string{"MES"},
				[]byte("PROTO_DATA"),
			},
			true,
			false,
		},
		{
			[]byte("PUBLISH 3 11\r\nddddd\r\nPROTO_DATA\r\n"),
			&Command{
				PUBLISH,
				[]string{"ddddd"},
				[]byte("PROTO_DATA"),
			},
			false,
			true,
		},
		{
			[]byte("PUBLISH 4 11\r\nddddd\r\nPROTO_DATA\r\n"),
			&Command{
				PUBLISH,
				[]string{"dddd"},
				[]byte("PROTO_DATA"),
			},
			false,
			true,
		},

		{
			[]byte("SUBSCRIBE 5\r\nddddd\r\n"),
			&Command{
				SUBSCRIBE,
				[]string{"ddddd"},
				nil,
			},
			true,
			false,
		},
		{
			[]byte("UNSUBSCRIBE 4\r\nxxxx\r\n"),
			&Command{
				UNSUBSCRIBE,
				[]string{"xxxx"},
				nil,
			},
			true,
			false,
		},
		//// multiple channel test
		//{
		//	[]byte("SUBSCRIBE 4,4\r\nxxxx xxxx\r\n"),
		//	&Command{
		//		UNSUBSCRIBE,
		//		[]string{"xxxx xxxx"},
		//		nil,
		//	},
		//	true,
		//	false,
		//},
	}
	for _, td := range testData {
		p := NewParser()
		bb := bytes.NewBuffer(td.send)

		scanner := bufio.NewScanner(bb)
		scanner.Split(bufio.ScanLines)

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
