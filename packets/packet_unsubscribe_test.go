package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestUnsubscribePacket(t *testing.T) {
	var buf bytes.Buffer
	up := CreateUnsubscribePacket(testTopic)
	err := up.Write(&buf)
	if err != nil {
		t.Errorf("could not write %s packet: %v", unsubscribeType.name, err)
	}
	fmt.Printf("unsubscribe packet: %s\n", up.String())

	packetType, err := decodeByte(&buf)
	if err != nil {
		t.Errorf("could not decode type from fixed header. got %v", packetType)
	}

	p := UnsubscribePacket{}
	err = p.Read(&buf)
	if err != nil {
		t.Errorf("could not read %s packet: %v", unsubscribeType.name, err)
	}
	fmt.Println("read that packet!", p)
}
