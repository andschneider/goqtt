package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestUnsubscribePacket(t *testing.T) {
	var buf bytes.Buffer
	var upRead, upWrite UnsubscribePacket

	upWrite.CreatePacket(testTopic)
	err := upWrite.Write(&buf)
	if err != nil {
		t.Errorf("could not write %s packet %v", upWrite.Name, err)
	}
	fmt.Printf("write %s packet: %+v\n", upWrite.Name, upWrite)

	packetType, err := decodeByte(&buf)
	if err != nil {
		t.Errorf("could not decode type from fixed header. got %v", packetType)
	}

	err = upRead.Read(&buf)
	if err != nil {
		t.Errorf("could not read %s packet: %v", upRead.Name, err)
	}
	fmt.Printf("read %s packet: %+v\n", upRead.Name, upRead)
}
