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

func TestConnackPacket(t *testing.T) {
	connack := bytes.NewBuffer(testConnackPacket)

	ca := ConnackPacket{}
	err := ca.Read(connack)
	if err != nil {
		t.Errorf("could not read connack packet: %v\n", err)
	}
	fmt.Printf("connack packet: %s\n", &ca)
}

func TestConnackPacket_Write(t *testing.T) {
	connack := bytes.NewBuffer(testConnackPacket)

	// correct packet
	ca := ConnackPacket{}
	err := ca.Read(connack)
	if err != nil {
		t.Errorf("could not read connack packet: %v\n", err)
	}
	fmt.Printf("connack packet: %s\n", &ca)

	// create packet
	var buf bytes.Buffer
	cp := CreateConnackPacket()
	err = cp.Write(&buf)
	if err != nil {
		t.Errorf("could not write CONNACK packet: %v", err)
	}
	fmt.Printf("connack packet: %s\n", &cp)

	if cp != ca {
		t.Errorf("connack packets don't match")
	}
}
