package goqtt

import (
	"fmt"
	"os"
	"testing"
)

func TestPingRespPacket(t *testing.T) {
	pr := CreatePingRespPacket()
	err := pr.Write(os.Stdout, true)
	if err != nil {
		t.Errorf("could not write PingReq packet %v", err)
	}
	fmt.Println()
}
