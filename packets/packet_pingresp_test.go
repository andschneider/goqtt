package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPingRespPacket(t *testing.T) {
	var buf bytes.Buffer
	var pr PingRespPacket

	pr.CreatePacket()
	err := pr.Write(&buf)
	if err != nil {
		t.Errorf("could not write %s packet %v", pr.Name, err)
	}
	fmt.Printf("%s packet: %+v\n", pr.Name, pr)
}
