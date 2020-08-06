package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestConnectPacket(t *testing.T) {
	var buf bytes.Buffer
	var cpWrite ConnectPacket

	cpWrite.CreatePacket()
	err := cpWrite.Write(&buf)
	if err != nil {
		t.Errorf("could not write %s packet: %v\n", cpWrite.name, err)
	}
	fmt.Printf("%s packet write: %+v\n", cpWrite.name, cpWrite)
}
