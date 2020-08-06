package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPublishPacket_Write(t *testing.T) {
	var buf bytes.Buffer
	var pp PublishPacket

	message := "hello world"
	pp.CreatePublishPacket(testTopic, message)

	err := pp.Write(&buf)
	if err != nil {
		t.Errorf("could not write %s packet %v", pp.name, err)
	}
	fmt.Printf("%s packet: %+v\n", pp.name, pp)

	packetType, err := decodeByte(&buf)
	if err != nil {
		t.Errorf("could not decode type from fixed header. got %v", packetType)
	}

	var ppRead PublishPacket
	err = ppRead.Read(&buf)
	if err != nil {
		t.Errorf("could not read %s packet %v", pp.name, err)
	}
	topic := ppRead.Topic
	if topic != testTopic {
		t.Errorf("topics do not match. expected %s, got %s", testTopic, topic)
	}
	ms := string(ppRead.Message)
	if ms != message {
		t.Errorf("messages do not match. expected %s, got %s", message, ms)
	}
	fmt.Printf("TOPIC: %s\nMESSAGE: %s\n", topic, ms)
}
