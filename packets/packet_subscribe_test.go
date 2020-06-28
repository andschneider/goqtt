package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSubscribePacket(t *testing.T) {
	var buf bytes.Buffer
	var spWrite, spRead SubscribePacket

	spWrite.CreatePacket(testTopic)
	err := spWrite.Write(&buf)
	if err != nil {
		t.Errorf("could not write %s packet: %v", spWrite.name, err)
	}
	fmt.Printf("write %s packet: %+v\n", spWrite.name, spWrite)

	packetType, err := decodeByte(&buf)
	if err != nil {
		t.Errorf("could not decode type from fixed header. got %v", packetType)
	}

	err = spRead.Read(&buf)
	if err != nil {
		t.Errorf("could not read %s packet: %v", spRead.name, err)
	}
	fmt.Printf("read %s packet: %+v\n", spRead.name, spRead)
}
