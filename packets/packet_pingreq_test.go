package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPingReqPacket(t *testing.T) {
	var buf bytes.Buffer
	pr := CreatePingReqPacket()
	err := pr.Write(&buf)
	if err != nil {
		t.Errorf("could not write PingReq packet %v", err)
	}
	fmt.Printf("pingreq packet: %s\n", &pr)

	packetType, err := decodeByte(&buf)
	if err != nil {
		t.Errorf("could not decode type from fixed header. got %v", packetType)
	}

	p := PingReqPacket{}
	err = p.Read(&buf)
	if err != nil {
		t.Errorf("could not read %s packet: %v", pingReqType.name, err)
	}
	fmt.Println("read that packet!", p)
}
