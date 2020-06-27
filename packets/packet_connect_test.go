package packets

import (
	"bytes"
	"fmt"
	"testing"
)

func TestConnectPacket(t *testing.T) {
	var buf bytes.Buffer
	var cpWrite, cpRead ConnectPacket

	cpWrite.CreatePacket()
	err := cpWrite.Write(&buf)
	if err != nil {
		t.Errorf("could not write %s packet: %v\n", cpWrite.Name, err)
	}
	fmt.Printf("%s packet write: %+v\n", cpWrite.Name, cpWrite)

	// have to read in the type from the fixed header for Read to work
	packetType, err := decodeByte(&buf)
	if err != nil {
		t.Errorf("could not decode type from fixed header. got %v", packetType)
	}

	err = cpRead.Read(&buf)
	if err != nil {
		t.Errorf("could not read %s packet: %v", cpRead.Name, err)
	}
	fmt.Printf("%s packet read: %+v\n", cpRead.Name, cpRead)
}
