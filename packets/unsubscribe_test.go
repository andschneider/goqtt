package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestUnsubscribePacket(t *testing.T) {
	var buf bytes.Buffer
	var upWrite UnsubscribePacket

	upWrite.CreateUnsubscribePacket(testTopic)
	err := upWrite.Write(&buf)
	if err != nil {
		t.Errorf("could not write %s packet %v", upWrite.name, err)
	}
	fmt.Printf("write %s packet: %+v\n", upWrite.name, upWrite)
}
