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
			t.Errorf("Command miss match expected publish got %v", got)
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
			[]byte("+PUBLISH MES\r\n$10\r\nPROTO_DATA\r\n"),
			&Command{
				PUBLISH,
				[]string{"MES"},
				[]byte("PROTO_DATA"),
			},
			true,
			false,
		},

		{
			[]byte("+SUBSCRIBE ddddd\r\n"),
			&Command{
				SUBSCRIBE,
				[]string{"ddddd"},
				nil,
			},
			true,
			false,
		},
		{
			[]byte("+UNSUBSCRIBE xxxx\r\n"),
			&Command{
				UNSUBSCRIBE,
				[]string{"xxxx"},
				nil,
			},
			true,
			false,
		},
		// multiple channel test
		{
			[]byte("+SUBSCRIBE xxxx xxxx\r\n"),
			&Command{
				SUBSCRIBE,
				[]string{"xxxx", "xxxx"},
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
