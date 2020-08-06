package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestConnackPacket_Write(t *testing.T) {
	var buf bytes.Buffer
	var ca ConnackPacket

	ca.CreatePacket()
	err := ca.Write(&buf)
	if err != nil {
		t.Errorf("could not write %s packet: %v\n", ca.Name(), err)
	}
	fmt.Printf("%s packet: %+v\n", ca.Name(), ca)
}
