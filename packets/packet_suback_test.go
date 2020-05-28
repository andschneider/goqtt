package packets

import (
	"bytes"
	"fmt"
	"testing"
)

var testSubackPacket = []byte{3, 0, 1, 0}

func TestSubackPacket(t *testing.T) {
	suback := bytes.NewBuffer(testSubackPacket)

	sa := SubackPacket{}
	err := sa.ReadSubackPacket(suback)
	if err != nil {
		t.Errorf("could not read suback packet: %v\n", err)
	}
	fmt.Printf("suback packet: %s\n", &sa)
}

func TestSubackPacket_Write(t *testing.T) {
	suback := bytes.NewBuffer(testSubackPacket)

	sa := SubackPacket{}
	err := sa.ReadSubackPacket(suback)
	if err != nil {
		t.Errorf("could not read suback packet: %v\n", err)
	}
	fmt.Printf("read suback packet:  %s\n", &sa)

	// create packet
	var buf bytes.Buffer
	sp := CreateSubackPacket()
	err = sp.Write(&buf)
	if err != nil {
		t.Errorf("could not write suback packet: %v", err)
	}
	fmt.Printf("write suback packet: %s\n", &sp)

	// verify they match
	if !bytes.Equal(sp.MessageId, sa.MessageId) {
		t.Errorf("suback messageids don't match")
	}
	if !bytes.Equal(sp.ReturnCodes, sa.ReturnCodes) {
		t.Errorf("suback returncodes don't match")
	}
	if sp.FixedHeader != sa.FixedHeader {
		t.Errorf("suback fixedheaders don't match")
	}
}
