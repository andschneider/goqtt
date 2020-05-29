package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSubscribePacket(t *testing.T) {
	var buf bytes.Buffer
	sp := CreateSubscribePacket(testTopic)
	err := sp.Write(&buf)
	if err != nil {
		t.Errorf("could not write SUBSCRIBE packet: %v", err)
	}
	fmt.Printf("subscribe packet: %s\n", &sp)

	packetType, err := decodeByte(&buf)
	if err != nil {
		t.Errorf("could not decode type from fixed header. got %v", packetType)
	}

	p := SubscribePacket{}
	err = p.Read(&buf)
	if err != nil {
		t.Errorf("could not read %s packet: %v", subscribeType.name, err)
	}
	fmt.Println("read that packet!", p)
}
