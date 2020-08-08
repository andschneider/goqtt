package goqtt_test

import (
	"testing"

	"github.com/andschneider/goqtt"
)

var (
	clientID  = "goqtt-testclient"
	keepAlive = 10
	server    = "mqtt.eclipse.org"
	port      = "1883"
	topic     = "goqtt-test"
)

func TestNewClientConfig(t *testing.T) {
	cc := &goqtt.ClientConfig{
		ClientId:  clientID,
		KeepAlive: keepAlive,
		Server:    server,
		Port:      port,
		Topic:     topic,
	}
	c, err := goqtt.NewClient(cc)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Connect()
	if err != nil {
		t.Fatal(err)
	}
}
func TestNewClientConfig_Default(t *testing.T) {
	cc := &goqtt.ClientConfig{}
	c, err := goqtt.NewClient(cc)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Connect()
	if err != nil {
		t.Fatal(err)
	}
}
