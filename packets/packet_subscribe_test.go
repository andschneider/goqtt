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
}
