package packets

import (
	"bytes"
	"fmt"
	"testing"
)

var testUnsubackPacket = []byte{2, 0, 1, 0}

func TestUnsubackPacket(t *testing.T) {
	unsuback := bytes.NewBuffer(testUnsubackPacket)

	ua := UnsubackPacket{}
	err := ua.Read(unsuback)
	if err != nil {
		t.Errorf("could not read suback packet: %v\n", err)
	}
	fmt.Printf("unsuback packet: %s\n", ua.String())
}

func TestUnsubackPacket_Write(t *testing.T) {
	suback := bytes.NewBuffer(testUnsubackPacket)

	ua := UnsubackPacket{}
	err := ua.Read(suback)
	if err != nil {
		t.Errorf("could not read suback packet: %v\n", err)
	}
	fmt.Printf("read suback packet:  %s\n", &ua)

	// create packet
	var buf bytes.Buffer
	sp := CreateUnsubackPacket()
	err = sp.Write(&buf)
	if err != nil {
		t.Errorf("could not write suback packet: %v", err)
	}
	fmt.Printf("write suback packet: %s\n", &sp)

	// verify they match
	if !bytes.Equal(sp.MessageId, ua.MessageId) {
		t.Errorf("unsuback messageids don't match")
	}
	if sp.FixedHeader != ua.FixedHeader {
		t.Errorf("unsuback fixedheaders don't match")
	}
}
