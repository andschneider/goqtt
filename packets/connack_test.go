package packets

import (
	"bytes"
	"fmt"
	"testing"
)

const (
	remainingLength = 2
	sessionPresent  = 0
	returnCode      = 0
)

// Test reading in a connack packet
func TestConnackPacket_Read(t *testing.T) {
	var ca ConnackPacket
	connack := bytes.NewBuffer(testConnackPacket)

	err := ca.Read(connack)
	if err != nil {
		t.Errorf("could not read %s packet: %v\n", ca.Name(), err)
	}
	fmt.Printf("%s packet: %+v\n", ca.Name(), ca)
}

// Test writing a packet to a buffer
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
