package packets

import (
	"fmt"
	"os"
	"testing"
)

func TestPingReqPacket(t *testing.T) {
	pr := CreatePingReqPacket()
	err := pr.Write(os.Stdout, true)
	if err != nil {
		t.Errorf("could not write PingReq packet %v", err)
	}
	fmt.Println()
}
