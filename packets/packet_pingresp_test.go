package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPingRespPacket(t *testing.T) {
	var buf bytes.Buffer
	pr := CreatePingRespPacket()
	err := pr.Write(&buf, true)
	if err != nil {
		t.Errorf("could not write PingResp packet %v", err)
	}
	fmt.Printf("pingresp packet: %s\n", &pr)
}
