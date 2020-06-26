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

var testConnackPacket = []byte{byte(remainingLength), byte(sessionPresent), byte(returnCode)}

// Test reading in a connack packet
func TestConnackPacket_Read(t *testing.T) {
	var ca ConnackPacket
	connack := bytes.NewBuffer(testConnackPacket)

	err := ca.Read(connack)
	if err != nil {
		t.Errorf("could not read %s packet: %v\n", ca.name, err)
	}
	fmt.Printf("%s packet: %+v\n", ca.name, ca)
}

// Test writing a packet to a buffer
func TestConnackPacket_Write(t *testing.T) {
	var buf bytes.Buffer
	var ca ConnackPacket

	ca.CreatePacket()
	err := ca.Write(&buf)
	if err != nil {
		t.Errorf("could not write %s packet: %v\n", ca.name, err)
	}
	fmt.Printf("%s packet: %+v\n", ca.name, ca)
}

// Test comparing a default packet to one read in
func TestConnackPacket_Compare(t *testing.T) {
	var caRead, caDefault ConnackPacket
	connack := bytes.NewBuffer(testConnackPacket)

	// correct packet
	err := caRead.Read(connack)
	if err != nil {
		t.Errorf("could not read %s packet: %v\n", caRead.name, err)
	}

	// default packet
	caDefault.CreatePacket()

	if caRead != caDefault {
		t.Errorf("connack packets don't match")
	}
}
