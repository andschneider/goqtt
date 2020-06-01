package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestConnectPacket(t *testing.T) {
	var buf bytes.Buffer
	cp := CreateConnectPacket()
	err := cp.Write(&buf)
	if err != nil {
		t.Errorf("could not write CONNECT packet: %v", err)
	}
	fmt.Printf("connect packet: %s\n", &cp)

	packetType, err := decodeByte(&buf)
	if err != nil {
		t.Errorf("could not decode type from fixed header. got %v", packetType)
	}

	p := ConnectPacket{}
	err = p.Read(&buf)
	if err != nil {
		t.Errorf("could not read %s packet: %v", connectType.name, err)
	}
	fmt.Println(p)
}
