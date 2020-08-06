package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSubscribePacket(t *testing.T) {
	var buf bytes.Buffer
	var spWrite SubscribePacket

	spWrite.CreateSubscribePacket(testTopic)
	err := spWrite.Write(&buf)
	if err != nil {
		t.Errorf("could not write %s packet: %v", spWrite.name, err)
	}
	fmt.Printf("write %s packet: %+v\n", spWrite.name, spWrite)
}
