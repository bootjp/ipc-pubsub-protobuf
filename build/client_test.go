package client_test

import (
	"os"
	"testing"

	client "github.com/bootjp/ipc-pubsub-protobuf/build"

	"github.com/golang/protobuf/jsonpb"
)

func TestName(t *testing.T) {
	container := client.MessageContainer{}
	m := jsonpb.Marshaler{
		EmitDefaults: true,
		Indent:       "    ",
		OrigName:     true,
	}
	m.Marshal(os.Stdout, &container)
}
