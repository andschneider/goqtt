package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPingReqPacket(t *testing.T) {
	var buf bytes.Buffer
	var prWrite PingReqPacket

	prWrite.CreatePacket()
	err := prWrite.Write(&buf)
	if err != nil {
		t.Errorf("could not write %s packet %v", prWrite.name, err)
	}
	fmt.Printf("%s packet: %+v\n", prWrite.name, prWrite)
}
