package goqtt_test

import (
	"testing"

	"github.com/andschneider/goqtt"
)

var (
	clientID  = "testClient"
	keepAlive = 10
	broker    = ""
	topic     = "test"
)

func TestNewClientConfig(t *testing.T) {
	cc := goqtt.NewClientConfig(clientID, keepAlive, broker, topic)
	c := goqtt.NewClient(cc)

	if c.GetClientId() != clientID {
		t.Errorf("clientIDs don't match. got %s, expected %s", c.GetClientId(), clientID)
	}
	if c.GetKeepAlive() != keepAlive {
		t.Errorf("keepAlives don't match. got %d, expected %d", c.GetKeepAlive(), keepAlive)
	}
	if c.GetBroker() != broker {
		t.Errorf("brokers don't match. got %s, expected %s", c.GetBroker(), broker)
	}
	if c.GetTopic() != topic {
		t.Errorf("topics don't match. got %s, expected %s", c.GetTopic(), topic)
	}
}
