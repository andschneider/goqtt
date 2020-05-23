package goqtt

import (
	"bytes"
	"fmt"
	"testing"
)

func TestConnackPacket(t *testing.T) {
	remainingLength := 2
	sessionPresent := 0
	returnCode := 0
	connack := bytes.NewBuffer([]byte{byte(remainingLength), byte(sessionPresent), byte(returnCode)})

	ca := ConnackPacket{}
	err := ca.ReadConnackPacket(connack)
	if err != nil {
		t.Errorf("could not read connack packet: %v\n", err)
	}
	fmt.Println(ca)
}
