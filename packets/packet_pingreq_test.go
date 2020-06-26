package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPingReqPacket(t *testing.T) {
	var buf bytes.Buffer
	var prWrite, prRead PingReqPacket

	prWrite.CreatePacket()
	err := prWrite.Write(&buf)
	if err != nil {
		t.Errorf("could not write %s packet %v", prWrite.name, err)
	}
	fmt.Printf("%s packet: %+v\n", prWrite.name, prWrite)

	packetType, err := decodeByte(&buf)
	if err != nil {
		t.Errorf("could not decode type from fixed header. got %v", packetType)
	}

	err = prRead.Read(&buf)
	if err != nil {
		t.Errorf("could not read %s packet: %v", prRead.name, err)
	}
	fmt.Printf("%s packet: %+v\n", prRead.name, prRead)
}
