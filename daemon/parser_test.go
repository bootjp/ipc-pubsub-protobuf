package daemon

import (
	"fmt"
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
	p := Parser{}

	testData := []struct {
		send   []byte
		want   *Command
		hasErr bool
	}{
		{
			[]byte("PUBLISH 3 12\r\nMES\r\nPROTO_DATA\r\n"),
			&Command{
				PUBLISH,
				"MES",
				[]byte("PROTO_DATA\r\n"),
			},
			false,
		},
	}

	for _, td := range testData {
		got, err := p.parse(td.send)
		fmt.Println(got)
		fmt.Println(td.want)
		fmt.Println(err)
		fmt.Println("--")

		if reflect.DeepEqual(got, td.want) || (err != nil) != td.hasErr {
			t.Errorf("channel name miss match expected %v got %v hasError %v", td.want, got, err != nil)
		}
		if got, err := p.parse(td.send); reflect.DeepEqual(got, *td.want) || (err != nil) != td.hasErr {
			t.Errorf("channel name miss match expected %v got %v hasError %v", td.want, got, err != nil)
		}
	}
}
