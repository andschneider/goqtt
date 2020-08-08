package goqtt_test

import (
	"fmt"
	"testing"

	"github.com/andschneider/goqtt"
)

const (
	clientID  = "goqtt-testclient"
	keepAlive = 10
	server    = "mqtt.eclipse.org"
	port      = "1883"
	topic     = "goqtt-test"
)

func TestNewClientConfig_Default(t *testing.T) {
	c := goqtt.NewClient(server)

	err := c.Connect()
	if err != nil {
		t.Fatal(err)
	}
}

func TestOptions(t *testing.T) {
	cid := goqtt.ClientId(clientID)
	ka := goqtt.KeepAlive(keepAlive)
	port := goqtt.Port(port)
	topic := goqtt.Topic(topic)

	c := goqtt.NewClient("test", cid, port, topic, ka)

	fmt.Printf("%+v\n", c.Config)
}
