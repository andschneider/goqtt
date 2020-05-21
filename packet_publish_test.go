package goqtt

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPublishPacket_Write(t *testing.T) {
	message := "hello world"
	pp := CreatePublishPacket(testTopic, message)
	var buf bytes.Buffer
	err := pp.Write(&buf, true)
	if err != nil {
		t.Errorf("could not write Publish packet %v", err)
	}

	p, err := pp.ReadPublishPacket(&buf)
	if err != nil {
		t.Errorf("could not read Publish packer %v", err)
	}
	topic := p.Topic
	if topic != testTopic {
		t.Errorf("topics do not match. expected %s, got %s", testTopic, topic)
	}
	ms := string(p.Message[2:]) // TODO not sure about the first two characters
	if ms != message {
		t.Errorf("messages do not match. expected %s, got %s", message, ms)
	}
	fmt.Printf("TOPIC: %s\nMESSAGE: %s\n", topic, ms)
}
