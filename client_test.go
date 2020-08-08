package goqtt

import (
	"testing"
)

const (
	clientID  = "goqtt-testclient"
	keepAlive = 10
	server    = "mqtt.eclipse.org"
	port      = "1883"
	topic     = "goqtt-test"
)

func TestNewClientConfig_Default(t *testing.T) {
	c := NewClient(server)

	if c.Config.topic != DefaultTopic {
		t.Fatalf("default topic error. got %s, want %s", c.Config.topic, DefaultTopic)
	}
	if c.Config.port != DefaultPort {
		t.Fatalf("default port error. got %s, want %s", c.Config.port, DefaultPort)
	}
	if c.Config.keepAlive != DefaultKeepAlive {
		t.Fatalf("default keep alive error. got %d, want %d", c.Config.keepAlive, DefaultKeepAlive)
	}
}

func TestOptions(t *testing.T) {
	// set config options
	cid := ClientId(clientID)
	ka := KeepAlive(keepAlive)
	p := Port(port)
	tp := Topic(topic)

	c := NewClient(server, cid, p, tp, ka)

	if c.Config.clientId != clientID {
		t.Fatalf("keep alive. got %s, want %s", c.Config.clientId, clientID)
	}
	if c.Config.keepAlive != keepAlive {
		t.Fatalf("keep alive. got %d, want %d", c.Config.keepAlive, keepAlive)
	}
	if c.Config.port != port {
		t.Fatalf("port error. got %s, want %s", c.Config.port, port)
	}
	if c.Config.topic != topic {
		t.Fatalf("topic error. got %s, want %s", c.Config.topic, topic)
	}
	if c.Config.server != server {
		t.Fatalf("server error. got %s, want %s", c.Config.server, server)
	}
}
