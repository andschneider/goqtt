package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPingReqPacket(t *testing.T) {
	var buf bytes.Buffer
	pr := CreatePingReqPacket()
	err := pr.Write(&buf, true)
	if err != nil {
		t.Errorf("could not write PingReq packet %v", err)
	}
	fmt.Printf("pingreq packet: %s\n", &pr)
}
